module github.com/efficientgo/tools/copyright

go 1.15

require (
	github.com/efficientgo/tools/core v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/protoconfig/protoconfig/go v0.0.0-20210106192113-733758adefac
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

// Pin this to exactly same version as rest of repo?
replace github.com/efficientgo/tools/core => ../core
