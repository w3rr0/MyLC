#!/bin/bash

# === Run Server ===

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


# 2. Launch Docker Compose
echo -e "${GREEN}🐳 Building and launching project inside Docker...${NC}"
docker compose up --build