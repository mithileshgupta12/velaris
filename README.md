# Velaris

A Trello-inspired task management application built with Go and Vue.

## Tech Stack

- **Backend**: Go
- **Frontend**: Vue.js / TypeScript
- **Database**: PostgreSQL
- **ORM**: XORM

## Prerequisites

- Go (1.25+)
- Node.js (22+)
- PostgreSQL
- Make

## Getting Started

### 1. Environment Setup

```bash
cp .env.example .env
```

### 2. Build & Run

**Build the backend:**

```bash
make build
```

**Run the API:**

```bash
make run
```

**In a new terminal, run the frontend:**

```bash
cd frontend
pnpm run dev
```

The API will be available at `http://localhost:3000` and frontend at `http://localhost:5173`.

## License

This project is licensed under the GNU Affero General Public License v3 (AGPL v3) - see the [LICENSE](LICENSE) file for details.

For a quick summary: you can use, modify, and distribute this software freely, but any modifications made to a publicly accessible version must also be made available under the same license.
