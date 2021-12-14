package signal

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func WaitAny(d time.Duration, wait, timeout func()) {
	ch := make(chan struct{})

	go func() {
		wait()
		ch <- struct{}{}
		close(ch)
	}()

	select {
	case <-time.After(d):
		timeout()
	case <-ch:
	}
}

func NoSignalArrived(t *testing.T) func() {
	return func() {
		t.Error("no signal arrived")
	}
}

func TestOnce(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	Once(SigUsr1).Notify(context.TODO(), func(sig os.Signal) {
		if assert.Equal(t, SigUsr1, sig) {
			wg.Done()
		}
	})

	if pid := os.Getpid(); assert.Greater(t, pid, 0) {
		if err := SendSignalUser1(pid); assert.NoError(t, err) {
			WaitAny(time.Second, wg.Wait, NoSignalArrived(t))
		}

		assert.NoError(t, SendSignalUser1(pid))
	}
}

func TestWhen(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	When(SigUsr2).Notify(context.TODO(), func(sig os.Signal) {
		if assert.Equal(t, SigUsr2, sig) {
			wg.Done()
		}
	})

	if pid := os.Getpid(); assert.Greater(t, pid, 0) {
		if assert.NoError(t, SendSignalUser2(pid)) {
			WaitAny(time.Second, func() {
				wg.Wait()
				wg.Add(1)
				if assert.NoError(t, SendSignalUser2(pid)) {
					wg.Wait()
				}
			}, NoSignalArrived(t))
		}
	}
}

func TestWith(t *testing.T) {
	ctx, cancel := With(context.TODO(), SigUsr1)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		<-ctx.Done()
		wg.Done()
	}()

	if pid := os.Getpid(); assert.Greater(t, pid, 0) {
		if err := SendSignalUser1(pid); assert.NoError(t, err) {
			WaitAny(time.Second, wg.Wait, NoSignalArrived(t))
		}
	}
}
