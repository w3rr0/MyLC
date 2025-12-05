#!/bin/bash

# --- CONFIGURATION ---
# Your repository URL
REPO_URL="https://github.com/w3rr0/MyLC.git"
# Folder name to clone into
DIR_NAME="server"

# --- colors ---
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Starting MyLC environment setup ===${NC}"

# 1. Check whether Docker is running
if ! docker info > /dev/null 2>&1; then
  echo -e "${RED}Error: Docker is not running!${NC}"
  echo "Launch Docker Desktop and try again."
  exit 1
fi

# 2. Clone or update the repository
if [ -d "$DIR_NAME" ]; then
    echo -e "${BLUE}Directory $DIR_NAME already exist.${NC}"
    echo "Downloading changes (git pull)..."
    cd "$DIR_NAME" || exit
    git pull
else
    echo -e "${GREEN}Cloning repository...${NC}"
    git clone "$REPO_URL" "$DIR_NAME"
    cd "$DIR_NAME" || exit
fi

# 3. Launch Docker Compose
echo -e "${GREEN}🐳 Building and launching project inside Docker...${NC}"
docker compose up --build