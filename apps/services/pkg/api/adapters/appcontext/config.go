package appcontext

import (
	"github.com/eser/aya.is/services/pkg/ajan"
	"github.com/eser/aya.is/services/pkg/api/adapters/arcade"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
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
	Auth      auth.Config     `conf:"auth"`
	Data      DataConfig      `conf:"data"`
	Externals ExternalsConfig `conf:"externals"`
	SiteURI   string          `conf:"site_uri"  default:"http://localhost:8080"`
	ajan.BaseConfig

	Features FeatureFlags `conf:"features"`
}
