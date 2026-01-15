package main

import (
	"context"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/processfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/appcontext"
	"github.com/eser/aya.is/services/pkg/api/adapters/http"
)

func main() {
	baseCtx := context.Background()

	appContext := appcontext.New()

	err := appContext.Init(baseCtx)
	if err != nil {
		panic(err)
	}

	process := processfx.New(baseCtx, appContext.Logger)

	process.StartGoroutine("http-server", func(ctx context.Context) error {
		cleanup, err := http.Run(
			ctx,
			appContext.Config.SiteURI,
			&appContext.Config.HTTP,
			appContext.Logger,
			appContext.AuthService,
			appContext.UserService,
			appContext.ProfileService,
			appContext.StoryService,
		)
		if err != nil {
			appContext.Logger.ErrorContext(
				ctx,
				"[Main] HTTP server run failed",
				slog.String("module", "main"),
				slog.Any("error", err))
		}

		defer cleanup()

		<-ctx.Done()

		return nil
	})

	process.Wait()
	process.Shutdown()
}
