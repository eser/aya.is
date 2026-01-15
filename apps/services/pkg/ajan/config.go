package ajan

import (
	"github.com/eser/aya.is/services/pkg/ajan/connfx"
	"github.com/eser/aya.is/services/pkg/ajan/httpclient"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

type BaseConfig struct {
	Conn       connfx.Config `conf:"conn"`
	AppName    string        `conf:"name"    default:"ajansvc"`
	AppEnv     string        `conf:"env"     default:"development"`
	AppVersion string        `conf:"version" default:"0.0.0"`

	// Security configuration
	JWTSecret string `conf:"jwt_secret"` // No default - validated at startup

	Log        logfx.Config      `conf:"log"`
	HTTP       httpfx.Config     `conf:"http"`
	HTTPClient httpclient.Config `conf:"http_client"`
}
