package version

var (
	// These are set via -ldflags at build
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)
