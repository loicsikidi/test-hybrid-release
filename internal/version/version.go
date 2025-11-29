package version

import (
	"fmt"
	"runtime/debug"
	"sync"
)

// Variables set via ldflags by build tools (GoReleaser, Nix).
// These serve as fallback when VCS info is not available in debug.ReadBuildInfo().
var (
	version  string // Set via: -X github.com/loicsikidi/test-hybrid-release/internal/version.version=v1.0.0
	revision string // Set via: -X github.com/loicsikidi/test-hybrid-release/internal/version.revision=abc123
)

type Version struct {
	Revision, Version, Time string
	Dirty                   bool
}

var ver = Version{
	Revision: "unknown",
	Version:  "unknown",
	Time:     "unknown",
	Dirty:    false,
}

func (v Version) String() string {
	return fmt.Sprintf(`Revision: %s
Version: %s
BuildTime: %s
Dirty: %t`, v.Revision, v.Version, v.Time, v.Dirty)
}

// Get retrieves version information from Go build info or ldflags.
//
// Priority order:
//  1. VCS info from [debug.ReadBuildInfo] (available when built with "go build" from a git repo)
//  2. Fallback to ldflags set by build tools (GoReleaser, Nix)
//  3. Default to "unknown"
func Get() Version {
	sync.OnceFunc(func() {
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			fallbackToLdflags()
			return
		}

		// Try to get version from build info
		if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			ver.Version = bi.Main.Version
		}

		// Try to get VCS info from build settings
		hasVCSInfo := false
		for _, setting := range bi.Settings {
			switch setting.Key {
			case "vcs.revision":
				ver.Revision = setting.Value
				hasVCSInfo = true
			case "vcs.time":
				ver.Time = setting.Value
			case "vcs.modified":
				ver.Dirty = setting.Value == "true"
			}
		}

		// If no VCS info found in build info, fallback to ldflags
		if !hasVCSInfo {
			fallbackToLdflags()
		}
	})()
	return ver
}

func fallbackToLdflags() {
	if version != "" {
		ver.Version = version
	}
	if revision != "" {
		ver.Revision = revision
	}
}
