package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T){
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})
	except := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, except){
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T){
	loadCounts := make(map[string]int, len(db))
	gee := NewGroup("score", 2 << 10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v,ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] = 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not existd", key)
	}))
	for k,v := range db{
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of Tom")
		}
		if _,err := gee.Get(k); err != nil|| loadCounts[k] > 1 {
			t.Fatalf("cache %s missing", k)
		}
	}
	if view,err := gee.Get("unknow");err == nil{
		t.Fatalf("the value of unknow should en empty, but %s got", view.String())
	}
}