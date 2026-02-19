package main

import (
	"context"
	"log/slog"

	"github.com/eser/aya.is/services/pkg/ajan/processfx"
	"github.com/eser/aya.is/services/pkg/ajan/workerfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/appcontext"
	githubadapter "github.com/eser/aya.is/services/pkg/api/adapters/github"
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
			appContext.Config.Features.DiscloseErrors,
			appContext.AuthService,
			appContext.UserService,
			appContext.ProfileService,
			appContext.ProfilePointsService,
			appContext.ProfileQuestionsService,
			appContext.StoryService,
			appContext.StoryInteractionService,
			appContext.SessionService,
			appContext.ProtectionService,
			appContext.UploadService,
			appContext.UnsplashClient,
			&http.ProfileLinkProviders{
				YouTube:                appContext.YouTubeProvider,
				GitHub:                 appContext.GitHubProvider,
				SiteImporter:           appContext.SiteImporterService,
				PendingConnectionStore: profiles.NewPendingConnectionStore(),
			},
			buildTelegramProviders(appContext),
			appContext.AIModels,
			appContext.RuntimeStateService,
			appContext.WorkerRegistry,
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

func buildTelegramProviders(appContext *appcontext.AppContext) *http.TelegramProviders {
	if appContext.TelegramBot == nil {
		return nil
	}

	return &http.TelegramProviders{
		Bot:           appContext.TelegramBot,
		Service:       appContext.TelegramService,
		WebhookSecret: appContext.Config.Telegram.WebhookSecret,
	}
}

func startWorkers(process *processfx.Process, appContext *appcontext.AppContext) { //nolint:funlen
	idGen := func() string {
		return string(users.DefaultIDGenerator())
	}

	// Create story processor (used by sync workers to create/reconcile stories after sync)
	storyProcessor := workers.NewYouTubeStoryProcessor(
		&appContext.Config.Workers.YouTubeSync,
		appContext.Logger,
		appContext.ProfileLinkSyncService,
		appContext.Repository,
		idGen,
	)

	// YouTube full sync worker
	if appContext.Config.Workers.YouTubeSync.FullSyncEnabled {
		fullSyncWorker := workers.NewYouTubeFullSyncWorker(
			&appContext.Config.Workers.YouTubeSync,
			appContext.Logger,
			appContext.ProfileLinkSyncService,
			appContext.YouTubeProvider,
			idGen,
			appContext.RuntimeStateService,
			storyProcessor,
		)

		runner := workerfx.NewRunner(fullSyncWorker, appContext.Logger)
		runner.SetStateKey("youtube.sync.full_sync_worker")
		appContext.WorkerRegistry.Register(runner)

		process.StartGoroutine("youtube-full-sync-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}

	// YouTube incremental sync worker
	if appContext.Config.Workers.YouTubeSync.IncrementalSyncEnabled {
		incrementalSyncWorker := workers.NewYouTubeIncrementalSyncWorker(
			&appContext.Config.Workers.YouTubeSync,
			appContext.Logger,
			appContext.ProfileLinkSyncService,
			appContext.YouTubeProvider,
			idGen,
			appContext.RuntimeStateService,
			storyProcessor,
		)

		runner := workerfx.NewRunner(incrementalSyncWorker, appContext.Logger)
		runner.SetStateKey("youtube.sync.incremental_sync_worker")
		appContext.WorkerRegistry.Register(runner)

		process.StartGoroutine("youtube-incremental-sync-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}

	// GitHub resource sync worker
	if appContext.Config.Workers.GitHubSync.Enabled {
		githubFetcher := githubadapter.NewResourceFetcherAdapter(appContext.GitHubClient)

		githubSyncWorker := workers.NewGitHubSyncWorker(
			&appContext.Config.Workers.GitHubSync,
			appContext.Logger,
			appContext.ProfileResourceSyncService,
			githubFetcher,
			appContext.RuntimeStateService,
		)

		runner := workerfx.NewRunner(githubSyncWorker, appContext.Logger)
		runner.SetStateKey("github.resource_sync_worker")
		appContext.WorkerRegistry.Register(runner)

		process.StartGoroutine("github-resource-sync-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}

	// SpeakerDeck sync worker
	if appContext.Config.Workers.SpeakerDeckSync.FullSyncEnabled {
		speakerDeckStoryProcessor := workers.NewSpeakerDeckStoryProcessor(
			&appContext.Config.Workers.SpeakerDeckSync,
			appContext.Logger,
			appContext.ProfileLinkSyncService,
			appContext.Repository,
			idGen,
		)

		speakerDeckSyncWorker := workers.NewSpeakerDeckSyncWorker(
			&appContext.Config.Workers.SpeakerDeckSync,
			appContext.Logger,
			appContext.ProfileLinkSyncService,
			appContext.SiteImporterService,
			speakerDeckStoryProcessor,
			appContext.RuntimeStateService,
			idGen,
		)

		runner := workerfx.NewRunner(speakerDeckSyncWorker, appContext.Logger)
		runner.SetStateKey("speakerdeck.sync.full_sync_worker")
		appContext.WorkerRegistry.Register(runner)

		process.StartGoroutine("speakerdeck-full-sync-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}

	// Queue worker
	if appContext.Config.Workers.Queue.Enabled {
		workerID := idGen()

		queueWorker := workers.NewQueueWorker(
			&appContext.Config.Workers.Queue,
			appContext.Logger,
			appContext.Repository,
			appContext.QueueRegistry,
			workerID,
			appContext.RuntimeStateService,
		)

		runner := workerfx.NewRunner(queueWorker, appContext.Logger)
		appContext.WorkerRegistry.Register(runner)

		process.StartGoroutine("queue-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}

	// Telegram bot polling worker (dev mode — alternative to webhooks)
	if appContext.TelegramBot != nil &&
		appContext.Config.Telegram.UsePolling &&
		appContext.Config.Workers.TelegramBot.Enabled {
		telegramPollWorker := workers.NewTelegramPollWorker(
			&appContext.Config.Workers.TelegramBot,
			appContext.Logger,
			appContext.TelegramClient,
			appContext.TelegramBot,
		)

		runner := workerfx.NewRunner(telegramPollWorker, appContext.Logger)
		runner.SetStateKey("telegram.bot.poll_worker")
		appContext.WorkerRegistry.Register(runner)

		process.StartGoroutine("telegram-bot-poll-worker", func(ctx context.Context) error {
			return runner.Run(ctx)
		})
	}
}
