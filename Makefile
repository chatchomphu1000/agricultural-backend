.PHONY: build clean deploy

# Build the Go binary for Lambda
build:
	@echo "Building Lambda function..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o lambda ./cmd/lambda/main.go

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -f lambda

# Deploy to AWS using SAM
deploy: build
	@echo "Deploying to AWS..."
	sam deploy --guided

# Deploy without prompts (after first deploy)
deploy-fast: build
	@echo "Fast deploying to AWS..."
	sam deploy

# Local testing with SAM
local: build
	@echo "Starting local API..."
	sam local start-api

# Validate SAM template
validate:
	@echo "Validating template..."
	sam validate

# Build and package
package: build
	@echo "Packaging application..."
	sam package --s3-bucket $(S3_BUCKET) --output-template-file packaged.yaml
