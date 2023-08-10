package hgcache

/*
只读缓存值
*/

type BytesView struct {
	b []byte
}

func (v BytesView) Len() int { 
	return len(v.b)
}

// copy
func (v BytesView) ByteSlice() []byte{
	return cloneBytes(v.b)
}

func (v BytesView) String() string {
	return string(v.b)
}

func cloneBytes(v []byte) []byte {
	b := make([]byte,len(v))
	copy(b,v)
	return b
}