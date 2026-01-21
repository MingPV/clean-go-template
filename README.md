# clean-go-template
<div><img width="600" alt="image" src="https://github.com/user-attachments/assets/5ff920c7-eccf-4fa2-8198-3cf2ec2dae6e" /></div>

**go-clean-template** is a clean and scalable starter template for building backend applications in Go, following Clean Architecture principles. This template uses:

- **Fiber v2** as a fast and lightweight web framework for building RESTful APIs 
- **GORM** as the ORM for PostgreSQL database access
- **gRPC** for high-performance RPC communication 
- **Docker Compose** for easy setup of PostgreSQL services

## Features

- Clear separation of concerns with Clean Architecture  
- High-performance HTTP handling with Fiber v2  
- Robust database integration using GORM with PostgreSQL  
- REST and gRPC APIs supported 
- Data Transfer Objects (DTO) to manage data structure transformations between layers  
- Swagger API documentation with automatic generation 
- Ready-to-use Docker Compose setup for dependencies  

## Getting Started

Follow the steps below to set up and run the project:

### Rename the Project (Optional but Recommended)

If you want to rename this project to your own project name, use the automated script:

```bash
./rename-project.sh
```

The script will guide you through renaming the project. For more details, see [docs/RENAME_PROJECT.md](docs/RENAME_PROJECT.md).

---

1. Clone the repository:

    ```bash
    git clone https://github.com/MingPV/clean-go-template.git
    cd clean-go-template
    ```

2. Install Go module dependencies:

    ```bash
    go mod tidy
    ```

3. Create and configure the environment file:

    ```bash
    cp .env.example .env.dev
    ```

    Open the `.env.dev` file and fill in all required configuration values:
    - **Development database** credentials: `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_PORT`
    - **Test database** credentials: `DB_TEST_NAME`, `DB_TEST_USER`, `DB_TEST_PASSWORD`, `DB_TEST_PORT`
    - Application settings: `APP_PORT`, `GRPC_PORT`, `JWT_SECRET`, etc.

    **Note:** This project uses a single `.env.dev` file for both development and testing environments.

4. Start PostgreSQL services using Docker Compose:

    ```bash
    # Start both development and test databases
    docker-compose --env-file .env.dev up -d
    
    # Or start only development database
    docker-compose --env-file .env.dev up -d postgres
    
    # Or start only test database
    docker-compose --env-file .env.dev up -d postgres-test
    ```

    The docker-compose file includes two PostgreSQL services:
    - `postgres`: Development database (uses `DB_NAME`, `DB_USER`, `DB_PASSWORD`, `DB_PORT`)
    - `postgres-test`: Test database (uses `DB_TEST_NAME`, `DB_TEST_USER`, `DB_TEST_PASSWORD`, `DB_TEST_PORT`)

5. Run the application:

    ```bash
    go run ./cmd/app
    ```

    The application will start on port `8000` by default (configurable via `APP_PORT` in `.env.dev`).

6. Run tests:

    ```bash
    # Run all tests
    go test ./...
    
    # Run tests with verbose output
    go test -v ./...
    
    # Run tests with coverage
    go test -v -coverprofile=coverage.out ./...
    ```

    **Note:** Tests use the same `.env.dev` file. Make sure to configure test database credentials (`DB_TEST_*`) in your `.env.dev` file.

Swagger UI for the API documentation is available at: `http://localhost:8000/api/v1/docs`

<img width="0" alt="image" src="https://github.com/user-attachments/assets/e38ff0e8-8fd1-4d39-baca-af30b85b353a" />
<img width="700" alt="image" src="https://github.com/user-attachments/assets/840f8d43-e07c-44a8-9b7d-3f4d62d912ce" />

## Environment Variables

The project uses a single `.env.dev` file for configuration. Key environment variables include:

### Application Settings
- `APP_PORT`: HTTP server port (default: `8000`)
- `GRPC_PORT`: gRPC server port (default: `50052`)
- `APP_ENV`: Application environment (default: `development`)
- `JWT_SECRET`: Secret key for JWT token generation
- `JWT_EXPIRATION`: JWT token expiration time in seconds (default: `3600`)

### Development Database
- `DB_HOST`: Database host (default: `localhost`)
- `DB_PORT`: Database port (default: `5432`)
- `DB_USER`: Database user (default: `postgres`)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

### Test Database
- `DB_TEST_NAME`: Test database name
- `DB_TEST_USER`: Test database user
- `DB_TEST_PASSWORD`: Test database password
- `DB_TEST_PORT`: Test database port

See `.env.example` for a complete list of available environment variables.

## Testing

For detailed testing information, see [docs/TESTING.md](docs/TESTING.md).

The test suite uses the same `.env.dev` file and automatically:
- Connects to test database using `DB_TEST_*` environment variables
- Runs database migrations automatically
- Cleans up test data before and after each test to ensure isolation

**Important:** Make sure to configure `DB_TEST_NAME`, `DB_TEST_USER`, `DB_TEST_PASSWORD`, and `DB_TEST_PORT` in your `.env.dev` file for tests to work properly.

## Project structure


```bash
/clean-go-template
├── cmd/
│   └── app/
│       └── main.go               
├── docs/
│   └── v1/                 
├── internal/               
│   ├── app/            
│   ├── entities/
│   ├── order/
│   │   ├── handler/
│   │   │   ├── grpc/
│   │   │   └── rest/
│   │   ├── usecase/
│   │   ├── repository/
│   │   └── dto/ 
│   └── user/               
├── pkg/
│   ├── config/
│   ├── database/
│   ├── middleware/
│   ├── responses/
│   └── routes/
├── proto/
│   └── order/
├── utils/                
├── .env.example             
├── .gitignore               
├── LICENSE                  
├── README.md             
├── docker-compose.yaml      
└── go.mod
```




