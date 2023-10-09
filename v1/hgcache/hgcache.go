package hgcache

import (
	"fmt"
	"log"
	"sync"

	"github.com/rshulabs/HgCache/v1/hgcache/rpc/pb"
	"github.com/rshulabs/HgCache/v1/hgcache/singleflight"
)

// 获取缓存数据接口，依赖注入，让用户自行实现获取方式
type Getter interface {
	Get(key string) ([]byte, error) // 根据key获取值
}

// 本地获取数据函数
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string              // 缓存组名字
	getter    Getter              // 本地获取数据
	mainCache cache               // 缓存
	peers     PeerPicker          // 远程节点获取数据
	loader    *singleflight.Group // 防止缓存击穿
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	// 判断getter是否存在
	if getter == nil {
		panic("Getter is nil")
	}
	// 获取锁
	mu.Lock()
	defer mu.Unlock()
	// 创建group实例
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	// 放入groups组
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	// 获取锁
	mu.RLock()
	// 从groups获取group
	if v, ok := groups[name]; ok {
		return v
	}
	mu.RUnlock()
	return nil
}

/*
主要函数

	获取缓存
	策略：1,2,3
	客户端访问-直接从本地缓存中取(本地缓存存在)
			 -(本地缓存不存在)是否从远程节点获取-(是)从远程获取-返回缓存
					                        -(否)回调函数getter，添加缓存-返回缓存
*/
func (g *Group) Get(key string) (BytesView, error) {
	// 验证key值
	if key == "" {
		return BytesView{}, fmt.Errorf("key is required")
	}
	// 先从本地缓存获取，策略1
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[HgCache] hit")
		return v, nil
	}
	// 本地没有，采用策略2,3
	return g.load(key)
}

func (g *Group) load(key string) (value BytesView, err error) {
	// 避免缓存击穿，同一时间内大量相同访问，导致服务崩溃，限制访问 10000-1
	view, err := g.loader.Do(key, func() (interface{}, error) {
		// 判断是否有采取远程调用
		if g.peers != nil {
			// 根据key得到远程节点调用方法
			if peer, ok := g.peers.PickPeer(key); ok {
				// 判断远程节点是否有数据
				if value, err := g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				
				log.Println("[HgCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return view.(BytesView), err
	}
	return
}

func (g *Group) getLocally(key string) (BytesView, error) {
	// 从本地获取数据
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
