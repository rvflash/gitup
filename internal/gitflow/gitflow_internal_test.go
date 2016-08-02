package gitflow

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	errPathTest   = "/home/error/path"
	okPathTest    = "/home/fake/path/to/git/repository"
	commitTest    = "9b7f1bbc8d82ef98bbb15e86f3ccb704ec35720a"
	remoteTagTest = "v1.2.4"
	tagTest       = "v1.2.3"
)

var errPathTests = []struct {
	path string // input
}{
	{""},                // empty path
	{" "},               // empty path (only space)
	{errPathTest},       // path
	{" " + errPathTest}, // path with leading space
	{errPathTest + " "}, // path ending by space
}

var okPathTests = []struct {
	path string // input
}{
	{okPathTest},       // path
	{" " + okPathTest}, // path with leading space
	{okPathTest + " "}, // path ending by space
}

var commitTests = []struct {
	id  string // input
	tag string // output
}{
	{"", tagTest},                     // no commit
	{" ", tagTest},                    // no commit (ony space)
	{" " + commitTest, remoteTagTest}, // commit with leading space
	{commitTest, remoteTagTest},       // commit
}

var tagTests = []struct {
	tag string // input
}{
	{""},
	{" " + gitTagFolder + remoteTagTest},
	{gitTagFolder + remoteTagTest + " "},
}

// fakeExecCommand returns a mock of the exec command.
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestNewRepo tests the methiod to return a structure for the Git repository.
func TestNewRepo(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with various incorrect paths.
	for _, tp := range errPathTests {
		if _, err := NewRepo(tp.path); err == nil {
			t.Errorf("Expected error with invalid path '%v'", tp.path)
		}
	}

	// Checks with valid paths
	for _, tp := range okPathTests {
		if r, err := NewRepo(tp.path); err != nil {
			t.Errorf("Expected no error with valid path '%v', got: %v", tp.path, err)
		} else if r.path != okPathTest {
			t.Errorf("Expected a valid repository with path '%v'", tp.path)
		}
	}
}

// TestRepo_CheckoutTag tests the method dedicated to checkout the given tag on the current repository.
func TestRepo_CheckoutTag(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with valid paths
	for _, tp := range okPathTests {
		if r, err := NewRepo(tp.path); err != nil {
			t.Errorf("Expected no error with valid path '%v', got: %v", tp.path, err)
		} else if err := r.CheckoutTag(""); err == nil {
			t.Error("Expected error with empty commit")
		} else if err := r.CheckoutTag(remoteTagTest); err != nil {
			t.Errorf("Expected no error with valid tag '%v', got: %v", remoteTagTest, err)
		} else if err := r.CheckoutTag(tagTest); err == nil {
			t.Error("Expected error with unknown tag on local repository")
		}
	}
}

// TestRepo_LocalTag tests the method dedicated to get the local tag of current repository.
func TestRepo_LocalTag(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with valid paths
	for _, tp := range okPathTests {
		if r, err := NewRepo(tp.path); err != nil {
			t.Errorf("Expected no error with valid path '%v', got: %v", tp.path, err)
		} else if tag, err := r.LocalTag(); err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		} else if tag != tagTest {
			t.Errorf("Expected tag '%v', got '%v'", tagTest, tag)
		}
	}
}

// TestRepo_LastTag tests the method dedicated to get the latest remote tag of current repository.
func TestRepo_LastTag(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with valid paths
	for _, tp := range okPathTests {
		if r, err := NewRepo(tp.path); err != nil {
			t.Errorf("Expected no error with valid path '%v', got: %v", tp.path, err)
		} else if tag, err := r.LastTag(); err != nil {
			t.Errorf("Expected no error, got '%v'", err)
		} else if tag != remoteTagTest {
			t.Errorf("Expected tag '%v', got '%v'", remoteTagTest, tag)
		}
	}
}

// TestGitCheck tests the internal method dedicated to verify if the given path is a Git repository.
func TestGitCheck(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with incorrect path.
	r := new(Repo)
	r.path = errPathTest
	if err := r.gitCheck(); err == nil {
		t.Errorf("Expected error with invalid path '%v'", errPathTest)
	}

	// Checks with valid path
	r = new(Repo)
	r.path = okPathTest
	if err := r.gitCheck(); err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}

// TestGitCheckout tests the internal method dedicated to git checkout.
func TestGitCheckout(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with incorrect path.
	r := new(Repo)
	r.path = errPathTest
	for _, c := range tagTests {
		if err := r.gitCheckout(c.tag); err == nil {
			t.Errorf("Expected error with branch '%v' on invalid Git path '%v'", c.tag, errPathTest)
		}
	}
	// Checks with valid path
	r = new(Repo)
	r.path = okPathTest
	for _, c := range tagTests {
		if err := r.gitCheckout(c.tag); err != nil {
			t.Errorf("Expected no error with valid path '%v' and branch '%v', got: %v", okPathTest, c.tag, err)
		}
	}
}

// TestGitDescribe tests the internal method dedicated to git describe.
func TestGitDescribe(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with incorrect path.
	r := new(Repo)
	r.path = errPathTest
	for _, c := range commitTests {
		if _, err := r.gitDescribe(c.id); err == nil {
			t.Errorf("Expected error with commit '%v' on invalid Git path '%v'", c.id, errPathTest)
		}
	}
	// Checks with valid path
	r = new(Repo)
	r.path = okPathTest
	for _, c := range commitTests {
		if out, err := r.gitDescribe(c.id); err != nil {
			t.Errorf("Expected no error with valid path '%v' and commit '%v', got: %v", errPathTest, c.id, err)
		} else if tag := string(out); tag != c.tag {
			t.Errorf("Expected tag '%v' for the local Git repository '%v', got: %v", c.tag, okPathTest, tag)
		}
	}
}

// TestGitFetch tests the internal method dedicated to git fetch.
func TestGitFetch(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with incorrect path.
	r := new(Repo)
	r.path = errPathTest
	if err := r.gitFetch(); err == nil {
		t.Errorf("Expected error on invalid Git path '%v'", errPathTest)
	}
	// Checks with valid path
	r = new(Repo)
	r.path = okPathTest
	if err := r.gitFetch(); err != nil {
		t.Errorf("Expected no error with valid path '%v', got: %v", okPathTest, err)
	}
}

// TestGitStatus tests the internal method dedicated to git status.
func TestGitStatus(t *testing.T) {
	execCommand = fakeExecCommand

	// Restore exec command behavior at the end of the test.
	defer func() { execCommand = exec.Command }()

	// Checks with incorrect path.
	r := new(Repo)
	r.path = errPathTest
	if _, err := r.gitStatus(); err == nil {
		t.Errorf("Expected error with invalid path '%v'", errPathTest)
	}

	// Checks with valid path
	r = new(Repo)
	r.path = okPathTest
	if out, err := r.gitStatus(); err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	} else if string(out) == "" {
		t.Errorf("Expected status message about the valid local Git repository: %v", okPathTest)
	}
}

// TestHelperProcess mocks exec commands and responds instead of the command git.
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	// Extract only exec arguments.
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No command\n")
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	// Manage only git command.
	if cmd != "git" && args[0] != "-C" {
		fmt.Fprintf(os.Stderr, "fatal: Not a Git command, received: %v\n", cmd)
		os.Exit(1)
	}
	// Manage exit status on error on "invalid" Git path.
	if strings.HasPrefix(args[1], errPathTest) {
		fmt.Fprintf(os.Stderr, "fatal: Not a git repository %v (or any of the parent directories): .git\n", args[1])
		os.Exit(1)
	}
	// Manage each git sub-commands.
	switch args[2] {
	case "checkout":
		switch len(args) {
		case 4:
			if args[3] == gitTagFolder+remoteTagTest {
				fmt.Fprintf(os.Stdout, "note: checking out '%v'.", remoteTagTest)
			} else {
				fmt.Fprintf(os.Stderr, "error: pathspec '%v' did not match any file(s) known to git.\n", args[3])
				os.Exit(1)
			}
		case 3:
			fmt.Fprintf(os.Stdout, "Your branch is up-to-date with '%v%v'.", gitTagFolder, tagTest)
		default:
			fmt.Fprintf(os.Stderr, "fatal: Not a git repository (or any of the parent directories): .git\n", args[1])
			os.Exit(1)
		}
	case "describe":
		if args[3] == "--abbrev=0" && args[4] == "--tags" {
			if len(args) == 6 {
				fmt.Fprint(os.Stdout, remoteTagTest+"\n")
			} else {
				fmt.Fprint(os.Stdout, tagTest+"\n")
			}
		}
	case "fetch":
		fmt.Fprint(os.Stdout, "\n")
	case "status":
		if len(args) == 3 {
			fmt.Fprint(os.Stdout, "On branch stable\n")
		}
	case "rev-list":
		if args[3] == "--tags" && args[4] == "--max-count=1" {
			fmt.Fprint(os.Stdout, commitTest+"\n")
		}
	default:
		fmt.Fprintf(os.Stderr, "fatal: Not a git sub-command (%v)\n", args[2])
		os.Exit(1)
	}
}
