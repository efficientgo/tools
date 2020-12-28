# tools

Set of tools, packages and libraries for debugging, tracking and drilling down resource usage regressions for Go code.

## Release model

Since this is meant to be critical, tiny import, multi module toolset, there are currently no semver releases planned. It's designed to pin modules via git commits, all commits to master should be stable and properly tests, vetted and linted.

## Modules

### `github.com/efficientgo/tools/core`

Main module containing set of useful, core packages for testing, closing, running and repeating.

This module is optimized for almost zero dependencies for ease of use.

### `github.com/efficientgo/tools/copyright`

This module is a very simple CLI for ensuring copyright header on code files.

```bash mdox-gen-exec="sh -c 'cd copyright && go run copyright.go --help || exit 0'"
usage: copyright [<flags>] [<files>...]

copyright

Flags:
  --help                        Show context-sensitive help (also try
                                --help-long and --help-man).
  --copyright-file=<file-path>  Path to Copyright content to apply to provided
                                files
  --copyright=<content>         Alternative to 'copyright-file' flag (lower
                                priority). Content of Copyright content to apply
                                to provided files

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

### `github.com/efficientgo/tools/kingpin`

Extra flag types for popular kingpin flag parsing library.
