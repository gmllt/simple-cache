package simple_cache

import (
	"time"
)

type Item struct {
	key       string `json:"key"`       // The key of the item
	data      any    `json:"data"`      // raw data
	expiresAt int64  `json:"expire_at"` // unix timestamp
	hit       bool   // true if the item was found in the pool
}

func (i *Item) GetKey() string {
	return i.key
}

func (i *Item) Set(data any) {
	i.data = data
}

func (i *Item) Get() any {
	return i.data
}

func (i *Item) IsHit() bool {
	return i.hit
}

func (i *Item) ExpiresAt(expiration int64) {
	i.expiresAt = expiration
}

func (i *Item) ExpiresAfter(ttl int64) {
	i.expiresAt = ttl + time.Now().Unix()
}
