/*
 * @Author: cnzf1
 * @Date: 2023-03-12 19:38:39
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-13 23:54:00
 * @Description:
 */
package cache

import "github.com/cnzf1/gocore/lang"

type LRUCache struct {
	size     int
	capacity int
	cache    map[string]*CacheNode
	head     *CacheNode
	tail     *CacheNode
}

type CacheNode struct {
	key        string
	value      lang.AnyType
	prev, next *CacheNode
}

func NewLRUCache(capacity int) LRUCache {
	l := LRUCache{
		cache:    make(map[string]*CacheNode),
		capacity: capacity,
		head:     createCacheNode("", ""),
		tail:     createCacheNode("", ""),
	}
	l.head.next = l.tail
	l.tail.prev = l.head
	return l
}

func createCacheNode(key string, value lang.AnyType) *CacheNode {
	return &CacheNode{
		key:   key,
		value: value,
	}
}

func (this *LRUCache) incrSize(inc int) {
	this.size += inc
	if this.size > this.capacity {
		removed := this.removeTail()
		delete(this.cache, removed.key)
	}
}

func (this *LRUCache) addToHead(node *CacheNode) {
	node.prev = this.head
	node.next = this.head.next
	this.head.next.prev = node
	this.head.next = node
	this.incrSize(1)
}

func (this *LRUCache) removeNode(node *CacheNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
	this.incrSize(-1)
}

func (this *LRUCache) moveToHead(node *CacheNode) {
	this.removeNode(node)
	this.addToHead(node)
}

func (this *LRUCache) removeTail() *CacheNode {
	node := this.tail.prev
	this.removeNode(node)
	return node
}

func (this *LRUCache) Get(key string) lang.AnyType {
	if _, ok := this.cache[key]; !ok {
		return -1
	}

	node := this.cache[key]
	this.moveToHead(node)
	return node.value
}

func (this *LRUCache) Put(key string, value lang.AnyType) {
	if _, ok := this.cache[key]; !ok {
		node := createCacheNode(key, value)
		this.cache[key] = node
		this.addToHead(node)
	} else {
		node := this.cache[key]
		node.value = value
		this.moveToHead(node)
	}
}
