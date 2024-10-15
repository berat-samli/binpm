#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[1;34m'
NC='\033[0m' # No Color

CHECKMARK="✅"
CROSS="❌"
HOURGLASS="⏳"
WARNING="⚠️"

echo -e "${HOURGLASS} ${YELLOW}Removing old Docker packages...${NC}"
for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do
    sudo apt-get remove -y $pkg
    if [ $? -eq 0 ]; then
        echo -e "${CHECKMARK} ${GREEN}Removed package: $pkg${NC}"
    else
        echo -e "${CROSS} ${RED}Failed to remove package: $pkg${NC}"
    fi
done

# apt-get update
echo -e "${HOURGLASS} ${YELLOW}Updating package index...${NC}"
sudo apt-get update
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Package index updated successfully!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to update package index.${NC}"
    exit 1
fi

# ca-certificates ve curl yükle
echo -e "${HOURGLASS} ${YELLOW}Installing ca-certificates and curl...${NC}"
sudo apt-get install -y ca-certificates curl
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}ca-certificates and curl installed successfully!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to install ca-certificates and curl.${NC}"
    exit 1
fi

# Docker keyring'i oluştur
echo -e "${HOURGLASS} ${YELLOW}Setting up Docker keyring...${NC}"
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Docker keyring setup successfully!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to set up Docker keyring.${NC}"
    exit 1
fi


echo -e "${HOURGLASS} ${YELLOW}Adding Docker repository...${NC}"
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Docker repository added successfully!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to add Docker repository.${NC}"
    exit 1
fi

echo -e "${HOURGLASS} ${YELLOW}Updating package index for Docker repository...${NC}"
sudo apt-get update
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Package index updated successfully for Docker repository!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to update package index for Docker repository.${NC}"
    exit 1
fi

# Docker ve ilgili paketleri yükle
echo -e "${HOURGLASS} ${YELLOW}Installing Docker and related packages...${NC}"
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Docker installed successfully!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to install Docker.${NC}"
    exit 1
fi

echo -e "${HOURGLASS} ${YELLOW}Starting Docker service...${NC}"
sudo systemctl start docker
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Docker service started successfully!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to start Docker service.${NC}"
    exit 1
fi

echo -e "${CHECKMARK} ${GREEN}Docker installation and setup completed!${NC}"
