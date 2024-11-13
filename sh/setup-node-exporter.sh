#!/bin/bash

# Node Exporter versiyonu
VERSION="1.7.0"

echo "1. Node Exporter indiriliyor..."
if ! sudo wget -q https://github.com/prometheus/node_exporter/releases/download/v$VERSION/node_exporter-$VERSION.linux-amd64.tar.gz; then
    echo "Node Exporter indirme başarısız, kurulumu atlanıyor."
    exit 0
fi

echo "2. Dosyalar açılıyor..."
if ! sudo tar xzf node_exporter-$VERSION.linux-amd64.tar.gz; then
    echo "Dosyalar açılamadı, kurulumu atlanıyor."
    exit 0
fi

echo "3. İndirilen dosya siliniyor..."
sudo rm -f node_exporter-$VERSION.linux-amd64.tar.gz

# Hedef dizini kontrol et ve eski dosyaları temizle
echo "4. Hedef dizin kontrol ediliyor..."
if [ -d "/etc/node_exporter" ]; then
    echo "Eski dosyalar temizleniyor..."
    sudo rm -rf /etc/node_exporter/*
else
    sudo mkdir -p /etc/node_exporter
fi

echo "5. Dosyalar taşınıyor..."
if ! sudo mv node_exporter-$VERSION.linux-amd64/* /etc/node_exporter/; then
    echo "Dosyalar taşınamadı, kurulumu atlanıyor."
    exit 0
fi
sudo rm -rf node_exporter-$VERSION.linux-amd64

echo "6. Servis dosyası oluşturuluyor..."
sudo tee /etc/systemd/system/node_exporter.service > /dev/null << EOL
[Unit]
Description=Node Exporter
Wants=network-online.target
After=network-online.target

[Service]
ExecStart=/etc/node_exporter/node_exporter
Restart=always

[Install]
WantedBy=multi-user.target
EOL

echo "7. Node Exporter başlatılıyor..."
if ! sudo systemctl daemon-reload || ! sudo systemctl enable node_exporter || ! sudo systemctl start node_exporter; then
    echo "Node Exporter servisi başlatılamadı, kurulumu atlanıyor."
    exit 0
fi

echo "8. Node Exporter durumu kontrol ediliyor..."
sudo systemctl status node_exporter --no-pager || echo "Node Exporter durumu kontrol edilemedi."