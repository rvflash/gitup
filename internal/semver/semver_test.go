package semver_test

import (
	"github.com/rvflash/gitup/internal/semver"
	"testing"
)

var errTests = []struct {
	tag string // input
}{
	{" "},
	{"v"},
	{"v1.2"},
	{"1.2.3"},
	{"va.b.c"},
	{"v1.2.a"},
	{"v1.a.3"},
	{"v1.2.3-"},
	{"v1.2.3.4"},
	{"v1.2.3-+2"},
	{"v1+5114f85"},
	{"v1.2.3.beta2"},
	{"v1.2.3-beta2+"},
}

var okTests = []struct {
	tag     string         // input
	version semver.Version // expected result
}{
	{"v1.2.3", semver.Version{1, 2, 3, "", ""}},
	{" v1.2.3", semver.Version{1, 2, 3, "", ""}},
	{"v1.2.3+92", semver.Version{1, 2, 3, "", "92"}},
	{"v1.2.3-alpha", semver.Version{1, 2, 3, "alpha", ""}},
	{"v1.2.3-alpha.1", semver.Version{1, 2, 3, "alpha.1", ""}},
	{"v1.2.3-0.3.7", semver.Version{1, 2, 3, "0.3.7", ""}},
	{"v1.2.3-x.7.z.92", semver.Version{1, 2, 3, "x.7.z.92", ""}},
	{"v1.2.3-alpha+001", semver.Version{1, 2, 3, "alpha", "001"}},
	{"v1.2.3+20130313144700", semver.Version{1, 2, 3, "", "20130313144700"}},
	{"v1.2.3-beta+exp.sha.5114f85", semver.Version{1, 2, 3, "beta", "exp.sha.5114f85"}},
}

var cpErrTests = []struct {
	tag1 string // input tag #1
	tag2 string // input tag #2
}{
	{"1.2.3", "v1.2.3"},
	{"v1.2.3", "1.2.3"},
}

var cpOkTests = []struct {
	tag1 string              // input tag #1
	tag2 string              // input tag #2
	diff semver.Relationship // expected result
}{
	{"v1.2.3", "v1.2.3", semver.Relationship{}},
	{"v1.2.3", "v2.2.3", semver.Relationship{-1, 0, 0, "", "", -1}},
	{"v1.2.3", "v1.3.3", semver.Relationship{0, -1, 0, "", "", -1}},
	{"v1.2.3", "v1.2.4", semver.Relationship{0, 0, -1, "", "", -1}},
	{"v1.2.3", "v1.2.3-beta", semver.Relationship{0, 0, 0, "<>beta", "", -1}},
	{"v1.2.3-alpha", "v1.2.3-beta", semver.Relationship{0, 0, 0, "alpha<>beta", "", -1}},
	{"v1.2.3+92", "v1.2.3", semver.Relationship{0, 0, 0, "", "92<>", 0}},
	{"v1.2.3-beta+92", "v0.2.3-beta", semver.Relationship{1, 0, 0, "", "92<>", 1}},
	{"v1.2.3", "v1.2.3-beta+92", semver.Relationship{0, 0, 0, "<>beta", "<>92", -1}},
}

// TestCompare tests Compare method with invalid or valid Semantic Version.
func TestCompare(t *testing.T) {
	// Checks with various incorrect tags
	for _, te := range cpErrTests {
		if _, err := semver.Compare(te.tag1, te.tag2); err == nil {
			t.Errorf("Expected error with invalid versions to compare : %v vs %v", te.tag1, te.tag2)
		}
	}
	// Checks with valid tags, expected well-formed relations between these tags
	for _, to := range cpOkTests {
		if tag, err := semver.Compare(to.tag1, to.tag2); err == nil {
			if tag != to.diff {
				t.Errorf("Expected relation %v for %v vs %v, received: %v", to.diff, to.tag1, to.tag2, tag)
			}
		} else {
			t.Errorf("Expected valid comparaison between %v and %v, received error: %v", to.tag1, to.tag2, err)
		}
	}
}

// TestParse tests Parse method with invalid or valid Semantic Version's tags.
func TestParse(t *testing.T) {
	// Checks with various incorrect tags
	for _, te := range errTests {
		if _, err := semver.Parse(te.tag); err == nil {
			t.Errorf("Expected error with invalid version : %v", te.tag)
		}
	}
	// Checks with valid tags, expected well-formed versions
	for _, to := range okTests {
		if tag, err := semver.Parse(to.tag); err == nil {
			if tag != to.version {
				t.Errorf("Expected valid version %v for %v, received: %v", to.version, to.tag, tag)
			}
		} else {
			t.Errorf("Expected valid version for %v, received error: %v", to.tag, err)
		}
	}
}
