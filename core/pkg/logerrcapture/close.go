// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

// Initially copied from Thanos
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

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
