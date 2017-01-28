package gitflow

import (
	"errors"
	"os/exec"
	"strings"
)

const (
	gitTagFolder        = "tags/"
	errMsgUndefinedPath = "directory path is undefined"
	errMsgUndefinedTag  = "tag name is undefined"
)

// Repo represents a Git repository.
type Repo struct {
	path  string
	valid bool
}

// Enable testing by mocking *exec.Cmd.
var execCommand = exec.Command

// NewRepo starts a new Git repository.
func NewRepo(path string) (*Repo, error) {
	if path = strings.TrimSpace(path); path == "" {
		return nil, errors.New(errMsgUndefinedPath)
	}
	r := &Repo{path: path}
	if err := r.gitCheck(); err != nil {
		return nil, err
	}
	return r, nil
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
		if commit, err = execCommand("git", "-C", r.path, "rev-list", "--tags", "--max-count=1").Output(); err == nil {
			// Get the tag name for this commit
			tag, err = r.gitDescribe(string(commit))
		}
	}
	return
}

// CheckoutTag returns an error if it can not switch the repository on the given tag.
func (r *Repo) CheckoutTag(tag string) error {
	if tag = strings.TrimSpace(tag); tag == "" {
		return errors.New(errMsgUndefinedTag)
	}
	return r.gitCheckout(gitTagFolder + tag)
}

// gitCheck returns err if path is not a valid Git repository.
func (r *Repo) gitCheck() (err error) {
	if r.valid == false {
		_, err = r.gitStatus()
	}
	return
}

// gitCheckout returns an error if it can switch to the given branch or restore it.
func (r *Repo) gitCheckout(branch string) (err error) {
	if err = r.gitCheck(); err != nil {
		return
	}
	args := []string{"-C", r.path, "checkout"}
	if branch = strings.TrimSpace(branch); branch != "" {
		args = append(args, branch)
	}
	return execCommand("git", args...).Run()
}

// gitDescribe returns the most recent tag reachable for this directory path.
func (r *Repo) gitDescribe(commit string) (tag string, err error) {
	if err = r.gitCheck(); err != nil {
		return
	}
	args := []string{"-C", r.path, "describe", "--abbrev=0", "--tags"}
	if commit = strings.TrimSpace(commit); commit != "" {
		args = append(args, commit)
	}
	var ref []byte
	if ref, err = execCommand("git", args...).Output(); err == nil {
		tag = strings.TrimSpace(string(ref))
	}
	return
}

// gitFetch returns in error if it fails to update local tag list.
func (r *Repo) gitFetch() (err error) {
	if err = r.gitCheck(); err != nil {
		return
	}
	return execCommand("git", "-C", r.path, "fetch", "--tags").Run()
}

// gitStatus returns the working tree status.
func (r *Repo) gitStatus() ([]byte, error) {
	return execCommand("git", "-C", r.path, "status").Output()
}
