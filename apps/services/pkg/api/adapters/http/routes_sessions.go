package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/middlewares"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/protection"
	"github.com/eser/aya.is/services/pkg/api/business/sessions"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// CreateSessionRequest is the request body for creating a new session.
type CreateSessionRequest struct {
	POWChallengeID string                      `json:"pow_challenge_id"`
	Nonce          string                      `json:"nonce"`
	Preferences    sessions.SessionPreferences `json:"preferences"`
}

// UpdatePreferencesRequest is the request body for updating session preferences.
type UpdatePreferencesRequest struct {
	Theme    *string `json:"theme,omitempty"`
	Locale   *string `json:"locale,omitempty"`
	Timezone *string `json:"timezone,omitempty"`
}

func RegisterHTTPRoutesForSessions( //nolint:funlen,cyclop
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	sessionService *sessions.Service,
	protectionService *protection.Service,
) {
	// GET /{locale}/sessions/_current - Consolidated session endpoint (cookie-based, no auth middleware)
	// Returns auth state, fresh JWT token, and preferences in a single response.
	// Used by frontend on app mount to establish session state without multiple round-trips.
	routes.Route(
		"GET /{locale}/sessions/_current",
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, err := GetSessionIDFromCookie(ctx.Request, authService.Config)
			if err != nil {
				return ctx.Results.JSON(map[string]any{
					"data": map[string]any{"authenticated": false, "preferences": nil},
				})
			}

			session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if err != nil || session == nil || session.Status != users.SessionStatusActive {
				ClearSessionCookie(ctx.ResponseWriter, authService.Config)

				return ctx.Results.JSON(map[string]any{
					"data": map[string]any{"authenticated": false, "preferences": nil},
				})
			}

			// Fetch preferences (available for both anonymous and authenticated sessions)
			prefs, prefsErr := sessionService.GetPreferences(ctx.Request.Context(), sessionID)
			if prefsErr != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to get preferences",
					slog.String("error", prefsErr.Error()),
					slog.String("session_id", sessionID))
			}

			// Anonymous session - return preferences only
			if session.LoggedInUserID == nil {
				return ctx.Results.JSON(map[string]any{
					"data": map[string]any{"authenticated": false, "preferences": prefs},
				})
			}

			// Authenticated session - generate fresh JWT and include user data
			tokenString, expiresAt, err := authService.GenerateSessionToken(
				session.ID,
				*session.LoggedInUserID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to generate session token",
					slog.String("error", err.Error()))

				return ctx.Results.JSON(map[string]any{
					"data": map[string]any{"authenticated": false, "preferences": prefs},
				})
			}

			user, _ := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)

			response := map[string]any{
				"authenticated":    true,
				"token":            tokenString,
				"expires_at":       expiresAt.Unix(),
				"user":             user,
				"selected_profile": nil,
				"preferences":      prefs,
			}

			// If user has an individual profile, fetch it and their accessible profiles
			if user != nil && user.IndividualProfileID != nil {
				locale := ctx.Request.PathValue("locale")
				profile, profileErr := profileService.GetBasicByID(
					ctx.Request.Context(),
					*user.IndividualProfileID,
				)
				if profileErr == nil && profile != nil {
					response["selected_profile"] = map[string]any{
						"id":                  profile.ID,
						"slug":                profile.Slug,
						"kind":                profile.Kind,
						"profile_picture_uri": profile.ProfilePictureURI,
					}
				}

				// Fetch profiles where user is a member
				memberships, membershipErr := profileService.GetMembershipsByUserProfileID(
					ctx.Request.Context(),
					locale,
					*user.IndividualProfileID,
				)
				if membershipErr == nil && len(memberships) > 0 {
					accessibleProfiles := make([]map[string]any, 0, len(memberships))
					for _, m := range memberships {
						if m.Profile == nil {
							continue
						}

						accessibleProfiles = append(accessibleProfiles, map[string]any{
							"id":                  m.Profile.ID,
							"slug":                m.Profile.Slug,
							"kind":                m.Profile.Kind,
							"title":               m.Profile.Title,
							"profile_picture_uri": m.Profile.ProfilePictureURI,
							"membership_kind":     m.Kind,
						})
					}

					response["accessible_profiles"] = accessibleProfiles
				}
			}

			return ctx.Results.JSON(map[string]any{
				"data": response,
			})
		}).
		HasSummary("Get Current Session").
		HasDescription("Returns auth state, token, and preferences via cookie-based session check.").
		HasResponse(http.StatusOK)

	// Register authenticated route with auth middleware
	routes.Route(
		"GET /{locale}/sessions/current",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get user ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get session data
			session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if err != nil {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithErrorMessage("Session not found"),
				)
			}

			// Prepare response with session data
			response := map[string]any{
				"id":               session.ID,
				"user":             nil,
				"selected_profile": nil,
			}

			if session.LoggedInUserID != nil { //nolint:nestif
				user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
				if userErr != nil {
					return ctx.Results.Error(
						http.StatusNotFound,
						httpfx.WithErrorMessage("User not found"),
					)
				}

				if user != nil {
					response["user"] = map[string]any{
						"id":                    user.ID,
						"kind":                  user.Kind,
						"name":                  user.Name,
						"email":                 user.Email,
						"phone":                 user.Phone,
						"github_handle":         user.GithubHandle,
						"github_remote_id":      user.GithubRemoteID,
						"bsky_handle":           user.BskyHandle,
						"x_handle":              user.XHandle,
						"individual_profile_id": user.IndividualProfileID,
						"created_at":            user.CreatedAt,
						"updated_at":            user.UpdatedAt,
					}

					// If user has an individual profile, fetch it
					if user.IndividualProfileID != nil {
						// Get locale from path
						locale := ctx.Request.PathValue("locale")
						logger.DebugContext(ctx.Request.Context(), "Fetching profile",
							slog.String("locale", locale),
							slog.String("profile_id", *user.IndividualProfileID))

						profile, profileErr := profileService.GetByID(
							ctx.Request.Context(),
							locale,
							*user.IndividualProfileID,
						)
						if profileErr != nil {
							logger.ErrorContext(ctx.Request.Context(), "Profile fetch error",
								slog.String("error", profileErr.Error()),
								slog.String("profile_id", *user.IndividualProfileID))
						}

						if profile != nil {
							response["selected_profile"] = map[string]any{
								"id":                  profile.ID,
								"slug":                profile.Slug,
								"kind":                profile.Kind,
								"title":               profile.Title,
								"description":         profile.Description,
								"profile_picture_uri": profile.ProfilePictureURI,
								"pronouns":            profile.Pronouns,
								"properties":          profile.Properties,
								"created_at":          profile.CreatedAt,
								"updated_at":          profile.UpdatedAt,
							}
						}
					}
				}
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  response,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get Current Session").
		HasDescription("Returns the current authenticated session with user and profile data.").
		HasResponse(http.StatusOK)

	// POST /{locale}/sessions - Create a new anonymous session (with PoW verification)
	routes.Route(
		"POST /{locale}/sessions",
		func(ctx *httpfx.Context) httpfx.Result {
			// Parse request body
			var req CreateSessionRequest
			if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Get client IP
			clientIP := middlewares.GetClientAddrs(ctx.Request)

			// Verify PoW challenge if enabled
			if protectionService.IsPOWChallengeEnabled() {
				if req.POWChallengeID == "" || req.Nonce == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("PoW challenge ID and nonce are required"),
					)
				}

				err := protectionService.VerifyAndConsumePOWChallenge(
					ctx.Request.Context(),
					req.POWChallengeID,
					req.Nonce,
					clientIP,
				)
				if err != nil {
					logger.WarnContext(ctx.Request.Context(), "PoW verification failed",
						"error", err.Error(),
						"challenge_id", req.POWChallengeID)

					switch {
					case errors.Is(err, protection.ErrPOWChallengeNotFound):
						return ctx.Results.Error(
							http.StatusForbidden,
							httpfx.WithErrorMessage("Invalid PoW challenge"),
						)
					case errors.Is(err, protection.ErrPOWChallengeExpired):
						return ctx.Results.Error(
							http.StatusForbidden,
							httpfx.WithErrorMessage("PoW challenge expired"),
						)
					case errors.Is(err, protection.ErrPOWChallengeUsed):
						return ctx.Results.Error(
							http.StatusForbidden,
							httpfx.WithErrorMessage("PoW challenge already used"),
						)
					case errors.Is(err, protection.ErrPOWChallengeInvalid):
						return ctx.Results.Error(
							http.StatusForbidden,
							httpfx.WithErrorMessage("Invalid PoW solution"),
						)
					default:
						return ctx.Results.Error(
							http.StatusInternalServerError,
							httpfx.WithErrorMessage("PoW verification failed"),
						)
					}
				}
			}

			// Create session with rate limiting
			ipHash := protection.HashIP(clientIP)
			session, err := sessionService.CreateSession(ctx.Request.Context(), ipHash)
			if err != nil {
				if errors.Is(err, sessions.ErrRateLimitExceeded) {
					return ctx.Results.Error(
						http.StatusTooManyRequests,
						httpfx.WithErrorMessage("Rate limit exceeded. Please try again later."),
					)
				}

				logger.ErrorContext(ctx.Request.Context(), "Failed to create session",
					slog.String("error", err.Error()))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to create session"),
				)
			}

			// Set initial preferences if provided
			if len(req.Preferences) > 0 {
				err := sessionService.SetPreferences(
					ctx.Request.Context(),
					session.ID,
					req.Preferences,
				)
				if err != nil {
					logger.WarnContext(ctx.Request.Context(), "Failed to set initial preferences",
						slog.String("error", err.Error()),
						slog.String("session_id", session.ID))
					// Don't fail the request - session was created successfully
				}
			}

			// Set session cookie
			expiresAt := time.Now().Add(24 * time.Hour) // Default 24 hours
			SetSessionCookie(ctx.ResponseWriter, session.ID, expiresAt, authService.Config)

			// Get preferences for response
			prefs, _ := sessionService.GetPreferences(ctx.Request.Context(), session.ID)

			return ctx.Results.JSON(map[string]any{
				"data": map[string]any{
					"session_id":  session.ID,
					"preferences": prefs,
				},
				"error": nil,
			})
		}).
		HasSummary("Create Session").
		HasDescription("Creates a new anonymous session with optional PoW verification.").
		HasResponse(http.StatusOK)

	// PATCH /{locale}/sessions/_current - Update session preferences
	routes.Route(
		"PATCH /{locale}/sessions/_current",
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from cookie
			sessionID, err := GetSessionIDFromCookie(ctx.Request, authService.Config)
			if err != nil {
				return ctx.Results.Unauthorized(
					httpfx.WithErrorMessage("No session found"),
				)
			}

			// Verify session exists
			session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if err != nil || session == nil || session.Status != users.SessionStatusActive {
				ClearSessionCookie(ctx.ResponseWriter, authService.Config)

				return ctx.Results.Unauthorized(
					httpfx.WithErrorMessage("Invalid session"),
				)
			}

			// Parse request body
			var req UpdatePreferencesRequest
			if err := json.NewDecoder(ctx.Request.Body).Decode(&req); err != nil {
				return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
			}

			// Update preferences
			if req.Theme != nil {
				err := sessionService.SetPreference(
					ctx.Request.Context(),
					sessionID,
					sessions.PrefKeyTheme,
					*req.Theme,
				)
				if err != nil {
					if errors.Is(err, sessions.ErrInvalidPreferenceValue) {
						return ctx.Results.BadRequest(
							httpfx.WithErrorMessage("Invalid theme value"),
						)
					}

					logger.ErrorContext(ctx.Request.Context(), "Failed to set theme preference",
						slog.String("error", err.Error()),
						slog.String("session_id", sessionID))

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithErrorMessage("Failed to update theme"),
					)
				}
			}

			if req.Locale != nil {
				err := sessionService.SetPreference(
					ctx.Request.Context(),
					sessionID,
					sessions.PrefKeyLocale,
					*req.Locale,
				)
				if err != nil {
					if errors.Is(err, sessions.ErrInvalidPreferenceValue) {
						return ctx.Results.BadRequest(
							httpfx.WithErrorMessage("Invalid locale value"),
						)
					}

					logger.ErrorContext(ctx.Request.Context(), "Failed to set locale preference",
						slog.String("error", err.Error()),
						slog.String("session_id", sessionID))

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithErrorMessage("Failed to update locale"),
					)
				}
			}

			if req.Timezone != nil {
				err := sessionService.SetPreference(
					ctx.Request.Context(),
					sessionID,
					sessions.PrefKeyTimezone,
					*req.Timezone,
				)
				if err != nil {
					if errors.Is(err, sessions.ErrInvalidPreferenceValue) {
						return ctx.Results.BadRequest(
							httpfx.WithErrorMessage("Invalid timezone value"),
						)
					}

					logger.ErrorContext(ctx.Request.Context(), "Failed to set timezone preference",
						slog.String("error", err.Error()),
						slog.String("session_id", sessionID))

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithErrorMessage("Failed to update timezone"),
					)
				}
			}

			// Get updated preferences for response
			prefs, _ := sessionService.GetPreferences(ctx.Request.Context(), sessionID)

			return ctx.Results.JSON(map[string]any{
				"data": map[string]any{
					"preferences": prefs,
				},
				"error": nil,
			})
		}).
		HasSummary("Update Session Preferences").
		HasDescription("Updates the current session preferences (theme, locale, timezone).").
		HasResponse(http.StatusOK)

	// GET /{locale}/sessions/list - List all sessions for current user
	routes.Route(
		"GET /{locale}/sessions/list",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get user ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get user from session
			session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if err != nil || session == nil || session.LoggedInUserID == nil {
				return ctx.Results.Unauthorized(
					httpfx.WithErrorMessage("User not logged in"),
				)
			}

			// List all sessions for user
			sessions, err := userService.ListSessionsByUserID(
				ctx.Request.Context(),
				*session.LoggedInUserID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to list sessions",
					slog.String("error", err.Error()),
					slog.String("user_id", *session.LoggedInUserID))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to list sessions"),
				)
			}

			// Map sessions to response format (hide sensitive fields)
			sessionList := make([]map[string]any, 0, len(sessions))
			for _, s := range sessions {
				sessionList = append(sessionList, map[string]any{
					"id":               s.ID,
					"status":           s.Status,
					"user_agent":       s.UserAgent,
					"last_activity_at": s.LastActivityAt,
					"logged_in_at":     s.LoggedInAt,
					"created_at":       s.CreatedAt,
					"is_current":       s.ID == sessionID,
				})
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]any{
					"sessions": sessionList,
				},
				"error": nil,
			})
		}).
		HasSummary("List User Sessions").
		HasDescription("Lists all sessions for the current user.").
		HasResponse(http.StatusOK)

	// POST /{locale}/sessions/{sessionId}/terminate - Terminate a specific session
	routes.Route(
		"POST /{locale}/sessions/{sessionId}/terminate",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get current session ID from context
			currentSessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			// Get target session ID from path
			targetSessionID := ctx.Request.PathValue("sessionId")
			if targetSessionID == "" {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage("Session ID is required"),
				)
			}

			// Prevent terminating current session (use logout instead)
			if targetSessionID == currentSessionID {
				return ctx.Results.BadRequest(
					httpfx.WithErrorMessage(
						"Cannot terminate current session. Use logout instead.",
					),
				)
			}

			// Get user from current session
			session, err := userService.GetSessionByID(ctx.Request.Context(), currentSessionID)
			if err != nil || session == nil || session.LoggedInUserID == nil {
				return ctx.Results.Unauthorized(
					httpfx.WithErrorMessage("User not logged in"),
				)
			}

			// Terminate the target session (only if it belongs to the same user)
			err = userService.TerminateSession(
				ctx.Request.Context(),
				targetSessionID,
				*session.LoggedInUserID,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to terminate session",
					slog.String("error", err.Error()),
					slog.String("target_session_id", targetSessionID),
					slog.String("user_id", *session.LoggedInUserID))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to terminate session"),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data": map[string]any{
					"status": "terminated",
				},
				"error": nil,
			})
		}).
		HasSummary("Terminate Session").
		HasDescription("Terminates a specific session for the current user.").
		HasResponse(http.StatusOK)
}
