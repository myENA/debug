package sync

import (
	"fmt"
	stdlog "log"
	"os"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultLockTTL = 2 * time.Second
	defaultTickTTL = time.Second
)

var (
	mid        uint64
	defaultLog = stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)
)

type Mutex struct {
	// uses actual sync Mutex to do the Mutex'ing
	mu sync.Mutex

	// number of times this mutex has been locked
	count uint64

	// duration to wait before considering this lock "stale"
	lockTTL time.Duration

	// duration to wait before logging again during "stale" state
	tickTTL time.Duration

	// will be populated once Unlock is called
	done chan struct{}

	// name of this mutex, will be generated for you if not set
	name string

	// logger to use, will use default defined above if not provided
	log Logger
}

func NewDebugMutex(name string, lockTTL, tickTTL time.Duration, log Logger) *Mutex {
	if name == "" {
		name = fmt.Sprintf("mutex-%d", atomic.AddUint64(&mid, 1))
	}
	if log == nil {
		log = defaultLog
	}
	if lockTTL == 0 {
		lockTTL = defaultLockTTL
	}
	if tickTTL == 0 {
		tickTTL = defaultTickTTL
	}
	return &Mutex{
		name:    name,
		log:     log,
		lockTTL: lockTTL,
		tickTTL: tickTTL,
		done:    make(chan struct{}),
	}
}

func (m *Mutex) Lock() {
	m.mu.Lock()
	if m.name == "" {
		m.name = fmt.Sprintf("mutex-%d", atomic.AddUint64(&mid, 1))
	}
	if m.log == nil {
		m.log = defaultLog
	}
	if m.lockTTL == 0 {
		m.lockTTL = defaultLockTTL
	}
	if m.tickTTL == 0 {
		m.tickTTL = defaultTickTTL
	}
	if m.done == nil {
		m.done = make(chan struct{})
	}
	atomic.AddUint64(&m.count, 1)
	go m.timeout(debug.Stack(), m.lockTTL, m.tickTTL)
}

func (m *Mutex) Unlock() {
	m.mu.Unlock()
	m.done <- struct{}{}
}

func (m *Mutex) timeout(stack []byte, lockTTL, tickTTL time.Duration) {
	start := time.Now() // TODO: not sure I like this...
	timer := time.NewTimer(lockTTL)
	select {
	case <-m.done:
		if !timer.Stop() {
			<-timer.C
		}
		return
	case <-timer.C:
	}

	count := atomic.LoadUint64(&m.count)
	m.log.Printf("[%s] (%d) Locked for >= %s. Trace:", m.name, count, lockTTL)
	m.log.Print(string(stack))

	ticker := time.NewTicker(tickTTL)
	for i := 0; ; i++ {
		select {
		case <-m.done:
			ticker.Stop()
			m.log.Printf("[%s] (%d) Unlocked after %s", m.name, count, time.Now().Sub(start))
			return
		case <-ticker.C:
			if i > 0 && i%5 == 0 {
				m.log.Printf("[%s] (%d) Still locked after %s. Trace:", m.name, count, time.Now().Sub(start))
				m.log.Print(string(stack))
			} else {
				m.log.Printf("[%s] (%d) Still locked after %s", m.name, count, time.Now().Sub(start))
			}
		}
	}
}
