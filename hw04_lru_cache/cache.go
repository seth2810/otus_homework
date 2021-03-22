package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   string
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()

	defer c.Unlock()

	item, exist := c.items[key]

	// I don't really understand why we use cacheItem type
	cachedValue := &cacheItem{key: string(key), value: value}

	if exist {
		item.Value = cachedValue
		c.queue.MoveToFront(item)
		return true
	}

	c.queue.PushFront(cachedValue)
	c.items[key] = c.queue.Front()

	if c.queue.Len() > c.capacity {
		back := c.queue.Back()
		c.queue.Remove(back)
		delete(c.items, Key(back.Value.(*cacheItem).key))
	}

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()

	defer c.Unlock()

	item, exist := c.items[key]

	if !exist {
		return nil, false
	}

	c.queue.MoveToFront(item)

	cachedItem := item.Value.(*cacheItem)

	return cachedItem.value, true
}

func (c *lruCache) Clear() {
	c.Lock()

	defer c.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	c := &lruCache{
		capacity: capacity,
	}

	c.Clear()

	return c
}
