package http

import (
	"encoding/json"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	telegrambiz "github.com/eser/aya.is/services/pkg/api/business/telegram"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func RegisterHTTPRoutesForProfileEnvelopes( //nolint:funlen,cyclop
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	mailboxService *mailbox.Service,
	telegramService *telegrambiz.Service,
) {
	// GET /{locale}/profiles/{slug}/_envelopes — list envelopes (inbox)
	routes.
		Route(
			"GET /{locale}/profiles/{slug}/_envelopes",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				slugParam := ctx.Request.PathValue("slug")

				sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
				if !ok {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Session ID not found"))
				}

				session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
				if err != nil || session == nil || session.LoggedInUserID == nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid session"))
				}

				canEdit, permErr := profileService.HasUserAccessToProfile(
					ctx.Request.Context(), *session.LoggedInUserID, slugParam,
					profiles.MembershipKindMaintainer,
				)
				if permErr != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(permErr),
					)
				}

				if !canEdit {
					return ctx.Results.Error(http.StatusForbidden,
						httpfx.WithErrorMessage("You do not have permission to view this inbox"))
				}

				profile, profileErr := profileService.GetBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
				)
				if profileErr != nil || profile == nil {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("profile not found"))
				}

				statusFilter := ctx.Request.URL.Query().Get("status")

				items, listErr := mailboxService.ListEnvelopes(
					ctx.Request.Context(), profile.ID, statusFilter,
				)
				if listErr != nil {
					logger.Error("failed to list envelopes", "error", listErr)

					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithSanitizedError(listErr),
					)
				}

				return ctx.Results.JSON(map[string]any{"data": items})
			},
		).
		HasSummary("List profile envelopes").
		HasDescription("List inbox envelopes for a profile. Requires maintainer permission.").
		HasResponse(http.StatusOK)

	// POST /{locale}/profiles/{slug}/_envelopes — create envelope
	routes.
		Route(
			"POST /{locale}/profiles/{slug}/_envelopes",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				slugParam := ctx.Request.PathValue("slug")

				user, userErr := getUserFromContext(ctx, userService)
				if userErr != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(userErr))
				}

				// Admin or Lead/Owner of the sender profile
				isAdmin := user.Kind == "admin"

				if !isAdmin {
					canSend, permErr := profileService.HasUserAccessToProfile(
						ctx.Request.Context(), user.ID, slugParam,
						profiles.MembershipKindLead,
					)
					if permErr != nil {
						return ctx.Results.Error(
							http.StatusInternalServerError, httpfx.WithSanitizedError(permErr),
						)
					}

					if !canSend {
						return ctx.Results.Error(
							http.StatusForbidden,
							httpfx.WithErrorMessage(
								"You must be an admin or lead/owner to send envelopes",
							),
						)
					}
				}

				var body struct {
					Kind              string `json:"kind"`
					TargetProfileID   string `json:"target_profile_id"`
					ConversationTitle string `json:"conversation_title"`
					Message           string `json:"message"`
					InviteCode        string `json:"invite_code"`
					Properties        any    `json:"properties"`
				}

				err := json.NewDecoder(ctx.Request.Body).Decode(&body)
				if err != nil {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
				}

				if body.Kind == "" || body.TargetProfileID == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("kind and target_profile_id are required"),
					)
				}

				if body.Message == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("message is required"),
					)
				}

				// If an invite code is provided, resolve it and merge properties
				if body.InviteCode != "" {
					if telegramService == nil {
						return ctx.Results.BadRequest(
							httpfx.WithErrorMessage("External code service is not configured"),
						)
					}

					externalCode, resolveErr := telegramService.ResolveGroupInviteCode(
						ctx.Request.Context(), body.InviteCode,
					)
					if resolveErr != nil {
						return ctx.Results.BadRequest(
							httpfx.WithErrorMessage("Invalid or expired invite code"),
						)
					}

					// Map external code properties into envelope properties.
					// The envelope format is the contract consumed by the bot.
					props := map[string]any{
						"external_system":  externalCode.ExternalSystem,
						"invitation_kind":  mailbox.InvitationKindTelegramGroup,
						"telegram_chat_id": externalCode.Properties["telegram_chat_id"],
						"group_name":       externalCode.Properties["telegram_chat_title"],
					}

					body.Properties = props

					_ = telegramService.ConsumeGroupInviteCode(
						ctx.Request.Context(), body.InviteCode,
					)
				}

				// Get the sender profile to get its ID
				senderProfile, senderErr := profileService.GetBySlug(
					ctx.Request.Context(), "en", slugParam,
				)
				if senderErr != nil || senderProfile == nil {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("sender profile not found"))
				}

				msgPtr := &body.Message

				envelope, createErr := mailboxService.SendSystemEnvelope(
					ctx.Request.Context(),
					&mailbox.SendMessageParams{
						TargetProfileID:    body.TargetProfileID,
						SenderProfileID:    senderProfile.ID,
						SenderUserID:       &user.ID,
						Kind:               body.Kind,
						ConversationTitle:  body.ConversationTitle,
						Message:            msgPtr,
						Properties:         body.Properties,
						SenderProfileTitle: senderProfile.Title,
						Locale:             localeParam,
					},
				)
				if createErr != nil {
					logger.Error("failed to create envelope", "error", createErr)

					return ctx.Results.Error(
						http.StatusInternalServerError, httpfx.WithSanitizedError(createErr),
					)
				}

				return ctx.Results.JSON(map[string]any{"data": envelope})
			},
		).
		HasSummary("Create profile envelope").
		HasDescription("Create an envelope for a target profile. Requires admin or lead/owner of sender profile.").
		HasResponse(http.StatusOK)

	// POST /{locale}/profiles/{slug}/_envelopes/{id}/accept
	routes.
		Route(
			"POST /{locale}/profiles/{slug}/_envelopes/{id}/accept",
			AuthMiddleware(authService, userService),
			envelopeActionHandler(logger, userService, profileService, mailboxService, "accept"),
		).
		HasSummary("Accept envelope").
		HasDescription("Accept a pending envelope.").
		HasResponse(http.StatusOK)

	// POST /{locale}/profiles/{slug}/_envelopes/{id}/reject
	routes.
		Route(
			"POST /{locale}/profiles/{slug}/_envelopes/{id}/reject",
			AuthMiddleware(authService, userService),
			envelopeActionHandler(logger, userService, profileService, mailboxService, "reject"),
		).
		HasSummary("Reject envelope").
		HasDescription("Reject a pending envelope.").
		HasResponse(http.StatusOK)

	// POST /{locale}/profiles/{slug}/_envelopes/{id}/revoke
	routes.
		Route(
			"POST /{locale}/profiles/{slug}/_envelopes/{id}/revoke",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				slugParam := ctx.Request.PathValue("slug")
				envelopeID := ctx.Request.PathValue("id")

				user, userErr := getUserFromContext(ctx, userService)
				if userErr != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(userErr))
				}

				// Admin or Lead/Owner of sender profile can revoke
				isAdmin := user.Kind == "admin"

				if !isAdmin {
					canRevoke, permErr := profileService.HasUserAccessToProfile(
						ctx.Request.Context(), user.ID, slugParam,
						profiles.MembershipKindLead,
					)
					if permErr != nil {
						return ctx.Results.Error(
							http.StatusInternalServerError, httpfx.WithSanitizedError(permErr),
						)
					}

					if !canRevoke {
						return ctx.Results.Error(
							http.StatusForbidden,
							httpfx.WithErrorMessage(
								"You must be an admin or lead/owner to revoke envelopes",
							),
						)
					}
				}

				err := mailboxService.RevokeEnvelope(ctx.Request.Context(), envelopeID)
				if err != nil {
					return ctx.Results.Error(
						http.StatusBadRequest, httpfx.WithErrorMessage(err.Error()),
					)
				}

				return ctx.Results.JSON(map[string]any{"data": "ok"})
			},
		).
		HasSummary("Revoke envelope").
		HasDescription("Revoke a pending envelope. Requires admin or lead/owner of sender profile.").
		HasResponse(http.StatusOK)
}

// envelopeActionHandler creates a handler for accept/reject actions on envelopes.
func envelopeActionHandler(
	logger *logfx.Logger,
	userService *users.Service,
	profileService *profiles.Service,
	mailboxService *mailbox.Service,
	action string,
) func(ctx *httpfx.Context) httpfx.Result {
	return func(ctx *httpfx.Context) httpfx.Result {
		localeParam, localeOk := validateLocale(ctx)
		if !localeOk {
			return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
		}

		slugParam := ctx.Request.PathValue("slug")
		envelopeID := ctx.Request.PathValue("id")

		sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
		if !ok {
			return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Session ID not found"))
		}

		session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
		if err != nil || session == nil || session.LoggedInUserID == nil {
			return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid session"))
		}

		canEdit, permErr := profileService.HasUserAccessToProfile(
			ctx.Request.Context(), *session.LoggedInUserID, slugParam,
			profiles.MembershipKindMaintainer,
		)
		if permErr != nil {
			return ctx.Results.Error(
				http.StatusInternalServerError,
				httpfx.WithSanitizedError(permErr),
			)
		}

		if !canEdit {
			return ctx.Results.Error(http.StatusForbidden,
				httpfx.WithErrorMessage("You do not have permission"))
		}

		profile, profileErr := profileService.GetBySlug(
			ctx.Request.Context(),
			localeParam,
			slugParam,
		)
		if profileErr != nil || profile == nil {
			return ctx.Results.NotFound(httpfx.WithErrorMessage("profile not found"))
		}

		var actionErr error

		switch action {
		case "accept":
			actionErr = mailboxService.AcceptEnvelope(
				ctx.Request.Context(),
				envelopeID,
				profile.ID,
			)
		case "reject":
			actionErr = mailboxService.RejectEnvelope(
				ctx.Request.Context(),
				envelopeID,
				profile.ID,
			)
		}

		if actionErr != nil {
			logger.Error("envelope action failed", "action", action, "error", actionErr)

			return ctx.Results.Error(
				http.StatusBadRequest, httpfx.WithErrorMessage(actionErr.Error()),
			)
		}

		return ctx.Results.JSON(map[string]any{"data": "ok"})
	}
}
