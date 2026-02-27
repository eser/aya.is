package resend

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/i18nfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	bulletinbiz "github.com/eser/aya.is/services/pkg/api/business/bulletin"
)

// BulletinSender sends bulletin digests via email using the Resend API.
type BulletinSender struct {
	client         *Client
	emailResolver  bulletinbiz.UserEmailResolver
	logger         *logfx.Logger
	localizer      *i18nfx.Localizer
	fromAddress    string
	frontendURI    string
	tmpl           *template.Template
	sandboxMode    bool
	sandboxAllowed map[string]bool
}

// NewBulletinSender creates a new email bulletin channel adapter.
func NewBulletinSender(
	client *Client,
	emailResolver bulletinbiz.UserEmailResolver,
	logger *logfx.Logger,
	localizer *i18nfx.Localizer,
	config *Config,
	frontendURI string,
	templatePath string,
) (*BulletinSender, error) {
	funcMap := template.FuncMap{
		"t": func(locale string, messageID string) string {
			return localizer.T(locale, messageID)
		},
		"dir": i18nfx.Dir,
		"pickSummary": func(story *bulletinbiz.DigestStory) string {
			if story.SummaryAI != nil && *story.SummaryAI != "" {
				return *story.SummaryAI
			}

			return story.Summary
		},
		"derefStr": func(s *string) string {
			if s == nil {
				return ""
			}

			return *s
		},
	}

	tmpl, err := template.New("bulletin_email.html.tmpl").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("parsing email template: %w", err)
	}

	return &BulletinSender{
		client:         client,
		emailResolver:  emailResolver,
		logger:         logger,
		localizer:      localizer,
		fromAddress:    config.FromAddress,
		frontendURI:    frontendURI,
		tmpl:           tmpl,
		sandboxMode:    config.SandboxMode,
		sandboxAllowed: config.SandboxAllowedEmails(),
	}, nil
}

// Kind returns the channel kind identifier.
func (s *BulletinSender) Kind() bulletinbiz.ChannelKind {
	return bulletinbiz.ChannelEmail
}

// Send delivers a digest to the recipient via email.
func (s *BulletinSender) Send(
	ctx context.Context,
	recipientProfileID string,
	digest *bulletinbiz.Digest,
) error {
	email, err := s.emailResolver.GetUserEmailByProfileID(ctx, recipientProfileID)
	if err != nil {
		s.logger.InfoContext(ctx, "Skipping bulletin email (no email found)",
			slog.String("profile_id", recipientProfileID),
			slog.String("error", err.Error()))

		return nil
	}

	// Sandbox mode: only send to explicitly allowed recipients
	if s.sandboxMode && !s.sandboxAllowed[strings.ToLower(email)] {
		s.logger.InfoContext(ctx, "Bulletin email dismissed (sandbox mode)",
			slog.String("profile_id", recipientProfileID),
			slog.String("to", email))

		return nil
	}

	htmlBody, renderErr := s.renderHTML(digest)
	if renderErr != nil {
		return fmt.Errorf("rendering email template: %w", renderErr)
	}

	subject := s.localizer.T(digest.Locale, "BulletinEmailSubject")

	sendErr := s.client.SendEmail(ctx, s.fromAddress, email, subject, htmlBody)
	if sendErr != nil {
		return fmt.Errorf("sending email: %w", sendErr)
	}

	s.logger.InfoContext(ctx, "Bulletin sent via email",
		slog.String("profile_id", recipientProfileID),
		slog.String("to", email))

	return nil
}

type templateData struct {
	FrontendURI   string
	Locale        string
	RecipientSlug string
	Groups        []*bulletinbiz.DigestGroup
}

func (s *BulletinSender) renderHTML(digest *bulletinbiz.Digest) (string, error) {
	data := templateData{
		FrontendURI:   s.frontendURI,
		Locale:        digest.Locale,
		RecipientSlug: digest.RecipientSlug,
		Groups:        digest.Groups,
	}

	var buf bytes.Buffer

	err := s.tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
