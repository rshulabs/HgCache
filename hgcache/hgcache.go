package hgcache

import (
	"fmt"
	"log"
	"sync"

	"github.com/rshulabs/HgCache/hgcache/rpc/pb"
	"github.com/rshulabs/HgCache/hgcache/singleflight"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	loader    *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("Getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	if v, ok := groups[name]; ok {
		return v
	}
	mu.RUnlock()
	return nil
}

func (g *Group) Get(key string) (BytesView, error) {
	if key == "" {
		return BytesView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[HgCache] hit")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value BytesView, err error) {
	view, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Panicln("[HgCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err != nil {
		return view.(BytesView), nil
	}
	return
}

func (g *Group) getLocally(key string) (BytesView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return BytesView{}, err
	}
	value := BytesView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value BytesView) {
	g.mainCache.add(key, value)
}

func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (BytesView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return BytesView{}, err
	}
	return BytesView{b: res.Value}, nil
}
