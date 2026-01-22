package appcontext

import (
	"github.com/eser/aya.is/services/pkg/ajan"
	"github.com/eser/aya.is/services/pkg/api/adapters/arcade"
	"github.com/eser/aya.is/services/pkg/api/adapters/s3client"
	"github.com/eser/aya.is/services/pkg/api/adapters/workers"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/protection"
	"github.com/eser/aya.is/services/pkg/api/business/sessions"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
)

type DataConfig struct {
	MigrationsPath string `conf:"migration_path" default:"etc/data/default/migrations"`
	SeedFilePath   string `conf:"seed_file_path" default:"etc/data/default/seed/seed.sql"`
}

type FeatureFlags struct {
	Dummy bool `conf:"dummy" default:"false"` // dummy feature flag
}

type ExternalsConfig struct {
	Arcade arcade.Config `conf:"arcade"`
}

type AppConfig struct {
	Auth       auth.Config       `conf:"auth"`
	Sessions   sessions.Config   `conf:"sessions"`
	Protection protection.Config `conf:"protection"`
	Profiles   profiles.Config   `conf:"profiles"`
	Stories    stories.Config    `conf:"stories"`
	Data       DataConfig        `conf:"data"`
	S3         s3client.Config   `conf:"s3"`
	Externals  ExternalsConfig   `conf:"externals"`
	Workers    workers.Config    `conf:"workers"`
	SiteURI    string            `conf:"site_uri"   default:"http://localhost:8080"`
	ajan.BaseConfig

	Features FeatureFlags `conf:"features"`
}
