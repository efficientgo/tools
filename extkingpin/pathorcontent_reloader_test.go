package extkingpin

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/efficientgo/core/testutil"
)

func TestPathContentReloader(t *testing.T) {
	type args struct {
		runSteps func(t *testing.T, testFile string, pathContent *StaticPathContent)
	}
	tests := []struct {
		name        string
		args        args
		wantReloads int
	}{
		{
			name: "Many operations, only rewrite triggers one reload",
			args: args{
				runSteps: func(t *testing.T, testFile string, pathContent *StaticPathContent) {
					testutil.Ok(t, os.Chmod(testFile, 0777))
					testutil.Ok(t, os.Remove(testFile))
					testutil.Ok(t, pathContent.Rewrite([]byte("test modified")))
				},
			},
			wantReloads: 1,
		},
		{
			name: "Many operations, only rename triggers one reload",
			args: args{
				runSteps: func(t *testing.T, testFile string, pathContent *StaticPathContent) {
					testutil.Ok(t, os.Chmod(testFile, 0777))
					testutil.Ok(t, os.Rename(testFile, testFile+".tmp"))
					testutil.Ok(t, os.Rename(testFile+".tmp", testFile))
				},
			},
			wantReloads: 1,
		},
		{
			name: "Many operations, two rewrites trigger two reloads",
			args: args{
				runSteps: func(t *testing.T, testFile string, pathContent *StaticPathContent) {
					testutil.Ok(t, os.Chmod(testFile, 0777))
					testutil.Ok(t, os.Remove(testFile))
					testutil.Ok(t, pathContent.Rewrite([]byte("test modified")))
					time.Sleep(2 * time.Second)
					testutil.Ok(t, pathContent.Rewrite([]byte("test modified again")))
				},
			},
			wantReloads: 1,
		},
		{
			name: "Chmod doesn't trigger reload",
			args: args{
				runSteps: func(t *testing.T, testFile string, pathContent *StaticPathContent) {
					testutil.Ok(t, os.Chmod(testFile, 0777))
				},
			},
			wantReloads: 0,
		},
		{
			name: "Remove doesn't trigger reload",
			args: args{
				runSteps: func(t *testing.T, testFile string, pathContent *StaticPathContent) {
					testutil.Ok(t, os.Remove(testFile))
				},
			},
			wantReloads: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := path.Join(t.TempDir(), "test")
			testutil.Ok(t, os.WriteFile(testFile, []byte("test"), 0666))
			pathContent, err := NewStaticPathContent(testFile)
			testutil.Ok(t, err)

			wg := &sync.WaitGroup{}
			wg.Add(tt.wantReloads)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			reloadCount := 0
			err = PathContentReloader(ctx, pathContent, newTestLogger("debug"), newTestLogger("error"), func() {
				reloadCount++
				wg.Done()
			}, 100*time.Millisecond)
			testutil.Ok(t, err)

			tt.args.runSteps(t, testFile, pathContent)
			wg.Wait()
			testutil.Equals(t, tt.wantReloads, reloadCount)
		})
	}
}

type testLogger struct {
	prefix string
}

func newTestLogger(prefix string) testLogger {
	return testLogger{prefix: prefix}
}

func (t testLogger) Log(keyvals ...interface{}) error {
	_, err := fmt.Printf("[%s] %s", t.prefix, keyvals)
	return err
}