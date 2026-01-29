# Project Banking Platform

Mini Banking Platform - handle financial transactions, proper accounting principles, user-friendly interface.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Migrate the database (requires Goose to be installed)
_Note: Maybe you shoud change the database host in the .env file to localhost_
```bash
make migrate-up
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

## API Documentation (Swagger)

The project includes Swagger/OpenAPI documentation that is automatically generated from code annotations.

### Accessing Swagger UI

After starting the server, you can access the Swagger UI at:
```
http://localhost:<PORT>/swagger/index.html
```

Replace `<PORT>` with the port number configured in your `.env` file (default is usually `8080`).

### Generating Documentation

The Swagger documentation is generated using `swag`. To regenerate the documentation after making changes to API annotations:

1. Install swag (if not already installed):
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

2. Generate documentation:
```bash
swag init -g ./cmd/api/main.go -o ./docs
```

This will update the files in the `docs/` directory:
- `docs.go` - Generated Go code
- `swagger.json` - JSON specification
- `swagger.yaml` - YAML specification