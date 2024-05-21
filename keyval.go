package main

import (
	"sync"
)

type KV struct {
	data map[string][]byte

	mu sync.RWMutex
}

func NewKV() *KV {
	return &KV{data: make(map[string][]byte)}
}

func (kv *KV) Set(key, value string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = []byte(value)

	return nil
}

func (kv *KV) Get(key string) ([]byte, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	value, ok := kv.data[key]
	return value, ok
}
