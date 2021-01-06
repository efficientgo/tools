// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

// Initially copied from Thanos
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

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
