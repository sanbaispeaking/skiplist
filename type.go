package skiplist

import "sync"

var (
	DefaultMaxLevel int = 21
	// Probability for an element to be inserted into a higher level approximates to 0.25
	DefaultProbability int64 = 0x3fff // float(0x3fff) / float(0xffff) = 0.25
)

type elementNode struct {
	next []*Element
}

// Element stores underlying data and pointers to next elements of different levels
type Element struct {
	elementNode
	Value interface{}
	key   uint64
}

type SkipList struct {
	elementNode
	maxLevel int
	length   int
	// caches search path for insertion and removal
	searchPathCache []*elementNode
	mutex           sync.Mutex
}
