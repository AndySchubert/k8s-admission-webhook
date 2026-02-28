# ── Config ────────────────────────────────────────────────────────────
APP          := webhook
MODULE       := github.com/AndySchubert/k8s-admission-webhook
IMAGE        := ghcr.io/andyschubert/k8s-admission-webhook
TAG          ?= dev
BIN_DIR      := bin
NAMESPACE    := platform-system
DEPLOY_NAME  := k8s-admission-webhook
CERTS_DIR    := certs

# ── Build ─────────────────────────────────────────────────────────────
.PHONY: build
build: ## Build the webhook binary
	CGO_ENABLED=0 go build -o $(BIN_DIR)/$(APP) ./cmd/webhook

.PHONY: run
run: build certs ## Run the webhook locally (generates certs if needed)
	TLS_CERT=$(CERTS_DIR)/server/tls.crt TLS_KEY=$(CERTS_DIR)/server/tls.key $(BIN_DIR)/$(APP)

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(BIN_DIR)

# ── Quality ───────────────────────────────────────────────────────────
.PHONY: test
test: ## Run unit tests
	go test -v -race ./...

.PHONY: cover
cover: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

.PHONY: lint
lint: ## Run golangci-lint
	@which golangci-lint > /dev/null 2>&1 || \
		{ echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	golangci-lint run ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Format Go source files
	gofmt -w .

.PHONY: check
check: vet lint test ## Run all quality checks

# ── Docker ────────────────────────────────────────────────────────────
.PHONY: docker-build
docker-build: ## Build the container image
	docker build -t $(IMAGE):$(TAG) .

.PHONY: docker-push
docker-push: ## Push the container image to GHCR
	docker push $(IMAGE):$(TAG)

.PHONY: docker
docker: docker-build docker-push ## Build and push the container image

# ── Certificates ──────────────────────────────────────────────────────
.PHONY: certs
certs: ## Generate self-signed CA + server certificates
	CERTS_DIR=$(CERTS_DIR) ./scripts/gen-certs.sh

.PHONY: certs-clean
certs-clean: ## Remove generated certificates
	rm -rf $(CERTS_DIR)

# ── Deploy ────────────────────────────────────────────────────────────
.PHONY: manifests
manifests: ## Render deploy manifests (inject CA bundle)
	./scripts/render-manifests.sh

.PHONY: secret
secret: ## Create the TLS secret in Kubernetes
	kubectl create namespace $(NAMESPACE) 2>/dev/null || true
	kubectl create secret tls webhook-tls \
		-n $(NAMESPACE) \
		--cert=$(CERTS_DIR)/server/tls.crt \
		--key=$(CERTS_DIR)/server/tls.key

.PHONY: deploy
deploy: manifests ## Apply manifests and wait for rollout
	kubectl apply -f deploy/manifests.yaml
	kubectl rollout status deploy/$(DEPLOY_NAME) -n $(NAMESPACE)

.PHONY: undeploy
undeploy: ## Delete all deployed resources
	kubectl delete -f deploy/manifests.yaml --ignore-not-found
	kubectl delete secret webhook-tls -n $(NAMESPACE) --ignore-not-found

.PHONY: logs
logs: ## Tail webhook pod logs
	kubectl logs -n $(NAMESPACE) -l app=$(DEPLOY_NAME) -f

# ── Setup (full first-time flow) ──────────────────────────────────────
.PHONY: setup
setup: certs secret manifests deploy ## Full first-time setup: certs → secret → manifests → deploy

# ── Help ──────────────────────────────────────────────────────────────
.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
