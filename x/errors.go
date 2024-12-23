package x

import "errors"

var (
	ErrForbidden                     = errors.New("forbidden")
	ErrUnauthorized                  = errors.New("unauthorized")
	ErrUserNotFoundInContext         = errors.New("user not found in context")
	ErrCookieNotFoundInContext       = errors.New("cookie not found in context")
	ErrOrganizationNotFoundInContext = errors.New("organization not found in context")
	ErrOrganizationExists            = errors.New("organization already exists")
	ErrNotOrganizationMember         = errors.New("not an organization member")
	ErrMetadataNotFound              = errors.New("metadata not found")
	ErrCookieNotFound                = errors.New("http cookie not found")
	ErrOAuthStateNotFoundInContext   = errors.New("oauth state not found in context")
)
