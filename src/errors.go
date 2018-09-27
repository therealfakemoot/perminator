package perminator

import (
	"errors"
)

var (
	// ErrBadFileType indicates a rule is not using one of the 3 accepted file type sigils: a, d, f.
	ErrBadFileType = errors.New("invalid filetype in rule")
	// ErrBadPerms indicates a rulle is using incorrectly formatted file permissions.
	ErrBadPerms = errors.New("invalid file permissions in rule")
	// ErrBadPattern indicates that filepath.Match has returned an error indicating incorrect glob pattern syntax.
	ErrBadPattern = errors.New("invalid file globbing pattern in rule")
)
