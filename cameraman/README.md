# The backend of the `golden-arm`

The person who operated the camera in a theater during the early days of film was the cameraman or cinematographer. The Lumière brothers, Auguste and Louis, invented the Cinématographe, a hand-cranked movie camera that could be used in a theater or carried by the cameraman.

## Developing

Stack: Go, PostgreSQL

Make sure you've [downloaded and installed](https://go.dev/doc/install) Go.

While in `cameraman`, execute `go mod tidy` to download missing modules and clean up any unnecessary dependencies.

Add the following to a `.env` file and replace "?" with your information:
```
DB_HOST="?"
DB_USER="?"
DB_PASS="?"
DB_PORT="5432"
DB_NAME="?"

API_KEY="?"
ADMIN_PASSKEY="?"

SMTP_USERNAME="?"
SMTP_PASSWORD="?"

AWS_ACCESS_KEY_ID="?"
AWS_SECRET_ACCESS_KEY="?"
AWS_REGION="?"
```

Execute `go run .` to start a local development server.