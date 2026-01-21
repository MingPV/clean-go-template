# Testing Guide

## Setup

### 1. Create Environment File

Copy the example environment file:

```bash
cp .env.example .env.dev
```

### 2. Configure Test Database

Edit `.env.dev` and update with your test database credentials:

```env
DB_HOST=localhost
DB_TEST_PORT=5432
DB_TEST_USER=postgres
DB_TEST_PASSWORD=postgres
DB_TEST_NAME=test
JWT_SECRET=test-secret-key-for-jwt-token-generation
```

**Note:** Tests use `DB_TEST_*` variables to connect to the test database, which is separate from the development database.

### 3. Start Test Database

Using Docker Compose:

```bash
docker-compose --env-file .env.dev up -d
```

Or use your existing PostgreSQL instance.

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Specific Test Package

```bash
# User repository tests
go test ./internal/user/repository/...

# User usecase tests
go test ./internal/user/usecase/...

# Order repository tests
go test ./internal/order/repository/...

# Order usecase tests
go test ./internal/order/usecase/...
```

### Run Specific Test

```bash
go test -v -run TestUserRepositoryTestSuite/TestSave ./internal/user/repository/...
```

## Test Database

The test suite automatically:
- Connects to a single test database (configured via `DB_TEST_*` variables)
- Runs migrations automatically
- Cleans up tables before and after each test to ensure isolation

Each test runs in isolation by truncating tables before and after execution, ensuring no data leaks between tests.

## CI/CD

The project includes GitHub Actions CI/CD workflow (`.github/workflows/ci.yml`) that:

1. **Runs tests** on every push and pull request
2. **Lints code** using golangci-lint
3. **Builds the application** to verify compilation
4. **Runs tests** with race detection and coverage reporting

### CI/CD Features

- Automatic PostgreSQL service setup
- Parallel job execution (test, lint, build)
- Code coverage reporting
- Race condition detection

### Viewing CI Results

1. Go to your GitHub repository
2. Click on "Actions" tab
3. View workflow runs and their results

## Troubleshooting

### Database Connection Issues

If tests fail with database connection errors:

1. Ensure PostgreSQL is running (either via Docker Compose or local instance)
2. Check `.env.dev` configuration, especially `DB_TEST_*` variables
3. Verify test database credentials are correct (`DB_TEST_NAME`, `DB_TEST_USER`, `DB_TEST_PASSWORD`, `DB_TEST_PORT`)
4. Ensure the test database exists and the user has proper permissions

### Test Database Not Cleaning Up

If test data is not being cleaned up between tests:

1. Check that `TearDownTest()` is being called (verify test output)
2. Check PostgreSQL logs for errors during table truncation
3. Ensure the test database user has permission to truncate tables
4. Manually clean tables if needed: `TRUNCATE TABLE users, orders RESTART IDENTITY CASCADE;`

### Environment Variables Not Loading

If `.env.dev` is not being loaded:

1. Ensure `.env.dev` exists in the project root
2. Check file permissions
3. Verify the path in test_helper.go matches your project structure

