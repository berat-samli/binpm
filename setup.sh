#!/bin/bash

# Renk kodları
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Emojiler
CHECKMARK="✅"
CROSS="❌"
HOURGLASS="⏳"

# .env dosyasının kaynağı ve hedefi
ENV_SOURCE="./.env"
ENV_DEST="/usr/local/bin/.env"

# Go dosyasını build et
echo -e "${HOURGLASS} ${YELLOW}Building the Go project...${NC}"
go build -o binpm
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Build successful!${NC}"
else
    echo -e "${CROSS} ${RED}Build failed!${NC}"
    exit 1
fi

# Binaries dizinine taşı
echo -e "${HOURGLASS} ${YELLOW}Moving the binary to /usr/local/bin...${NC}"
sudo mv binpm /usr/local/bin/
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}binpm is now available in /usr/local/bin!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to move binpm to /usr/local/bin!${NC}"
    exit 1
fi


# Shell script dosyalarını da aynı dizine kopyalayalım
echo -e "${HOURGLASS} ${YELLOW}Copying shell scripts to /usr/local/bin...${NC}"
sudo cp ./docker/docker_install_macos.sh /usr/local/bin/
sudo cp ./docker/docker_install_linux.sh /usr/local/bin/
if [ $? -eq 0 ]; then
    echo -e "${CHECKMARK} ${GREEN}Shell scripts copied to /usr/local/bin!${NC}"
else
    echo -e "${CROSS} ${RED}Failed to copy shell scripts!${NC}"
    exit 1
fi

# /etc/binpm dizinini oluştur ve .env dosyasını kopyala
echo -e "${HOURGLASS} ${YELLOW}Setting up .env file in /etc/binpm...${NC}"
if [ ! -d "/etc/binpm" ]; then
    sudo mkdir /etc/binpm
fi

if [ -f "$ENV_SOURCE" ]; then
    sudo cp "$ENV_SOURCE" "$ENV_DEST"
    if [ $? -eq 0 ]; then
        echo -e "${CHECKMARK} ${GREEN}.env file copied to /etc/binpm!${NC}"
    else
        echo -e "${CROSS} ${RED}Failed to copy .env file!${NC}"
        exit 1
    fi
else
    echo -e "${CROSS} ${RED}.env file not found!${NC}"
    exit 1
fi

echo -e "${CHECKMARK} ${GREEN}Setup completed successfully! You can now use 'binpm' from anywhere.${NC}"
