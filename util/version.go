package util

import "fmt"

type BuildInfo struct {
	Version string
	Commit  string
	Date    string
}

func (v *BuildInfo) String() string {
	return fmt.Sprintf("version: %s | commit: %s | date: %s", v.Version, v.Commit, v.Date)
}
