package siteimporter

import "context"

// SiteProvider defines the interface that site import adapters must implement.
type SiteProvider interface {
	// Kind returns the provider kind (e.g., "speakerdeck", "medium").
	Kind() string
	// Check validates a URL and returns connection info.
	Check(ctx context.Context, url string) (*CheckResult, error)
	// FetchAll fetches all importable items from the source.
	FetchAll(ctx context.Context, username string) ([]*ImportItem, error)
}
