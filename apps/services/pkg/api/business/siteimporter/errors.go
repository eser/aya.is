package siteimporter

import "errors"

var (
	ErrProviderNotFound = errors.New("siteimporter: provider not found")
	ErrInvalidURL       = errors.New("siteimporter: invalid URL")
	ErrSiteNotFound     = errors.New("siteimporter: site not found")
	ErrAlreadyConnected = errors.New("siteimporter: already connected")
)
