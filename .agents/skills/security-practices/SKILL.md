---
name: security-practices
description: "Security rules for aya.is in the Go services and the webclient: secrets, content-free logs, input bounds, authorization, injection, SSRF guards, sanitized errors, sessions, cookies, CORS, uploads, crypto, production settings, model output and the MCP server. Use when handling secrets, input, auth, sessions, cookies, crypto, files, uploads, outbound URLs, AI prompts or a security review."
---

# Security Practices

Security rules for the Go services (`apps/services`) and the webclient
(`apps/webclient`). Privacy behavior (export, hide, delete) is in AGENTS.md:
Privacy (export, hide, delete).

## Always

- Secrets arrive through the environment (a `conf` field with no default,
  platform secret variables, git-ignored `.env.local`), never from code or
  committed configuration; a `VITE_` variable is public
- Secrets, session ids, content and personal data never reach logs, audit
  records, errors or telemetry; log IDs, codes and counts
- Validate every input where it enters, after bounding it; parse unknown
  input into a typed value (Go struct, zod schema)
- Every protected route runs behind `AuthMiddleware` and authorizes the
  resource for the acting user; deny by default, the UI is not access control
- No SQL, command, path or markup built from input strings; SQL goes through
  sqlc, paths through `os.OpenRoot`, user text never into
  `dangerouslySetInnerHTML` unescaped
- URLs from users or settings pass `lib.ValidateExternalURL`; redirects are
  refused or re-checked
- Answers carry a code or a fixed message (`httpfx.WithSanitizedError`); Go
  error text stays in the log
- Secrets, tokens and one-time codes come from `crypto/rand`; JWTs pin the
  signing method; session cookies are `HttpOnly` and `Secure`
- TLS is never relaxed: no `InsecureSkipVerify`, 1.2 minimum
- Settings that weaken security default to the safe value and are named
- Model output is untrusted input, escaped like user content; MCP tools return
  only what an anonymous caller may read

## References

| File                                                  | Read when                                                                                                |
| ----------------------------------------------------- | -------------------------------------------------------------------------------------------------------- |
| [security-rules.md](references/security-rules.md)     | Code that handles secrets, logs, input, authorization, injection, outbound URLs, errors, production, models |
| [http-and-auth.md](references/http-and-auth.md)       | HTTP listeners, limits, sessions, tokens, cookies, CSRF, CORS, headers, redirects                        |
| [crypto-and-files.md](references/crypto-and-files.md) | Hashing, encryption, TLS, uploads, files, decoding untrusted data                                        |
| [review-checklist.md](references/review-checklist.md) | Reviewing a change for security                                                                          |
