package http

import (
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	envelopes "github.com/eser/aya.is/services/pkg/api/business/profile_envelopes"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// mailboxEnvelope wraps an envelope with the owning profile metadata.
type mailboxEnvelope struct {
	*envelopes.Envelope

	OwningProfileSlug  string `json:"owning_profile_slug"`
	OwningProfileTitle string `json:"owning_profile_title"`
	OwningProfileKind  string `json:"owning_profile_kind"`
}

func RegisterHTTPRoutesForMailbox(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	envelopeService *envelopes.Service,
) {
	// GET /{locale}/mailbox — aggregated inbox across all maintainer+ profiles
	routes.
		Route(
			"GET /{locale}/mailbox",
			AuthMiddleware(authService, userService),
			func(ctx *httpfx.Context) httpfx.Result {
				localeParam, localeOk := validateLocale(ctx)
				if !localeOk {
					return ctx.Results.BadRequest(httpfx.WithErrorMessage("unsupported locale"))
				}

				sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
				if !ok {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Session ID not found"))
				}

				session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
				if err != nil || session == nil || session.LoggedInUserID == nil {
					return ctx.Results.Unauthorized(httpfx.WithErrorMessage("Invalid session"))
				}

				// Get user to find their individual profile
				user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
				if userErr != nil || user == nil || user.IndividualProfileID == nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithErrorMessage("Could not resolve user profile"),
					)
				}

				// Fetch the individual profile to get its slug/title
				individualProfile, ipErr := profileService.GetByID(
					ctx.Request.Context(),
					localeParam,
					*user.IndividualProfileID,
				)
				if ipErr != nil || individualProfile == nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithErrorMessage("Could not resolve individual profile"),
					)
				}

				statusFilter := ctx.Request.URL.Query().Get("status")

				var allEnvelopes []mailboxEnvelope

				// 1. Fetch envelopes for user's own individual profile
				items, listErr := envelopeService.ListEnvelopes(
					ctx.Request.Context(), individualProfile.ID, statusFilter,
				)
				if listErr != nil {
					logger.ErrorContext(
						ctx.Request.Context(),
						"failed to list individual envelopes",
						slog.String("error", listErr.Error()),
					)
				} else {
					for _, item := range items {
						allEnvelopes = append(allEnvelopes, mailboxEnvelope{
							Envelope:           item,
							OwningProfileSlug:  individualProfile.Slug,
							OwningProfileTitle: individualProfile.Title,
							OwningProfileKind:  individualProfile.Kind,
						})
					}
				}

				// 2. Fetch memberships to discover maintainer+ profiles
				memberships, membershipErr := profileService.GetMembershipsByUserProfileID(
					ctx.Request.Context(),
					localeParam,
					*user.IndividualProfileID,
				)
				if membershipErr != nil {
					logger.ErrorContext(ctx.Request.Context(), "failed to list memberships",
						slog.String("error", membershipErr.Error()))
				} else {
					for _, m := range memberships {
						if m.Profile == nil {
							continue
						}

						level := profiles.MembershipKindLevel[profiles.MembershipKind(m.Kind)]
						minLevel := profiles.MembershipKindLevel[profiles.MembershipKindMaintainer]

						if level < minLevel {
							continue
						}

						mItems, mErr := envelopeService.ListEnvelopes(
							ctx.Request.Context(), m.Profile.ID, statusFilter,
						)
						if mErr != nil {
							logger.ErrorContext(ctx.Request.Context(), "failed to list envelopes for profile",
								slog.String("error", mErr.Error()),
								slog.String("profile_id", m.Profile.ID))

							continue
						}

						for _, item := range mItems {
							allEnvelopes = append(allEnvelopes, mailboxEnvelope{
								Envelope:           item,
								OwningProfileSlug:  m.Profile.Slug,
								OwningProfileTitle: m.Profile.Title,
								OwningProfileKind:  m.Profile.Kind,
							})
						}
					}
				}

				return ctx.Results.JSON(map[string]any{"data": allEnvelopes})
			},
		).
		HasSummary("List aggregated mailbox").
		HasDescription("List inbox envelopes across all profiles where the user is maintainer+.").
		HasResponse(http.StatusOK)
}
