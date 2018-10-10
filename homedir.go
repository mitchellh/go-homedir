// Package homedir detects the user's home directory without the use of cgo, for use in cross-compilation environments.
package homedir

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
)

// DisableCache will disable caching of the home directory.
// Caching is enabled by default.
var DisableCache bool

var homedirCache string
var cacheLock sync.RWMutex

// Dir returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Dir() (string, error) {
	if !DisableCache {
		cacheLock.RLock()
		cached := homedirCache
		cacheLock.RUnlock()
		if cached != "" {
			return cached, nil
		}
	}

	cacheLock.Lock()
	defer cacheLock.Unlock()

	if !DisableCache && homedirCache != "" {
		return homedirCache, nil
	}

	result, err := dir()
	if err != nil {
		return "", err
	}

	homedirCache = result
	return result, nil
}

// Expand expands the path to include the home directory if the path is prefixed with `~`.
// If it isn't prefixed with `~`, the path is returned as-is.
func Expand(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] != '~' {
		return path, nil
	}

	if len(path) == 1 {
		return Dir()
	}

	if !os.IsPathSeparator(path[1]) {
		return "", errors.New("cannot expand user-specific home dir")
	}

	dir, err := Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, path[2:]), nil
}
