package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterHTTPRoutesForProfileCandidates registers routes for managing profile membership candidates.
func RegisterHTTPRoutesForProfileCandidates( //nolint:gocognit,gocyclo,cyclop,funlen,maintidx
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	mailboxService *mailbox.Service,
) {
	// List candidates for a profile (member+ only)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_candidates",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			candidates, err := profileService.ListCandidates(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				slugParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to list candidates",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrInsufficientAccess) {
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			if candidates == nil {
				candidates = []*profiles.ProfileMembershipCandidate{}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  candidates,
				"error": nil,
			})
		},
	).HasDescription("List profile membership candidates")

	// Create a new candidate (member+ only)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_candidates",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			var input struct {
				ReferredProfileSlug string   `json:"referred_profile_slug"`
				TeamIDs             []string `json:"team_ids"`
			}

			err := json.NewDecoder(ctx.Request.Body).Decode(&input)
			if err != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			if input.ReferredProfileSlug == "" {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("referred_profile_slug is required"),
				)
			}

			candidate, err := profileService.CreateCandidate(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				input.ReferredProfileSlug,
				input.TeamIDs,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to create candidate",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError

				switch {
				case errors.Is(err, profiles.ErrInsufficientAccess):
					statusCode = http.StatusForbidden
				case errors.Is(err, profiles.ErrCandidateAlreadyExists):
					statusCode = http.StatusConflict
				case errors.Is(err, profiles.ErrCannotReferSelf):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrCannotReferExistingMember):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrCannotReferNonIndividual):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrProfileNotFound):
					statusCode = http.StatusNotFound
				case errors.Is(err, profiles.ErrInvalidInput):
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  candidate,
				"error": nil,
			})
		},
	).HasDescription("Create a membership candidate")

	// Get votes for a candidate (member+ only)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_candidates/{id}/votes",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			idParam := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			votes, err := profileService.GetCandidateVotes(
				ctx.Request.Context(),
				localeParam,
				*session.LoggedInUserID,
				slugParam,
				idParam,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to get candidate votes",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("candidateId", idParam))

				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrInsufficientAccess) {
					statusCode = http.StatusForbidden
				} else if errors.Is(err, profiles.ErrCandidateNotFound) {
					statusCode = http.StatusNotFound
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			if votes == nil {
				votes = []*profiles.CandidateVote{}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  votes,
				"error": nil,
			})
		},
	).HasDescription("Get votes for a candidate")

	// Cast or update a vote on a candidate (member+ only)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_candidates/{id}/votes",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			idParam := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			var input struct {
				Score   int16   `json:"score"`
				Comment *string `json:"comment"`
			}

			err := json.NewDecoder(ctx.Request.Body).Decode(&input)
			if err != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			vote, err := profileService.VoteCandidate(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				idParam,
				input.Score,
				input.Comment,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to vote on candidate",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("candidateId", idParam))

				statusCode := http.StatusInternalServerError

				switch {
				case errors.Is(err, profiles.ErrInsufficientAccess):
					statusCode = http.StatusForbidden
				case errors.Is(err, profiles.ErrCandidateNotFound):
					statusCode = http.StatusNotFound
				case errors.Is(err, profiles.ErrInvalidVoteScore):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrCandidateNotVoting):
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  vote,
				"error": nil,
			})
		},
	).HasDescription("Cast or update a vote on a candidate")

	// Update candidate status (lead+ only)
	routes.Route(
		"PATCH /{locale}/profiles/{slug}/_candidates/{id}/status",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam, localeOk := validateLocale(ctx)
			if !localeOk {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid locale"),
				)
			}

			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			idParam := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			var input struct {
				Status string `json:"status"`
			}

			err := json.NewDecoder(ctx.Request.Body).Decode(&input)
			if err != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			newStatus := profiles.CandidateStatus(input.Status)

			err = profileService.UpdateCandidateStatus(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				idParam,
				newStatus,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to update candidate status",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam),
					slog.String("candidateId", idParam),
					slog.String("newStatus", input.Status))

				statusCode := http.StatusInternalServerError

				switch {
				case errors.Is(err, profiles.ErrInsufficientAccess):
					statusCode = http.StatusForbidden
				case errors.Is(err, profiles.ErrCandidateNotFound):
					statusCode = http.StatusNotFound
				case errors.Is(err, profiles.ErrInvalidStatusTransition):
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			// When transitioning to invitation_pending_response, send a mailbox invitation.
			if newStatus == profiles.CandidateStatusInvitationPendingResponse {
				sendErr := sendCandidateInvitation(
					ctx.Request.Context(),
					logger,
					profileService,
					mailboxService,
					localeParam,
					slugParam,
					idParam,
				)
				if sendErr != nil {
					logger.ErrorContext(
						ctx.Request.Context(),
						"Failed to send candidate invitation",
						slog.String("error", sendErr.Error()),
						slog.String("candidateId", idParam),
					)
				} else {
					candidate, candErr := profileService.GetCandidateByID(ctx.Request.Context(), idParam)
					if candErr == nil {
						profileService.RecordCandidateInvitationSent(
							ctx.Request.Context(),
							*session.LoggedInUserID,
							idParam,
							candidate.ProfileID,
							candidate.ReferredProfileID,
						)
					}
				}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  map[string]string{"status": input.Status},
				"error": nil,
			})
		},
	).HasDescription("Update candidate status")

	// --- Application Form Endpoints ---

	// List application presets (public, locale-aware)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_application-presets",
		func(ctx *httpfx.Context) httpfx.Result {
			localeParam := ctx.Request.PathValue("locale")
			presets := profiles.ListApplicationPresets(localeParam)

			return ctx.Results.JSON(map[string]any{
				"data":  presets,
				"error": nil,
			})
		},
	).HasDescription("List application form presets")

	// Get the active application form for a profile (public)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_application-form",
		func(ctx *httpfx.Context) httpfx.Result {
			slugParam := ctx.Request.PathValue("slug")

			form, err := profileService.GetApplicationForm(
				ctx.Request.Context(),
				slugParam,
			)
			if err != nil {
				statusCode := http.StatusInternalServerError

				switch {
				case errors.Is(err, profiles.ErrApplicationsNotEnabled):
					statusCode = http.StatusNotFound
				case errors.Is(err, profiles.ErrNoApplicationForm):
					statusCode = http.StatusNotFound
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  form,
				"error": nil,
			})
		},
	).HasDescription("Get application form for a profile")

	// Upsert application form (lead+ only)
	routes.Route(
		"PUT /{locale}/profiles/{slug}/_application-form",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			var input struct {
				PresetKey           *string                              `json:"preset_key"`
				Fields              []profiles.ApplicationFormFieldInput `json:"fields"`
				ResponsesVisibility string                               `json:"responses_visibility"`
			}

			err := json.NewDecoder(ctx.Request.Body).Decode(&input)
			if err != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			form, err := profileService.UpsertApplicationForm(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				input.PresetKey,
				input.Fields,
				input.ResponsesVisibility,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to upsert application form",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError

				switch {
				case errors.Is(err, profiles.ErrInsufficientAccess):
					statusCode = http.StatusForbidden
				case errors.Is(err, profiles.ErrInvalidFieldType):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrInvalidResponsesVisibility):
					statusCode = http.StatusBadRequest
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  form,
				"error": nil,
			})
		},
	).HasDescription("Create or update application form")

	// Submit an application (authenticated, non-member)
	routes.Route(
		"POST /{locale}/profiles/{slug}/_candidates/apply",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			var input struct {
				ApplicantMessage *string           `json:"applicant_message"`
				FormResponses    map[string]string `json:"form_responses"`
			}

			err := json.NewDecoder(ctx.Request.Body).Decode(&input)
			if err != nil {
				return ctx.Results.Error(
					http.StatusBadRequest,
					httpfx.WithErrorMessage("Invalid request body"),
				)
			}

			candidate, err := profileService.CreateApplication(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				input.ApplicantMessage,
				input.FormResponses,
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Failed to create application",
					slog.String("error", err.Error()),
					slog.String("slug", slugParam))

				statusCode := http.StatusInternalServerError

				switch {
				case errors.Is(err, profiles.ErrApplicationsNotEnabled):
					statusCode = http.StatusForbidden
				case errors.Is(err, profiles.ErrNoApplicationForm):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrAlreadyApplied):
					statusCode = http.StatusConflict
				case errors.Is(err, profiles.ErrCannotReferExistingMember):
					statusCode = http.StatusConflict
				case errors.Is(err, profiles.ErrMissingRequiredField):
					statusCode = http.StatusBadRequest
				case errors.Is(err, profiles.ErrInsufficientAccess):
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			return ctx.Results.JSON(map[string]any{
				"data":  candidate,
				"error": nil,
			})
		},
	).HasDescription("Submit a membership application")

	// Get my application status (authenticated)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_candidates/my-application",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			candidate, err := profileService.GetMyApplication(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
			)
			if err != nil {
				if errors.Is(err, profiles.ErrCandidateNotFound) {
					return ctx.Results.JSON(map[string]any{
						"data":  nil,
						"error": nil,
					})
				}

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithSanitizedError(err),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"data":  candidate,
				"error": nil,
			})
		},
	).HasDescription("Get current user's application status")

	// Get form responses for a candidate (member+/lead+ depending on visibility)
	routes.Route(
		"GET /{locale}/profiles/{slug}/_candidates/{id}/responses",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Session ID not found in context"),
				)
			}

			slugParam := ctx.Request.PathValue("slug")
			idParam := ctx.Request.PathValue("id")

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithErrorMessage("Failed to get session information"),
				)
			}

			responses, err := profileService.GetCandidateFormResponses(
				ctx.Request.Context(),
				*session.LoggedInUserID,
				slugParam,
				idParam,
			)
			if err != nil {
				statusCode := http.StatusInternalServerError
				if errors.Is(err, profiles.ErrInsufficientAccess) {
					statusCode = http.StatusForbidden
				}

				return ctx.Results.Error(statusCode, httpfx.WithSanitizedError(err))
			}

			if responses == nil {
				responses = []*profiles.CandidateFormResponse{}
			}

			return ctx.Results.JSON(map[string]any{
				"data":  responses,
				"error": nil,
			})
		},
	).HasDescription("Get form responses for a candidate")
}

// sendCandidateInvitation sends a mailbox invitation to the referred profile.
func sendCandidateInvitation(
	ctx context.Context,
	logger *logfx.Logger,
	profileService *profiles.Service,
	mailboxService *mailbox.Service,
	localeCode string,
	profileSlug string,
	candidateID string,
) error {
	_ = localeCode

	candidate, err := profileService.GetCandidateByID(ctx, candidateID)
	if err != nil {
		return fmt.Errorf("failed to get candidate: %w", err)
	}

	profile, profileErr := profileService.GetBySlugInternal(ctx, profileSlug)
	if profileErr != nil {
		return fmt.Errorf("failed to get profile: %w", profileErr)
	}

	message := "You have been invited to join " + profile.Title
	props := mailbox.InvitationProperties{
		InvitationKind:   mailbox.InvitationKindProfileJoin,
		TelegramChatID:   0,
		GroupProfileSlug: "",
		GroupName:        "",
		InviteLink:       nil,
		CandidateID:      candidateID,
		ProfileID:        candidate.ProfileID,
		ProfileSlug:      profileSlug,
	}

	_, envelopeErr := mailboxService.SendSystemEnvelope(ctx, &mailbox.SendMessageParams{
		SenderProfileID:    candidate.ProfileID,
		TargetProfileID:    candidate.ReferredProfileID,
		SenderUserID:       nil,
		Kind:               mailbox.KindInvitation,
		ConversationTitle:  "Membership Invitation",
		Message:            &message,
		Properties:         props,
		ReplyToID:          nil,
		SenderProfileTitle: profile.Title,
		Locale:             "",
	})
	if envelopeErr != nil {
		return fmt.Errorf("failed to send envelope: %w", envelopeErr)
	}

	logger.InfoContext(ctx, "Candidate invitation sent",
		slog.String("candidate_id", candidateID),
		slog.String("referred_profile_id", candidate.ReferredProfileID))

	return nil
}
