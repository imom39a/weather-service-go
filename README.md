# weather-service-go

This Weather service project hosts an http server that uses the Open Weather API that exposes an endpoint that takes in lat/long coordinates and return what the weather condition is outside in that area (snow, rain, etc), whether it’s hot, cold, or moderate outside.

## Project Structure

```
/weather-service-go
├── cmd                     # Source code for application's executables
│   └── api                 # API server directory
│       ├── main.go         # Main file for the API server
├── e2e                     # End-to-end tests directory
│   └── app_test.go         # API server end-to-end tests
├── internal
│   ├── api                 # API specific code, OpenAPI specifications
│   │   ├── server.go       # Generated server initialization and routing
│   │   ├── types.go        # Generated types and models from the OpenAPI specification
│   │   └── handler         # Handlers implementing the server logic
│   │       ├── handler.go  # Handlers for each API endpoint
│   ├── service             # Business logic layer
│   ├── entity              # Domain entities. In this case, the Weather entity from the Open Weather API
├── spec                    # OpenAPI Specifications for the project
├── tools                   # Tools and utilities
├── Makefile                # Makefile for common tasks
├── go.mod                  # Go module file
├── go.sum                  # Go sum file
└── README.md               # Project README
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- Make

### Running the application

1. Clone the repository

```bash
git clone
```

2. Change into the project directory

```bash
cd weather-service-go
```

3. Run the application

```bash
make run-api
```

The API server will start on port 8080.

4. To Run tests

```bash
make test
```

### OpenAPI Specifications

- The OpenAPI specifications for the project can be found in the `spec` directory.
- The server code is generated from the OpenAPI specifications using `oapi-codegen`.
- To generate the server code, run the following command:

```bash
make generate-api
```

### Features

- [x] Get weather by lat/long

### API Endpoints

- `GET /weather?lat={latitude}&long={longitude}` - Get weather by lat/long
