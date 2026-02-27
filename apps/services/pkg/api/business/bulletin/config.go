package bulletin

// Config holds business-layer configuration for the bulletin module.
type Config struct {
	MinStoryThreshold   int    `conf:"min_story_threshold"    default:"5"`
	MaxStoriesPerDigest int    `conf:"max_stories_per_digest" default:"15"`
	FrontendURI         string `conf:"frontend_uri"           default:"https://aya.is"`
}
