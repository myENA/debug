package sync

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type RWMutex struct {
	// uses actual sync RWMutex to do the RWMutex'ing
	mu sync.RWMutex

	// number of times this mutex has been locked
	count uint64

	// duration to wait before considering this lock "stale"
	lockTTL time.Duration

	// duration to wait before logging again during "stale" state
	tickRate time.Duration

	// name of this mutex, will be generated for you if not set
	name string

	// this will be used to track full Locks
	done chan struct{}
}

func NewDebugRWMutex(name string, lockTTL, tickRate time.Duration, log Logger) *RWMutex {
	if name == "" {
		name = fmt.Sprintf("rwmutex-%d", atomic.AddUint64(&rwmid, 1))
	}
	if lockTTL == 0 {
		lockTTL = defaultLockTTL
	}
	if tickRate == 0 {
		tickRate = defaultTickRate
	}
	return &RWMutex{
		name:     name,
		lockTTL:  lockTTL,
		tickRate: tickRate,
		done:     make(chan struct{}),
	}
}

func (rw *RWMutex) Lock() {
	rw.mu.Lock()
	if rw.name == "" {
		rw.name = fmt.Sprintf("mutex-%d", atomic.AddUint64(&mid, 1))
	}
	if rw.lockTTL == 0 {
		rw.lockTTL = defaultLockTTL
	}
	if rw.tickRate == 0 {
		rw.tickRate = defaultTickRate
	}
	if rw.done == nil {
		rw.done = make(chan struct{})
	}
	go lockTimeout(rw.name, rw.lockTTL, rw.tickRate, atomic.AddUint64(&rw.count, 1), rw.done, debug.Stack())
}

func (rw *RWMutex) Unlock() {
	rw.mu.Unlock()
	rw.done <- struct{}{}
}

func (rw *RWMutex) RLock() {
	rw.mu.RLock()
}

func (rw *RWMutex) RUnlock() {
	rw.mu.RUnlock()
}
