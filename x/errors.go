package x

import "errors"

var (
	ErrForbidden                   = errors.New("forbidden")
	ErrUnauthorized                = errors.New("unauthorized")
	ErrUserNotFoundInContext       = errors.New("user not found in context")
	ErrCookieNotFoundInContext     = errors.New("cookie not found in context")
	ErrProjectNotFoundInContext    = errors.New("project not found in context")
	ErrProjectExists               = errors.New("project already exists")
	ErrNotProjectMember            = errors.New("not an project member")
	ErrMetadataNotFound            = errors.New("metadata not found")
	ErrCookieNotFound              = errors.New("http cookie not found")
	ErrOAuthStateNotFoundInContext = errors.New("oauth state not found in context")
)
