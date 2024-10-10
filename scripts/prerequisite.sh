#!/bin/bash

# Function to log messages
log() {
    echo "$(date +'%Y-%m-%d %H:%M:%S') - $1"
}

# Function to install essential tools
install_tools() {
    log "Updating package lists..."
    if ! apt-get update; then
        log "Failed to update package lists."
        exit 1
    fi
    log "Installing wget, curl, and other essential tools..."
    if ! apt-get install -y wget curl; then
        log "Failed to install essential tools."
        exit 1
    fi
}

# Function to install Go
install_go() {
    log "Installing Go..."
    wget -q https://dl.google.com/go/go1.22.5.linux-amd64.tar.gz
    if [ -f go1.22.5.linux-amd64.tar.gz ]; then
        tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
        rm go1.22.5.linux-amd64.tar.gz
    else
        log "Failed to download Go."
        exit 1
    fi
}

# Function to install IPFS
install_ipfs() {
    log "Installing IPFS..."
    wget -q https://github.com/ipfs/kubo/releases/download/v0.29.0/kubo_v0.29.0_linux-amd64.tar.gz
    if [ -f kubo_v0.29.0_linux-amd64.tar.gz ]; then
        tar -C /usr/local/bin -xzf kubo_v0.29.0_linux-amd64.tar.gz --strip-components=1
        rm kubo_v0.29.0_linux-amd64.tar.gz
    else
        log "Failed to download IPFS."
        exit 1
    fi
}

# Function to install Python and required dependencies
install_python() {
    log "Installing Python and required dependencies..."
    if ! apt-get install -y python3 python3-pip python3-dev build-essential; then
        log "Failed to install Python and dependencies."
        exit 1
    fi
    log "Installing Python packages with pip..."
    if ! pip3 install --break-system-packages llama_index==0.11.5 openai==1.43.0; then
        log "Failed to install Python packages."
        exit 1
    fi
}

# Function to set up environment variables for Go
setup_go_env() {
    log "Setting up environment variables for Go..."
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
}

# First, install essential tools
install_tools

# Start other installations in parallel
install_go &
install_ipfs &
install_python &

# Wait for all background jobs to finish
wait

# Setup Go environment variables
setup_go_env

log "Prerequisites installation completed."
