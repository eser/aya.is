package main

import (
	"context"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/processfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/appcontext"
	"github.com/eser/aya.is/services/pkg/api/adapters/http"
	"github.com/eser/aya.is/services/pkg/api/adapters/workers"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func main() {
	baseCtx := context.Background()

	appContext := appcontext.New()

	err := appContext.Init(baseCtx)
	if err != nil {
		panic(err)
	}

	process := processfx.New(baseCtx, appContext.Logger)

	startHTTPServer(process, appContext)
	startWorkers(process, appContext)

	process.Wait()
	process.Shutdown()
}

func startHTTPServer(process *processfx.Process, appContext *appcontext.AppContext) {
	process.StartGoroutine("http-server", func(ctx context.Context) error {
		cleanup, err := http.Run(
			ctx,
			appContext.Config.SiteURI,
			&appContext.Config.HTTP,
			appContext.Logger,
			appContext.AuthService,
			appContext.UserService,
			appContext.ProfileService,
			appContext.ProfilePointsService,
			appContext.StoryService,
			appContext.SessionService,
			appContext.ProtectionService,
			appContext.UploadService,
			&http.ProfileLinkProviders{
				YouTube:                appContext.YouTubeProvider,
				GitHub:                 appContext.GitHubProvider,
				PendingConnectionStore: profiles.NewPendingConnectionStore(),
			},
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
}

func startWorkers(process *processfx.Process, appContext *appcontext.AppContext) {
	idGen := func() string {
		return string(users.DefaultIDGenerator())
	}

	// YouTube full sync worker
	if appContext.Config.Workers.YouTubeSync.FullSyncEnabled {
		fullSyncWorker := workers.NewYouTubeFullSyncWorker(
			&appContext.Config.Workers.YouTubeSync,
			appContext.Logger,
			appContext.LinkSyncService,
			appContext.YouTubeProvider,
			idGen,
			appContext.RuntimeStateService,
		)

		runner := workerfx.NewRunner(fullSyncWorker, appContext.Logger)

		process.StartGoroutine("youtube-full-sync-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}

	// YouTube incremental sync worker
	if appContext.Config.Workers.YouTubeSync.IncrementalSyncEnabled {
		incrementalSyncWorker := workers.NewYouTubeIncrementalSyncWorker(
			&appContext.Config.Workers.YouTubeSync,
			appContext.Logger,
			appContext.LinkSyncService,
			appContext.YouTubeProvider,
			idGen,
			appContext.RuntimeStateService,
		)

		runner := workerfx.NewRunner(incrementalSyncWorker, appContext.Logger)

		process.StartGoroutine("youtube-incremental-sync-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}

	// Event queue worker
	if appContext.Config.Workers.EventQueue.Enabled {
		workerID := idGen()

		eventQueueWorker := workers.NewEventQueueWorker(
			&appContext.Config.Workers.EventQueue,
			appContext.Logger,
			appContext.Repository,
			appContext.EventRegistry,
			workerID,
		)

		runner := workerfx.NewRunner(eventQueueWorker, appContext.Logger)

		process.StartGoroutine("event-queue-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}
}
