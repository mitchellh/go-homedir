// Package homedir detects the user's home directory without the use of cgo, for use in cross-compilation environments.

// +build plan9

package homedir

const homeEnv = "home"
