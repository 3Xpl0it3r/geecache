package geecache

import (
	"fmt"
	"geecache/geecache/singlefight"
	"log"
	"sync"
)

// Getter loads data for a key
type Getter interface {
	Get(key string)([]byte, error)
}

// GetterFun implements Getter with a function
type GetterFunc func(key string)([]byte, error)

func (f GetterFunc)Get(key string)([]byte, error){
	return f(key)
}


type Group struct {
	name string
	getter Getter
	mainCache cache
	peers PeerPicker
	//
	loader *singlefight.Group
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter)*Group{
	if getter == nil{
		panic("getter is nil")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
		loader: &singlefight.Group{},
	}
	groups[name] = g
	return g
}
func GetGroup(name string)*Group{
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

func(g *Group)RegisterPeers(peers PeerPicker){
	if g.peers != nil{
		panic("RegisterPeerPicker  called more than one")
	}
	g.peers = peers
}

func (g *Group)Get(key string)(ByteView, error){
	if key == ""{
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v,ok := g.mainCache.get(key);ok {
		log.Println("[Gcache hit]")
		return v, nil
	}
	return g.load(key)
}


func(g *Group)load(key string)(value ByteView, err error){
	viewi,err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil{
			if peer,ok := g.peers.PickPeer(key);ok {
				if value,err = g.getFromPeer(peer, key); err == nil{
					return value, nil
				}
				log.Println("[GeeCache] Failed to get from peer", err)
			}
		}

		return g.getLocally(key)
	})
	if err == nil{
		return viewi.(ByteView), nil
	}
	return
}



func(g *Group)getLocally(key string)(value ByteView, err error){
	bytes,err := g.getter.Get(key)
	if err != nil{
		return ByteView{}, err
	}
	value = ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return  value ,nil
}

func(g *Group)getFromPeer(peer PeerGetter, key string)(ByteView, error){
	bytes,err := peer.Get(g.name, key)
	if err != nil{
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

func(g *Group)populateCache(key string, value ByteView){
	g.mainCache.add(key, value)
}