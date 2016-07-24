package gitflow

import (
	"errors"
	"os/exec"
	"strings"
)

const (
	errMsgUndefinedPath = "Directory path is undefined"
	errMsgUndefinedTag  = "Tag name is undefined"
)

type Repo struct {
	Path  string
	valid bool
}

// gitCheck returns err if path is not a valid Git repository.
func (r *Repo) gitCheck() (err error) {
	if r.valid == false {
		_, err = r.gitStatus()
	}
	return
}

// gitDescribe returns the most recent tag reachable for this directory path.
func (r *Repo) gitDescribe(commit string) (tag string, err error) {
	if err = r.gitCheck(); err != nil {
		return
	}
	args := []string{"-C", r.Path, "describe", "--abbrev=0", "--tags"}
	if commit = strings.TrimSpace(commit); commit != "" {
		args = append(args, commit)
	}
	var ref []byte
	if ref, err = exec.Command("git", args...).Output(); err == nil {
		tag = strings.TrimSpace(string(ref))
	}
	return
}

// gitFetch returns in error if it fails to update local tag list
func (r *Repo) gitFetch() (err error) {
	if err = r.gitCheck(); err != nil {
		return
	}
	return exec.Command("git", "-C", r.Path, "fetch", "--tags").Run()
}

// gitStatus returns the working tree status.
func (r *Repo) gitStatus() (status []byte, err error) {
	if r.Path = strings.TrimSpace(r.Path); r.Path == "" {
		err = errors.New(errMsgUndefinedPath)
	} else {
		status, err = exec.Command("git", "-C", r.Path, "status").Output()
	}
	return
}

// LocalTag returns the most recent tag reachable from the defined repository.
func (r *Repo) LocalTag() (string, error) {
	return r.gitDescribe("")
}

// LastTag returns the last available tag of the Git repository.
func (r *Repo) LastTag() (tag string, err error) {
	// Get new tags from the remote
	if err = r.gitFetch(); err == nil {
		// Get the latest commit of tag list
		var commit []byte
		if commit, err = exec.Command("git", "-C", r.Path, "rev-list", "--tags", "--max-count=1").Output(); err == nil {
			// Get the tag name for this commit
			tag, err = r.gitDescribe(string(commit))
		}
	}
	return
}

// CheckoutTag returns an error if it can not switch the repository on this tag
func (r *Repo) CheckoutTag(tag string) (err error) {
	if err = r.gitCheck(); err != nil {
		return
	} else if tag = strings.TrimSpace(tag); tag == "" {
		return errors.New(errMsgUndefinedTag)
	}
	return exec.Command("git", "-C", r.Path, "checkout", tag).Run()
}

// IsRepository return true if the repo is a valid Git repository.
func (r *Repo) IsRepository() bool {
	if err := r.gitCheck(); err != nil {
		return r.valid
	}
	return false
}
