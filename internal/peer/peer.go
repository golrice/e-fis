package peer

import pb "github.com/golrice/e-fis/internal/protocal"

// we can use PickPeer function to get the peergetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// peergetter function can return the value according to the key
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
