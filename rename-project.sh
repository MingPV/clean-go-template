#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Old project name
OLD_MODULE="github.com/MingPV/clean-go-template"
OLD_API_TITLE="CleanGO API"
OLD_API_DESC="CleanGO project"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Project Renaming Script${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Get new module name
read -p "Enter your new module path (e.g., github.com/username/project-name): " NEW_MODULE

if [ -z "$NEW_MODULE" ]; then
    echo -e "${RED}Error: Module path cannot be empty!${NC}"
    exit 1
fi

# Get new project name for API title (optional)
read -p "Enter your project name for API title (press Enter to use module name): " NEW_PROJECT_NAME

if [ -z "$NEW_PROJECT_NAME" ]; then
    # Extract project name from module path (last part after /)
    NEW_PROJECT_NAME=$(echo "$NEW_MODULE" | awk -F'/' '{print $NF}')
fi

NEW_API_TITLE="${NEW_PROJECT_NAME} API"
NEW_API_DESC="${NEW_PROJECT_NAME} project"

echo ""
echo -e "${YELLOW}Summary:${NC}"
echo -e "  Old module: ${OLD_MODULE}"
echo -e "  New module: ${NEW_MODULE}"
echo -e "  New API title: ${NEW_API_TITLE}"
echo ""
read -p "Continue? (y/n): " CONFIRM

if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo -e "${YELLOW}Cancelled.${NC}"
    exit 0
fi

echo ""
echo -e "${BLUE}Starting renaming process...${NC}"
echo ""

# Detect OS for sed command
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    SED_IN_PLACE="sed -i ''"
else
    # Linux
    SED_IN_PLACE="sed -i"
fi

# Step 1: Update go.mod
echo -e "${GREEN}[1/8]${NC} Updating go.mod..."
$SED_IN_PLACE "s|$OLD_MODULE|$NEW_MODULE|g" go.mod
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✓${NC} go.mod updated"
else
    echo -e "  ${RED}✗${NC} Failed to update go.mod"
    exit 1
fi

# Step 2: Update all .go files
echo -e "${GREEN}[2/8]${NC} Updating import statements in .go files..."
find . -type f -name "*.go" -not -path "./vendor/*" -exec $SED_IN_PLACE "s|$OLD_MODULE|$NEW_MODULE|g" {} +
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✓${NC} All .go files updated"
else
    echo -e "  ${RED}✗${NC} Failed to update .go files"
    exit 1
fi

# Step 3: Update proto file
echo -e "${GREEN}[3/8]${NC} Updating proto/order/order.proto..."
if [ -f "proto/order/order.proto" ]; then
    $SED_IN_PLACE "s|$OLD_MODULE|$NEW_MODULE|g" proto/order/order.proto
    if [ $? -eq 0 ]; then
        echo -e "  ${GREEN}✓${NC} proto file updated"
    else
        echo -e "  ${RED}✗${NC} Failed to update proto file"
        exit 1
    fi
else
    echo -e "  ${YELLOW}⚠${NC}  proto/order/order.proto not found, skipping..."
fi

# Step 4: Update docs.go
echo -e "${GREEN}[4/8]${NC} Updating Swagger documentation..."
if [ -f "docs/v1/docs.go" ]; then
    $SED_IN_PLACE "s|$OLD_API_TITLE|$NEW_API_TITLE|g" docs/v1/docs.go
    $SED_IN_PLACE "s|$OLD_API_DESC|$NEW_API_DESC|g" docs/v1/docs.go
    if [ $? -eq 0 ]; then
        echo -e "  ${GREEN}✓${NC} Swagger docs updated"
    else
        echo -e "  ${RED}✗${NC} Failed to update Swagger docs"
        exit 1
    fi
else
    echo -e "  ${YELLOW}⚠${NC}  docs/v1/docs.go not found, skipping..."
fi

# Step 5: Update README.md
echo -e "${GREEN}[5/8]${NC} Updating README.md..."
if [ -f "README.md" ]; then
    # Update project name in title (first line)
    $SED_IN_PLACE "1s|.*|# $NEW_PROJECT_NAME|" README.md
    
    # Update git clone URL
    $SED_IN_PLACE "s|git clone https://github.com/MingPV/clean-go-template.git|git clone https://github.com/$(echo $NEW_MODULE | cut -d'/' -f2-3 | tr '/' '/').git|g" README.md
    
    # Update directory name in README
    $SED_IN_PLACE "s|cd clean-go-template|cd $NEW_PROJECT_NAME|g" README.md
    
    # Update project structure path
    $SED_IN_PLACE "s|/clean-go-template|/$NEW_PROJECT_NAME|g" README.md
    
    if [ $? -eq 0 ]; then
        echo -e "  ${GREEN}✓${NC} README.md updated"
    else
        echo -e "  ${YELLOW}⚠${NC}  Some parts of README.md may not have been updated"
    fi
else
    echo -e "  ${YELLOW}⚠${NC}  README.md not found, skipping..."
fi

# Step 6: Run go mod tidy
echo -e "${GREEN}[6/8]${NC} Running go mod tidy..."
go mod tidy
if [ $? -eq 0 ]; then
    echo -e "  ${GREEN}✓${NC} go mod tidy completed"
else
    echo -e "  ${RED}✗${NC} go mod tidy failed"
    exit 1
fi

# Step 7: Regenerate proto files
echo -e "${GREEN}[7/8]${NC} Regenerating proto files..."
if command -v protoc &> /dev/null; then
    # Check if protoc-gen-go is installed
    if ! command -v protoc-gen-go &> /dev/null; then
        echo -e "  ${YELLOW}⚠${NC}  protoc-gen-go not found, installing..."
        go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    fi
    
    # Check if protoc-gen-go-grpc is installed
    if ! command -v protoc-gen-go-grpc &> /dev/null; then
        echo -e "  ${YELLOW}⚠${NC}  protoc-gen-go-grpc not found, installing..."
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    fi
    
    if [ -f "proto/order/order.proto" ]; then
        protoc --go_out=. --go_opt=paths=source_relative \
               --go-grpc_out=. --go-grpc_opt=paths=source_relative \
               proto/order/order.proto 2>/dev/null
        
        if [ $? -eq 0 ]; then
            echo -e "  ${GREEN}✓${NC} Proto files regenerated"
        else
            echo -e "  ${YELLOW}⚠${NC}  Failed to regenerate proto files (protoc may not be installed)"
        fi
    fi
else
    echo -e "  ${YELLOW}⚠${NC}  protoc not found, skipping proto regeneration"
    echo -e "  ${YELLOW}   ${NC}  Install protoc and run: protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/order/order.proto"
fi

# Step 8: Regenerate Swagger docs
echo -e "${GREEN}[8/8]${NC} Regenerating Swagger documentation..."
if command -v swag &> /dev/null; then
    swag init -g cmd/app/main.go -o docs/v1 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "  ${GREEN}✓${NC} Swagger docs regenerated"
    else
        echo -e "  ${YELLOW}⚠${NC}  Failed to regenerate Swagger docs"
    fi
else
    echo -e "  ${YELLOW}⚠${NC}  swag not found, skipping Swagger regeneration"
    echo -e "  ${YELLOW}   ${NC}  Install swag: go install github.com/swaggo/swag/cmd/swag@latest"
    echo -e "  ${YELLOW}   ${NC}  Then run: swag init -g cmd/app/main.go -o docs/v1"
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  Project Renaming Completed!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo -e "  1. Review the changes in git: ${YELLOW}git diff${NC}"
echo -e "  2. Test the project: ${YELLOW}go test ./...${NC}"
echo -e "  3. Build the project: ${YELLOW}go build ./cmd/app${NC}"
echo -e "  4. Update your Git remote URL if needed"
echo ""
echo -e "${YELLOW}Note:${NC} If proto files or Swagger docs weren't regenerated,"
echo -e "      install the required tools and run the commands shown above."
echo ""

