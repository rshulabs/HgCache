package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // 哈希函数，默认crc32.ChecksumIEEE
	replicas int            // 虚拟节点数
	keys     []int          // 哈希环
	hashMap  map[int]string // 节点映射 虚拟-物理
}

func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 哈希环添加节点
// 根据物理节点建立映射
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 对哈希环排序，便于用二分搜索
	sort.Ints(m.keys)
}

// 获取对应物理节点
func (m *Map) Get(key string) string {
	// 验证客户端ip
	if len(key) == 0 {
		return ""
	}
	// 对客户端ip哈希
	hash := int(m.hash([]byte(key)))
	// 找到对应虚拟节点索引
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 根据索引找到命中缓存真实节点
	return m.hashMap[m.keys[idx%len(m.keys)]] // 哈希环状，取模
}
