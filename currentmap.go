package fxk

import (
	"sync"
)

var shardSize = 32

type Cmap []*CurrentMapShard
type CurrentMapShard struct {
	items map[string]interface{}
	sync.RWMutex
}

func NewSyncMap() Cmap {
	m := make(Cmap, shardSize)
	for i := 0; i < shardSize; i++ {
		m[i] = &CurrentMapShard{items: make(map[string]interface{})}
	}
	return m
}

func (m Cmap) GetShard(key string) *CurrentMapShard {
	return m[uint(fnv(key))%uint(shardSize)]
}

func (m Cmap) Set(key string, value interface{}) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

func (m Cmap) Get(key string) (interface{}, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (m Cmap) Count() int {
	count := 0
	for i := 0; i < shardSize; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

func (m Cmap) Remove(key string) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

type Tuple struct {
	Key string
	Val interface{}
}

func snapshot(m Cmap) []chan Tuple {
	chans := make([]chan Tuple, shardSize)
	wg := sync.WaitGroup{}
	wg.Add(shardSize)
	for index, shard := range m {
		go func(index int, shard *CurrentMapShard) {
			shard.RLock()
			chans[index] = make(chan Tuple, len(shard.items))
			wg.Done()
			for key, val := range shard.items {
				chans[index] <- Tuple{Key: key, Val: val}
			}
			shard.RUnlock()
			close(chans[index])
		}(index, shard)
	}
	wg.Wait()
	return chans
}

func fanIn(chans []chan Tuple, out chan Tuple) {
	wg := sync.WaitGroup{}
	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch chan Tuple) {
			for t := range ch {
				out <- t
			}
			wg.Done()
		}(ch)
	}
	wg.Wait()
	close(out)
}

func (m Cmap) IterBuffered() <-chan Tuple {
	chans := snapshot(m)
	total := 0
	for _, c := range chans {
		total += cap(c)
	}
	ch := make(chan Tuple, total)
	go fanIn(chans, ch)
	return ch
}

func (m Cmap) Clear() {
	for item := range m.IterBuffered() {
		m.Remove(item.Key)
	}
}

func fnv(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
