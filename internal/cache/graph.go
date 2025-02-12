package cache

import (
	"fmt"
	"log"
	"sync"

	"github.com/golrice/e-fis/internal/cache/flowcontrol"
	"github.com/golrice/e-fis/internal/peer"
	pb "github.com/golrice/e-fis/internal/protocal"
)

// we define a namespace
type Graph struct {
	mu      sync.RWMutex
	records map[string]*Node
}

func DefaultGraph() *Graph {
	return &Graph{
		mu:      sync.RWMutex{},
		records: map[string]*Node{},
	}
}

func (g *Graph) AddNode(node *Node) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.records[node.name] = node
}

type Getter interface {
	Get(key string) ([]byte, error)
}

// implement Getter all function which "looks like Get"
type GetterLikeFunc func(key string) ([]byte, error)

func (f GetterLikeFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// define a basic node
type Node struct {
	name          string
	cache         *cache
	getter        Getter
	peers         peer.PeerPicker
	flowcontroler *flowcontrol.Controler
}

func NewNode(name string, capacity int64, getter GetterLikeFunc) *Node {
	if getter == nil {
		panic("need a good getter")
	}

	node := &Node{
		name:          name,
		cache:         NewCache(capacity),
		getter:        getter,
		peers:         nil,
		flowcontroler: &flowcontrol.Controler{},
	}

	return node
}

func GetNode(graph *Graph, name string) (*Node, error) {
	graph.mu.RLock()
	defer graph.mu.RUnlock()

	if node, ok := graph.records[name]; ok {
		return node, nil
	}

	return nil, fmt.Errorf("no such a node")
}

func (n *Node) RegisterPeers(peers peer.PeerPicker) {
	if n.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	n.peers = peers
}

func (n *Node) Get(key string) (ByteView, error) {
	if key == "" {
		return NewByteView(nil), nil
	}

	if v, ok := n.cache.get(key); ok {
		return v, nil
	}

	// cache miss, fix it
	return n.load(key)
}

func (n *Node) load(key string) (value ByteView, err error) {
	// we load data from local or remote, it depends.
	v, err := n.flowcontroler.Do(key, func() (any, error) {
		if n.peers != nil {
			if peer, ok := n.peers.PickPeer(key); ok {
				if value, err = n.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[Cache] Failed to get from peer", err)
			}
		}

		return n.loadLocally(key)
	})

	if err == nil {
		return v.(ByteView), err
	}

	return
}

func (n *Node) getFromPeer(peer peer.PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		NodeName: n.name,
		Key:      key,
	}
	resp := &pb.Response{}
	err := peer.Get(req, resp)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: resp.Value}, nil
}

func (n *Node) loadLocally(key string) (ByteView, error) {
	vb, err := n.getter.Get(key)

	if err != nil {
		return NewByteView(nil), err
	}

	bv := ByteView{b: cloneBytes(vb)}

	// add new item in cache
	n.addCache(key, bv)

	return bv, nil
}

func (n *Node) addCache(key string, value ByteView) {
	n.cache.add(key, value)
}
