# Blog Management System

## Overview

A robust and scalable RESTful API for managing blog content, built with Go and PostgreSQL. This system provides comprehensive functionality for managing posts, categories, tags, user authentication, and media uploads.

## Postman Collection

You can find the Postman collection for this API in the `https://documenter.getpostman.com/view/10995686/2sB2x9jqfj`.

## Technology Stack

- **Language**: Go (Golang)
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Media Storage**: Cloudinary
- **Configuration**: Viper

## Design Patterns

### Repository Pattern

The application implements the Repository Pattern to abstract data access logic from business logic. Each entity (Post, Category, Tag, User) has its own repository that handles database operations.

### Service Layer Pattern

Business logic is encapsulated in service layers that sit between controllers (handlers) and repositories. This promotes separation of concerns and makes the codebase more maintainable.

### Dependency Injection

The application uses manual dependency injection to provide dependencies to components. This approach enhances testability and modularity.

### Middleware Pattern

Gin middleware is used for cross-cutting concerns like authentication, authorization, and request logging.

### Model-View-Controller (MVC)

The application follows a modified MVC pattern where:

- Models represent data structures and database entities
- Controllers (handlers) handle HTTP requests and responses
- Views are represented as JSON responses

## Lifecycle Framework

### Request Lifecycle

1. **HTTP Request**: Client sends a request to the server
2. **Middleware Processing**: Request passes through middleware chain (logging, authentication, etc.)
3. **Route Matching**: Gin router matches the request to the appropriate handler
4. **Handler Processing**: Handler processes the request, interacts with services
5. **Service Logic**: Service layer applies business logic
6. **Repository Access**: Repositories interact with the database
7. **Response Generation**: Handler generates and returns the response
8. **Middleware Post-Processing**: Response passes through middleware chain
9. **HTTP Response**: Response is sent back to the client

### Application Lifecycle

1. **Configuration Loading**: Load environment variables and configuration files
2. **Database Connection**: Establish connection to PostgreSQL
3. **Migration**: Run database migrations
4. **Dependency Setup**: Initialize repositories, services, and handlers
5. **Router Setup**: Configure routes and middleware
6. **Server Start**: Start the HTTP server
7. **Request Handling**: Process incoming requests
8. **Graceful Shutdown**: Handle shutdown signals and close resources

## Project Structure

```
├── cmd/
│   └── server/           # Application entry point
├── configs/              # Configuration files and loading logic
├── internal/             # Internal application code
│   ├── database/         # Database connection and migration
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── models/           # Data models and DTOs
│   ├── repositories/     # Data access layer
│   ├── services/         # Business logic layer
│   │   └── cloudinary/   # Cloudinary integration
│   ├── setup/            # Application setup and initialization
│   └── utils/            # Utility functions
└── pkg/                  # Public packages
```

## Core Features

### Authentication

- User registration and login
- JWT-based authentication
- Token refresh mechanism
- Role-based authorization

### Post Management

- Create, read, update, delete posts
- Post publishing workflow
- Post categorization
- Post tagging

### Category Management

- Create, read, update, delete categories
- Associate posts with categories

### Tag Management

- Create, read, update, delete tags
- Associate posts with multiple tags
- Filter posts by tag

### Media Management

- Upload user avatars to Cloudinary
- Secure media storage and retrieval

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and get access token
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout and invalidate token

### Categories

- `GET /api/categories` - List all categories
- `GET /api/categories/:id` - Get category by ID
- `GET /api/categories/slug/:slug` - Get category by slug
- `POST /api/admin/categories` - Create a new category (admin only)
- `PUT /api/admin/categories/:id` - Update a category (admin only)
- `DELETE /api/admin/categories/:id` - Delete a category (admin only)

### Posts

- `GET /api/posts` - List all published posts
- `GET /api/posts/:id` - Get post by ID
- `GET /api/posts/slug/:slug` - Get post by slug
- `POST /api/admin/posts` - Create a new post (admin only)
- `PUT /api/admin/posts/:id` - Update a post (admin only)
- `DELETE /api/admin/posts/:id` - Delete a post (admin only)
- `PUT /api/admin/posts/:id/publish` - Publish a post (admin only)

### Tags

- `GET /api/tags` - List all tags
- `GET /api/tags/:id` - Get tag by ID
- `GET /api/tags/slug/:slug` - Get tag by slug
- `GET /api/tags/:id/posts` - Get posts by tag ID
- `GET /api/posts/:id/tags` - Get tags by post ID
- `POST /api/admin/tags` - Create a new tag (admin only)
- `PUT /api/admin/tags/:id` - Update a tag (admin only)
- `DELETE /api/admin/tags/:id` - Delete a tag (admin only)

### User Profile

- `POST /api/profile/avatar` - Upload user avatar

### Admin Dashboard

- `GET /api/admin/dashboard` - Get admin dashboard data

## Setup and Installation

### Prerequisites

- Go 1.16+
- PostgreSQL 12+
- Cloudinary account

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Server Configuration
SERVER_PORT=8080
SERVER_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=blog_db
DB_SSLMODE=disable

# JWT Configuration
JWT_ACCESS_SECRET=your_access_secret
JWT_REFRESH_SECRET=your_refresh_secret
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=7d

# Cloudinary Configuration
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
CLOUDINARY_FOLDER=avatars
```

### Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/blog-management.git
cd blog-management
```

2. Install dependencies

```bash
go mod download
```

3. Run the application

```bash
go run cmd/server/main.go
```

Or use the provided Makefile:

```bash
make run
```

## Testing

Run the tests with:

```bash
go test ./...
```

Or use the provided Makefile:

```bash
make test
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
