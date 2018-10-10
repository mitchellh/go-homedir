// Package homedir detects the user's home directory without the use of cgo, for use in cross-compilation environments.

// +build darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package homedir

const homeEnv = "HOME"
