// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package profiles

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/efficientgo/tools/core/pkg/errcapture"
	"github.com/felixge/fgprof"
	"github.com/pkg/errors"
)

func Heap(dir string) (err error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(dir, "mem.pprof"))
	if err != nil {
		return err
	}
	defer errcapture.Do(&err, f.Close, "close")
	return pprof.WriteHeapProfile(f)
}

type CPUType string

const (
	CPUTypeBuiltIn CPUType = "built-in"
	// CPUTypeFGProf represents enhanced https://github.com/felixge/fgprof CPU profiling.
	CPUTypeFGProf CPUType = "fgprof"
)

// StartCPU starts CPU profiling. If no error, it returns close function that stops and flushes profile to provided
// directory.
// NOTE(bwplotka): It does not make sense to run more than one of those.
func StartCPU(dir string, typ CPUType) (closeFn func() error, err error) {
	fileName := "cpu.pprof"
	switch typ {
	case CPUTypeBuiltIn:
	case CPUTypeFGProf:
		fileName = "cpu.fgprof.pprof"
	default:
		return nil, errors.Errorf("unknown CPU profile type %v", typ)
	}

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	f, err := os.Create(filepath.Join(dir, fileName))
	if err != nil {
		return nil, err
	}

	switch typ {
	case CPUTypeBuiltIn:
		if err = pprof.StartCPUProfile(f); err != nil {
			errcapture.Do(&err, f.Close, fmt.Sprintf("close %v", filepath.Join(dir, fileName)))
			return nil, err
		}
		closeFn = func() (ferr error) {
			pprof.StopCPUProfile()
			return errors.Wrapf(f.Close(), "close %v", filepath.Join(dir, fileName))
		}
	case CPUTypeFGProf:
		closeFGProfFn := fgprof.Start(f, fgprof.FormatPprof)
		closeFn = func() (ferr error) {
			defer errcapture.Do(&ferr, f.Close, fmt.Sprintf("close %v", filepath.Join(dir, fileName)))
			return closeFGProfFn()
		}
	}
	return closeFn, nil
}
