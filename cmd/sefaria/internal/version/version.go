package version

var (
	Version   = "0.0.0-unknown"
	Commit    = "unknown"
	BuildDate = "unknown"
	Tag       = "0.0.0"
)

func String() string {
	return Version
}
