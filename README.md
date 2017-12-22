# go-homedir

This is a Go library for detecting the user's home directory without
the use of cgo, so the library can be used in cross-compilation environments.

Usage is incredibly simple, just call `homedir.Dir()` to get the home directory
for a user, and `homedir.Expand()` to expand the `~` in a path to the home
directory.

**Why not just use `os/user`?** The built-in `os/user` package requires
cgo on Darwin systems. This means that any Go code that uses that package
cannot cross compile. But 99% of the time the use for `os/user` is just to
retrieve the home directory, which we can do for the current user without
cgo. This library does that, enabling cross-compilation.

## Installing with Glide

```
glide get github.com/mitchellh/go-homedir
```

## Using go-homedir

```go
import (
  "github.com/mitchellh/go-homedir"
)

// get the home dir 
home, err := homedir.Dir()

// expand a directory if its path starts with ~
someDir, err : = homedir.Expand("~/some-dir")
```
