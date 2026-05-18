# Book Service API

A RESTful book management service built with Go, Echo framework, PostgreSQL, and JWT authentication using hexagonal architecture with Viper configuration management.

## Features

- **Hexagonal Architecture**: Clean separation of concerns with domain, application, and infrastructure layers
- **Viper Configuration**: Centralized YAML configuration with environment variable override support
- **JWT Authentication**: Secure token-based authentication for protected endpoints
- **PostgreSQL Database**: Persistent storage with GORM ORM
- **Request Validation**: Automatic validation using go-playground/validator
- **CRUD Operations**: Complete book management with pagination and filtering

## Endpoints

### Public Endpoints
- `GET /ping` - Health check
- `POST /echo` - Echo request body
- `POST /auth/token` - Generate JWT token

### Book Endpoints
- `POST /books` - Create a new book
- `GET /books` - Get all books (protected with JWT)
- `GET /books?author=X` - Filter books by author
- `GET /books?page=1&limit=2` - Get paginated books
- `GET /books/:id` - Get book by ID
- `PUT /books/:id` - Update book (protected with JWT)
- `DELETE /books/:id` - Delete book (protected with JWT)

## Project Structure

```
decent/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Viper configuration management
│   ├── domain/
│   │   ├── entity/
│   │   │   └── book.go          # Book entity
│   │   ├── repository/
│   │     │   └── book.go          # Repository interface
│   │   └── service/
│   │       └── book.go          # Service interface
│   ├── application/
│   │   ├── dto/
│   │   │   ├── book.go          # Book DTOs
│   │   │   └── auth.go          # Auth DTOs
│   │   └── usecase/
│   │       └── book.go          # Book use cases
│   └── infrastructure/
│       ├── database/
│       │   ├── postgres.go       # Database connection
│       │   └── book.go          # Book repository implementation
│       ├── auth/
│       │   └── jwt.go           # JWT service
│       └── http/
│           ├── handler/
│           │   ├── book.go      # Book handlers
│           │   └── auth.go      # Auth handlers
│           ├── middleware/
│           │   ├── auth.go      # Auth middleware
│           │   └── validation.go # Validation middleware
│           └── router.go        # Router setup
├── config.yaml                   # Application configuration
├── config.example.yaml           # Example configuration template
├── go.mod                        # Go module file
└── go.sum                        # Go module checksums
```

## Configuration

The application uses Viper for configuration management. You can configure the service using:
1. `config.yaml` file (primary configuration)
2. Environment variables (override YAML values)

### Configuration Priority (highest to lowest)
1. Environment variables with `BOOK_` prefix
2. YAML configuration file
3. Default values

### YAML Configuration

Create `config.yaml` with the following structure:

```yaml
database:
  host: localhost
  port: "5432"
  user: postgres
  password: password
  dbname: bookdb

server:
  port: "8080"

jwt:
  secret: your-secret-key-here
```

### Environment Variables

You can override any configuration value using environment variables with the `BOOK_` prefix:

```bash
export BOOK_DATABASE_HOST=localhost
export BOOK_DATABASE_PORT=5432
export BOOK_DATABASE_USER=postgres
export BOOK_DATABASE_PASSWORD=password
export BOOK_DATABASE_DBNAME=bookdb
export BOOK_SERVER_PORT=8080
export BOOK_JWT_SECRET=your-secret-key
```

### Configuration Options

| Path | Default | Description |
|------|---------|-------------|
| `database.host` | `localhost` | Database host |
| `database.port` | `5432` | Database port |
| `database.user` | `postgres` | Database user |
| `database.password` | `password` | Database password |
| `database.dbname` | `bookdb` | Database name |
| `server.port` | `8080` | Server port |
| `jwt.secret` | auto-generated | JWT secret key |

## Setup

### Prerequisites
- Go 1.25+ (for local development)
- PostgreSQL 12+ (for local development)
- Docker & Docker Compose (for containerized deployment)
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd decent
```

2. Install dependencies:
```bash
go mod download
```

3. Create configuration file:
```bash
cp config.example.yaml config.yaml
```

4. Edit `config.yaml` with your settings:
```yaml
database:
  host: localhost
  port: "5432"
  user: postgres
  password: your_password
  dbname: bookdb

server:
  port: "8080"

jwt:
  secret: your-secret-key-here
```

5. Create PostgreSQL database:
```sql
CREATE DATABASE bookdb;
```

6. Build and run:
```bash
go build -o server ./cmd/server
./server
```

## Usage

### Authentication

Get JWT token:
```bash
curl -X POST http://localhost:8080/auth/token \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=user&password=pass"
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Book Operations

Create a book:
```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "description": "A classic American novel",
    "published_year": 1925
  }'
```

Get all books (protected):
```bash
curl -X GET http://localhost:8080/books \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Get book by ID:
```bash
curl -X GET http://localhost:8080/books/1
```

Filter by author:
```bash
curl -X GET "http://localhost:8080/books?author=F. Scott Fitzgerald"
```

Get paginated books:
```bash
curl -X GET "http://localhost:8080/books?page=1&limit=10"
```

Update book (protected):
```bash
curl -X PUT http://localhost:8080/books/1 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Great Gatsby (Updated)",
    "description": "An updated description"
  }'
```

Delete book (protected):
```bash
curl -X DELETE http://localhost:8080/books/1 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Book Model

```json
{
  "id": 1,
  "title": "The Great Gatsby",
  "author": "F. Scott Fitzgerald",
  "description": "A classic American novel",
  "published_year": 1925,
  "created_at": "2026-05-18T13:21:00Z",
  "updated_at": "2026-05-18T13:21:00Z"
}
```

### Validation Rules
- `title`: required, max 200 characters
- `author`: required, max 100 characters
- `description`: max 1000 characters
- `published_year`: between 1000 and 9999

## Development

### Run tests:
```bash
go test ./...
```

### Run with live reload:
```bash
go run ./cmd/server
```

### Build for production:
```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server
```

### Override config with environment variables:
```bash
BOOK_DATABASE_PASSWORD=mypass BOOK_SERVER_PORT=3000 ./server
```

## Docker Deployment

### Quick Start with Docker Compose

The easiest way to run the application is using Docker Compose, which will automatically:
- Build the Docker image
- Start PostgreSQL database
- Configure environment variables
- Set up network connectivity

```bash
# Clone the repository
git clone <repository-url>
cd decent

# Start all services (app + database)
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

The application will be available at `http://localhost:8080`

### Docker Compose Configuration

The `docker-compose.yml` file includes:

**Services:**
- **app**: The book service application
- **postgres**: PostgreSQL database

**Environment Variables:**
- Database credentials
- Server configuration
- JWT secret

**Volumes:**
- PostgreSQL data persistence

### Detailed Docker Compose Usage

#### Start the Application:

```bash
# Start in background (detached mode)
docker-compose up -d

# Start with logs in foreground
docker-compose up

# Start specific services only
docker-compose up -d app
```

#### View Logs:

```bash
# All services
docker-compose logs -f

# App service only
docker-compose logs -f app

# Last 100 lines
docker-compose logs --tail=100 app
```

#### Check Service Status:

```bash
# All services
docker-compose ps

# App service only
docker-compose ps app
```

#### Stop and Remove Services:

```bash
# Stop services but keep data
docker-compose stop

# Stop and remove containers
docker-compose down

# Stop and remove containers + volumes (data will be lost!)
docker-compose down -v
```

#### Restart Services:

```bash
# Restart all services
docker-compose restart

# Restart specific service
docker-compose restart app
```

#### Access Database:

```bash
# Connect to PostgreSQL database
docker-compose exec postgres psql -U postgres -d bookdb

# Execute SQL commands
docker-compose exec postgres psql -U postgres -d bookdb -c "SELECT * FROM books;"

# Backup database
docker-compose exec postgres pg_dump -U postgres bookdb > backup.sql

# Restore database
cat backup.sql | docker-compose exec -T postgres psql -U postgres bookdb
```

#### Update the Application:

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build

# Or rebuild specific service
docker-compose up -d --build app
```

#### Production Docker Compose:

For production deployments, create a `docker-compose.prod.yml` file:

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - BOOK_DATABASE_HOST=postgres
      - BOOK_DATABASE_PORT=5432
      - BOOK_DATABASE_USER=postgres
      - BOOK_DATABASE_PASSWORD=postgres
      - BOOK_DATABASE_DBNAME=bookdb
      - BOOK_SERVER_PORT=8080
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=bookdb
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:

networks:
  default:
    name: book-service-network
```

Run with:
```bash
export POSTGRES_PASSWORD=your_secure_password
docker-compose -f docker-compose.prod.yml up -d
```

### Manual Docker Usage

#### Build Docker Image:

```bash
# Build with tag
docker build -t decent-api:latest .

# Build with custom tag
docker build -t decent-api:v1.0.0 .

# Build with build args
docker build --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) -t decent-api:latest .
```

#### Run Docker Container:

```bash
# Basic run
docker run -d --name book-service decent-api

# With port mapping
docker run -d \
  --name book-service \
  -p 8080:8080 \
  decent-api

# With configuration file
docker run -d \
  --name book-service \
  -p 8080:8080 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  decent-api

# With all environment variables
docker run -d \
  --name book-service \
  -p 8080:8080 \
  -e BOOK_DATABASE_HOST=postgres \
  -e BOOK_DATABASE_PORT=5432 \
  -e BOOK_DATABASE_USER=postgres \
  -e BOOK_DATABASE_PASSWORD=postgres \
  -e BOOK_DATABASE_DBNAME=bookdb \
  -e BOOK_SERVER_PORT=8080 \
  decent-api

# With network access
docker run -d \
  --name book-service \
  --network book-service-network \
  decent-api
```

#### Docker Container Management:

```bash
# View running containers
docker ps

# View all containers (including stopped)
docker ps -a

# View container logs
docker logs book-service

# View logs with follow
docker logs -f book-service

# Enter container shell
docker exec -it book-service sh

# View container processes
docker top book-service

# Stop container
docker stop book-service

# Start container
docker start book-service

# Restart container
docker restart book-service

# Remove container
docker rm book-service

# Remove container and volumes
docker rm -v book-service
```

#### Docker Images Management:

```bash
# List images
docker images

# Remove image
docker rmi decent-api:latest

# Remove all unused images
docker image prune

# Remove all images (DANGEROUS)
docker rmi $(docker images -q)
```

#### Docker Networks:

```bash
# List networks
docker network ls

# Create network
docker network create book-service-network

# Connect container to network
docker network connect book-service-network book-service

# Inspect network
docker network inspect book-service-network
```

### Troubleshooting

#### Common Issues:

1. **Port already in use**:
```bash
# Find process using port 8080
lsof -i :8080
# Or Windows:
netstat -ano | findstr :8080

# Kill the process or use a different port
docker run -p 8081:8080 decent-api
```

2. **Database connection failed**:
```bash
# Check if database is ready
docker-compose ps postgres

# Wait for database to be ready
docker-compose exec app sleep 5

# Check database logs
docker-compose logs postgres
```

3. **Volume permissions issues**:
```bash
# Fix permissions on Linux/Mac
sudo chown -R $USER:$USER ./data

# Or use a named volume (recommended)
```

4. **Container won't start**:
```bash
# Check logs
docker-compose logs app

# Check container status
docker-compose ps

# Rebuild container
docker-compose up -d --build
```

5. **Memory issues**:
```bash
# Check memory usage
docker stats

# Limit memory usage
docker run -d --memory="512m" --memory-swap="1g" decent-api
```

#### Debug Mode:

```bash
# Run container with debug output
docker run -it --rm --entrypoint sh decent-api

# Inside container, run the app manually
/app/server
```

#### Performance Optimization:

```bash
# Use build cache
docker build --cache-from decent-api:latest -t decent-api:latest .

# Use multi-stage build (already configured in Dockerfile)

# Use Docker layer caching

# Monitor performance
docker stats
```

## Technologies Used

- **Echo**: High-performance web framework
- **Viper**: Configuration management
- **GORM**: ORM library for database operations
- **PostgreSQL**: Robust relational database
- **golang-jwt**: JWT token generation and validation
- **go-playground/validator**: Struct validation

## License

MIT License