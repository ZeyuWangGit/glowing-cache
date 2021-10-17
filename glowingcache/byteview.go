package glowingcache

// A ByteView holds an immutable view of bytes.
type ByteView struct {
	byteView []byte
}

func (v ByteView) Len() int {
	 return len(v.byteView)
}

func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.byteView)
}

func (v ByteView) String() string {
	return string(v.byteView)
}

func cloneBytes(b []byte) []byte  {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}