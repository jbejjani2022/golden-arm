# golden-arm

Web application for The Golden Arm, Harvard's student-run movie theater in Eliot House.

Find us on Instagram: [@eliotgoldenarm](https://www.instagram.com/eliotgoldenarm?utm_source=ig_web_button_share_sheet&igsh=ZDNlZDc0MzIxNw==)

## Structure

`/cameraman` contains the backend, written in Go and using a PostgreSQL database.

`/screen` is the frontend, written with SvelteKit and making requests to `/cameraman`.

## About

A key goal of this app is to streamline and simplify the workflow of a Golden Arm "operator"—any student who wishes to plan, schedule, and manage film screenings—while providing regular users with a friendly interface for viewing the theater calendar, finding information about future and past screenings, booking seats, submitting suggestions, and engaging with the Harvard Film Festival.

Our app features a content management system and interface for Golden Arm operators to perform admin-only functions including adding film screenings, posters, menus, and interacting with user comments, reservations, and email lists. An admin can log in at `/admin` to view and engage with the management dashboard.

The user-facing site dynamically retrieves content uploaded by admins. For instance, the landing page displays the film screening that is closest in the future (i.e. the upcoming screening). If no such film exists in the database, it displays the most recently screened film. The landing page also shows the movie calendar that contains today's date; if no such calendar exists, it displays the calendar closest in the future, and if there are no future calendars, it displays the calendar closest in the past.

## Developing

Make sure you've [downloaded and installed](https://go.dev/doc/install) Go and psql.

While in `/cameraman`, execute `go mod tidy` to download missing modules and clean up any unnecessary dependencies.

Add the following to a `.env` file in `/cameraman` and replace "?" with your information:
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

While in `/screen`, ensure you've installed dependencies with `npm install` (or `pnpm install` or `yarn`). Then, start a development server with

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

## Authors

- Joey Bejjani [@jbejjani2022](https://github.com/jbejjani2022)
- Renée Perpignan [@reneeperpignan](https://github.com/reneeperpignan)

In collaboration with
- Xander Patton
- Alexandre Benoit
- Karen Choi
