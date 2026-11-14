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
)

// exportReadyTemplate keeps the HTML on long lines on purpose: inline e-mail markup does not survive wrapping.
//
//nolint:lll
const exportReadyTemplate = `<!doctype html>
<html lang="{{.Locale}}" dir="{{.Dir}}">
<body style="margin:0;padding:24px;background:#f6f6f4;font-family:-apple-system,Segoe UI,Helvetica,Arial,sans-serif;color:#1c1c1a;">
<table role="presentation" width="100%" style="max-width:520px;margin:0 auto;background:#ffffff;border-radius:8px;padding:32px;">
<tr><td>
<h1 style="margin:0 0 16px;font-size:20px;">{{.Heading}}</h1>
<p style="margin:0 0 24px;font-size:15px;line-height:1.5;">{{.Body}}</p>
<p style="margin:0 0 24px;"><a href="{{.URL}}" style="display:inline-block;padding:10px 18px;background:#1c1c1a;color:#ffffff;text-decoration:none;border-radius:6px;font-size:14px;">{{.Button}}</a></p>
<p style="margin:0;font-size:12px;line-height:1.5;color:#6b6b66;">{{.Footer}}</p>
</td></tr>
</table>
</body>
</html>`

// PrivacyMailer sends privacy-related notices (a personal data archive being ready) through Resend.
type PrivacyMailer struct {
	client         *Client
	logger         *logfx.Logger
	localizer      *i18nfx.Localizer
	tmpl           *template.Template
	sandboxAllowed map[string]bool
	fromAddress    string
	sandboxMode    bool
}

// NewPrivacyMailer creates the mailer. It honours the same sandbox settings as the bulletin sender, so
// non-production environments never e-mail real users.
func NewPrivacyMailer(
	client *Client,
	logger *logfx.Logger,
	localizer *i18nfx.Localizer,
	config *Config,
) (*PrivacyMailer, error) {
	tmpl, err := template.New("privacy_export_ready").Parse(exportReadyTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing privacy e-mail template: %w", err)
	}

	return &PrivacyMailer{
		client:         client,
		logger:         logger,
		localizer:      localizer,
		tmpl:           tmpl,
		sandboxAllowed: config.SandboxAllowedEmails(),
		fromAddress:    config.FromAddress,
		sandboxMode:    config.SandboxMode,
	}, nil
}

// SendExportReady tells a user that their personal data archive can be downloaded.
func (m *PrivacyMailer) SendExportReady(
	ctx context.Context,
	recipient string,
	locale string,
	settingsURL string,
	retentionDays int,
) error {
	if m.sandboxMode && !m.sandboxAllowed[strings.ToLower(recipient)] {
		m.logger.InfoContext(
			ctx,
			"sandbox mode: skipping export-ready e-mail",
			slog.String("locale", locale),
		)

		return nil
	}

	locale = strings.TrimSpace(locale)

	data := struct {
		Locale  string
		Dir     string
		Heading string
		Body    string
		Button  string
		Footer  string
		URL     string
	}{
		Locale:  locale,
		Dir:     i18nfx.Dir(locale),
		Heading: m.localizer.T(locale, "PrivacyExportReadyHeading"),
		Body: m.localizer.TWithData(locale, "PrivacyExportReadyBody",
			map[string]any{"Days": retentionDays}),
		Button: m.localizer.T(locale, "PrivacyExportReadyButton"),
		Footer: m.localizer.T(locale, "PrivacyExportReadyFooter"),
		URL:    settingsURL,
	}

	var body bytes.Buffer

	err := m.tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("rendering privacy e-mail: %w", err)
	}

	return m.client.SendEmail(
		ctx,
		m.fromAddress,
		recipient,
		m.localizer.T(locale, "PrivacyExportReadySubject"),
		body.String(),
	)
}
