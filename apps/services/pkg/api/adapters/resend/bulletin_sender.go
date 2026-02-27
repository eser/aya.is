package resend

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"strings"
	"time"

	"github.com/eser/aya.is/services/pkg/ajan/i18nfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	bulletinbiz "github.com/eser/aya.is/services/pkg/api/business/bulletin"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

// localeToMessageID maps a locale code to a go-i18n message ID (e.g. "en" → "LocaleEn").
var localeToMessageID = map[string]string{ //nolint:gochecknoglobals
	"tr":    "LocaleTr",
	"en":    "LocaleEn",
	"fr":    "LocaleFr",
	"de":    "LocaleDe",
	"es":    "LocaleEs",
	"pt-PT": "LocalePtPT",
	"it":    "LocaleIt",
	"nl":    "LocaleNl",
	"ja":    "LocaleJa",
	"ko":    "LocaleKo",
	"ru":    "LocaleRu",
	"zh-CN": "LocaleZhCN",
	"ar":    "LocaleAr",
}

// kindEmoji maps story kinds to Unicode symbols for email rendering.
var kindEmoji = map[string]string{ //nolint:gochecknoglobals
	"news":         "📰",
	"article":      "✏️",
	"announcement": "📢",
	"status":       "ℹ️",
	"content":      "🖼️",
	"presentation": "📊",
	"activity":     "📅",
}

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
func NewBulletinSender( //nolint:funlen
	client *Client,
	emailResolver bulletinbiz.UserEmailResolver,
	logger *logfx.Logger,
	localizer *i18nfx.Localizer,
	config *Config,
	frontendURI string,
	templatePath string,
) (*BulletinSender, error) {
	md := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)

	funcMap := template.FuncMap{
		"t": func(locale string, messageID string) string {
			return localizer.T(locale, messageID)
		},
		"dir": i18nfx.Dir,
		"mdToHTML": func(text string) template.HTML {
			var buf bytes.Buffer

			err := md.Convert([]byte(text), &buf)
			if err != nil {
				return template.HTML(template.HTMLEscapeString(text)) //nolint:gosec
			}

			return template.HTML(buf.String()) //nolint:gosec
		},
		"localeBadge": func(emailLocale string, storyLocale string) string {
			trimmed := strings.TrimRight(storyLocale, " ")
			if trimmed == emailLocale {
				return ""
			}

			msgID, ok := localeToMessageID[trimmed]
			if !ok {
				return ""
			}

			return localizer.T(emailLocale, msgID)
		},
		"kindIcon": func(kind string) string {
			if emoji, ok := kindEmoji[kind]; ok {
				return emoji
			}

			return ""
		},
		"formatDate": func(locale string, t *time.Time) string {
			if t == nil {
				return ""
			}

			month := localizer.T(locale, fmt.Sprintf("MonthShort%d", t.Month()))

			// English: "Dec 22, 2025" — all others: "22 Ara 2025"
			if locale == "en" {
				return fmt.Sprintf("%s %d, %d", month, t.Day(), t.Year())
			}

			return fmt.Sprintf("%d %s %d", t.Day(), month, t.Year())
		},
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
