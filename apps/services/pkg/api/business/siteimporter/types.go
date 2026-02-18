package siteimporter

import "time"

// ImportItem represents a single item from an external site.
type ImportItem struct {
	RemoteID     string
	Title        string
	Description  string
	PublishedAt  time.Time
	Link         string
	ThumbnailURL string
	StoryKind    string
	Properties   map[string]any
}

// CheckResult represents the result of checking an external site connection.
type CheckResult struct {
	Valid    bool
	Username string
	URI      string
	Title    string
}
