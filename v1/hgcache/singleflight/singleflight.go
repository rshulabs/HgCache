package singleflight

import (
	"sync"
)

type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	// 获取锁
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	// 如果map有值，wait第一个请求完成，再一起返回第一个map值
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		// wait控制
		c.wg.Wait()
		return c.val, c.err
	}
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	// 第一次请求赋完值，释放锁
	g.mu.Unlock()
	// 返回请求结果
	c.val, c.err = fn()
	// 请求完成
	c.wg.Done()
	// 再加锁，删除存在的key
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}
