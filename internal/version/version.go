package version

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

var (
	version string
	tag     string 
	commit  string
	dirty   bool
)

func init() {
	// Always attempt to get build info from the Go runtime, regardless of build type.
	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "" {
			// If version wasn't set via ldflags, try to get it from VCS info
			version = getVersionFromVCS(info)
		}
		commit = getCommit(info)
		dirty = getDirty(info)
	}

	// Fallback version if nothing else is available
	if version == "" {
		version = "dev"
	}
}

func getVersionFromVCS(info *debug.BuildInfo) string {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			// Use commit hash as version if no proper version is available
			if len(setting.Value) >= 7 {
				return "dev-" + setting.Value[:7]
			}
		}
	}
	return ""
}

func getDirty(info *debug.BuildInfo) bool {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.modified" {
			return setting.Value == "true"
		}
	}
	return false
}

func getCommit(info *debug.BuildInfo) string {
	for _, setting := range info.Settings {
		if setting.Key == "vcs.revision" {
			return setting.Value[:7]
		}
	}
	return ""
}

// GetVersion returns the version of pulsebridge. This should be injected at build time
// using: -ldflags="-X 'wavezync/pulse-bridge/internal/version.version=vX.X.X'".
// If not provided, it will attempt to derive a version from VCS info.
func GetVersion() string {
	return version
}

// GetVersionWithBuildInfo returns detailed version information including
// commit hash, dirty status, and Go version. This will work when built
// within a Git checkout.
func GetVersionWithBuildInfo() string {
	var parts []string

	// Add the main version
	parts = append(parts, fmt.Sprintf("Version: %s", version))

	// Add build metadata
	var buildMetadata []string
	if tag != "" {
		buildMetadata = append(buildMetadata, tag)
	}
	if commit != "" {
		buildMetadata = append(buildMetadata, commit)
	}
	if dirty {
		buildMetadata = append(buildMetadata, "dirty")
	}

	if len(buildMetadata) > 0 {
		parts = append(parts, fmt.Sprintf("Build: %s", strings.Join(buildMetadata, ".")))
	}

	// Add Go version
	parts = append(parts, fmt.Sprintf("Go: %s", runtime.Version()))

	// Add architecture
	parts = append(parts, fmt.Sprintf("Platform: %s/%s", runtime.GOOS, runtime.GOARCH))

	return strings.Join(parts, "\n")
}
