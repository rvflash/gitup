package gitup

import (
	"errors"
	"github.com/rvflash/gitup/internal/gitflow"
	"io/ioutil"
	"os"
	"testing"
)

const (
	errMsgFake = "fake error"
	errValue   = "error"
)

// FakeGitFlow implements GitFlow to order to mock *gitflow.Repo.
type FakeGitFlow struct {
	localError, remoteError, checkoutError bool
	localTag, remoteTag                    string
}

var repoTests = []struct {
	git      *FakeGitFlow
	strategy UpdateStrategy
	inDemand bool
	stdin    string
}{
	{&FakeGitFlow{true, true, false, "", ""}, UpdateStrategy{}, false, ""}, // Invalid entries, errors
	{&FakeGitFlow{false, true, false, "v1.0.0", ""}, UpdateStrategy{}, false, ""},
	{&FakeGitFlow{true, false, false, "", "v1.0.0"}, UpdateStrategy{}, false, ""},
	{&FakeGitFlow{false, false, true, "v1.0.0", "v1.0.0"}, UpdateStrategy{}, false, ""},
	{&FakeGitFlow{false, false, false, "v1.0", "v1.0.0"}, UpdateStrategy{}, false, ""},
	{&FakeGitFlow{false, false, true, "v1.0.0", "v1.1.0"}, UpdateStrategy{[4]uint8{Auto}}, true, ""},
	{&FakeGitFlow{false, false, false, "v1.0.0", "v1.0.0"}, UpdateStrategy{}, false, ""}, // Valid entries, fails
	{&FakeGitFlow{false, false, false, "v1.0.0", "v2.0.0"}, UpdateStrategy{}, false, ""},
	{&FakeGitFlow{false, false, false, "v2.0.0", "v1.0.0"}, UpdateStrategy{}, false, ""},
	{&FakeGitFlow{false, false, false, "v1.0.0", "v1.0.1"}, UpdateStrategy{[4]uint8{Auto}}, true, ""}, // Valid entries, successful
	{&FakeGitFlow{false, false, false, "v1.0.0-alpha", "v1.0.0-beta"}, UpdateStrategy{[4]uint8{Auto}}, true, ""},
	{&FakeGitFlow{false, false, false, "v1.0.0", "v2.0.0"}, UpdateStrategy{[4]uint8{Manual}}, true, "y"},
	{&FakeGitFlow{false, false, false, "v1.0.0", "v2.0.0"}, UpdateStrategy{[4]uint8{Manual}}, true, "n"},
}

var confirmTests = []struct {
	str       string // input
	ok, onErr bool   // expected result
}{
	{"y", true, false}, // yes
	{"Y", true, false},
	{"yes", true, false},
	{"Yes", true, false},
	{"YES", true, false},
	{"n", false, false}, // no
	{"N", false, false},
	{"no", false, false},
	{"No", false, false},
	{"NO", false, false},
	{" yes", false, true}, // invalid
	{"n ", false, true},
}

var strategyTests = []struct {
	version, action uint8 // input
	onErr           bool  // expected result
}{
	{MajorVersion, Noop, false},
	{MinorVersion, Noop, false},
	{PatchVersion, Noop, false},
	{PreReleaseVersion, Noop, false},
	{BuildMetadata, Noop, true}, // Build metadata is ignored when determining version precedence, expected error
	{5, Noop, true},             // Unknown version type, expected error
	{MinorVersion, Manual, false},
	{MinorVersion, Auto, false},
	{MinorVersion, 3, true},      // Unknown action, expected error
	{PatchVersion, Manual, true}, // Due to the previous setting in success, we can not downgrade behavior for minor versions
}

var actionTests = []struct {
	strategy, actions [4]uint8
}{
	{[4]uint8{}, [4]uint8{}},
	{[4]uint8{Noop}, [4]uint8{Noop}},
	{[4]uint8{Manual}, [4]uint8{Manual, Manual, Manual, Manual}},
	{[4]uint8{Noop, Manual}, [4]uint8{Noop, Manual, Manual, Manual}},
	{[4]uint8{Manual, Auto}, [4]uint8{Manual, Auto, Auto, Auto}},
	{[4]uint8{Noop, Manual, Auto}, [4]uint8{Noop, Manual, Auto, Auto}},
	{[4]uint8{Noop, Manual, Auto, Noop}, [4]uint8{Noop, Manual, Auto, Auto}},
	{[4]uint8{Noop, Auto, Manual, Noop}, [4]uint8{Noop, Auto, Auto, Auto}},
	{[4]uint8{Noop, Noop, Manual, Auto}, [4]uint8{Noop, Noop, Manual, Auto}},
	{[4]uint8{Noop, Noop, Auto, Noop}, [4]uint8{Noop, Noop, Auto, Auto}},
	{[4]uint8{Noop, Noop, Noop, Manual}, [4]uint8{Noop, Noop, Noop, Manual}},
	{[4]uint8{Auto, Auto, Noop, Noop}, [4]uint8{Auto, Auto, Auto, Auto}},
}

// LocalTag mocks the gitflow's method LocalTag() on FakeGitFlow struct.
func (r FakeGitFlow) LocalTag() (string, error) {
	if r.localError {
		return "", errors.New(errMsgFake)
	}
	return r.localTag, nil
}

// LastTag mocks the gitflow's method LastTag() on FakeGitFlow struct.
func (r FakeGitFlow) LastTag() (string, error) {
	if r.remoteError {
		return "", errors.New(errMsgFake)
	}
	return r.remoteTag, nil
}

// CheckoutTag mocks the gitflow's method CheckoutTag() on FakeGitFlow struct.
func (r FakeGitFlow) CheckoutTag(string) error {
	if r.checkoutError {
		return errors.New(errMsgFake)
	}
	return nil
}

// fakeStdin returns a temporary file with required content to mock stdin.
// os.Stdin as a file implements *os.File interface.
func fakeStdin(str string) (stdin *os.File, err error) {
	// Creates temporary file.
	if stdin, err = ioutil.TempFile(os.TempDir(), "stdin"); err != nil {
		return
	}
	// Mocks user interaction by writing required string into it.
	if _, err = stdin.WriteString(str); err != nil {
		return
	}
	// Opens it in order to scan it as stdin.
	return os.Open(stdin.Name())
}

// TestNewRepo tests NewRepo method with various values and uses mock to get a fake Git repository.
func TestNewRepo(t *testing.T) {
	gitRepo = func(path string) (*gitflow.Repo, error) {
		if path == errValue {
			return nil, errors.New(errMsgFake)
		}
		return &gitflow.Repo{}, nil
	}
	// Restore new repo behavior at the end of the test.
	defer func() { gitRepo = gitflow.NewRepo }()

	// Checks with various type of path
	for _, pt := range []string{errValue, "/repo"} {
		if _, err := NewRepo(pt); err != nil {
			if pt != errValue {
				t.Errorf("Expected no error with valid path '%v', got: %v", pt, err)
			}
		} else if pt == errValue {
			t.Errorf("Expected an error with invalid path '%v'", pt)
		}
	}
}

// TestRepo_InDemand tests InDemand method with various valid or invalid values
func TestRepo_InDemand(t *testing.T) {
	for _, rt := range repoTests {
		r := &Repo{git: rt.git}
		if r.InDemand(rt.strategy) {
			if !rt.inDemand {
				t.Errorf("Expected no update with repository %#v and strategy %#v", rt.git, rt.strategy)
			}
		} else if rt.inDemand {
			t.Errorf("Expected an update with repository %#v and strategy %#v", rt.git, rt.strategy)
		}
	}
}

// TestRepo_Update tests Update method with various valid or invalid values
func TestRepo_Update(t *testing.T) {
	// Restore stdin source file at the end of the test.
	defer func() { stdin = os.Stdin }()

	for _, rt := range repoTests {
		r := &Repo{git: rt.git}
		if rt.stdin != "" {
			stdin, _ = fakeStdin(rt.stdin)
		}
		if err := r.Update(rt.strategy); err == nil {
			if !rt.inDemand || rt.git.checkoutError {
				t.Errorf("Expected error and no update with repository %#v and strategy %#v", rt.git, rt.strategy)
			}
		} else if rt.inDemand && !rt.git.checkoutError {
			t.Errorf("Expected an update with repository %#v and strategy %#v", rt.git, rt.strategy)
		}
	}
}

// TestAddStrategy tests AddStrategy method with various values.
func TestAddStrategy(t *testing.T) {
	s := new(UpdateStrategy)
	for _, st := range strategyTests {
		if err := s.AddStrategy(st.version, st.action); err == nil {
			if st.onErr {
				t.Errorf("Expected error for the version %v with action %v", st.version, st.action)
			} else if s.until[st.version] != st.action {
				t.Errorf("Expected result %v for the version %v, received: %v", st.action, st.version, s.until[st.version])
			}
		} else if !st.onErr {
			t.Errorf("Expected no error for the version %v with action %v, received: %v", st.version, st.action, err)
		}
	}
}

// TestGetStrategy tests getStrategy method.
func TestGetStrategy(t *testing.T) {
	s := new(UpdateStrategy)
	// Checks with invalids version's type but valid entries for this method.
	if Noop != s.getStrategy(-1) || Noop != s.getStrategy(BuildMetadata) {
		t.Errorf("Expected no operation with unknown version type")
	}
	// Checks with valid bounces
	for _, ac := range actionTests {
		s.until = ac.strategy
		for i := MajorVersion; i < BuildMetadata; i++ {
			if a := s.getStrategy(int8(i)); a != ac.actions[i] {
				t.Errorf("Expected action %v for version %v with strategy %v, received: %v", ac.actions[i], i, ac.strategy, a)
			}
		}
	}
}

// TestConfirmUpdate tests confirmUpdate method by mocking stdin
func TestConfirmUpdate(t *testing.T) {
	// Restore stdin source file at the end of the test.
	defer func() { stdin = os.Stdin }()

	var err error
	for _, cf := range confirmTests {
		if stdin, err = fakeStdin(cf.str); err != nil {
			t.Fatalf("Unable to mock stdin, received error: %v", err)
		}
		if ok := confirmUpdate(); cf.ok != ok {
			t.Errorf("Expected result %t for confirm with %v, received: %t", cf.ok, cf.str, ok)
		}
		if err = os.Remove(stdin.Name()); err != nil {
			t.Fatalf("Unable to remove stdin mock file, received error: %v", err)
		}
	}
}

// TestParseConfirm tests parseConfirm method with invalid or valid user responses.
func TestParseConfirm(t *testing.T) {
	// Checks with valid or invalids tags, expected true for yes, false if response is no and an error otherwise
	for _, cf := range confirmTests {
		if ok, err := parseConfirm(cf.str); err == nil {
			if cf.onErr {
				t.Errorf("Expected error for confirm with %v, received for %t: %t", cf.str, cf.ok, ok)
			} else if ok != cf.ok {
				t.Errorf("Expected result %t for confirm with %v, received: %t", cf.ok, cf.str, ok)
			}
		} else if !cf.onErr {
			t.Errorf("Expected no error for confirm with %v, received: %v", cf.str, err)
		}
	}
}
