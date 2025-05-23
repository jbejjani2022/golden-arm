# golden-arm

Web application for The Golden Arm, Harvard's student-run movie theater in Eliot House.

Find us on Instagram: [@eliotgoldenarm](https://www.instagram.com/eliotgoldenarm?utm_source=ig_web_button_share_sheet&igsh=ZDNlZDc0MzIxNw==)

## Structure

`/cameraman` contains the backend, written in Go and using a PostgreSQL database.

`/screen` is the frontend, written with SvelteKit and making requests to `/cameraman`.

## About

A key goal of this app is to streamline and simplify the workflow of a Golden Arm "operator"—any student who wishes to plan, schedule, and manage film screenings—while providing regular users with a friendly interface for viewing the theater calendar, booking seats, finding information about future and past screenings, submitting suggestions, and ordering movie posters and merch.

![The Golden Arm past screenings](assets/UI/past_screenings_carousel.png)

![The Golden Arm reservation page](assets/UI/reservation_page.png)

Our app features a content management system and interface for Golden Arm operators to perform admin-only functions including adding film screenings, posters, menus, merch items, and interacting with user comments, reservations, shop orders, and email lists. An admin can log in at `/admin` to view and engage with the management dashboard.

![The Golden Arm admin dashboard](assets/UI/admin_dashboard.png)

The user-facing site dynamically retrieves content uploaded by admins. For instance, the landing page displays the film screening that is closest in the future (i.e. the upcoming screening). If no such film exists in the database, it displays the most recently screened film. The landing page also shows the movie calendar that contains today's date; if no such calendar exists, it displays the calendar closest in the future, and if there are no future calendars, it displays the calendar closest in the past.

![The Golden Arm calendar](assets/UI/calendar.png)

In the admin dashboard, an admin can add a movie that will be screened. This involves uploading a movie poster and concessions menu as image files. These image files are stored in an AWS S3 bucket and made visible via a public URL. In the admin merch dashboard, they can view and manage merch shop items, inventory, and orders.

## Developing

Make sure you've [downloaded and installed](https://go.dev/doc/install) Go and psql.

While in `/cameraman`, execute `go mod tidy` to download missing modules and clean up any unnecessary dependencies.

Add the following to a `.env` file in `/cameraman` and replace "?" with your information:
```
DB_HOST="?"
DB_USER="?"
DB_PASS="?"
DB_PORT="?"  # default: "5432"
DB_NAME="?"

API_KEY="?"
ADMIN_PASSKEY="?"

RESERVATIONS_SENDER="?"  # address from which reservation confirmation emails are sent
ORDERS_SENDER="?"  # address from which order confirmation emails are sent
REPLYTO="?" # monitored admin inbox

AWS_ACCESS_KEY_ID="?"
AWS_SECRET_ACCESS_KEY="?"
AWS_REGION="?"

S3_BUCKET_NAME="?"
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
- [Xander Patton](https://xanderdraven.weebly.com)
- Alexandre Benoit
- Karen Choi
