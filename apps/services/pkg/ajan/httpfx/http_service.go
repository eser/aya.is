package httpfx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/lib"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
)

var (
	ErrFailedToLoadCertificate        = errors.New("failed to load certificate")
	ErrFailedToGenerateSelfSignedCert = errors.New("failed to generate self-signed certificate")
	ErrFailedToCreateHTTPMetrics      = errors.New("failed to create HTTP metrics")
	ErrHTTPServiceNetListenError      = errors.New("HTTP service net listen error")
)

type HTTPService struct {
	InnerServer  *http.Server
	InnerRouter  *Router
	InnerMetrics *Metrics

	Config *Config
	logger *logfx.Logger
}

func NewHTTPService(
	config *Config,
	router *Router,
	logger *logfx.Logger,
) *HTTPService {
	server := &http.Server{ //nolint:exhaustruct
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		ReadTimeout:       config.ReadTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,

		Addr: config.Addr,

		Handler: router.GetMux(),
	}

	metricsBuilder := logger.NewMetricsBuilder("httpfx")
	metrics := NewMetrics(metricsBuilder)

	return &HTTPService{
		InnerServer:  server,
		InnerRouter:  router,
		InnerMetrics: metrics,
		Config:       config,
		logger:       logger,
	}
}

func (hs *HTTPService) Server() *http.Server {
	return hs.InnerServer
}

func (hs *HTTPService) Router() *Router {
	return hs.InnerRouter
}

func (hs *HTTPService) SetupTLS(ctx context.Context) error {
	switch {
	case hs.Config.CertString != "" && hs.Config.KeyString != "":
		cert, err := tls.X509KeyPair([]byte(hs.Config.CertString), []byte(hs.Config.KeyString))
		if err != nil {
			return fmt.Errorf("%w: %w", ErrFailedToLoadCertificate, err)
		}

		hs.InnerServer.TLSConfig = &tls.Config{ //nolint:exhaustruct
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
	case hs.Config.SelfSigned:
		cert, err := lib.GenerateSelfSignedCert()
		if err != nil {
			return fmt.Errorf("%w: %w", ErrFailedToGenerateSelfSignedCert, err)
		}

		hs.InnerServer.TLSConfig = &tls.Config{ //nolint:exhaustruct
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
	default:
		hs.logger.WarnContext(
			ctx,
			"HTTPService is starting without TLS, this will cause HTTP/2 support to be disabled",
		)
	}

	return nil
}

func (hs *HTTPService) Start(ctx context.Context) (func(), error) {
	hs.logger.InfoContext(ctx, "HTTPService is starting...", slog.String("addr", hs.Config.Addr))

	if hs.InnerMetrics != nil {
		err := hs.InnerMetrics.Init()
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrFailedToCreateHTTPMetrics, err)
		}
	}

	err := hs.SetupTLS(ctx)
	if err != nil {
		return nil, err
	}

	listener, lnErr := net.Listen("tcp", hs.InnerServer.Addr)
	if lnErr != nil {
		return nil, fmt.Errorf("%w: %w", ErrHTTPServiceNetListenError, lnErr)
	}

	go func() {
		var sErr error

		if hs.InnerServer.TLSConfig != nil {
			sErr = hs.InnerServer.ServeTLS(listener, "", "")
		} else {
			sErr = hs.InnerServer.Serve(listener)
		}

		if sErr != nil && !errors.Is(sErr, http.ErrServerClosed) {
			hs.logger.ErrorContext(ctx, "HTTPService ServeTLS error: %w", slog.Any("error", sErr))
		}
	}()

	cleanup := func() {
		hs.logger.InfoContext(ctx, "Shutting down server...")

		newCtx, cancel := context.WithTimeout(ctx, hs.Config.GracefulShutdownTimeout)
		defer cancel()

		err := hs.InnerServer.Shutdown(newCtx)
		if err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			hs.logger.ErrorContext(ctx, "HTTPService forced to shutdown", slog.Any("error", err))

			return
		}

		hs.logger.InfoContext(ctx, "HTTPService has gracefully stopped.")
	}

	return cleanup, nil
}
