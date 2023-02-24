package mgr

import (
	"fmt"
	"sync"
	"time"
)

type item struct {
	ts   time.Time
	used bool
}

type Mgr struct {
	store      map[string]item
	m          sync.Mutex
	expireTime time.Duration
}

func New(expireDur time.Duration) *Mgr {
	m := Mgr{
		store:      make(map[string]item),
		expireTime: expireDur,
	}
	return &m
}

func (m *Mgr) Push(code string) error {
	m.m.Lock()
	defer m.m.Unlock()
	if _, ok := m.store[code]; ok {
		return fmt.Errorf("code already exists: %s", code)
	}
	m.store[code] = item{ts: time.Now().UTC()}
	return nil
}

func (m *Mgr) Check(code string) error {
	m.m.Lock()
	defer m.m.Unlock()
	item, ok := m.store[code]
	if !ok {
		return fmt.Errorf("no code found: %s", code)
	}
	if item.used {
		return fmt.Errorf("code already in use: %s", code)
	}
	if time.Now().UTC().Add(-m.expireTime).After(item.ts) {
		return fmt.Errorf("code expired: %s", code)
	}
	return nil
}

func (m *Mgr) Use(code string) error {
	m.m.Lock()
	defer m.m.Unlock()
	item, ok := m.store[code]
	if !ok {
		return fmt.Errorf("no code found: %s", code)
	}
	if item.used {
		return fmt.Errorf("code already in use: %s", code)
	}
	if time.Now().UTC().Add(-m.expireTime).After(item.ts) {
		return fmt.Errorf("code expired: %s", code)
	}
	item.used = true
	m.store[code] = item
	return nil
}
