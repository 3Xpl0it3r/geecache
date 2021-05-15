package lru

import "container/list"

// Lru 算法，
type Cache struct {
	maxBytes int64		// 所允许最大内存
	nbytes int64		// 已经使用了的内存
	ll *list.List
	cache map[string]*list.Element

	OnEvicted func(key string, value Value)
}

type entry struct {
	key string
	value Value
}


// Value use Len to count how many bytes it takes
type Value interface {
	Len()int
}


func New(maxBytes int64, onEvicted func(string, Value))*Cache{
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}


// 查找功能，如果查找到说明命中，则将元素移动到队列头部
func(c *Cache)Get(key string)(value Value,ok bool){
	if value,ok := c.cache[key]; ok {
		c.ll.MoveToFront(value)
		kv := value.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

// 删除功能，移除访问最小的节点
func(c *Cache)RemoveOldest(){
	ele := c.ll.Back()
	if ele!= nil{
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil{
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 新增/修改
func(c *Cache)Add(key string,value Value){
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	if c.maxBytes != 0 && c.maxBytes < c.nbytes{
		c.RemoveOldest()
	}
}

func(c *Cache)Len()int{
	return c.ll.Len()
}