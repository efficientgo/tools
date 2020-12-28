// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

// Initially copied from Thanos
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.
//
// Package errcapture is useful when you want to close a `Closer` interface. As we all know, we should close all implements of `Closer`, such as *os.File. Commonly we will use:
//
// 	defer closer.Close()
//
// The problem is that Close() usually can return important error e.g for os.File the actual file flush might happen (and fail) on `Close` method. It's important to *always* check error. Thanos provides utility functions to log every error like those, allowing to put them in convenient `defer`:
//
// 	defer errcapture.CloseWithLog(logger, closer, "log format message")
//
// The errcapture.Exhaust* family of functions provide the same functionality but
// they take an io.ReadCloser and they exhaust the whole reader before closing
// them. They are useful when trying to use http keep-alive connections because
// for the same connection to be re-used the whole response body needs to be
// exhausted.
package logerrcapture

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

// Logger interface compatible with go-kit/logger.
type Logger interface {
	Log(keyvals ...interface{}) error
}

type closerFunc func() error

// CloseWithLog is making sure we log every error, even those from best effort tiny closers.
func Close(logger Logger, closer closerFunc, format string, a ...interface{}) {
	err := closer()
	if err == nil {
		return
	}

	// Not a problem if it has been closed already.
	if errors.Is(err, os.ErrClosed) {
		return
	}

	_ = logger.Log("msg", "detected close error", "err", errors.Wrap(err, fmt.Sprintf(format, a...)))
}

// ExhaustClose closes the io.ReadCloser with a log message on error but exhausts the reader before.
func ExhaustClose(logger Logger, r io.ReadCloser, format string, a ...interface{}) {
	_, err := io.Copy(ioutil.Discard, r)
	if err != nil {
		_ = logger.Log("msg", "failed to exhaust reader, performance may be impeded", "err", err)
	}

	Close(logger, r.Close, format, a...)
}
