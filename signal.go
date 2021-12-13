/*
 * Copyright (c) 2021 Jayson Wang <jayson@laboys.org>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// Package signal provides simple, semantic manipulation of the operating system's
// signal processing.
package signal

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

const (
	// SigCtx represents signal from cancelled context
	SigCtx = syscall.Signal(0xff)
)

// Notifier creates and listens to a specified system signal and calls a handler
// when the signal is received or when the context is cancelled.
//
// Notifier will be closed when the context is cancelled and the
// handler can get an instance of SigCtx (number = 0xff).
//
// methods of the Notifier can be safely called multiple times.
type Notifier interface {
	Notify(context.Context, func(sig os.Signal))
}

// When perform the handler function when one of the listed signals arrives.
//
// Notifier created by When can only be closed by canceling the context object.
func When(signals ...os.Signal) Notifier {
	return &notifier{
		once:    false,
		signals: signals,
	}
}

// Once perform the handler once only when the context is cancelled or
// when one of the listed signals arrives, after which the Notifier
// will be closed directly.
func Once(signals ...os.Signal) Notifier {
	return &notifier{
		once:    true,
		signals: signals,
	}
}

// notifier implements the Notifier interface
//
// The internal once state determines whether the handler
// can be performed multiple.
//
// signals indicate the list of operating system's signal
// to be listened to.
type notifier struct {
	once    bool
	signals []os.Signal
}

// Notify creates a channel to receive signals from the operating system, passed
// the signal to the handler when it is received.
//
// when the context object is cancelled, the SigCtx is passed to the handler and
// the channel and goroutine are cleaned up and exited.
func (n *notifier) Notify(ctx context.Context, handler func(sig os.Signal)) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, n.signals...)

	go func() {
		defer func() { signal.Stop(signals); close(signals) }()

		for {
			select {
			case sig := <-signals:
				if handler(sig); n.once {
					return
				}
			case <-ctx.Done():
				handler(SigCtx)
				return
			}
		}
	}()
}

// With returns a copy of the parent context that is marked done
// when one of the listed signals arrives, when the returned cancel
// function is called, or when the parent context's canceled,
// whichever happens first.
func With(parent context.Context, signals ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	Once(signals...).Notify(parent, func(sig os.Signal) {
		cancel()
	})
	return ctx, cancel
}
