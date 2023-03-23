/*
 * @Author: cnzf1
 * @Date: 2023-03-16 01:21:31
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-03-16 01:21:44
 * @Description: 字典树
 */

package tree

const MAXCAP = 26 // a-z

type Trie struct {
	next   map[rune]*Trie
	isWord bool
}

/** Initialize your data structure here. */
func NewTrie() Trie {
	root := new(Trie)
	root.next = make(map[rune]*Trie, MAXCAP)
	root.isWord = false
	return *root
}

/** Inserts a word into the trie. */
func (this *Trie) Insert(word string) {
	for _, v := range word {
		if this.next[v] == nil {
			node := new(Trie)
			//子节点数量为26
			node.next = make(map[rune]*Trie, MAXCAP)
			//初始化节点单词标志为假
			node.isWord = false
			this.next[v] = node
		}
		this = this.next[v]
	}
	this.isWord = true
}

/** Returns if the word is in the trie. */
func (this *Trie) Search(word string) bool {
	for _, v := range word {
		if this.next[v] == nil {
			return false
		}
		this = this.next[v]
	}
	return this.isWord
}

/** Returns if there is any word in the trie that starts with the given prefix. */
func (this *Trie) StartsWith(prefix string) bool {
	for _, v := range prefix {
		if this.next[v] == nil {
			return false
		}
		this = this.next[v]
	}
	return true
}
