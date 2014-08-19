# go-homedir

This is a Go library for detecting the user's home directory without
the use of cgo, so the library can be used in cross-compilation environments.

Usage is incredibly simple, just call `homedir.Dir()` to get the home directory
for a user, and `homedir.Expand()` to expand the `~` in a path to the home
directory.
