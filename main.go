package main

import (
	"flag"
	"fmt"
	"geecache/geecache"
	"log"
	"net/http"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup()*geecache.Group{
	return geecache.NewGroup("scores",2 <<10,geecache.GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDb] search key:", key)
		if v,ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not existed", key)
	}))
}

func startCacheServer(addr string, addrs []string, gee *geecache.Group){
	peers := geecache.NewHTTPPool(addr)
	peers.Set(addrs...)
	gee.RegisterPeers(peers)
	log.Println("geecache is running at ", addrs)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startApiServer(apiAddr string,gee *geecache.Group){
	http.Handle("/api", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		key := request.URL.Query().Get("key")
		view,err := gee.Get(key)
		if err != nil{
			http.Error(writer,err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Write(view.ByteSlice())
	}))
	log.Println("fronted server is running at ",apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main(){
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "GeeCache server port")
	flag.BoolVar(&api, "api", false, "Start a api server")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	var addrs []string
	for _, v := range addMap{
		addrs = append(addrs, v)
	}
	gee := createGroup()
	if api{
		go startApiServer(apiAddr, gee)
	}
	startCacheServer(addMap[port], addrs, gee)
}
