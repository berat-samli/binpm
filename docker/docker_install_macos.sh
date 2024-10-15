#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

CHECKMARK="✅"
CROSS="❌"
HOURGLASS="⏳"
WARNING="⚠️"

echo -e "${HOURGLASS} ${YELLOW}Checking Brew...${NC}"
sleep 5
if command -v brew &> /dev/null
then
    echo -e "${CHECKMARK} ${GREEN}Brew is installed.${NC}"
else
    echo -e "${CROSS} ${RED}Brew is not installed.${NC}"
    echo -e "${WARNING} ${YELLOW}Please install Homebrew first: https://brew.sh${NC}"
    exit 1
fi

echo -e "${HOURGLASS} ${YELLOW}Installing Docker...${NC}"
brew install --cask docker

if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Docker has been successfully installed!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to install Docker.${NC}"
    exit 1
fi
