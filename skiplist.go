package skiplist

import "math/rand"

// New create a new skip list with the default max level.
func New() *SkipList {
	return &SkipList{
		elementNode:     elementNode{make([]*Element, DefaultMaxLevel)},
		maxLevel:        DefaultMaxLevel,
		searchPathCache: make([]*elementNode, DefaultMaxLevel),
	}

}

// randLevel generates a random level in [0, maxLevel) for new element.
func (list *SkipList) randLevel() (level int) {

	for ((rand.Int63()>>32)&0xffff < DefaultProbability) && (level < list.maxLevel) {
		level++
	}
	return level
}

// Get searches for an element by key and returns a pointer to it.
// Returns nil if not found.
func (list *SkipList) Get(key uint64) *Element {
	var prev *elementNode = &list.elementNode
	var next *Element

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]
		// search in level i
		for next != nil && next.key < key {
			prev = &next.elementNode
			next = next.next[i]
		}
	}

	if next != nil && next.key == key {
		return next
	}
	return nil
}

func (list *SkipList) searchPath(key uint64) []*elementNode {
	var prev *elementNode = &list.elementNode
	var next *Element

	path := list.searchPathCache

	for i := list.maxLevel - 1; i >= 0; i-- {
		next = prev.next[i]
		for next != nil && next.key < key {
			prev = &next.elementNode
			next = next.next[i]
		}
		path[i] = prev
	}
	return path
}

// Set updates the value of the element specified by key.
// If the key does not exist, insert a new node into the list ordered by key.
// Returns a pointer to the element node
func (list *SkipList) Set(key uint64, value interface{}) (el *Element) {
	path := list.searchPath(key)

	if el = path[0].next[0]; el != nil && el.key == key {
		el.Value = value
		return
	}

	// key not found
	el = &Element{
		elementNode: elementNode{
			next: make([]*Element, list.randLevel()),
		},
		key:   key,
		Value: value,
	}

	for i := range el.next {
		el.next[i] = path[i].next[i]
		path[i].next[i] = el
	}

	list.length++
	return
}

// Remove delete the element specified by key from the list.
// Returns a pointer to that element, or nil if key not in the list.
func (list *SkipList) Remove(key uint64) (el *Element) {
	path := list.searchPath(key)

	if el = path[0].next[0]; el != nil && el.key == key {
		for level, element := range el.next {
			path[level].next[level] = element
		}
		list.length--
		return el
	}
	return
}
