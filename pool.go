package simple_cache

import (
	"sync"
	"time"
)

type Pool struct {
	stop      chan struct{}
	waitGroup sync.WaitGroup
	rwMutex   sync.RWMutex
	items     map[string]Item
}

func NewPool(cleanupInterval time.Duration) *Pool {
	p := &Pool{
		items: make(map[string]Item),
		stop:  make(chan struct{}),
	}
	p.waitGroup.Add(1)
	go func(cleanupInterval time.Duration) {
		defer p.waitGroup.Done()
		p.cleanupLoop(cleanupInterval)
	}(cleanupInterval)
	return p
}

func (p *Pool) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-p.stop:
			return
		case <-ticker.C:
			p.rwMutex.Lock()
			for key, item := range p.items {
				if item.expiresAt <= time.Now().Unix() {
					delete(p.items, key)
				}
			}
			p.rwMutex.Unlock()
		}
	}
}

func (p *Pool) StopCleanup() {
	close(p.stop)
	p.waitGroup.Wait()
}

func (p *Pool) GetItem(key string) Item {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	item, ok := p.items[key]
	if !ok {
		item.key = key
		item.hit = false
	}
	return item
}

func (p *Pool) GetItems(keys []string) []Item {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	items := make([]Item, len(keys))
	for i, key := range keys {
		items[i] = p.GetItem(key)
	}
	return items
}

func (p *Pool) HasItem(key string) bool {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()
	_, ok := p.items[key]
	return ok
}

func (p *Pool) Clear() {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.items = make(map[string]Item)
}

func (p *Pool) DeleteItem(key string) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	delete(p.items, key)
}

func (p *Pool) DeleteItems(keys []string) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	for _, key := range keys {
		delete(p.items, key)
	}
}

func (p *Pool) Save(item Item) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()
	p.items[item.key] = item
}
