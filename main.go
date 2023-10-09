package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/rshulabs/HgCache/v1/hgcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup(name string) *hgcache.Group {
	return hgcache.NewGroup(name, 2<<10, hgcache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
}

func startCacheServer(addr string, addrs []string, g *hgcache.Group) {
	peers := hgcache.NewHttpPool(addr)
	peers.Set(addrs...)
	g.RegisterPeers(peers)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startApiServer(apiAddr string, g *hgcache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, err := g.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())
		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "hgcache server port")
	flag.BoolVar(&api, "api", false, "start hgcache")
	flag.Parse()
	apiAddr := "http://0.0.0.0:8791"
	addrMap := map[int]string{
		8001: "http://0.0.0.0:8001",
		8002: "http://0.0.0.0:8002",
		8003: "http://0.0.0.0:8003",
	}
	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}
	g := createGroup("test")
	if api {
		go startApiServer(apiAddr, g)
	}
	startCacheServer(addrMap[port], addrs, g)
}
