# tools

[![copyright module docs](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/efficientgo/tools/copyright) [![core module docs](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/efficientgo/tools/core) [![e2e module docs](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/efficientgo/tools/e2e) [![extkingpin module docs](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/efficientgo/tools/extkingpin)

Set of lightweight tools, packages and modules that every open-source Go project always needs with almost no dependencies.

## Release model

Since this is meant to be critical, tiny import, multi module toolset, there are currently no semver releases planned. It's designed to pin modules via git commits, all commits to master should be stable and properly tested, vetted and linted.

API is considered stable, but rare API changes might occur. If they do - they will cause compilation error, so it will be easy to spot.

## Modules

### Module `github.com/efficientgo/tools/core`

The main module containing set of useful, core packages for testing, closing, running and repeating.

**This module is optimized for almost zero dependencies for ease of use.**

This module contains:

* `pkg/clilog`:

```go mdox-gen-exec="sh -c 'tail -n +6 core/pkg/clilog/doc.go'"
// Logging formatter that transforms structure log entry into human readable, clean friendly entry
// suitable more for CLI tools.
//
// In details this means:
//
// * No special sign escaping.
// * No key printing.
// * Values separated with ': '
// * Support for pretty printing multi errors (including nested ones) in format of (<something>: <err1>; <err2>; ...; <errN>)
// * TODO(bwplotka): Support for multiple multilines.
//
// Compatible with `github.com/go-kit/kit/log.Logger`
```

* `pkg/errcapture`

```go mdox-gen-exec="sh -c 'tail -n +6 core/pkg/errcapture/doc.go'"
// Close a `io.Closer` interface or execute any function that returns error safely while capturing error.
// It's often forgotten but it's a caller responsibility to close all implementations of `Closer`,
// such as *os.File or io.ReaderCloser. Commonly we would use:
//
// 	defer closer.Close()
//
// This is wrong. Close() usually return important error (e.g for os.File the actual file flush might happen and fail on `Close` method).
// It's very important to *always* check error. `errcapture` provides utility functions to capture error and add to provided one,
// still allowing to put them in a convenient `defer` statement:
//
// 	func <...>(...) (err error) {
//  	...
//  	defer errcapture.Do(&err, closer.Close, "log format message")
//
// 		...
// 	}
//
// If Close returns error, `errcapture.Do` will capture it, add to input error if not nil and return by argument.
//
// The errcapture.ExhaustClose function provide the same functionality but takes an io.ReadCloser and exhausts the whole
// reader before closing. This is useful when trying to use http keep-alive connections because for the same connection
// to be re-used the whole response body needs to be exhausted.
//
// Check https://pkg.go.dev/github.com/efficientgo/tools/pkg/logerrcapture if you want to just log an error instead.
```

* `pkg/logerrcapture`

```go mdox-gen-exec="sh -c 'tail -n +6 core/pkg/logerrcapture/doc.go'"
// Close a `io.Closer` interface or execute any function that returns error safely while logging error.
// It's often forgotten but it's a caller responsibility to close all implementations of `Closer`,
// such as *os.File or io.ReaderCloser. Commonly we would use:
//
// 	defer closer.Close()
//
// This is wrong. Close() usually return important error (e.g for os.File the actual file flush might happen and fail on `Close` method).
// It's very important to *always* check error. `logerrcapture` provides utility functions to capture error and log it via provided
// logger, while still allowing to put them in a convenient `defer` statement:
//
// 	func <...>(...) (err error) {
//  	...
//  	defer logerrcapture.Do(logger, closer.Close, "log format message")
//
// 		...
// 	}
//
// If Close returns error, `logerrcapture.Do` will capture it, add to input error if not nil and return by argument.
//
// The logerrcapture.ExhaustClose function provide the same functionality but takes an io.ReadCloser and exhausts the whole
// reader before closing. This is useful when trying to use http keep-alive connections because for the same connection
// to be re-used the whole response body needs to be exhausted.
//
// Recommended: Check https://pkg.go.dev/github.com/efficientgo/tools/pkg/errcapture if you want to return error instead of just logging (causing
// hard error).
```

* `pkg/merrors`

```go mdox-gen-exec="sh -c 'tail -n +6 core/pkg/merrors/doc.go'"
// Safe multi error implementation that chains errors on the same level. Supports errors.As and errors.Is functions.
//
// Example 1:
//
//  return merrors.New(err1, err2).Err()
//
// Example 2:
//
//  merr := merrors.New(err1)
//  merr.Add(err2, errOrNil3)
//  for _, err := range errs {
//    merr.Add(err)
//  }
//  return merr.Err()
//
```

* `pkg/runutil`

```go mdox-gen-exec="sh -c 'tail -n +6 core/pkg/runutil/doc.go'"
// Helpers for advanced function scheduling control like repeat or retry.
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
```

* `pkg/testutil`

```go mdox-gen-exec="sh -c 'tail -n +6 core/pkg/testutil/doc.go'"
// Simplistic assertion helpers for testing code. TestOrBench utils for union of testing and benchmarks.
```

### Module `github.com/efficientgo/tools/e2e`

This module is a fully featured e2e suite allowing utilizing `go test` for setting hermetic up complex microservice testing scenarios using docker.

```go mdox-gen-exec="sh -c 'tail -n +6 e2e/doc.go'"
// This module is a fully featured e2e suite allowing utilizing `go test` for setting hermetic up complex microservice integration testing scenarios using docker.
// Example usages:
//  * https://github.com/cortexproject/cortex/tree/master/integration
//  * https://github.com/thanos-io/thanos/tree/master/test/e2e
//
// Check github.com/efficientgo/tools/e2e/db for common DBs services you can run out of the box.
```

Credits:

* [Cortex Team](https://github.com/cortexproject/cortex/tree/f639b1855c9f0c9564113709a6bce2996d151ec7/integration)
* Initial Author: [@pracucci](https://github.com/pracucci)

### Module `github.com/efficientgo/tools/copyright`

This module is a very simple CLI for ensuring copyright header on code files.

```bash mdox-gen-exec="sh -c 'cd copyright && go run copyright.go --help || exit 0'"
usage: copyright [<flags>] [<files>...]

copyright

Flags:
      --help                 Show context-sensitive help (also try --help-long
                             and --help-man).
      --copyright-file=<file-path>  
                             Path to Copyright content to apply to provided
                             files
      --copyright=<content>  Alternative to 'copyright-file' flag (lower
                             priority). Content of Copyright content to apply to
                             provided files
  -v, --verbose              Enable verbose prints.

Args:
  [<files>]  Files to apply copyright to.

```

Install via standard Go installation pattern:

```shell
go install github.com/efficientgo/tools/copyright
```

or via [bingo](https://github.com/bwplotka/bingo) if want to pin it:

```shell
go install github.com/bwplotka/bingo
bingo get -u github.com/efficientgo/tools/copyright
```

### Module `github.com/efficientgo/tools/extkingpin`

This module provides the PathOrContent flag type which defines two flags to fetch bytes. Either from file (\*-file flag) or content (\* flag). Also returns the content of a YAML file with substituted environment variables.

```go mdox-gen-exec="sh -c 'tail -n +6 extkingpin/doc.go'"
// PathOrContent is a flag type that defines two flags to fetch bytes. Either from file (*-file flag) or content (* flag).
// Also returns content of YAML file with substituted environment variables.
// Follows K8s convention, i.e $(...), as mentioned here https://kubernetes.io/docs/tasks/inject-data-application/define-interdependent-environment-variables/.

// RegisterPathOrContent registers PathOrContent flag in kingpinCmdClause.

// Content returns the content of the file when given or directly the content that has been passed to the flag.
// It returns an error when:
// * The file and content flags are both not empty.
// * The file flag is not empty but the file can't be read.
// * The content is empty and the flag has been defined as required.

// Option is a functional option type for PathOrContent objects.
// WithRequired allows you to override default required option.
// WithEnvSubstitution allows you to override default envSubstitution option.
```
