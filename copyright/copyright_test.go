// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestCopyrightApplier(t *testing.T) {
	tmpDir, err := ioutil.TempDir(os.TempDir(), "copyright")
	testutil.Ok(t, err)
	t.Cleanup(func() { testutil.Ok(t, os.RemoveAll(tmpDir)) })

	for _, a := range []*copyrightApplier{
		NewCopyrightApplier([]byte(`Copyright (c) The EfficientGo Authors.

Licensed under the Apache License 2.0.`)),
		NewCopyrightApplier([]byte(`Copyright (c) The EfficientGo Authors.

Licensed under the Apache License 2.0.`)),
		NewCopyrightApplier([]byte(`
Copyright (c) The EfficientGo Authors.

Licensed under the Apache License 2.0.
`)),
	} {
		t.Run("", func(t *testing.T) {
			for _, tcase := range []struct {
				applier *copyrightApplier

				filename    string
				input       string
				expected    string
				expectedErr error
			}{
				{
					applier: a, filename: "abc", input: `package x

// yolo
,`, expectedErr: errors.New("unsupported file extension "),
				},
				{
					applier: a, filename: "abc.go", input: `package x

// yolo
,`, expected: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					applier: a, filename: "abc.go", input: `
package x

// yolo
,`, expected: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					applier: a, filename: "abc.go", input: `Copyright (c) The EfficientGo Authors.
Licensed under the Apache License 2.0.

package x

// yolo
,`, expected: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

Copyright (c) The EfficientGo Authors.
Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					applier: a, filename: "abc.go", input: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

package x

// yolo
,`, expected: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					applier: a, filename: "empty.go", input: ``, expected: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

`,
				},
				{
					applier: a, filename: "abc.c", input: `package x

// yolo
,`, expected: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					applier: a, filename: "abc.cpp", input: `package x

// yolo
,`, expected: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					applier: a, filename: "abc.sh", input: `package x

// yolo
,`, expected: `# Copyright (c) The EfficientGo Authors.
#
# Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					applier: a, filename: "abc.py", input: `package x

// yolo
,`, expected: `# Copyright (c) The EfficientGo Authors.
#
# Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					applier: a, filename: "abc.proto", input: `package x

// yolo
,`, expected: `// Copyright (c) The EfficientGo Authors.
//
// Licensed under the Apache License 2.0.

package x

// yolo
,`,
				},
				{
					// Generated.
					applier: a, filename: "abc.pb.go", input: `package x

// yolo
,`, expected: `package x

// yolo
,`,
				},
			} {
				t.Run("", func(t *testing.T) {
					file := filepath.Join(tmpDir, tcase.filename)
					testutil.Ok(t, ioutil.WriteFile(file, []byte(tcase.input), os.ModePerm))
					defer func() { testutil.Ok(t, os.RemoveAll(file)) }()

					err := tcase.applier.Apply(file)
					if tcase.expectedErr != nil {
						testutil.NotOk(t, err)
						testutil.Equals(t, tcase.expectedErr.Error(), err.Error())
						return
					}
					testutil.Ok(t, err)

					b, err := ioutil.ReadFile(file)
					testutil.Ok(t, err)
					testutil.Equals(t, tcase.expected, string(b))
				})
			}
		})
	}
}
