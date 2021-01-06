// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package errcapture

// Close a `Closer` interface safely while capturing error. It's often forgotten but it's a caller responsibility
// to close all implementations of `Closer`, such as *os.File or io.ReaderCloser. Commonly we would use:
//
// 	defer closer.Close()
//
// This is wrong. Close() usually return important error (e.g for os.File the actual file flush might happen and fail on `Close` method).
// It's very important to *always* check error. `errcapture` provides utility functions to capture error and add to provided one,
// still allowing to put them in a convenient `defer` statement:
//
// 	func <...>(...) (err error) {
//  	...
//  	defer errcapture.Close(&err, closer, "log format message")
//
// 		...
// 	}
//
//If closer returns error, `errcapture.Close` will capture it, add to input error if not nil and return by argument.
//
// The errcapture.ExhaustClose function provide the same functionality but takes an io.ReadCloser and exhausts the whole
// reader before closing. This is useful when trying to use http keep-alive connections because for the same connection
// to be re-used the whole response body needs to be exhausted.
//
// Check https://pkg.go.dev/github.com/efficientgo/tools/pkg/logerrcapture if you want to just log an error instead.
