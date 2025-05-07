# Movies Service

**`movies_service`** is a simple, modular Go backend for managing movies with user authentication. It uses:

* **Uber Fx** for dependency injection
* **Gin** for HTTP routing
* **GORM** with **PostgreSQL** for ORM
* **JWT** for authentication/authorization
* **sql-migrate** for database migrations
* **Swaggo** for auto-generated Swagger docs
* **Docker** for containerized deployment

---

## Features

* User registration and login (JWT-based)
* Secure CRUD endpoints for movies:

  * Create a movie: `POST /movies`
  * List all movies: `GET /movies`
  * Retrieve a movie: `GET /movies/:id`
  * Update a movie: `PUT /movies/:id`
  * Delete a movie: `DELETE /movies/:id`
* Input validation and consistent error responses
* Swagger UI at `/docs` for interactive API documentation (http://localhost:8080/docs/index.html)
* Automatic SQL migrations on container startup
* Unit tests for service and handler layers
* Coverage reporting via Go `cover` tool

---

## Project Structure

```text
movies_service/
├── cmd/
│   └── movies-service/      # Composition root (main.go)
├── config/                  # Env-driven configuration loader
├── model/                   # Domain models (Movie, User, Responses)
├── repository/              # GORM-based data access
├── service/                 # Business logic (user + movie services)
├── auth/                    # JWT generation and middleware
├── handlers/                # Gin handlers (controllers)
├── migrations/              # SQL migration files (sql-migrate)
├── docs/                    # Swagger spec and docs.go
├── Dockerfile               # Multi-stage build + migrations
├── docker-entrypoint.sh     # Run migrations, start server
├── dbconfig.yml             # sql-migrate configuration
├── Makefile                 # Build, run, migrate, test, coverage
├── go.mod / go.sum          # Go modules
└── README.md                # This file
```

---

## Prerequisites

* Go 1.23+
* PostgreSQL (local or remote)
* `sql-migrate` CLI (for local migrations):

  ```bash
  go install github.com/rubenv/sql-migrate/...@latest
  ```
* `swag` CLI (for docs):

  ```bash
  go install github.com/swaggo/swag/cmd/swag@latest
  ```
* Docker (or Colima) for containerization

---

## Configuration

Copy `.env.sample` to `.env` and update values:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=movies_user
DB_PASSWORD=yourpassword
DB_NAME=movies_db
JWT_SECRET=supersecretkey
PORT=8080
```

Environment variables are loaded by `config.NewConfig()` at startup.

---

## Local Development

1. **Create database & user** (if not already):

   ```sql
   CREATE ROLE movies_user WITH LOGIN PASSWORD 'yourpassword';
   CREATE DATABASE movies_db OWNER movies_user;
   ```
2. **Set env vars**:

   ```bash
   source .env
   ```
3. **Run migrations**:

   ```bash
   make migrate-up
   ```
4. **Generate Swagger docs**:

   ```bash
   make swagger
   ```
5. **Run tests + coverage**:

   ```bash
   make test
   make coverage
   open coverage.html
   ```
6. **Start the service**:

   ```bash
   make run
   ```

Visit `http://localhost:8080/docs` for the Swagger UI.

---

## Docker Deployment

1. **Build the image**:

   ```bash
   make docker-build
   ```
2. **Run the container**:

   ```bash
   docker run --rm -p 8080:8080 \
     -e DB_HOST=host.docker.internal \
     -e DB_PORT=5432 \
     -e DB_USER=$DB_USER \
     -e DB_PASSWORD=$DB_PASSWORD \
     -e DB_NAME=$DB_NAME \
     -e JWT_SECRET=$JWT_SECRET \
     movies_service
   ```

By default, the entrypoint script will apply migrations before launching the server.

---

## API Usage Examples

### Register & Login

```bash
# Register
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret"}'

# Login
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"secret"}'
# → {"token":"<JWT_TOKEN>"}
```

### CRUD Movies

```bash
TOKEN=<JWT_TOKEN>
# Create
curl -X POST http://localhost:8080/movies \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Inception","director":"Nolan","year":2010,"plot":"Dream heist"}'

# List
curl http://localhost:8080/movies -H "Authorization: Bearer $TOKEN"

# Get by ID
curl http://localhost:8080/movies/1 -H "Authorization: Bearer $TOKEN"

# Update
curl -X PUT http://localhost:8080/movies/1 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"Inception (2010)"}'

# Delete
curl -X DELETE http://localhost:8080/movies/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

## Contributing

Pull requests are welcome! Please open an issue first to discuss any significant changes.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
