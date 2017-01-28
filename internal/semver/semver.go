package semver

import (
	"errors"
	"strconv"
	"strings"
)

const (
	errMsgSemanticVersion = "not a valid semantic version"
)

// Version represents a semantic versioning like 2.0.0.
// @see http://semver.org/spec/v2.0.0.html
// @example v1.2.3
type Version struct {
	Major, Minor, Patch uint8
	PreRelease, Build   string
}

// Relationship represents the difference between two versions.
type Relationship struct {
	Major, Minor, Patch int8
	PreRelease, Build   string
	Upstream            int8
}

// Compare returns the difference for each type and the relation between two versions.
func Compare(tag1 string, tag2 string) (version Relationship, err error) {
	var v1, v2 Version
	if v1, err = Parse(tag1); err != nil {
		return
	}
	if v2, err = Parse(tag2); err != nil {
		return
	}
	// Build metadata should be ignored when determining version precedence, so ...
	if bp1 := strings.Index(tag1, "+"); -1 < bp1 {
		tag1 = tag1[:bp1]
	}
	if bp2 := strings.Index(tag2, "+"); -1 < bp2 {
		tag2 = tag2[:bp2]
	}
	if tag1 > tag2 {
		version.Upstream = 1
	} else if tag1 < tag2 {
		version.Upstream = -1
	}
	version.Major = int8(v1.Major) - int8(v2.Major)
	version.Minor = int8(v1.Minor) - int8(v2.Minor)
	version.Patch = int8(v1.Patch) - int8(v2.Patch)
	if v1.PreRelease != v2.PreRelease {
		version.PreRelease = v1.PreRelease + "<>" + v2.PreRelease
	}
	if v1.Build != v2.Build {
		version.Build = v1.Build + "<>" + v2.Build
	}
	return
}

// Parse returns all parts of a tag in a Version's struct.
func Parse(tag string) (version Version, err error) {
	if tag = strings.Trim(tag, " "); tag == "" {
		err = errors.New(errMsgSemanticVersion)
		return
	}
	// Tag must start with "v"
	if !strings.HasPrefix(tag, "v") {
		err = errors.New(errMsgSemanticVersion)
		return
	}
	// Tag must have 3 parts: major.minor.path
	ver := strings.SplitN(tag, ".", 3)
	if len(ver) != 3 {
		err = errors.New(errMsgSemanticVersion)
		return
	}
	// Major
	if version.Major, err = toVersionNumber(ver[0][1:]); err != nil {
		err = errors.New(errMsgSemanticVersion)
		return
	}
	// Minor
	if version.Minor, err = toVersionNumber(ver[1]); err != nil {
		err = errors.New(errMsgSemanticVersion)
		return
	}
	// Last member may contain path, pre-release and build metadata
	var ep, pr, bm int
	ep = len(ver[2])
	pr = strings.Index(ver[2], "-")
	bm = strings.Index(ver[2], "+")
	if pr > -1 && (bm == -1 || pr < bm) {
		// With pre-release version
		pre := bm
		if pre == -1 {
			pre = ep
		}
		ep = pr
		pr++
		if version.PreRelease = ver[2][pr:pre]; version.PreRelease == "" {
			err = errors.New(errMsgSemanticVersion)
			return
		}
	}
	if bm > -1 {
		// With build metadata
		if ep > bm {
			ep = bm
		}
		bm++
		if version.Build = ver[2][bm:]; version.Build == "" {
			err = errors.New(errMsgSemanticVersion)
			return
		}
	}
	// Path
	if version.Patch, err = toVersionNumber(ver[2][:ep]); err != nil {
		err = errors.New(errMsgSemanticVersion)
	}
	return
}

// toVersionNumber returns a uint8 number for a string version.
func toVersionNumber(version string) (number uint8, err error) {
	var num uint64
	if num, err = strconv.ParseUint(version, 10, 8); err == nil {
		number = uint8(num)
	}
	return
}
