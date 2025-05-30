# go-clean-template
<div><img width="600" alt="image" src="https://github.com/user-attachments/assets/5ff920c7-eccf-4fa2-8198-3cf2ec2dae6e" /></div>

**go-clean-template** is a clean and scalable starter template for building backend applications in Go, following Clean Architecture principles. This template uses:

- **Fiber v2** as a fast and lightweight web framework for building RESTful APIs 
- **GORM** as the ORM for PostgreSQL database access
- **Redis** for caching to improve performance and reduce database load
- **gRPC** for high-performance RPC communication 
- **Docker Compose** for easy setup of PostgreSQL and Redis services

## Features

- Clear separation of concerns with Clean Architecture  
- High-performance HTTP handling with Fiber v2  
- Robust database integration using GORM with PostgreSQL  
- Caching support via Redis
- REST and gRPC APIs supported 
- Data Transfer Objects (DTO) to manage data structure transformations between layers  
- Swagger API documentation with automatic generation 
- Ready-to-use Docker Compose setup for dependencies  

## Getting Started

Follow the steps below to set up and run the project:

1. Clone the repository:

    ```bash
    git clone https://github.com/MingPV/clean-go-template.git
    cd clean-go-template
    ```

2. Install Go module dependencies:

    ```bash
    go mod tidy
    ```

3. Rename the environment file and configure it:

    ```bash
    cp .env.example .env
    ```

    Open the `.env` file and fill in all required configuration values such as PostgreSQL credentials, Redis connection details, and any other environment-specific settings.

4. Start PostgreSQL and Redis services using Docker Compose:

    ```bash
    docker-compose up -d
    ```

5. Run the application:

    ```bash
    go run ./cmd/app
    ```

6. Test:

    ```bash
    go test ./pkg/routes
    ```

Swagger UI for the API documentation is available at: localhost:8080/api/v1/docs

<img width="0" alt="image" src="https://github.com/user-attachments/assets/e38ff0e8-8fd1-4d39-baca-af30b85b353a" />
<img width="700" alt="image" src="https://github.com/user-attachments/assets/840f8d43-e07c-44a8-9b7d-3f4d62d912ce" />


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
│   ├── redisclient/
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




