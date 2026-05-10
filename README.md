# go-integration-email

`github.com/go-web-services/go-integration-email`

HTTP integration service that sends transactional email on behalf of other microservices. Callers submit an email type, a list of recipients, and a params map to a single endpoint; the service selects the matching HTML/text/subject templates, renders them with the provided params, and delivers via SMTP. Scope is outbound transactional mail only — no inbox, no webhooks.

---

## Responsibilities

- Receive send requests from backend services via `POST /api/v1/send`.
- Resolve the correct template set (HTML body, plain-text body, subject) for each email type.
- Render templates with caller-supplied params and deliver via SMTP.
- Expose a typed `pkg/client` module so consuming services import DTOs and an HTTP client wrapper instead of building raw requests.

---

## Configuration

| Variable | Purpose | Default |
|----------|---------|---------|
| `APP_PORT` | HTTP listen port | — |
| `APP_ENV` | Environment (`dev` / `prod`) | — |
| `EMAIL_SERVER` | SMTP host | — |
| `EMAIL_PORT` | SMTP port | `587` |
| `EMAIL_USERNAME` | SMTP auth username | — |
| `EMAIL_PASSWORD` | SMTP auth password | — |
| `EMAIL_FROM` | Sender address shown to recipients | — |

---

## Run locally

```bash
git clone git@github.com:go-web-services/go-integration-email.git
cd go-integration-email
cp .env.sample .env
# Set EMAIL_SERVER, EMAIL_USERNAME, EMAIL_PASSWORD, EMAIL_FROM
docker compose up -d
```

---

## Docker

- **Dev** (hot reload, sources bind-mounted via `debug/Dockerfile`):
  ```bash
  docker compose up -d
  ```
- **Prod**:
  ```bash
  docker compose -f docker-compose-prod.yml up --build
  ```

---

## API surface

### Send email

`POST /api/v1/send`

**Request**

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

**Response `200 OK`**

```json
{
  "message": "Email sent successfully"
}
```

Built-in email type constants (defined in `pkg/client/constants`):

| Constant | Description |
|----------|-------------|
| `AuthForgotPassword` | Password-reset link |
| `AuthEmailConfirm` | Email activation link |
| `AuthOTPSignin` | One-time password for sign-in |

Swagger UI is available at `/swagger` (dev environment only). Regenerate after changes:

```bash
gocheck -d
# or: swag init -g cmd/app/main.go -o docs --parseDependency --parseInternal
```

Validation errors follow the go-web-platform format: `error_code: VALIDATION_ERROR` with `errors[].field` names matching JSON tags (`emailType`, not `EmailType`).

---

## Client module (`pkg/client`)

Other services import `github.com/go-web-services/go-integration-email/pkg/client` for typed DTOs and the `EmailAPIService` HTTP client.

```go
import (
    clientapi "github.com/go-web-services/go-integration-email/pkg/client/service"
    "github.com/go-web-services/go-integration-email/pkg/client/dto"
)

svc := clientapi.NewEmailAPIService("http://localhost:8010")
err := svc.SendForgotPasswordV1(ctx, dto.ForgotPasswordEmailDTO{
    Recipients: []string{"user@example.com"},
    Params: dto.ForgotPasswordParams{
        ForgotPasswordLink:      "https://app.example.com/reset?token=abc",
        ExpirationTimeInMinutes: 30,
    },
})
```

For local development in a consuming service:

```bash
go mod edit -replace github.com/go-web-services/go-integration-email=/path/to/go-integration-email
```

---

## Adding a new email type

1. Add a constant in `internal/constants/constants.go` (and mirror it in `pkg/client/constants/`).
2. Create `internal/templates/<email-type>/` with `subject.txt`, `body.txt`, and `body.html`.
3. Add input DTOs in `pkg/client/dto/` and map the type in `internal/mappings/email_mapping.go`.
4. Regenerate Swagger: `gocheck -d`.

---

## Private dependencies

```bash
export GOPRIVATE='github.com/go-web-services/*'
```

---

## Author

[Lomank](https://lomank.com)
