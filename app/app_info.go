package app

var (
	Version = "unknown"
	Commit  = "unknown"
	Created = "unknown"
)

type BuildInfo struct {
	Version string `json:"version"`
	Created string `json:"created"`
	Commit  string `json:"commit"`
}

func NewBuildInfo() *BuildInfo {
	return &BuildInfo{
		Version: Version,
		Created: Created,
		Commit:  Commit,
	}
}
