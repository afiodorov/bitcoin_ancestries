package main

import "container/list"

// RecentAncestry is LRU for recent ancestries. This way if certain blocks keep relying on the same
// ancestry, we don't need to keep fetching them
type RecentAncestry struct {
	Capacity int
	Size     int
	Cache    map[string]*list.Element
	Lru      *list.List
}

// NewRecentAncestry implements LRU for most recently needed ancestries
func NewRecentAncestry(capacity int) *RecentAncestry {
	ra := RecentAncestry{
		Capacity: capacity,
		Size:     0,
		Lru:      list.New(),
		Cache:    make(map[string]*list.Element),
	}

	return &ra
}

// Set add a value to cache & removes value that wasn't used
func (rc *RecentAncestry) Set(a *Ancestry) {
	if elem, present := rc.Cache[a.Hash]; present {
		rc.Lru.MoveToFront(elem)
	} else {
		elem := rc.Lru.PushFront(a)
		rc.Size++
		rc.Cache[a.Hash] = elem

		if rc.Size > rc.Capacity {
			lruItem := rc.Lru.Back()
			rc.Lru.Remove(lruItem)
			rc.Size--
			delete(rc.Cache, lruItem.Value.(*Ancestry).Hash)
		}
	}
}

// Get gets value and makes sure it's at the front of the LRU doubly-linked list
func (rc *RecentAncestry) Get(hash string) (*Ancestry, bool) {
	if elem, present := rc.Cache[hash]; present {
		rc.Lru.MoveToFront(elem)
		return elem.Value.(*Ancestry), true
	}
	return nil, false
}

// LookUp is immutalbe Get
func (rc RecentAncestry) LookUp(hash string) (*Ancestry, bool) {
	if elem, present := rc.Cache[hash]; present {
		return elem.Value.(*Ancestry), true
	}

	return nil, false
}
