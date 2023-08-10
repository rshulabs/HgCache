package lru

import (
	"reflect"
	"testing"
)

type Str string

func (s Str) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := NewCache(int64(0), nil)
	lru.Add("key1", Str("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(Str)) != "1234" {
		t.Fatal("cache hit key1 = 1234 failes")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatal("cache miss key2 failed")
	}
}

func TestRemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "key3"
	v1, v2, v3 := "value1", "value2", "value3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewCache(int64(cap), nil)
	lru.Add(k1, Str(v1))
	lru.Add(k2, Str(v2))
	lru.Add(k3, Str(v3))
	if _, ok := lru.Get("key1"); ok || lru.Len() != 2 {
		t.Fatal("remove failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, val Value) {
		keys = append(keys, key)
	}
	lru := NewCache(int64(10), callback)
	lru.Add("k1", Str("12345"))
	lru.Add("k2", Str("ere"))
	lru.Add("k3", Str("k3"))
	lru.Add("k4", Str("k4"))
	expect := []string{"k1", "k2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatal("callback failed")
	}
}
