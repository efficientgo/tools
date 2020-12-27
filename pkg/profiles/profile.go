package profiles

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/pprof"

	"github.com/efficientgo/tools/pkg/errcapture"
	"github.com/felixge/fgprof"
	"github.com/pkg/errors"
)

func Heap(dir string) (err error) {
	f, err := os.Open(filepath.Join(dir, "mem.pprof"))
	if err != nil {
		return err
	}
	defer errcapture.Close(&err, f, "close")
	return pprof.WriteHeapProfile(f)
}

type CPUProfileType string

const (
	CPUProfileTypeBuiltIn CPUProfileType = "built-in"
	// CPUProfileTypeFGProf represents enhanced https://github.com/felixge/fgprof profiling.
	CPUProfileTypeFGProf CPUProfileType = "fgprof"
)

// StartCPU starts CPU profiling. If no error, it returns close function that stops and flushes profile to provided
// directory.
func StartCPU(dir string, typ CPUProfileType) (closeFn func() error, err error) {
	fileName := "cpu.pprof"
	switch typ {
	case CPUProfileTypeBuiltIn:
	case CPUProfileTypeFGProf:
		fileName = "cpu.fgprof.pprof"
	default:
		return nil, errors.Errorf("unknown CPU profile type %v", typ)

	}
	f, err := os.Open(filepath.Join(dir, fileName))
	if err != nil {
		return nil, err
	}

	switch typ {
	case CPUProfileTypeBuiltIn:
		if err = pprof.StartCPUProfile(f); err != nil {
			errcapture.Close(&err, f, fmt.Sprintf("close %v", filepath.Join(dir, fileName)))
			return nil, err
		}
		closeFn = func() (ferr error) {
			pprof.StopCPUProfile()
			return errors.Wrapf(f.Close(), "close %v", filepath.Join(dir, fileName))
		}
	case CPUProfileTypeFGProf:
		closeFGProfFn := fgprof.Start(f, fgprof.FormatPprof)
		closeFn = func() (ferr error) {
			defer errcapture.Close(&ferr, f, fmt.Sprintf("close %v", filepath.Join(dir, fileName)))
			return closeFGProfFn()
		}
	}
	return closeFn, nil
}
