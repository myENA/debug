package sync

import "time"

const (
	defaultLockTTL  = 2 * time.Second
	defaultTickRate = time.Second
)

var (
	mid   uint64
	rwmid uint64
)

func lockTimeout(name string, lockTTL, tickTTL time.Duration, count uint64, done chan struct{}, stack []byte) {
	start := time.Now() // TODO: not sure I like this...
	timer := time.NewTimer(lockTTL)
	select {
	case <-done:
		if !timer.Stop() {
			<-timer.C
		}
		return
	case <-timer.C:
	}

	log.Printf("[%s] (%d) Locked for >= %s. Trace:", name, count, lockTTL)
	log.Print(string(stack))

	ticker := time.NewTicker(tickTTL)
	for i := 0; ; i++ {
		select {
		case <-done:
			ticker.Stop()
			log.Printf("[%s] (%d) Unlocked after %s", name, count, time.Now().Sub(start))
			return
		case <-ticker.C:
			if i > 0 && i%5 == 0 {
				log.Printf("[%s] (%d) Still locked after %s. Trace:", name, count, time.Now().Sub(start))
				log.Print(string(stack))
			} else {
				log.Printf("[%s] (%d) Still locked after %s", name, count, time.Now().Sub(start))
			}
		}
	}
}
