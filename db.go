package main

import (
	"sync"
)

// Db represent our storage
type Db struct {
	mu    sync.RWMutex
	name  string
	items map[string][]byte
}

func (db *Db) Get(key string) ([]byte, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	item, exists := db.items[key]

	if !exists {
		return nil, false
	}
	return item, true
}

func (db *Db) Set(key string, value []byte) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.items[key] = value
}

func (db *Db) Delete(key string) ([]byte, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	_, exists := db.items[key]

	if !exists {
		return nil, false
	}

	delete(db.items, key)
	return db.items[key], true
}

func NewDb(name string) *Db {
	db := new(Db)
	db.name = name
	db.items = make(map[string][]byte)

	return db
}
