package siteimporter

import "time"

// ImportItem represents a single item from an external site.
type ImportItem struct {
	PublishedAt  time.Time
	Properties   map[string]any
	RemoteID     string
	Title        string
	Description  string
	Link         string
	ThumbnailURL string
	StoryKind    string
}

// CheckResult represents the result of checking an external site connection.
type CheckResult struct {
	Username string
	URI      string
	Title    string
	Valid    bool
}
