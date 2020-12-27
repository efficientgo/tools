// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

// Initially copied from Thanos
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.
//
// Package runutil provides helpers to advanced function scheduling control like repeat or retry.
//
// It's very often the case when you need to excutes some code every fixed intervals or have it retried automatically.
// To make it reliably with proper timeout, you need to carefully arrange some boilerplate for this.
// Below function does it for you.
//
// For repeat executes, use Repeat:
//
// 	err := runutil.Repeat(10*time.Second, stopc, func() error {
// 		// ...
// 	})
//
// Retry starts executing closure function f until no error is returned from f:
//
// 	err := runutil.Retry(10*time.Second, stopc, func() error {
// 		// ...
// 	})
//
// For logging an error on each f error, use RetryWithLog:
//
// 	err := runutil.RetryWithLog(logger, 10*time.Second, stopc, func() error {
// 		// ...
// 	})
package runutil

import (
	"time"
)

// Repeat executes f every interval seconds until stopc is closed or f returns an error.
// It executes f once right after being called.
func Repeat(interval time.Duration, stopc <-chan struct{}, f func() error) error {
	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		if err := f(); err != nil {
			return err
		}
		select {
		case <-stopc:
			return nil
		case <-tick.C:
		}
	}
}

// Logger interface compatible with go-kit/logger.
type Logger interface {
	Log(keyvals ...interface{}) error
}

// Retry executes f every interval seconds until timeout or no error is returned from f.
func Retry(interval time.Duration, stopc <-chan struct{}, f func() error) error {
	return RetryWithLog(nil, interval, stopc, f)
}

// RetryWithLog executes f every interval seconds until timeout or no error is returned from f. It logs an error on each f error.
func RetryWithLog(logger Logger, interval time.Duration, stopc <-chan struct{}, f func() error) error {
	tick := time.NewTicker(interval)
	defer tick.Stop()

	var err error
	for {
		if err = f(); err == nil {
			return nil
		}
		if logger != nil {
			_ = logger.Log("msg", "function failed. Retrying in next tick", "err", err)
		}
		select {
		case <-stopc:
			return err
		case <-tick.C:
		}
	}
}
