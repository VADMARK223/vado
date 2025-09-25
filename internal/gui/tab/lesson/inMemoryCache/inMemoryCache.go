package inMemoryCache

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Cache interface {
	Set(key string, value string)
	Get(key string) (string, bool)
}

type InMemoryCache struct {
	shards []Shard
}

type Shard struct {
	data map[string]string
	mtx  sync.RWMutex
}

func (i *InMemoryCache) Set(key string, value string) {
	shardID := hasher(key) % len(i.shards)
	i.shards[shardID].Set(key, value)
}

func (i *InMemoryCache) Get(key string) (string, bool) {
	shardID := hasher(key) % len(i.shards)
	return i.shards[shardID].Get(key)
}

func (s *Shard) Set(key string, value string) {
	s.data[key] = value
}

func (s *Shard) Get(key string) (string, bool) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	data, ok := s.data[key] // Под капотом есть перегрузка, значит явно надо заводить переменные
	return data, ok
}

func NewInMemoryCache(shardsCount int) *InMemoryCache {
	shards := make([]Shard, 0, shardsCount)
	for i := 0; i < shardsCount; i++ {
		shards = append(shards, Shard{data: make(map[string]string)})
	}

	return &InMemoryCache{shards: shards}
}

func RunInMemoryCache() {
	fmt.Println("RunInMemoryCache")
	cache := NewInMemoryCache(5) // Работаем с указателей, потому что есть поле RWMutex, которое копировать нельзя (оно содержит внутренние поля состояния). При копировании мьютекса получаются две независимые копии, которые "думают", что они синхронизируют доступ к одной и той же map.
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		wg.Done()
		cache.Set("name", "vadim")
	}()

	go func() {
		wg.Done()
		cache.Set("surname", "markitanov")
	}()
	wg.Wait()

	fmt.Println("Cache", cache.shards)
}

func hasher(key string) int {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return int(h.Sum32())
}
