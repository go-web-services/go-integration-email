# go-integration-email

HTTP integration service that sends transactional email (auth flows, notifications) on behalf of other microservices. Callers use a single `POST /api/v1/send` endpoint with an **email type** string, **recipients**, and a **params** map; the service picks HTML/text/subject templates, renders them with those params, and delivers via SMTP.

**Scope:** template-backed outbound mail only (no inbox, no webhooks). **Consumers:** backend services that import the `pkg/client` module and call the API or use the thin client wrapper.

## API examples (JSON)

### Send email (generic payload)

`emailType` values are string constants such as `AuthForgotPassword`, `AuthEmailConfirm`, and `AuthOTPSignin` (see `pkg/client/constants`).

**Request** `POST /api/v1/send`

```json
{
  "emailType": "AuthForgotPassword",
  "recipients": ["user@example.com"],
  "params": {
    "forgotPasswordLink": "https://app.example.com/reset?token=abc",
    "expirationTimeInMinutes": 30
  }
}
```

**Response** `200 OK`

```json
{
  "message": "Email sent successfully"
}
```

Typed client DTOs wrap the same shape with explicit `params` structs, for example forgot-password:

```json
{
  "emailType": "AuthForgotPassword",
  "recipients": ["user@example.com"],
  "params": {
    "forgotPasswordLink": "https://app.example.com/reset?token=abc",
    "expirationTimeInMinutes": 30
  }
}
```

Email confirmation:

```json
{
  "emailType": "AuthEmailConfirm",
  "recipients": ["user@example.com"],
  "params": {
    "emailConfirmLink": "https://app.example.com/confirm?token=xyz",
    "expirationTimeInMinutes": 60
  }
}
```

OTP sign-in:

```json
{
  "emailType": "AuthOTPSignin",
  "recipients": ["user@example.com"],
  "params": {
    "otpCode": "482915",
    "expirationTimeInMinutes": 10
  }
}
```

Validation and platform errors follow the shared `go-web-platform` error format (for example `VALIDATION_ERROR` with `errors[].field` using **JSON tag** names such as `emailType`, not Go struct field names).

## Run locally

- Clone: `git clone git@github.com:go-web-services/go-integration-email.git`
- Copy env: `cp .env.sample .env` and set SMTP-related variables (`EMAIL_SERVER`, `EMAIL_PORT`, `EMAIL_USERNAME`, `EMAIL_PASSWORD`, `EMAIL_FROM`, etc.).
- With Docker (`debug/Dockerfile`, bind-mounted sources): `docker compose up -d`
- Application port defaults from `APP_PORT` (see `.env.sample`).

## Client module (`pkg/client`)

Other services depend on `github.com/go-web-services/go-integration-email/pkg/client` for DTOs, constants, and `EmailAPIService` (HTTP client using platform `SendRequest`). For local development, use a replace in the consuming `go.mod`:

```bash
go mod edit -replace github.com/go-web-services/go-integration-email/pkg/client=/path/to/go-integration-email/pkg/client
```

Private modules: set `GOPRIVATE=github.com/go-web-services/*` when pulling from GitHub.

## Adding a new email type

1. Add a constant in `pkg/client/constants/email_constants.go`.
2. Add templates under `internal/templates/<name>/`: `subject.txt`, `body.txt`, `body.html`.
3. Add DTOs in `pkg/client/dto` if the payload is typed; map the type in `internal/mappings/email_mapping.go`.
4. Regenerate Swagger from the repo root: `gocheck -d` (or `swag init` with your project’s usual flags).

## Swagger

After code or comment changes, regenerate docs (e.g. `gocheck -d`). Open `http://127.0.0.1:<APP_PORT>/swagger/index.html` when the service is running.

## Author

[Lomank](https://lomank.com)
