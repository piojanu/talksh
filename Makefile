# Makefile for talksh development

# Variables
BINARY_NAME=talksh
CONTAINER_NAME=talksh-dev
IMAGE_NAME=talksh-dev
# Build command that will run inside the container
GO_BUILD_COMMAND=go build -v -o $(BINARY_NAME) .

# Build the Docker development image
build-docker:
	@echo "Building $(IMAGE_NAME) Docker image..."
	@docker build -t $(IMAGE_NAME) .
	@echo "$(IMAGE_NAME) image built successfully."

# Build the binary inside a container for Linux compatibility
build-talksh-linux:
	@echo "Building $(BINARY_NAME) for Linux inside Docker container..."
	@docker run --rm \
		-v "$(shell pwd):/app" \
		-w /app \
		golang:1.24.4 \
		$(GO_BUILD_COMMAND)
	@echo "$(BINARY_NAME) built successfully for Linux."

# Run the development container interactively with native shell experience
# - Runs as your actual user/group to preserve file permissions
# - Sets proper hostname, HOME, SHELL, TERM, and USER environment variables
# - Mounts your home directory for access to shell configs, aliases, and dotfiles.
# - Mounts the talksh binary for easy execution from anywhere.
# - Mounts the project directory (read-only) for development context.
run-dev: build-talksh-linux # Build Linux binary instead of host binary
	@echo "Running $(IMAGE_NAME) container interactively..."
	@docker rm -f $(CONTAINER_NAME) || true
	@docker run -it \
		--name $(CONTAINER_NAME) \
		--user "$(shell id -u):$(shell id -g)" \
		--hostname "talksh-dev" \
		-e HOME="/home/$(shell whoami)" \
		-e SHELL="/bin/zsh" \
		-e TERM="$(TERM)" \
		-e USER="$(shell whoami)" \
		-v "$(HOME):/home/$(shell whoami)" \
		-v "$(shell pwd)/$(BINARY_NAME):/usr/local/bin/$(BINARY_NAME)" \
		-v "$(shell pwd):/app:ro" \
		$(IMAGE_NAME) /bin/zsh

# Rerun the existing development container
start-dev:
	@echo "Restarting existing $(CONTAINER_NAME) container..."
	@docker start -ai $(CONTAINER_NAME) || (echo "Container $(CONTAINER_NAME) not found. Use 'make run-dev' to create a new one." && exit 1)

# Clean the built binary from your host
clean:
	@echo "Cleaning built binary $(BINARY_NAME)..."
	@rm -f $(BINARY_NAME)
	@echo "Clean complete."

.PHONY: build-docker build-talksh-linux run-dev start-dev clean
