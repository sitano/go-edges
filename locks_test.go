package main

import "testing"
import "sync"
import "sync/atomic"

func Test_RLockReentrancy(t *testing.T) {
	rw := sync.RWMutex{}
	rw.RLock()
	// if ask WLock() by other thread in this point, it will deadlock
	rw.RLock()
}

type SVal *int64

type Config struct {
	sync.RWMutex
	endpoint string
}

func BenchmarkPMutexSet(b *testing.B) {
	config := Config{}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.Lock()
			config.endpoint = "api.example.com"
			config.Unlock()
		}
	})
}

func BenchmarkPMutexGet(b *testing.B) {
	config := Config{endpoint: "api.example.com"}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.RLock()
			_ = config.endpoint
			config.RUnlock()
		}
	})
}

func BenchmarkPAtomicSet(b *testing.B) {
	var config atomic.Value
	c := Config{endpoint: "api.example.com"}
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			config.Store(c)
		}
	})
}

func BenchmarkPAtomicGet(b *testing.B) {
	var config atomic.Value
	config.Store(Config{endpoint: "api.example.com"})
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = config.Load().(Config)
		}
	})
}
