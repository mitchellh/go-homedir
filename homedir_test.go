package homedir

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func patchEnv(key, value string) func() {
	bck := os.Getenv(key)
	deferFunc := func() {
		os.Setenv(key, bck)
	}

	os.Setenv(key, value)
	return deferFunc
}

func BenchmarkDir(b *testing.B) {
	// We do this for any "warmups"
	for i := 0; i < 10; i++ {
		Dir()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Dir()
	}
}

func TestUser(t *testing.T) {
	DisableCache = true

	u, err := user.Current()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if runtime.GOOS == "windows" {
		a := strings.Split(u.Username, "\\")
		if len(a) != 2 {
			t.Fatalf("could not split username")
		}
		u.Username = a[1] // patch up username
	}

	user, err := User()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if u.Username != user {
		t.Fatalf("invalid user %v != %v", u.Username, user)
	}

	if runtime.GOOS == "windows" {
		return
	}

	// force whoami
	patchEnv("USER", "")
	user, err = User()
	if err != nil {
		t.Fatalf("err whoami: %s", err)
	}

	if u.Username != user {
		t.Fatalf("invalid user %v != %v", u.Username, user)
	}

	// force id
	whoamiBypass = true
	user, err = User()
	if err != nil {
		t.Fatalf("err id: %s", err)
	}

	if u.Username != user {
		t.Fatalf("invalid user %v != %v", u.Username, user)
	}
}

func TestDir(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	dir, err := Dir()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if u.HomeDir != dir {
		t.Fatalf("%#v != %#v", u.HomeDir, dir)
	}
}

func TestExpand(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	cases := []struct {
		Input  string
		Output string
		Err    bool
	}{
		{
			"/foo",
			"/foo",
			false,
		},

		{
			"~/foo",
			filepath.Join(u.HomeDir, "foo"),
			false,
		},

		{
			"",
			"",
			false,
		},

		{
			"~",
			u.HomeDir,
			false,
		},

		{
			"~foo/foo",
			"",
			true,
		},
	}

	for _, tc := range cases {
		actual, err := Expand(tc.Input)
		if (err != nil) != tc.Err {
			t.Fatalf("Input: %#v\n\nErr: %s", tc.Input, err)
		}

		if actual != tc.Output {
			t.Fatalf("Input: %#v\n\nOutput: %#v", tc.Input, actual)
		}
	}

	DisableCache = true
	defer func() { DisableCache = false }()
	defer patchEnv("HOME", "/custom/path/")()
	expected := filepath.Join("/", "custom", "path", "foo/bar")
	actual, err := Expand("~/foo/bar")

	if err != nil {
		t.Errorf("No error is expected, got: %v", err)
	} else if actual != expected {
		t.Errorf("Expected: %v; actual: %v", expected, actual)
	}
}
