package common

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// 2^32-1
type units []uint32

func (x units) Len() int {
	return len(x)
}

func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("Hash is empty")

type Consistent struct {
	// hash circle, key is hash value, value is the name of the node
	circle map[uint32]string
	// sort the hash circle
	sortedHashes units
	// virtual node number
	VirtualNode int
	//lock
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		circle:      make(map[uint32]string),
		VirtualNode: 20,
	}
}

// generate key
func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var srcatch [64]byte
		copy(srcatch[:], key)
		return crc32.ChecksumIEEE(srcatch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

func (c *Consistent) updateSortedHash() {
	hashes := c.sortedHashes[:0]
	if cap(c.sortedHashes)/(c.VirtualNode) > len(c.circle) {
		hashes = nil
	}

	for k := range c.circle {
		hashes = append(hashes, k)
	}

	sort.Sort(hashes)

	c.sortedHashes = hashes
}

func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	c.updateSortedHash()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHash()
}

func (c *Consistent) search(key uint32) int {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	key := c.hashKey(name)
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
