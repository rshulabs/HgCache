package lru

/*
lru淘汰算法
	淘汰最近访问最少的那次
		数据结构：map[string]dqueue
		节点被访问，放入队尾，淘汰队首
*/
import "container/list"

type Cache struct {
	maxBytes int64                    // 允许使用最大内存
	nbytes   int64                    // 当前已使用内存
	ll       *list.List               // 双向链表
	cache    map[string]*list.Element // map
	// optional and excuted when an entyr is purged.
	OnEvicted func(key string, value Value) // 移除记录时的回调函数
}

// 链表节点数据类型
type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int // 返回占用内存大小
}

func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

/**
 查找
	从字典找到对应链表节点
	将该节点移到队尾
*/

func (c *Cache) Get(key string) (value Value, ok bool) {
	// 先查后取
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

/*
 删除
	取队首，删除节点和对应map
	更新当前使用内存
	有回调，调用回调
*/
func (c *Cache) RemoveOldest() {
	// 取出队头
	ele := c.ll.Back()
	if ele != nil {
		// 移除队头
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		// 更新当前使用内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

/*
增加/修改
	查找key，存在，更新节点并移到队尾
	不存在，队尾新增节点，map添加记录
	更新当前使用内存，若超过最大内存，移除最少访问节点
*/
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

/*
添加了几条数据
*/
func (c *Cache) Len() int {
	return c.ll.Len()
}
