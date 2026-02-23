package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/adapters/github"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterHTTPRoutesForProfileResources registers the routes for managing profile resources.
func RegisterHTTPRoutesForProfileResources(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	providers *ProfileLinkProviders,
) {
	// List profile resources
	routes.Route(
		"GET /{locale}/profiles/{slug}/_resources",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			resources, err := profileService.ListProfileResources(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to list profile resources",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  resources,
				"error": nil,
			})
		}).
		HasSummary("List Profile Resources").
		HasDescription("List all resources for a profile.").
		HasResponse(http.StatusOK)

	// List accessible GitHub repositories for adding as resources
	routes.Route(
		"GET /{locale}/profiles/{slug}/_resources/github/repos",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			// Check permission
			canEdit, permErr := profileService.HasUserAccessToProfile(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				profiles.MembershipKindMaintainer,
			)
			if permErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(permErr),
				)
			}

			if !canEdit {
				return ctx.Results.Error(
					http.StatusForbidden,
					httpfx.WithErrorMessage("You do not have permission to edit this profile"),
				)
			}

			// Use the viewer's individual profile's GitHub link to list their repos.
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			var gitHubLink *profiles.ManagedGitHubLink

			if user.IndividualProfileID != nil {
				gitHubLink, _ = profileService.GetManagedGitHubLink(
					ctx.Request.Context(),
					*user.IndividualProfileID,
				)
			}

			if gitHubLink == nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage(
						"No GitHub access token available. Please connect GitHub on your profile.",
					),
				)
			}

			// Fetch repos from GitHub API using the viewer's token
			page := 1
			pageParam := ctx.Request.URL.Query().Get("page")
			if pageParam != "" {
				if p, parseErr := strconv.Atoi(pageParam); parseErr == nil && p > 0 {
					page = p
				}
			}

			repos, err := providers.GitHub.Client().FetchUserRepos(
				ctx.Request.Context(),
				gitHubLink.AuthAccessToken,
				"owner,collaborator,organization_member",
				page,
				100,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to fetch GitHub repos",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to fetch GitHub repositories"),
				)
			}

			// Convert to response format
			type repoResponse struct {
				ID          string `json:"id"`
				FullName    string `json:"full_name"`
				Name        string `json:"name"`
				Description string `json:"description"`
				HTMLURL     string `json:"html_url"`
				Language    string `json:"language"`
				Stars       int    `json:"stars"`
				Forks       int    `json:"forks"`
				Private     bool   `json:"private"`
			}

			repoList := make([]repoResponse, 0, len(repos))
			for _, repo := range repos {
				repoList = append(repoList, repoResponse{
					ID:          strconv.FormatInt(repo.ID, 10),
					FullName:    repo.FullName,
					Name:        repo.Name,
					Description: repo.Description,
					HTMLURL:     repo.HTMLURL,
					Language:    repo.Language,
					Stars:       repo.Stars,
					Forks:       repo.Forks,
					Private:     repo.Private,
				})
			}

			return ctx.Results.JSON(map[string]any{
				"data":  repoList,
				"error": nil,
			})
		}).
		HasSummary("List GitHub Repositories").
		HasDescription("List accessible GitHub repositories for adding as profile resources.").
		HasResponse(http.StatusOK)

	// Create profile resource
	routes.Route(
		"POST /{locale}/profiles/{slug}/_resources",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			var reqBody struct {
				Kind        string         `json:"kind"`
				RemoteID    string         `json:"remote_id"`
				PublicID    string         `json:"public_id"`
				URL         string         `json:"url"`
				Title       string         `json:"title"`
				Description *string        `json:"description"`
				Properties  map[string]any `json:"properties"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&reqBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			if reqBody.Kind == "" || reqBody.Title == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Kind and title are required"),
				)
			}

			var remoteID *string
			if reqBody.RemoteID != "" {
				remoteID = &reqBody.RemoteID
			}

			var publicID *string
			if reqBody.PublicID != "" {
				publicID = &reqBody.PublicID
			}

			var resourceURL *string
			if reqBody.URL != "" {
				resourceURL = &reqBody.URL
			}

			resource, err := profileService.CreateProfileResource(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				reqBody.Kind,
				true, // is_managed: always true for API-added resources
				remoteID,
				publicID,
				resourceURL,
				reqBody.Title,
				reqBody.Description,
				reqBody.Properties,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to create profile resource",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  resource,
				"error": nil,
			})
		}).
		HasSummary("Create Profile Resource").
		HasDescription("Add a new resource to a profile.").
		HasResponse(http.StatusOK)

	// Set resource teams
	routes.Route(
		"PUT /{locale}/profiles/{slug}/_resources/{id}/teams",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			_, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")
			resourceID := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			var reqBody struct {
				TeamIDs []string `json:"team_ids"`
			}

			if err := json.NewDecoder(ctx.Request.Body).Decode(&reqBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			err := profileService.SetResourceTeams(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				resourceID,
				reqBody.TeamIDs,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to set resource teams",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("resource_id", resourceID))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "updated"},
				"error": nil,
			})
		}).
		HasSummary("Set Resource Teams").
		HasDescription("Assign teams to a resource.").
		HasResponse(http.StatusOK)

	// Delete profile resource
	routes.Route(
		"DELETE /{locale}/profiles/{slug}/_resources/{id}",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
			}
			slugParam := ctx.Request.PathValue("slug")
			resourceID := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get user information"),
				)
			}

			err := profileService.DeleteProfileResource(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				user.Kind,
				slugParam,
				resourceID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to delete profile resource",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("resource_id", resourceID))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": "deleted"},
				"error": nil,
			})
		}).
		HasSummary("Delete Profile Resource").
		HasDescription("Remove a resource from a profile.").
		HasResponse(http.StatusOK)
}

// getGitHubRepoInfoForResource is a helper that validates and returns GitHub repo info.
func getGitHubRepoInfoForResource(
	ctx *httpfx.Context,
	githubClient *github.Client,
	accessToken string,
	owner string,
	repo string,
) (*github.GitHubRepoInfo, error) {
	return githubClient.FetchRepoInfo(ctx.Request.Context(), accessToken, owner, repo)
}
