module github.com/efficientgo/tools/copyright

go 1.15

require (
	github.com/efficientgo/tools/core v0.0.0-00010101000000-000000000000
	github.com/efficientgo/tools/kingpin v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

replace (
	github.com/efficientgo/tools/core => ../core
	// Pin those to exactly same version as rest of repo.
	github.com/efficientgo/tools/kingpin => ../kingpin
)
