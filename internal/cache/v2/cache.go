package v2

import (
	"sync"

	"github.com/rshulabs/HgCache/v1/hgcache/lru"
)

type cache struct {
	mu         sync.Mutex // 锁，并发
	lru        *lru.Cache // 缓存
	cacheBytes int64      // 内存容量
}

func (c *cache) add(key string, value BytesView) {
	// 获取锁
	c.mu.Lock()
	defer c.mu.Unlock()
	// 判断缓存是否创建
	if c.lru == nil {
		c.lru = lru.NewCache(c.cacheBytes, nil)
	}
	// 将值插入缓存
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value BytesView, ok bool) {
	// 获取锁
	c.mu.Lock()
	defer c.mu.Unlock()
	// 判断缓存是否创建
	if c.lru == nil {
		return
	}
	// 查询缓存是否存在该值
	if v, ok := c.lru.Get(key); ok {
		return v.(BytesView), ok
	}
	return
}
