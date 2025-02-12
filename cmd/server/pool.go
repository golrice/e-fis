package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/golrice/e-fis/internal/cache"
	"github.com/golrice/e-fis/internal/consistenthash"
	"github.com/golrice/e-fis/internal/peer"
	pb "github.com/golrice/e-fis/internal/protocal"
	"google.golang.org/protobuf/proto"
)

const defaultReplicas = 50

type HttpPool struct {
	info  HttpInfo
	graph *cache.Graph

	mu          sync.Mutex
	peers       *consistenthash.DHTMap
	httpGetters map[string]*peer.HttpGetter
}

func NewHttpPool(addr string) *HttpPool {
	return &HttpPool{
		info:        *NewHttpInfo(addr),
		graph:       cache.DefaultGraph(),
		mu:          sync.Mutex{},
		peers:       nil,
		httpGetters: nil,
	}
}

func (p *HttpPool) Log(format string, v ...any) {
	log.Printf("[Server %s] %s", p.info.addr, fmt.Sprintf(format, v...))
}

func (p *HttpPool) NewNode(name string, capacity int64, getter cache.GetterLikeFunc) *cache.Node {
	node := cache.NewNode(name, capacity, getter)
	p.graph.AddNode(node)

	return node
}

func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check whether it is a valid request
	if !strings.HasPrefix(r.URL.Path, p.info.basePath) {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	p.Log("%s %s", r.Method, r.URL.Path)

	// path -> <base>/<node_name>/<key>
	s := strings.SplitN(r.URL.Path[len(p.info.basePath):], "/", 2)

	if len(s) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	node_name, key := s[0], s[1]

	node, err := cache.GetNode(p.graph, node_name)
	if err != nil {
		http.Error(w, "no such node", http.StatusNotFound)
		return
	}

	v, err := node.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := proto.Marshal(&pb.Response{Value: v.ByteSlice()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(body); err != nil {
		p.Log("fail to write to the requester, err: %s", err.Error())
	}
}

func (p *HttpPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*peer.HttpGetter, len(peers))

	for _, eachPeer := range peers {
		p.httpGetters[eachPeer] = &peer.HttpGetter{BaseURL: eachPeer + p.info.basePath}
	}
}

func (p *HttpPool) PickPeer(key string) (peer.PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if target := p.peers.Get(key); target != "" && target != p.info.addr {
		p.Log("Pick peer %s", target)
		return p.httpGetters[target], true
	}

	return nil, false
}

var _ peer.PeerPicker = (*HttpPool)(nil)
