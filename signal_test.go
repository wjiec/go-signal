package signal

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func WaitUntil(duration time.Duration, process, timout func()) {
	stop := make(chan struct{})

	go func() {
		process()
		close(stop)
	}()

	select {
	case <-time.After(duration):
		timout()
	case <-stop:
	}
}

func NoSignalArrived(t *testing.T) func() {
	return func() {
		t.Error("no signal arrived")
	}
}

func TestOnce(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(1)
		Once(SigUsr1).Notify(context.Background(), func(sig os.Signal) {
			if assert.Equal(t, SigUsr1, sig) {
				wg.Done()
			}
		})

		if pid := os.Getpid(); assert.Greater(t, pid, 0) {
			if err := SendSignalUser1(pid); assert.NoError(t, err) {
				WaitUntil(time.Second, wg.Wait, NoSignalArrived(t))
			}

			assert.NoError(t, SendSignalUser1(pid))
		}
	})

	t.Run("canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		Once(SigUsr1).Notify(ctx, func(sig os.Signal) {
			assert.Equal(t, SigCtx, sig)
		})

		cancel()
	})
}

func TestWhen(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(1)
		When(SigUsr2).Notify(context.TODO(), func(sig os.Signal) {
			if assert.Equal(t, SigUsr2, sig) {
				wg.Done()
			}
		})

		if pid := os.Getpid(); assert.Greater(t, pid, 0) {
			if assert.NoError(t, SendSignalUser2(pid)) {
				WaitUntil(time.Second, func() {
					wg.Wait()
					wg.Add(1)
					if assert.NoError(t, SendSignalUser2(pid)) {
						wg.Wait()
					}
				}, NoSignalArrived(t))
			}
		}
	})

	t.Run("canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		When(SigUsr2).Notify(ctx, func(sig os.Signal) {
			assert.Equal(t, SigCtx, sig)
		})

		cancel()
	})
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
			WaitUntil(time.Second, wg.Wait, NoSignalArrived(t))
		}
	}
}
