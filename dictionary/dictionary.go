// Package dictionary implements a hash table. The hash table accepts strings,
// uint16 values and bytes as keys, which are the data types used by the
// compression that require the dictionary.
package dictionary

import (
	"github.com/mjjs/gompressor/linkedlist"
	"github.com/mjjs/gompressor/vector"
)

const defaultSize uint = 32

type dictionaryNode struct {
	key   interface{}
	value interface{}
}

// Dictionary implements a hashtable.
type Dictionary struct {
	buckets []*linkedlist.LinkedList
	size    int
}

// New returns a pointer to a new Dictionary with an initial size.
func New() *Dictionary {
	return NewWithSize(defaultSize)
}

// NewWithSize returns a pointer to a new Dictionary with the given size. If
// zero is passed in as size, the default size is used.
func NewWithSize(size uint) *Dictionary {
	buckets := make([]*linkedlist.LinkedList, size)
	for i := range buckets {
		buckets[i] = new(linkedlist.LinkedList)
	}

	return &Dictionary{buckets: buckets}
}

// Set maps key to value in the dictionary. If key is present in the map, the
// value is updated.
func (d *Dictionary) Set(key interface{}, value interface{}) {
	bucket := d.getBucket(key)

	updated := false

	bucket.ForEach(func(iNode interface{}) {
		if updated {
			return
		}

		node := iNode.(*dictionaryNode)
		if node.key == key {
			node.value = value
			updated = true
		}
	})

	if updated {
		return
	}

	bucket.Append(&dictionaryNode{key: key, value: value})
	d.size++

	if float32(d.size/len(d.buckets)) > 0.75 {
		d.grow()
	}
}

// Get returns the value associated with key. An additional boolean value is
// returned to indicate whether or not the key exists in the map.
func (d *Dictionary) Get(key interface{}) (interface{}, bool) {
	bucket := d.getBucket(key)

	var value interface{} = nil
	found := false

	bucket.ForEach(func(iNode interface{}) {
		if found {
			return
		}

		node := iNode.(*dictionaryNode)
		if node.key == key {
			value = node.value
			found = true
		}
	})

	return value, found
}

// Remove removes the the given key-value pair from the dictionary.
func (d *Dictionary) Remove(key interface{}) {
	bucket := d.getBucket(key)

	removed := false

	bucket.ForEach(func(iNode interface{}) {
		if removed {
			return
		}

		node := iNode.(*dictionaryNode)
		if node.key == key {
			removed = true
			bucket.Remove(iNode)
		}
	})

	d.size--
}

// Size returns the amount of unique values present in the dictionary.
func (d *Dictionary) Size() int {
	return d.size
}

// Keys returns a vector containing all the keys in the dictionary.
func (d *Dictionary) Keys() *vector.Vector {
	keys := vector.New()

	for _, bucket := range d.buckets {
		bucket.ForEach(func(node interface{}) {
			keys.Append(node.(*dictionaryNode).key)
		})
	}

	return keys
}

func (d *Dictionary) getBucket(key interface{}) *linkedlist.LinkedList {
	hash := hash(key)
	n := int64(len(d.buckets))

	return d.buckets[((hash%n)+n)%n]
}

func (d *Dictionary) grow() {
	newBuckets := make([]*linkedlist.LinkedList, d.size*2)
	for i := range newBuckets {
		newBuckets[i] = new(linkedlist.LinkedList)
	}

	n := int64(len(newBuckets))

	for _, bucket := range d.buckets {
		bucket.ForEach(func(node interface{}) {
			hash := hash((node.(*dictionaryNode).key))
			newBucket := newBuckets[((hash%n)+n)%n]
			newBucket.Append(node)
		})
	}

	d.buckets = newBuckets
}

func hash(key interface{}) int64 {
	switch v := key.(type) {
	case string:
		const (
			prime  int64 = 31
			modulo int64 = 1e9 + 9
		)

		var (
			hash       int64 = 0
			primePower int64 = 1
		)

		for _, c := range v {
			hash = (hash + (int64(c-'a')+1)*primePower) % modulo
			primePower = (primePower * prime) % modulo
		}

		return hash
	case uint16:
		return int64(v)
	case byte:
		return int64(v)
	default:
		panic("Unsupported key type")
	}
}
