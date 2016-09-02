package gitup

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/rvflash/gitup/internal/gitflow"
	"github.com/rvflash/gitup/internal/semver"
	"os"
)

const (
	Noop   = iota // 0
	Manual        // 1
	Auto          // 2
)

const (
	MajorVersion      = iota // 0
	MinorVersion             // 1
	PatchVersion             // 2
	PreReleaseVersion        // 3
	BuildMetadata            // 4
)

const (
	errMsgVersion         = "unknown type of version"
	errMsgAction          = "unknown action's type"
	errMsgInDemand        = "no available update"
	errMsgConfirm         = "only accepts yes or no as valid response"
	errMsgDowngradeAction = "unable to downgrade behavior on minor versions"
)

type GitFlow interface {
	LocalTag() (string, error)
	LastTag() (string, error)
	CheckoutTag(string) error
}

type Repo struct {
	git           GitFlow
	diff          semver.Relationship
	local, remote string
	upStrategy    uint8
}

type UpdateStrategy struct {
	until [4]uint8
	// soon, we will also manage retryLater.
}

// Enable testing by mocking *os.File.
var stdin = os.Stdin

// Enable testing by mocking *gitflow.Repo.
var gitRepo = gitflow.NewRepo

// NewRepo starts a new Git repository.
func NewRepo(path string) (*Repo, error) {
	if git, err := gitRepo(path); err != nil {
		return nil, err
	} else {
		return &Repo{git: git}, nil
	}
}

// AddStrategy starts a new Git repository.
func (s *UpdateStrategy) AddStrategy(version, action uint8) (err error) {
	if version >= BuildMetadata {
		err = errors.New(errMsgVersion)
	} else if action < Noop || action > Auto {
		err = errors.New(errMsgAction)
	} else if action < s.getStrategy(int8(version)) {
		err = errors.New(errMsgDowngradeAction)
	} else {
		s.until[version] = action
	}
	return
}

// InDemand returns true if the Git repository needs to be updated because it is not on the latest tag.
func (r *Repo) InDemand(s UpdateStrategy) bool {
	var err error
	// Gets local version
	if r.local == "" {
		if r.local, err = r.git.LocalTag(); err != nil {
			return false
		}
	}
	// Gets latest remote version
	if r.remote == "" {
		if r.remote, err = r.git.LastTag(); err != nil {
			return false
		}
	}
	// Gets differences between local and remote tags
	if r.diff, err = semver.Compare(r.local, r.remote); err != nil {
		return false
	}
	// Defines strategy to use by type of difference: major strategy by passing minor, etc.
	if r.diff.Major < 0 {
		if r.upStrategy = s.getStrategy(MajorVersion); r.upStrategy > Noop {
			return true
		}
	} else if r.diff.Minor < 0 {
		if r.upStrategy = s.getStrategy(MinorVersion); r.upStrategy > Noop {
			return true
		}
	} else if r.diff.Patch < 0 {
		if r.upStrategy = s.getStrategy(PatchVersion); r.upStrategy > Noop {
			return true
		}
	} else if r.diff.PreRelease != "" {
		if r.upStrategy = s.getStrategy(PreReleaseVersion); r.upStrategy > Noop {
			return true
		}
	}
	return false
}

// Update returns an error if it can not to update Git repository with the latest tag.
func (r *Repo) Update(s UpdateStrategy) error {
	if !r.InDemand(s) {
		return errors.New(errMsgInDemand)
	}
	// Manual update required, demands authorisation to user
	if r.upStrategy == Manual {
		// Display a message in order to inform about the available update.
		fmt.Printf("You are currenly on the '%v', a new version is available.\n", r.local)
		fmt.Printf("Do you want to update and move on '%v'?\n", r.remote)
		if !confirmUpdate() {
			return nil
		}
	}
	// Checkout it on the local repository
	return r.git.CheckoutTag(r.remote)
}

// getStrategy returns for the type of version (major, minor, etc.), the action to perform.
func (s *UpdateStrategy) getStrategy(version int8) (action uint8) {
	if version < MajorVersion || version > PreReleaseVersion {
		return
	}
	for i := version; i >= MajorVersion; i-- {
		if s.until[i] > action {
			action = uint8(s.until[i])
		}
	}
	return
}

// confirmUpdate reads in console input the user response.
// It returns true if user responds yes or y, false otherwise.
// It deliberately ignores input errors
func confirmUpdate() bool {
	input := bufio.NewScanner(stdin)
	for input.Scan() {
		if ok, err := parseConfirm(input.Text()); err != nil {
			fmt.Println(err)
		} else {
			return ok
		}
	}
	return false
}

// parseConfirm returns the boolean value represented by the string.
// It accepts y, Y, yes, Yes, YES, n, N, no, No, NO.
// Any other value returns an error.
func parseConfirm(str string) (bool, error) {
	switch str {
	case "y", "Y", "yes", "Yes", "YES":
		return true, nil
	case "n", "N", "no", "No", "NO":
		return false, nil
	}
	return false, errors.New(errMsgConfirm)
}
