# Makefile for Go applicationset-dynamodb-plugin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt


AWS_ACCOUNT_ID ?= $(shell aws sts get-caller-identity --query Account --output text)

DIST_FOLDER := dist
BINARY_NAME := applicationset-dynamodb-plugin
REGISTRY_OWNER := tsahiduek
REGISTRY := public.ecr.aws/$(REGISTRY_OWNER)

IMAGE_NAME := $(BINARY_NAME)
IMAGE_REPO := $(BINARY_NAME)
IMAGE_TAG := latest
KO_DOCKER_REPO ?= public.ecr.aws/tsahiduek



# Application details
APP_NAME=applicationset-dynamodb-plugin

# Source files
SRC_FILES=$(wildcard *.go)

.PHONY: all build run test fmt compile

all: build

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	$(GOCMD) mod tidy
	$(GOBUILD) -o $(DIST_FOLDER)/$(APP_NAME) $(SRC_FILES)

.PHONY: run
run: 
	$(MAKE) build
	@echo "Running $(APP_NAME)..."
	$(DIST_FOLDER)/$(APP_NAME) --debug

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

compile: build
	@echo "Compiling for different platforms..."
	GOOS=linux $(GOBUILD) -o $(APP_NAME)_linux $(SRC_FILES)

# image:
# 	aws ecr-public get-login-password --region us-east-1 | podman login --username AWS --password-stdin public.ecr.aws/tsahiduek
# 	KO_DOCKER_REPO=ko.local DOCKER_HOST=unix:///run/podman/podman.sock ko publish . --local --tags $(IMAGE_NAME):$(IMAGE_TAG)

.PHONY: image-build
image-build:
	# podman build -t $(IMAGE_NAME):$(IMAGE_TAG) .
	podman image prune -a -f
	podman image rm $(IMAGE_NAME):$(IMAGE_TAG) -f
	podman manifest create -a $(IMAGE_NAME):$(IMAGE_TAG) 
	podman build --platform linux/amd64,linux/arm64 --manifest  $(IMAGE_NAME):$(IMAGE_TAG) .
	# podman build --os linux --arch amd64 --manifest  $(IMAGE_NAME):$(IMAGE_TAG) .


.PHONY: image-push
image-push:
	aws ecr-public get-login-password --region us-east-1 | podman login --username AWS --password-stdin $(REGISTRY)
	podman manifest push --all $(IMAGE_NAME):$(IMAGE_TAG) $(REGISTRY)/$(IMAGE_NAME)

.PHONY: image-release
image-release: image-build image-push