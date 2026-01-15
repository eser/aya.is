package httpfx

import (
	"time"
)

type Config struct {
	Addr string `conf:"addr" default:":8080"`

	CertString        string        `conf:"cert_string"`
	KeyString         string        `conf:"key_string"`
	ReadHeaderTimeout time.Duration `conf:"read_header_timeout" default:"5s"`
	ReadTimeout       time.Duration `conf:"read_timeout"        default:"10s"`
	WriteTimeout      time.Duration `conf:"write_timeout"       default:"10s"`
	IdleTimeout       time.Duration `conf:"idle_timeout"        default:"120s"`

	InitializationTimeout   time.Duration `conf:"init_timeout"     default:"25s"`
	GracefulShutdownTimeout time.Duration `conf:"shutdown_timeout" default:"5s"`

	SelfSigned bool `conf:"self_signed" default:"false"`

	HealthCheckEnabled bool `conf:"health_check" default:"true"`
	OpenAPIEnabled     bool `conf:"openapi"      default:"true"`
	ProfilingEnabled   bool `conf:"profiling"    default:"false"`
}
