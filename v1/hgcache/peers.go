package hgcache

import "github.com/rshulabs/HgCache/v1/hgcache/rpc/pb"

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
