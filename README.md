# gld-integration-email

## Usage

- Clone the repo:

```shell
git clone https://git.hexens.com/gld/integration/email.git
```

- Create and fill in `.env` file:

```shell
cp .env.sample .env
```

- Build and start the container:

```shell
docker compose up -d
```

---

## How to add new email

1. Add email type as constant to `/pkg/client/constants/email_types.go`.

Example:

```go
const WELCOME_EMAIL types.EmailType = "3"
```

2. Create 3 template files for email subject, text and html representations.

By convention a separate folder must be created for each email type and files must have these names:

- `subject.txt` - Subject
- `body.txt` - Text representation
- `body.html` - HTML representation

Example:

```shell
cd internal/templates
mkdir welcome
cd welcome
touch subject.txt
touch body.txt
touch body.html
```

3. If you need to insert variables into templates, do so using `{{.paramFirst}}` notation.

4. Create DTO for your email payload in `/pkg/client/dto/dto.go`.

Example:

```go
type WelcomeInputDTO struct {
	BaseSendEmailDTO
	Params WelcomeParams `json:"params" binding:"required"`
}

type WelcomeParams struct {
	Username  string `json:"username" binding:"required"`
}
```

5. Commit & Push changes made to `/pkg/client` to separate branch.
6. Update module dependancy in email microservice.

Example:

```shell
source .env
go get git.hexens.com/gld/integration/email.git/pkg/client@new-branch
```

Where `@new-branch` is name of your branch.

7. Add template filename mapping to `/internal/mappings/mappings.go`.

Example:

```go
var EmailTypeMapping = map[clientTypes.EmailType]EmailInfo{
  ...
	clientConsts.WELCOME_EMAIL: {
		HTMLTemplateFileName: buildTemplatePath("welcome/body.html"),
		TextTemplateFileName: buildTemplatePath("welcome/body.txt"),
		SubjectFileName:      buildTemplatePath("welcome/subject.txt"),
	},
  ...
}
```

8.  Restart docker container.

```shell
docker compose down
docker compose up -d
```

9. Adjust `/pkg/client/examples/main.go` and run it to test new email type.

```shell
go run pkg/client/examples/main.go
```

---

## Swagger

1. Install `swag`:

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate docs based on comments:

```shell
swag init --pd --parseInternal -g cmd/app/main.go
```

Note: This will generate files in `/docs` dir. `--pd --parseInternal` flags are needed to resolve module dependancies if they are used in swagger schema.

3. Go to http://127.0.0.1:8000/swagger/index.html

---

## Client module

Client is a shared module which contains:

- Constants
- DTOs for other microservices
- Service methods to make a call to email microservice API

---

## Shared modules

To setup a shared module:

1. Create folder inside main module dir and init a new one.

```shell
mkdir client
cd client
go mod init <module_name>
```

Example module name: `git.hexens.com/gld/integration/email.git/pkg/client`

2. Push your module to GitLab.
3. Create `.env` file in the repository you want to install shared modules and add these variables.

```
# Must be exported
export GOPRIVATE=git.hexens.com/path/to/repo/dir
export GITLAB_TOKEN=<gitlab_token>
```

GitLab token generation: https://git.hexens.com/-/user_settings/personal_access_tokens

4. Install the module.

Example:

```shell
source .env
go get git.hexens.com/gld/integration/email.git/pkg/client@initial
```

where `@initial` is the name of the branch (if you need it during development)

5. Add shortcut to local repository so you could use shared module without upgrading the version each time:

```shell
go mod edit -replace git.hexens.com/gld/integration/email.git/pkg/client=pkg/client
```

Reference:
- https://go.dev/doc/tutorial/call-module-code
