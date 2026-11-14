package appcontext

import (
	"github.com/eser/aya.is/services/pkg/ajan"
	"github.com/eser/aya.is/services/pkg/api/adapters/arcade"
	"github.com/eser/aya.is/services/pkg/api/adapters/coolify"
	"github.com/eser/aya.is/services/pkg/api/adapters/resend"
	"github.com/eser/aya.is/services/pkg/api/adapters/s3client"
	telegramadapter "github.com/eser/aya.is/services/pkg/api/adapters/telegram"
	"github.com/eser/aya.is/services/pkg/api/adapters/unsplash"
	"github.com/eser/aya.is/services/pkg/api/adapters/workers"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	bulletinbiz "github.com/eser/aya.is/services/pkg/api/business/bulletin"
	"github.com/eser/aya.is/services/pkg/api/business/privacy"
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
	DiscloseErrors bool `conf:"disclose_errors" default:"false"` // show real error messages in HTTP responses
}

type ExternalsConfig struct {
	Arcade   arcade.Config   `conf:"arcade"`
	Unsplash unsplash.Config `conf:"unsplash"`
	Resend   resend.Config   `conf:"resend"`
}

type AppConfig struct {
	ajan.BaseConfig

	Auth      auth.Config     `conf:"auth"`
	Externals ExternalsConfig `conf:"externals"`
	S3        s3client.Config `conf:"s3"`
	// S3Export is the private bucket for personal data archives. Endpoint, region and credentials fall back to the
	// public S3 settings when empty, so only bucket_name is required to enable exports.
	S3Export s3client.Config    `conf:"s3_export"`
	Privacy  privacy.Config     `conf:"privacy"`
	Profiles profiles.Config    `conf:"profiles"`
	Data     DataConfig         `conf:"data"`
	Bulletin bulletinbiz.Config `conf:"bulletin"`
	Stories  stories.Config     `conf:"stories"`
	SiteURI  string             `conf:"site_uri"  default:"http://localhost:8080"`

	Telegram   telegramadapter.Config `conf:"telegram"`
	Coolify    coolify.Config         `conf:"coolify"`
	Workers    workers.Config         `conf:"workers"`
	Protection protection.Config      `conf:"protection"`
	Sessions   sessions.Config        `conf:"sessions"`

	Features FeatureFlags `conf:"features"`
}
