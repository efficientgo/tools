// Copyright (c) The EfficientGo Authors.
// Licensed under the Apache License 2.0.

package extflag

import (
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Clause is a fluid interface used to build extended flags.
// This is useful for flags that are a bit of uncommon logic extension like fileOrContent.
type Clause struct {
	cmd      FlagClause
	name     string
	help     string
	required bool

	hiddenPath  bool
	defaultPath string

	hiddenContent  bool
	defaultContent string
}

type FlagClause interface {
	Flag(name, help string) *kingpin.FlagClause
}

type AllClause interface {
	Required() AllClause
	HiddenPath() AllClause
	HiddenContent() AllClause
	DefaultPath(values string) PathOrContentClause
	DefaultContent(values string) PathOrContentClause
	PathOrContentClause
}

type PathOrContentClause interface {
	PathOrContent() *PathOrContent
}

func Flag(cmd FlagClause, flagName string, help string) AllClause {
	return &Clause{cmd: cmd, name: flagName, help: help}
}

func (f *Clause) DefaultPath(value string) PathOrContentClause {
	f.defaultPath = value
	return f
}

func (f *Clause) DefaultContent(value string) PathOrContentClause {
	f.defaultContent = value
	return f
}

// HiddenContent hides a flag from usage but still allows it to be used.
func (f *Clause) HiddenContent() AllClause {
	f.hiddenContent = true
	return f
}

// HiddenPath hides a -file flag from usage but still allows it to be used.
func (f *Clause) HiddenPath() AllClause {
	f.hiddenPath = true
	return f
}

// Required makes the flag required. You can not provide a Default() value to a Required() flag.
func (f *Clause) Required() AllClause {
	f.required = true
	return f
}

// PathOrContent is a flag type that defines two flags to fetch bytes. Either from file (*-file flag) or content (* flag).
type PathOrContent struct {
	flagName string

	required bool

	path    *string
	content *string
}

func (f *Clause) PathOrContent() *PathOrContent {
	pathFlagName := fmt.Sprintf("%s-file", f.name)
	contentFlagName := f.name

	c := f.cmd.Flag(pathFlagName, fmt.Sprintf("Path to %s", f.help))
	if f.hiddenPath {
		c = c.Hidden()
	}
	if f.required {
		c = c.Required()
	}
	if f.defaultPath != "" {
		// TODO(bwplotka): Add default to help.
		c = c.Default(f.defaultPath)
	}
	pathFlag := c.PlaceHolder("<file-path>").String()

	c = f.cmd.Flag(contentFlagName, fmt.Sprintf("Alternative to '%s' flag (lower priority). Content of %s", pathFlagName, f.help))
	if f.hiddenContent {
		c = c.Hidden()
	}
	if f.required {
		c = c.Required()
	}
	if f.defaultContent != "" {
		// TODO(bwplotka): Add default to help.
		c = c.Default(f.defaultContent)
	}
	contentFlag := c.PlaceHolder("<content>").String()

	return &PathOrContent{
		flagName: f.name,
		required: f.required,
		path:     pathFlag,
		content:  contentFlag,
	}
}

// Content returns content of the file. Flag that specifies path has priority.
// It returns error if the content is empty and required flag is set to true.
func (p *PathOrContent) Content() ([]byte, error) {
	contentFlagName := p.flagName
	fileFlagName := fmt.Sprintf("%s-file", p.flagName)

	if len(*p.path) > 0 && len(*p.content) > 0 {
		return nil, errors.Errorf("both %s and %s flags set.", fileFlagName, contentFlagName)
	}

	var content []byte
	if len(*p.path) > 0 {
		c, err := ioutil.ReadFile(*p.path)
		if err != nil {
			return nil, errors.Wrapf(err, "loading YAML file %s for %s", *p.path, fileFlagName)
		}
		content = c
	} else {
		content = []byte(*p.content)
	}

	if len(content) == 0 && p.required {
		return nil, errors.Errorf("flag %s or %s is required for running this command and content cannot be empty.", fileFlagName, contentFlagName)
	}

	return content, nil
}
