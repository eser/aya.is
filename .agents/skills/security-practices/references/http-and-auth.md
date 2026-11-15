# HTTP Services and Authentication

How the HTTP server of the Go services is hardened, and how sessions, tokens,
cookies and CORS are handled. General rules (secrets, input, injection, errors)
are in [security-rules.md](security-rules.md).

## Contents

- Layered Controls
- Client Identity Comes From the Server
- Tokens and Credentials
- Passwords
- Cookies and CSRF
- CORS
- Security Headers
- Redirects and Debug Endpoints

---

## Layered Controls

Scope: Every HTTP listener

Rule: No single check protects a service, so each layer assumes the one before
it failed. Use the existing layer instead of writing a new one:

| Layer | Control | In this repository |
| --- | --- | --- |
| Transport | Header and read timeouts, header size, connection cap | `httpfx.Config` in `pkg/ajan/httpfx/config.go`: `ReadHeaderTimeout`, `ReadTimeout`, `WriteTimeout`, `IdleTimeout`, `MaxHeaderBytes`, `MaxConnections` |
| Transport | TLS 1.2 minimum for serving | `MinVersion` in `pkg/ajan/httpfx/http_service.go` |
| Volume | Request budgets per client | `RateLimitMiddleware` (`HTTP__RATE_LIMIT_REQUESTS`) |
| Volume | Body limits | `RequestSizeLimitMiddleware` (`HTTP__MAX_REQUEST_SIZE_MB`), `http.MaxBytesReader` |
| Identity | Session or bearer token on every protected request | `AuthMiddleware` in `pkg/api/adapters/http/auth_middleware.go` |
| Permission | Per-resource authorization | security-rules.md: Authorization in the Handler |
| Browser | CORS allowlist, security headers, `Cache-Control: no-store` | `CorsMiddlewareWithCustomDomains`, `SecurityHeadersMiddleware` |
| Data | Redaction on account deletion, private export bucket | AGENTS.md: Privacy (export, hide, delete) |

Login and other credential endpoints and expensive operations take the tighter
budgets. A limiter keyed per client evicts idle keys; a map that only grows is a
memory exhaustion hole (coding-practices: Bound Every Resource). The middleware
and settings exist in `pkg/ajan`; a listener that does not apply them is
unprotected by that layer.

**Why:** a handler that forgot a check is then bounded by the other layers
instead of becoming a breach.

---

## Client Identity Comes From the Server

Scope: Anything that reads who the caller is or where it connects from

Rule: Identity and permissions come from a verified session, never from a
header or a body field the client can set (`X-User-Id`, `X-Is-Admin`, a
`user_id` in the JSON). The session is looked up on every request and must be
active; the token only names it (the JWT carries `session_id`, and the user is
derived from the session). The client address that
`ResolveAddressMiddleware` reads (`True-Client-IP`, `X-Forwarded-For`,
`X-Real-IP`) is set by whoever sends the request unless a trusted proxy
overwrites it, so it is for logging and rate-limit keys, never an input to
authorization. Hidden URLs and UI checks are presentation, not access control.

Incorrect:

```go
if ctx.Request.Header.Get("X-Is-Admin") == "true" { // any client can send this
    return adminView(ctx)
}
```

**Why:** everything a client sends can be forged; only what the server holds or
verified can be trusted.

---

## Tokens and Credentials

Scope: Session JWTs, OAuth sign-in, one-time codes and any credential the
service issues or accepts

Rule:

- Session JWTs: pin the signing method when parsing and reject any other
  (`jwt_token_service.go` and `AuthMiddleware` accept only HMAC, the Apple
  client only `RS256` through `jwt.WithValidMethods`). The signing secret is
  required, has no default and is at least as strong as a generated key
  (`AUTH__JWT_SECRET`). Tokens expire (`exp` is checked) and carry as little as
  possible: the session ID, not the user or the permissions. Use an asymmetric
  algorithm when a service other than the issuer must verify a token.
- OAuth sign-in: a random `state` per login (`auth.GenerateRandomState`), PKCE
  for the code exchange (`pkg/api/business/profiles/pkce_store.go`), and the redirect URI checked against
  the allowed origins before and after the round trip (Redirects and Debug
  Endpoints).
- One-time codes (the Telegram codes of `pkg/api/business/telegram`): from
  `crypto/rand`, short-lived and single-use (security-rules.md: Secure
  Randomness).
- A new credential type follows the same design: random secret, expiry, one
  refusal for every cause so an answer does not tell which part was wrong.

Correct:

```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("%w (method=%s)", ErrInvalidSigningMethod, token.Method.Alg())
    }

    return jwtSecretBytes, nil
})
```

Incorrect:

```go
token, _ := jwt.Parse(tokenString, func(*jwt.Token) (any, error) { return key, nil }) // accepts whatever method the token names
```

**Why:** an unpinned algorithm lets an attacker pick one the server did not
intend, and a refusal that differs by cause tells an attacker which part of a
guess was right.

---

## Passwords

Scope: Any feature that stores a password

Rule: Users sign in through OAuth providers and the services store no
passwords. If a feature ever needs one: store only an Argon2id hash with a
random salt and the parameters encoded with it, compare in constant time, spend
the same work for an unknown account, answer every failed login with one error,
and bound the length of the input. Fast hashes (SHA-256, MD5) are never
password storage.

**Why:** a stolen hash database is attacked offline; only a slow, memory-hard
hash keeps that expensive, and equal work for unknown accounts hides which
accounts exist.

---

## Cookies and CSRF

Scope: Session, sign-in state and CSRF handling

Rule: Session cookies are `HttpOnly`, `Secure`, with a bounded lifetime and the
narrowest `SameSite` and `Path` that work (`SetSessionCookie` in
`pkg/api/adapters/http/cookie_helpers.go`). `Secure` follows the
`SecureCookie` setting, which production sets to true. The session cookie uses
`SameSite=None` on purpose, so that a profile's custom domain can use the
single sign-on; a cookie with `SameSite=None` is sent on cross-site requests,
so it does not protect against CSRF by itself. A state-changing endpoint
therefore must not be reachable by a plain cross-site form post: accept JSON
bodies only (a content type that forces a CORS preflight), or require the
`Authorization` header, and check the `Origin` against the allowed origins
(CORS). A cookie that scripts must read (the theme cookie) is the exception and
holds nothing but a preference, with the reason in a comment. Logout ends the
session on the server and clears the cookie (`ClearSessionCookie`).

Correct:

```go
http.SetCookie(w, &http.Cookie{
    Name:     config.CookieName,
    Value:    sessionID,
    Path:     "/",
    Expires:  expiresAt,
    Secure:   config.SecureCookie,
    HttpOnly: true,
    SameSite: http.SameSiteNoneMode, // cross-domain sign-in on custom domains
})
```

Incorrect:

```go
http.SetCookie(w, &http.Cookie{Name: "session", Value: sessionID}) // readable by scripts, sent over HTTP
```

**Why:** `HttpOnly` keeps an XSS bug from reading the session, `Secure` keeps it
off plain HTTP, and the rest limits where the browser attaches it.

---

## CORS

Scope: Cross-origin browser access to the API

Rule: Allow exact origins from a list (the configured origins plus the custom
domains held in the database), never `*` together with credentials and never an
origin reflected without a check. Set `Access-Control-Allow-Origin`,
`Access-Control-Allow-Credentials`, the allowed headers and methods only for an
allowed origin, and vary the response on `Origin`, `Cookie` and
`Authorization` so a cache does not serve one caller's answer to another
(`CorsMiddlewareWithCustomDomains` does this and sets `Cache-Control:
no-store`). CORS governs what a browser lets a page read; it does not
authenticate anything (Authorization in the Handler).

Correct:

```go
allowed := slices.Contains(allowedOrigins, requestOrigin) // exact match, not a prefix or substring
```

Incorrect:

```go
headers.Set("Access-Control-Allow-Origin", "*")
headers.Set("Access-Control-Allow-Credentials", "true") // browsers refuse it, and it signals the wrong intent
```

**Why:** a loose origin check lets any site act with the visitor's cookies.

---

## Security Headers

Scope: Pages and API answers that browsers read

Rule: `SecurityHeadersMiddleware` in `pkg/ajan/httpfx/middlewares` sets
`X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY` and
`Strict-Transport-Security` (one year, `includeSubDomains`) on every response;
register it on every listener. API responses that can differ per caller are
`Cache-Control: no-store`. A change that needs a weaker header (framing, an
inline script) is redesigned or decided on purpose, with the reason in a
comment.

---

## Redirects and Debug Endpoints

Scope: Redirect targets and diagnostic handlers

Rule:

- Redirect only to relative paths or to an exact origin on an allowlist. Parse
  the target with `url.Parse` and compare scheme and host exactly, as
  `validateRedirectURI` in `pkg/api/business/auth/service.go` does
  (`ErrUnsafeRedirectURI`); a prefix or substring check lets
  `https://aya.is.evil.test` through.
- The pprof handlers are mounted only when `profiling` is enabled (default
  false) and are protected by `PPROF_TOKEN`; without the token they are open and
  the service logs a warning, so production either leaves profiling off or sets
  the token. A token in a query string lands in logs and browser history; prefer
  a header.
- Diagnostics go to logs through `logfx`, and nothing mounts a debug handler
  that is reachable without a decision (security-rules.md: Safe Defaults and
  Explicit Outcomes).

Incorrect:

```go
if strings.HasPrefix(redirectURI, "https://aya.is") { // matches https://aya.is.evil.test
```

**Why:** an open redirect turns a trusted domain into a phishing link, and a
profiling endpoint exposes heap contents.
