# Project Renaming Guide

When you clone this project, you'll need to rename it in several places to ensure everything works correctly.

## Quick Start (Recommended)

The easiest way to rename the project is to use the automated script:

```bash
./rename-project.sh
```

The script will:
1. Ask for your new module path (e.g., `github.com/your-username/your-project-name`)
2. Ask for your project name for API title (optional)
3. Automatically update all files
4. Regenerate proto files and Swagger docs (if tools are installed)
5. Run `go mod tidy`

**That's it!** The script handles everything for you.

---

## Manual Steps (Alternative)

If you prefer to rename manually or the script doesn't work for you, follow these steps:

### 1. Change Module Name in `go.mod`

Open `go.mod` and change the module name from:
```go
module github.com/MingPV/clean-go-template
```

To your new name, for example:
```go
module github.com/your-username/your-project-name
```

### 2. Update All Import Statements

After changing the module name in `go.mod`, you need to update all import statements in the project. You can do this automatically or manually:

**Automatic (Recommended):**
```bash
go mod tidy
```

**Manual - Using VS Code:**
- Press `Cmd+Shift+H` (Mac) or `Ctrl+Shift+H` (Windows/Linux)
- Find: `github.com/MingPV/clean-go-template`
- Replace with: `github.com/your-username/your-project-name`
- Select "Replace All"

**Manual - Using Command Line:**
```bash
# Using sed (Mac/Linux)
find . -type f -name "*.go" -exec sed -i '' 's|github.com/MingPV/clean-go-template|github.com/your-username/your-project-name|g' {} +

# Or using grep + sed
grep -rl "github.com/MingPV/clean-go-template" . --include="*.go" | xargs sed -i '' 's|github.com/MingPV/clean-go-template|github.com/your-username/your-project-name|g'
```

### 3. Update Proto File

Open `proto/order/order.proto` and change the `go_package` option:

```protobuf
option go_package = "github.com/your-username/your-project-name/proto/orderpb";
```

After that, you need to regenerate the proto files:

```bash
# Make sure you have protoc and protoc-gen-go-grpc installed
# Install if you don't have them:
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Regenerate proto files
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/order/order.proto
```

### 4. Update Swagger Documentation

Open `docs/v1/docs.go` and change the Title and Description:

```go
Title:            "Your Project Name API",
Description:      "This is the backend API for Your Project Name project.",
```

Or if you want to regenerate Swagger docs:

```bash
# Install swag if you don't have it
go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
swag init -g cmd/app/main.go -o docs/v1
```

### 5. Update README.md

Edit `README.md`:
- Change the project name in the title
- Update the URL in the git clone command
- Update the project name in other relevant sections

### 6. Update Docker Compose (Optional)

If you want to change container names in `docker-compose.yaml`:

```yaml
container_name: your-project-name-postgres
container_name: your-project-name-postgres-test
```

### 7. Run go mod tidy

After making all changes, run:

```bash
go mod tidy
```

This will update dependencies and verify that all import statements are correct.

### 8. Test the Project

Run tests to make sure everything works:

```bash
# Build the project
go build ./cmd/app

# Run tests
go test ./...

# Run the project
go run ./cmd/app
```

## Summary of Files to Update

1. ✅ `go.mod` - Change module name
2. ✅ `proto/order/order.proto` - Change go_package option
3. ✅ `docs/v1/docs.go` - Change Title and Description (or regenerate)
4. ✅ `README.md` - Update project name and URLs
5. ✅ All `.go` files - Import statements (auto-updated with go mod tidy or find/replace)
6. ⚠️ `docker-compose.yaml` - (Optional: if you want to change container names)

## Using the Rename Script

The `rename-project.sh` script is located in the project root. It's an interactive script that guides you through the renaming process.

### Prerequisites

- Bash shell (available on macOS and Linux by default)
- Go installed
- (Optional) `protoc` for regenerating proto files
- (Optional) `swag` for regenerating Swagger docs

### Usage

1. Make sure the script is executable:
   ```bash
   chmod +x rename-project.sh
   ```

2. Run the script:
   ```bash
   ./rename-project.sh
   ```

3. Follow the prompts:
   - Enter your new module path (e.g., `github.com/your-username/your-project-name`)
   - Enter your project name for API title (or press Enter to use the module name)

4. Review the changes and test:
   ```bash
   git diff
   go test ./...
   go build ./cmd/app
   ```

The script automatically:
- Updates `go.mod`
- Updates all import statements in `.go` files
- Updates proto files
- Updates Swagger documentation
- Updates README.md
- Runs `go mod tidy`
- Attempts to regenerate proto files (if `protoc` is installed)
- Attempts to regenerate Swagger docs (if `swag` is installed)

### Troubleshooting

If the script fails or you need to regenerate files manually:

1. **Proto files not regenerated:**
   ```bash
   # Install protoc and plugins
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   
   # Regenerate
   protoc --go_out=. --go_opt=paths=source_relative \
          --go-grpc_out=. --go-grpc_opt=paths=source_relative \
          proto/order/order.proto
   ```

2. **Swagger docs not regenerated:**
   ```bash
   # Install swag
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # Regenerate
   swag init -g cmd/app/main.go -o docs/v1
   ```

## Important Notes

- After renaming, you must regenerate proto files using `protoc`
- After renaming, you may need to regenerate swagger docs using `swag init`
- Make sure your Git remote URL is correct if you're pushing to a new repository
- Test the project thoroughly after renaming to ensure everything works correctly
