# SnippetBox

<img width="1400" height="843" alt="screen113207" src="https://github.com/user-attachments/assets/19cb72b5-fc55-4468-9230-922328ea6b36" />

<p align="center">
  <em>A web application for saving, sharing, and editing text snippets.</em>
</p>

<p align="center">
  <a href="https://snippetbox-cbxz.onrender.com/">Live Demo 🚀</a>
</p>

## About

SnippetBox is a full-stack web application built with Go. Users can create an
account, log in, and save short text snippets (such as poems or phrases),
which they can later edit or delete.

Visitors who are not logged in can browse and read snippets, but only
authenticated users can create and manage their own content.

The project is based on the book **"Let's Go" by Alex Edwards**, implementing
its concepts and architecture as a hands-on learning exercise.

## Features

- 📝 Create, read, update, and delete text snippets (CRUD)
- 👤 User registration and login
- 🔒 Session-based authentication
- 🌐 Read-only access for anonymous visitors
- 🧩 Form validation with user-friendly error messages
- ⚡ Middleware: logging, panic recovery, and secure headers

## Tech Stack

| Category | Technology |
|----------|-----------|
| Language | Go |
| Database | SQLite |
| Templates | Go `html/template` + embedded UI assets |
| Styling | CSS |
| Deployment | Render |

> **Note:** The original book project uses MySQL, but this deployment was
> adapted to run with SQLite for easier hosting on Render.

## Project Structure

```
.
├── cmd/
│   └── web/          # Application entrypoint (main.go)
├── data/             # Database files
├── internal/
│   ├── assert/       # Testing helpers
│   ├── config/       # Application configuration
│   ├── models/       # Data models (snippets, users, sessions)
│   └── validator/    # Form validation helpers
├── tls/              # TLS certificates
├── ui/
│   └── html/         # HTML templates
├── go.mod
└── LICENSE
```

## Getting Started

### Prerequisites

- Go 1.x installed on your machine

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/Zyz7/snippetbox.git
   cd snippetbox
   ```

2. Run the application:

   ```bash
   go run ./cmd/web
   ```

3. Open your browser and navigate to:

   ```
   http://localhost:4000
   ```

The application uses a local SQLite database located at `data/app.db`,
so no additional database server is required.

## Routes

| Method | Route | Description |
|--------|-------|-------------|
| GET    | `/`                  | Home — latest snippets |
| GET    | `/snippet/new`       | Create snippet form *(auth required)* |
| POST   | `/snippet/create`    | Create a new snippet *(auth required)* |
| GET    | `/snippet/view/:id`  | View a single snippet |
| GET    | `/user/signup`       | Registration form |
| POST   | `/user/signup`       | Register a new user |
| GET    | `/user/login`        | Login form |
| POST   | `/user/login`        | Log in |
| POST   | `/user/logout`       | Log out |

## Credits

This project was built following the book
[Let's Go](https://lets-go.alexedwards.net/) by Alex Edwards.

## License

This project is licensed under the terms of the [LICENSE](./LICENSE) file in
this repository.

