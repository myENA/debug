package sync_test

import (
	"github.com/myENA/debug/sync"
	"math/rand"
	gosync "sync"
	"testing"
	"time"
)

func TestMutex(t *testing.T) {
	t.Parallel()
	m := sync.Mutex{}
	hit := ""

	wg := new(gosync.WaitGroup)
	wg.Add(2)
	t.Run("Routine1", func(t *testing.T) {
		m.Lock()
		if hit == "" {
			hit = "r1"
		}
		m.Unlock()
		wg.Done()
	})
	t.Run("Routine2", func(t *testing.T) {
		m.Lock()
		if hit == "" {
			hit = "r2"
		}
		m.Unlock()
		wg.Done()
	})

	wg.Wait()

	if hit == "" {
		t.Log("Expected either routine 1 or 2 to hit")
		t.FailNow()
	} else {
		t.Logf("%s hit first", hit)
	}
}

func TestMutex_Stale(t *testing.T) {
	t.Parallel()
	m := sync.Mutex{}

	wg := new(gosync.WaitGroup)
	wg.Add(3)
	t.Run("routine-1", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			m.Lock()
			s := rand.Intn(5)
			t.Logf("routine-1 acquired lock, sleeping for %d seconds", s)
			time.Sleep(time.Duration(s) * time.Second)
			m.Unlock()
		}
		wg.Done()
	})

	t.Run("routine-2", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			m.Lock()
			s := rand.Intn(5)
			t.Logf("routine-2 acquired lock, sleeping for %d seconds", s)
			time.Sleep(time.Duration(s) * time.Second)
			m.Unlock()
		}
		wg.Done()
	})

	t.Run("routine-3", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			m.Lock()
			s := rand.Intn(5)
			t.Logf("routine-3 acquired lock, sleeping for %d seconds", s)
			time.Sleep(time.Duration(s) * time.Second)
			m.Unlock()
		}
		wg.Done()
	})

	wg.Wait()
}
