
# === General Variables ===


BACKEND_DIR=backend
FRONTEND_DIR=frontend
INFRA_DIR=infra
DOCS_DIR=$(BACKEND_DIR)/docs
LAMBDA_ZIP=$(BACKEND_DIR)/bin/backend.zip

# === Targets ===

.PHONY: all build build-backend build-frontend swagger deploy update-lambda update-site destroy clean

# Default target (alias for build)
all: build

# === Build Section ===

# Full build (backend + frontend)
build: build-backend build-frontend

# Build Go backend
build-backend: swagger
	@echo "Compiling Go backend..."
	@cd $(BACKEND_DIR) && GOOS=linux GOARCH=amd64 go build -o bin/bootstrap cmd/api/main.go
	@echo "Packaging Lambda deployment artifact..."
	@cd $(BACKEND_DIR)/bin && zip -q -j backend.zip bootstrap

# Build Vue frontend
build-frontend:
	@echo "Building Vue frontend..."
	@cd $(FRONTEND_DIR) && npm install
	@cd $(FRONTEND_DIR) && npm run build

# === Swagger ===

swagger:
	@echo "Generating Swagger documentation..."
	@cd $(BACKEND_DIR) && swag init --generalInfo cmd/api/main.go --output docs || echo "Skipping docs generation (swag not installed or docs dir missing)"

# === Deployment ===

# Full deploy
deploy: build
	@echo "Deploying infrastructure with Terraform..."
	@cd $(INFRA_DIR) && terraform init && terraform apply -auto-approve
	@$(MAKE) update-site

# Lambda update
update-lambda: build-backend
	@echo "Updating Lambda function code only..."
	aws lambda update-function-code \
		--function-name $$(terraform -chdir=$(INFRA_DIR) output -raw lambda_function_name) \
		--zip-file fileb://$(LAMBDA_ZIP) \
		--no-cli-pager

# Site update
update-site: build-frontend
	@echo "Uploading frontend to S3..."
	aws s3 sync frontend/dist/ s3://$$(terraform -chdir=$(INFRA_DIR) output -raw frontend_bucket_name) --delete
	@echo "Invalidating CloudFront cache..."
	aws cloudfront create-invalidation \
		--distribution-id $$(terraform -chdir=$(INFRA_DIR) output -raw cloudfront_distribution_id) \
		--paths "/*" \
		--no-cli-pager

# === Destroy ===

destroy:
	@echo "Destroying all Terraform-managed infrastructure..."
	@cd $(INFRA_DIR) && terraform destroy -auto-approve

# === Clean ===

clean:
	@echo "Cleaning compiled artifacts..."
	@rm -rf $(BACKEND_DIR)/bin/*
	@rm -rf $(DOCS_DIR)/*
	@rm -rf $(FRONTEND_DIR)/dist/*
	@rm -rf $(FRONTEND_DIR)/node_modules/*
