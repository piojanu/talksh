# Use the official Golang image as a base
FROM golang:1.24.4

# Install zsh, python 3.12 and any other common tools
RUN apt-get update && apt-get install -y \
    curl \
    git \
    python3 \
    python3-venv \
    python3-pip \
    tree \
    vim \
    zsh \
    && rm -rf /var/lib/apt/lists/*

# Install zoxide system-wide
RUN curl -sSfL https://raw.githubusercontent.com/ajeetdsouza/zoxide/main/install.sh \
    | sh -s -- --bin-dir /usr/local/bin

# Create a playground directory that can be used by any user
RUN mkdir -p /playground && chmod 777 /playground

# Set the working directory inside the container
WORKDIR /playground
