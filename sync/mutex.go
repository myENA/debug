package sync

import (
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

type Mutex struct {
	// uses actual sync Mutex to do the Mutex'ing
	mu sync.Mutex

	// name of this mutex, will be generated for you if not set
	name string

	// number of times this mutex has been locked
	count uint64

	// duration to wait before considering this lock "stale"
	lockTTL time.Duration

	// duration to wait before logging again during "stale" state
	tickRate time.Duration

	// will be populated once Unlock is called
	done chan struct{}
}

func NewDebugMutex(name string, lockTTL, tickRate time.Duration, log Logger) *Mutex {
	if name == "" {
		name = fmt.Sprintf("mutex-%d", atomic.AddUint64(&mid, 1))
	}
	if lockTTL == 0 {
		lockTTL = defaultLockTTL
	}
	if tickRate == 0 {
		tickRate = defaultTickRate
	}
	return &Mutex{
		name:     name,
		lockTTL:  lockTTL,
		tickRate: tickRate,
		done:     make(chan struct{}),
	}
}

func (m *Mutex) Lock() {
	m.mu.Lock()
	if m.name == "" {
		m.name = fmt.Sprintf("mutex-%d", atomic.AddUint64(&mid, 1))
	}
	if m.lockTTL == 0 {
		m.lockTTL = defaultLockTTL
	}
	if m.tickRate == 0 {
		m.tickRate = defaultTickRate
	}
	if m.done == nil {
		m.done = make(chan struct{})
	}
	go lockTimeout(m.name, m.lockTTL, m.tickRate, atomic.AddUint64(&m.count, 1), m.done, debug.Stack())
}

func (m *Mutex) Unlock() {
	m.mu.Unlock()
	m.done <- struct{}{}
}
