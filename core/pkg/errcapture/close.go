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
// 	var err error
// 	defer errcapture.Close(&err, closer, "log format message")
//
// 	// ...
//
// If Close() returns error, err will capture it and return by argument.
//
// The errcapture.Exhaust* family of functions provide the same functionality but
// they take an io.ReadCloser and they exhaust the whole reader before closing
// them. They are useful when trying to use http keep-alive connections because
// for the same connection to be re-used the whole response body needs to be
// exhausted.
package errcapture

import (
	"io"
	"io/ioutil"

	"github.com/efficientgo/tools/core/pkg/merrors"
	"github.com/pkg/errors"
)

type closerFunc func() error

// Close runs function and on error return error by argument including the given error (usually
// from caller function).
func Close(err *error, closer closerFunc, format string, a ...interface{}) {
	*err = merrors.New(*err, errors.Wrapf(closer(), format, a...)).Err()
}

// ExhaustClose closes the io.ReadCloser with error capture but exhausts the reader before.
func ExhaustClose(err *error, r io.ReadCloser, format string, a ...interface{}) {
	_, copyErr := io.Copy(ioutil.Discard, r)

	Close(err, r.Close, format, a...)

	// Prepend the io.Copy error.
	*err = merrors.New(copyErr, *err).Err()
}
