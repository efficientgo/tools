// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/efficientgo/tools/core/pkg/errcapture"
	"github.com/pkg/errors"
	"github.com/protoconfig/protoconfig/go/kingpinv2"
	"gopkg.in/alecthomas/kingpin.v2"
)

type langSpec struct {
	lineCommentChars []byte
	isGenererated    func(file string) (bool, error)
}

// TODO(bwplotka): Support more, fill gaps.
var specByExt = map[string]langSpec{
	".go": {
		lineCommentChars: []byte("//"),
		isGenererated: func(file string) (bool, error) {
			// Yolo
			return strings.HasSuffix(file, ".pb.go"), nil
		},
	},
	".proto": {
		lineCommentChars: []byte("//"),
	},
	".c": {
		lineCommentChars: []byte("//"),
	},
	".cpp": {
		lineCommentChars: []byte("//"),
	},
	".py": {
		lineCommentChars: []byte("#"),
	},
	".sh": {
		lineCommentChars: []byte("#"),
	},
}

type copyrightApplier struct {
	copyright []byte

	copyrightBuff map[string][]byte
	fileBuff      bytes.Buffer

	verbose bool
}

func NewCopyrightApplier(content []byte) *copyrightApplier {
	return &copyrightApplier{
		copyright:     bytes.Trim(content, "\n"),
		copyrightBuff: map[string][]byte{},
		verbose:       false,
	}
}

func (c *copyrightApplier) EnableVerbose() {
	c.verbose = true
}

// TODO(bwplotka): Make this concurrently.
func (c *copyrightApplier) Apply(file string) (err error) {
	spec, ok := specByExt[filepath.Ext(file)]
	if !ok {
		return errors.Errorf("unsupported file extension %v", filepath.Ext(file))
	}

	if spec.isGenererated != nil {
		is, err := spec.isGenererated(file)
		if err != nil {
			return err
		}
		if is {
			return nil
		}
	}

	cb, ok := c.copyrightBuff[filepath.Ext(file)]
	if !ok {
		split := bytes.Split(c.copyright, []byte("\n"))
		for i, s := range split {
			if len(s) > 0 {
				split[i] = append(
					append(
						append([]byte{}, spec.lineCommentChars...),
						' ', // Add space only if there is some content on the line.
					),
					s...,
				)
				continue
			}
			split[i] = append([]byte{}, spec.lineCommentChars...)
		}
		cb = bytes.Join(split, []byte("\n"))
		c.copyrightBuff[filepath.Ext(file)] = cb
	}

	f, err := os.OpenFile(file, os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer errcapture.Close(&err, f.Close, "close")

	hdr := make([]byte, len(cb))
	n, err := f.Read(hdr)
	if err != io.EOF && err != nil {
		return errors.Wrapf(err, "read first %v bytes", len(cb))
	}

	if !bytes.Equal(hdr, cb) {
		// TODO(bwplotka): Figure out how to help on renames?
		if c.verbose {
			log.Println("file", file, "is missing Copyright header or header changed (got ", string(hdr), "). Adding.")
		} else {
			log.Println("file", file, "is missing Copyright header or header changed. Adding.")
		}

		c.fileBuff.Reset()
		c.fileBuff.Write(cb)
		c.fileBuff.WriteByte('\n')
		if len(hdr) == 0 || hdr[0] != '\n' {
			c.fileBuff.WriteByte('\n')
		}
		c.fileBuff.Write(hdr[:n])
		if n >= len(cb) {
			// Read rest of file to the buffer.
			if _, err := io.Copy(&c.fileBuff, f); err != nil {
				return errors.Wrap(err, "read")
			}
		}

		// TODO(bwplotka): Not atomic and safest ever, do it better (tmp file?)
		if _, err = f.Seek(0, 0); err != nil {
			return errors.Wrap(err, "seek")
		}

		if _, err = c.fileBuff.WriteTo(f); err != nil {
			return errors.Wrap(err, "write")
		}
	}
	return nil
}

func main() {
	app := kingpin.New(filepath.Base(os.Args[0]), `copyright`)
	copyright := kingpinv2.Flag(app, "copyright", "Copyright content to apply to provided files").DefaultPath("./COPYRIGHT").PathOrContent()
	files := app.Arg("files", "Files to apply copyright to.").ExistingFiles()
	verbose := app.Flag("verbose", "Enable verbose prints.").Short('v').Bool()
	if _, err := app.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	content, err := copyright.Content()
	if err != nil {
		log.Fatal(err)
	}

	a := NewCopyrightApplier(content)
	if *verbose {
		a.EnableVerbose()
	}
	for _, f := range *files {
		if err := a.Apply(f); err != nil {
			log.Fatal(err)
		}
	}
}
