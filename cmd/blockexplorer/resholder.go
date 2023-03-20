package main

import (
	"log"

	"github.com/go-pg/pg"
)

const syncStep = 100000
const cacheSize = 22000
const fetchStep = 1

type ResHolder struct {
	m       map[string]Ancestry
	fetched map[string]Ancestry
	cache   *RecentAncestry
}

func (r *ResHolder) Clear() {
	r.fetched = make(map[string]Ancestry, len(r.m))
	for k, v := range r.m {
		r.fetched[k] = v
	}
	r.m = make(map[string]Ancestry, syncStep)
}

func (r *ResHolder) Fetch(db *pg.DB, nodes []string) {
	var l []Ancestry

	var filtered []string
	for _, n := range nodes {
		_, ok := r.GetCached(n)
		if !ok {
			filtered = append(filtered, n)
		}
	}
	err := db.Model(&l).Column().WhereIn("hash IN (?)", filtered).Select()
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range l {
		r.fetched[string(a.Hash)] = a
	}
}

// LookUp without changing the cache
func (r *ResHolder) LookUp(n string) (Ancestry, bool) {
	if a, cached := r.cache.LookUp(n); cached {
		return *a, true
	}

	a, found := r.Get(n)
	return a, found
}

// GetCached gets the item & changes the cache
func (r *ResHolder) GetCached(n string) (Ancestry, bool) {
	if a, cached := r.cache.Get(n); cached {
		return *a, true
	}
	a, found := r.Get(n)
	if found {
		r.cache.Set(&a)
	}
	return a, found
}

// Get ignores the cache
func (r *ResHolder) Get(n string) (Ancestry, bool) {
	res, ok := r.fetched[n]
	if !ok {
		res, ok = r.m[n]
	}
	return res, ok
}

func NewResHolder() ResHolder {
	m := make(map[string]Ancestry, syncStep)
	fetched := make(map[string]Ancestry, syncStep)
	cache := NewRecentAncestry(cacheSize)

	return ResHolder{
		m:       m,
		fetched: fetched,
		cache:   cache,
	}
}

func (r ResHolder) Insert(db *pg.DB) {
	outs := make([]Ancestry, 0, len(r.m))

	for _, v := range r.m {
		outs = append(outs, v)
	}

	if _, err := db.Model(&outs).Insert(); err != nil {
		log.Fatalf("%v\n%v\n", err, outs)
	}
}
