package gitup

import (
	"io/ioutil"
	"os"
	"testing"
)

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
		t.Errorf("Expected no operation witn unknown version type")
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
			t.Errorf("Expected no error for confirm with %v, received:", cf.str, err)
		}
	}
}
