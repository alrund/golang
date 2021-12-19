package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	key   Key
	value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := l.items[key]; ok {
		item.Value = cacheItem{key, value}
		l.queue.MoveToFront(item)
		return true
	}

	l.items[key] = l.queue.PushFront(cacheItem{key, value})
	l.checkCapacity()

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := l.items[key]; ok {
		itemValue := item.Value.(cacheItem)
		l.queue.MoveToFront(item)
		return itemValue.value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

func (l *lruCache) checkCapacity() {
	if l.queue.Len() <= l.capacity {
		return
	}

	lastItem := l.queue.Back()
	lastItemValue := lastItem.Value.(cacheItem)

	l.queue.Remove(lastItem)
	delete(l.items, lastItemValue.key)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
