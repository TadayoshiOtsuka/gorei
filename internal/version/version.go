package version

import (
	"fmt"
	"runtime/debug"
)

var version = ""

func Exec() error {
	fmt.Printf("%v\n", getBuildVersion())

	return nil
}

func getBuildVersion() string {
	if version != "" {
		return version
	}
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	return i.Main.Version
}
