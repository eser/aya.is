package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/mailbox"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

var (
	ErrMailboxSessionNotFound = errors.New("session ID not found")
	ErrMailboxInvalidSession  = errors.New("invalid session")
	ErrMailboxUserNotFound    = errors.New("could not resolve user profile")
)

func RegisterHTTPRoutesForMailbox( //nolint:funlen,cyclop
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	mailboxService *mailbox.Service,
) {
	// GET /{locale}/mailbox/conversations — list conversations
	routes.
		Route(
			"GET /{locale}/mailbox/conversations",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				user, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
				}

				_ = user
				archivedParam := ctx.Request.URL.Query().Get("archived")
				includeArchived := archivedParam == "true"

				var allConversations []*mailbox.Conversation

				seen := make(map[string]bool)

				for _, profileID := range profileIDs {
					convs, listErr := mailboxService.ListConversations(
						ctx.Request.Context(), profileID, includeArchived, 50,
					)
					if listErr != nil {
						logger.ErrorContext(ctx.Request.Context(), "failed to list conversations",
							slog.String("error", listErr.Error()),
							slog.String("profile_id", profileID))

						continue
					}

					for _, conv := range convs {
						if seen[conv.ID] {
							continue
						}

						seen[conv.ID] = true

						// Populate all participants for each conversation.
						participants, pErr := mailboxService.GetConversationParticipants(
							ctx.Request.Context(), conv.ID,
						)
						if pErr == nil {
							conv.Participants = participants
						}

						allConversations = append(allConversations, conv)
					}
				}

				return ctx.Results.JSON(map[string]any{"data": allConversations})
			},
		).
		HasSummary("List mailbox conversations").
		HasDescription("List conversations across all profiles where the user is maintainer+.").
		HasResponse(http.StatusOK)

	// GET /{locale}/mailbox/conversations/{id} — get conversation detail
	routes.
		Route(
			"GET /{locale}/mailbox/conversations/{id}",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				conversationID := ctx.Request.PathValue("id")

				_, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
				}

				// Try each profile until one is a participant.
				var (
					conv      *mailbox.Conversation
					envelopes []*mailbox.Envelope
					found     bool
				)

				for _, profileID := range profileIDs {
					c, e, getErr := mailboxService.GetConversation(
						ctx.Request.Context(), conversationID, profileID,
					)
					if getErr == nil {
						conv = c
						envelopes = e
						found = true

						// Mark as read for this profile.
						_ = mailboxService.MarkConversationRead(
							ctx.Request.Context(), conversationID, profileID,
						)

						break
					}
				}

				if !found {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("conversation not found"))
				}

				return ctx.Results.JSON(map[string]any{
					"data": map[string]any{
						"conversation": conv,
						"envelopes":    envelopes,
					},
				})
			},
		).
		HasSummary("Get conversation detail").
		HasDescription("Get a conversation with its messages.").
		HasResponse(http.StatusOK)

	// POST /{locale}/mailbox/conversations/{id}/read — mark as read
	routes.
		Route(
			"POST /{locale}/mailbox/conversations/{id}/read",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				conversationID := ctx.Request.PathValue("id")

				_, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
				}

				for _, profileID := range profileIDs {
					_ = mailboxService.MarkConversationRead(
						ctx.Request.Context(), conversationID, profileID,
					)
				}

				return ctx.Results.JSON(map[string]any{"data": "ok"})
			},
		).
		HasSummary("Mark conversation as read").
		HasResponse(http.StatusOK)

	// POST /{locale}/mailbox/conversations/{id}/archive
	routes.
		Route(
			"POST /{locale}/mailbox/conversations/{id}/archive",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				conversationID := ctx.Request.PathValue("id")

				_, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
				}

				for _, profileID := range profileIDs {
					_ = mailboxService.ArchiveConversation(
						ctx.Request.Context(), conversationID, profileID,
					)
				}

				return ctx.Results.JSON(map[string]any{"data": "ok"})
			},
		).
		HasSummary("Archive conversation").
		HasResponse(http.StatusOK)

	// POST /{locale}/mailbox/conversations/{id}/unarchive
	routes.
		Route(
			"POST /{locale}/mailbox/conversations/{id}/unarchive",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				conversationID := ctx.Request.PathValue("id")

				_, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
				}

				for _, profileID := range profileIDs {
					_ = mailboxService.UnarchiveConversation(
						ctx.Request.Context(), conversationID, profileID,
					)
				}

				return ctx.Results.JSON(map[string]any{"data": "ok"})
			},
		).
		HasSummary("Unarchive conversation").
		HasResponse(http.StatusOK)

	// POST /{locale}/mailbox/messages — send a message
	routes.
		Route(
			"POST /{locale}/mailbox/messages",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				user, userErr := getUserFromContext(ctx, userService)
				if userErr != nil {
					return ctx.Results.Unauthorized(httpfx.WithSanitizedError(userErr))
				}

				var body struct {
					SenderProfileSlug string  `json:"sender_profile_slug"`
					TargetProfileSlug string  `json:"target_profile_slug"`
					Kind              string  `json:"kind"`
					ConversationTitle string  `json:"conversation_title"`
					Message           *string `json:"message"`
					ReplyToID         *string `json:"reply_to_id"`
				}

				err := json.NewDecoder(ctx.Request.Body).Decode(&body)
				if err != nil {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("Invalid request body"))
				}

				if body.SenderProfileSlug == "" || body.TargetProfileSlug == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage(
							"sender_profile_slug and target_profile_slug are required",
						),
					)
				}

				if body.Message == nil || *body.Message == "" {
					return ctx.Results.BadRequest(
						httpfx.WithErrorMessage("message is required"),
					)
				}

				if body.Kind == "" {
					body.Kind = mailbox.KindMessage
				}

				// Verify sender access.
				canSend, permErr := profileService.HasUserAccessToProfile(
					ctx.Request.Context(), user.ID, body.SenderProfileSlug,
					profiles.MembershipKindMaintainer,
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
							"You do not have permission to send from this profile",
						),
					)
				}

				senderProfile, sErr := profileService.GetBySlug(
					ctx.Request.Context(), localeParam, body.SenderProfileSlug,
				)
				if sErr != nil || senderProfile == nil {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("sender profile not found"))
				}

				targetProfile, tErr := profileService.GetBySlug(
					ctx.Request.Context(), localeParam, body.TargetProfileSlug,
				)
				if tErr != nil || targetProfile == nil {
					return ctx.Results.NotFound(httpfx.WithErrorMessage("target profile not found"))
				}

				envelope, createErr := mailboxService.SendMessage(
					ctx.Request.Context(),
					&mailbox.SendMessageParams{
						SenderProfileID:    senderProfile.ID,
						TargetProfileID:    targetProfile.ID,
						SenderUserID:       &user.ID,
						Kind:               body.Kind,
						ConversationTitle:  body.ConversationTitle,
						Message:            body.Message,
						ReplyToID:          body.ReplyToID,
						SenderProfileTitle: senderProfile.Title,
						Locale:             localeParam,
					},
				)
				if createErr != nil {
					logger.Error("failed to send message", "error", createErr)

					if errors.Is(createErr, mailbox.ErrInvalidEnvelopeKind) ||
						errors.Is(createErr, mailbox.ErrConversationPending) {
						return ctx.Results.BadRequest(
							httpfx.WithErrorMessage(createErr.Error()),
						)
					}

					return ctx.Results.Error(
						http.StatusInternalServerError, httpfx.WithSanitizedError(createErr),
					)
				}

				return ctx.Results.JSON(map[string]any{"data": envelope})
			},
		).
		HasSummary("Send a mailbox message").
		HasDescription("Send a message, creating or finding a direct conversation.").
		HasResponse(http.StatusOK)

	// POST /{locale}/mailbox/messages/{id}/accept
	routes.
		Route(
			"POST /{locale}/mailbox/messages/{id}/accept",
			AuthMiddleware(authService, userService),
			mailboxMessageActionHandler(
				logger,
				userService,
				profileService,
				mailboxService,
				"accept",
			),
		).
		HasSummary("Accept mailbox message").
		HasResponse(http.StatusOK)

	// POST /{locale}/mailbox/messages/{id}/reject
	routes.
		Route(
			"POST /{locale}/mailbox/messages/{id}/reject",
			AuthMiddleware(authService, userService),
			mailboxMessageActionHandler(
				logger,
				userService,
				profileService,
				mailboxService,
				"reject",
			),
		).
		HasSummary("Reject mailbox message").
		HasResponse(http.StatusOK)

	// POST /{locale}/mailbox/messages/{id}/revoke
	routes.
		Route(
			"POST /{locale}/mailbox/messages/{id}/revoke",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				envelopeID := ctx.Request.PathValue("id")

				err := mailboxService.RevokeEnvelope(ctx.Request.Context(), envelopeID)
				if err != nil {
					return ctx.Results.Error(
						http.StatusBadRequest, httpfx.WithErrorMessage(err.Error()),
					)
				}

				return ctx.Results.JSON(map[string]any{"data": "ok"})
			},
		).
		HasSummary("Revoke mailbox message").
		HasResponse(http.StatusOK)

	// POST /{locale}/mailbox/messages/{id}/reactions — add reaction
	routes.
		Route(
			"POST /{locale}/mailbox/messages/{id}/reactions",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				envelopeID := ctx.Request.PathValue("id")

				_, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
				}

				var body struct {
					Emoji string `json:"emoji"`
				}

				decErr := json.NewDecoder(ctx.Request.Body).Decode(&body)
				if decErr != nil || body.Emoji == "" {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("emoji is required"))
				}

				// Use the individual profile for reactions.
				if len(profileIDs) == 0 {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("no profile"))
				}

				addErr := mailboxService.AddReaction(
					ctx.Request.Context(), envelopeID, profileIDs[0], body.Emoji,
				)
				if addErr != nil {
					return ctx.Results.Error(
						http.StatusBadRequest, httpfx.WithErrorMessage(addErr.Error()),
					)
				}

				return ctx.Results.JSON(map[string]any{"data": "ok"})
			},
		).
		HasSummary("Add reaction to message").
		HasResponse(http.StatusOK)

	// DELETE /{locale}/mailbox/messages/{id}/reactions/{emoji} — remove reaction
	routes.
		Route(
			"DELETE /{locale}/mailbox/messages/{id}/reactions/{emoji}",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				envelopeID := ctx.Request.PathValue("id")
				emoji := ctx.Request.PathValue("emoji")

				_, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
				}

				if len(profileIDs) == 0 {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("no profile"))
				}

				rmErr := mailboxService.RemoveReaction(
					ctx.Request.Context(), envelopeID, profileIDs[0], emoji,
				)
				if rmErr != nil {
					return ctx.Results.Error(
						http.StatusBadRequest, httpfx.WithErrorMessage(rmErr.Error()),
					)
				}

				return ctx.Results.JSON(map[string]any{"data": "ok"})
			},
		).
		HasSummary("Remove reaction from message").
		HasResponse(http.StatusOK)

	// GET /{locale}/mailbox/unread-count — total unread count
	routes.
		Route(
			"GET /{locale}/mailbox/unread-count",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				_, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				_, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
				}

				totalUnread := 0

				for _, profileID := range profileIDs {
					count, cErr := mailboxService.CountPendingEnvelopes(
						ctx.Request.Context(),
						profileID,
					)
					if cErr == nil {
						totalUnread += count
					}
				}

				return ctx.Results.JSON(map[string]any{"data": totalUnread})
			},
		).
		HasSummary("Get unread count").
		HasResponse(http.StatusOK)
}

// resolveMailboxProfileIDs returns the user and all profile IDs where the user has maintainer+ access.
func resolveMailboxProfileIDs(
	ctx *httpfx.Context,
	userService *users.Service,
	profileService *profiles.Service,
) (*users.User, []string, error) {
	sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
	if !ok {
		return nil, nil, ErrMailboxSessionNotFound
	}

	session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
	if err != nil || session == nil || session.LoggedInUserID == nil {
		return nil, nil, ErrMailboxInvalidSession
	}

	user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
	if userErr != nil || user == nil || user.IndividualProfileID == nil {
		return nil, nil, ErrMailboxUserNotFound
	}

	profileIDs := []string{*user.IndividualProfileID}

	memberships, membershipErr := profileService.GetMembershipsByUserProfileID(
		ctx.Request.Context(),
		"en",
		*user.IndividualProfileID,
	)
	if membershipErr == nil {
		for _, m := range memberships {
			if m.Profile == nil {
				continue
			}

			level := profiles.MembershipKindLevel[profiles.MembershipKind(m.Kind)]
			minLevel := profiles.MembershipKindLevel[profiles.MembershipKindMaintainer]

			if level < minLevel {
				continue
			}

			profileIDs = append(profileIDs, m.Profile.ID)
		}
	}

	return user, profileIDs, nil
}

// mailboxMessageActionHandler creates a handler for accept/reject actions on mailbox messages.
func mailboxMessageActionHandler(
	logger *logfx.Logger,
	userService *users.Service,
	profileService *profiles.Service,
	mailboxService *mailbox.Service,
	action string,
) func(ctx *httpfx.Context) httpfx.Result {
	return func(ctx *httpfx.Context) httpfx.Result {
		_, localeOk := validateLocale(ctx)
		if !localeOk {
			return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
		}

		envelopeID := ctx.Request.PathValue("id")

		_, profileIDs, err := resolveMailboxProfileIDs(ctx, userService, profileService)
		if err != nil {
			return ctx.Results.Unauthorized(httpfx.WithErrorMessage(err.Error()))
		}

		// Find which profile owns this envelope and perform the action.
		var actionErr error

		actionDone := false

		for _, profileID := range profileIDs {
			switch action {
			case "accept":
				actionErr = mailboxService.AcceptEnvelope(
					ctx.Request.Context(),
					envelopeID,
					profileID,
				)
			case "reject":
				actionErr = mailboxService.RejectEnvelope(
					ctx.Request.Context(),
					envelopeID,
					profileID,
				)
			}

			if actionErr == nil {
				actionDone = true

				break
			}
		}

		if !actionDone {
			logger.Error("mailbox message action failed", "action", action, "error", actionErr)

			return ctx.Results.Error(
				http.StatusBadRequest,
				httpfx.WithErrorMessage("action failed or permission denied"),
			)
		}

		return ctx.Results.JSON(map[string]any{"data": "ok"})
	}
}
