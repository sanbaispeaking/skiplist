package skiplist

import (
	"sync"
	"testing"
)

func sancheck(list *SkipList, t *testing.T) {
	for k, element := range list.next {
		if element == nil {
			continue
		}
		t.Logf("Level %d, element level %d\n", k, len(element.next))
		if k > len(element.next) {

			t.Fatal("1st element level < current level")
		}

		var cur *Element = element
		// element count of current level
		var count int = 1
		// Each level should be ordered (cur.key < next.key)
		for cur.next[k] != nil {
			if cur.key > cur.next[k].key {
				t.Fatalf("next key < cur key (%v < %v)", cur.next[k].key, cur.key)
			}
			if k > len(cur.next) {
				t.Fatalf("element level < current level (%v < %v)", k, cur.next)
			}
			cur = cur.next[k]
			count++
		}

		if k == 0 {
			if count != list.length {
				t.Fatalf("list length != level 0 node count (%v != %v)", list.length, count)
			}
		}
	}
}

func TestBasicCRUD(t *testing.T) {
	var list *SkipList = New()

	kv := []struct {
		key uint64
		val int
	}{{10, 1}, {60, 2}, {30, 3}, {20, 4}, {90, 5}}

	// insert some int elements
	for _, item := range kv {
		list.Set(item.key, item.val)
	}
	sancheck(list, t)
	for _, item := range kv {
		element := list.Get(item.key)
		if element == nil {
			t.Fatalf("element(key=%d) not found", item.key)
		}
		if val := element.Value.(int); val != item.val {
			t.Fatalf("expect (key=%d) val :%d != got :%d", item.key, item.val, val)
		}
	}

	// update value of an existing key
	list.Set(30, 9)
	sancheck(list, t)
	element := list.Get(30)
	if element.Value.(int) != 9 {
		t.Fatal("element value not updated")
	}

	// remove existing key
	list.Remove(20)
	// key does not exist
	list.Remove(11)
	sancheck(list, t)

	keyNotExist := []uint64{11, 20}
	for _, key := range keyNotExist {
		element := list.Get(key)
		if element != nil {
			t.Fatalf("element(key=%d) should not exist", key)
		}
	}

}

func TestConcurrency(t *testing.T) {
	var list *SkipList = New()

	var wg sync.WaitGroup
	wg.Add(3)

	repeat := 10000
	go func() {
		for i := 0; i < repeat; i++ {
			list.Set(uint64(i), i)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < repeat; i++ {
			list.Get(uint64(i))
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < repeat; i++ {
			list.Get(uint64(i))
		}
		wg.Done()
	}()

	wg.Wait()
	if list.length != repeat {
		t.Fail()
	}
}
