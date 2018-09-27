package perminator

import (
	"errors"
)

var (
	ErrBadFileType = errors.New("invalid filetype in rule")
	ErrBadPerms    = errors.New("invalid file permissions in rule")
	ErrBadPattern  = errors.New("invalid file globbing pattern in rule")
)
