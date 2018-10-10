// Package homedir detects the user's home directory without the use of cgo, for use in cross-compilation environments.

// +build dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris plan9

package homedir

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
)

func dir() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv(homeEnv); home != "" {
		return home, nil
	}

	stdout := new(bytes.Buffer)

	cmd := exec.Command("getent", "passwd", strconv.Itoa(os.Getuid()))
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		// If the error is ErrNotFound, we ignore it. Otherwise, return it.
		if err != exec.ErrNotFound {
			return "", err
		}
	} else {
		if passwd := bytes.TrimSpace(stdout.Bytes()); len(passwd) > 0 {
			// username:password:uid:gid:gecos:home:shell
			passwdParts := bytes.SplitN(passwd, []byte(":"), 7)
			if len(passwdParts) > 5 {
				return string(passwdParts[5]), nil
			}
		}
	}

	// If all else fails, try the shell
	stdout.Reset()

	cmd = exec.Command("sh", "-c", "cd && pwd")
	cmd.Stdout = stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := string(bytes.TrimSpace(stdout.Bytes()))
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}
