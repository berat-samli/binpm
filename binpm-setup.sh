#!/bin/bash

# Hedef dizin
TARGET_DIR="/tmp/binpm"

# Hedef dizini oluştur
echo "Kurulum dizini oluşturuluyor: $TARGET_DIR"
sudo mkdir -p "$TARGET_DIR/sh"
sudo mkdir -p "$TARGET_DIR/dashboards"

# config.ini dosyasını taşı
echo "config.ini dosyası taşınıyor..."
sudo mv config.ini "$TARGET_DIR/"

# dashboards dizinindeki dosyaları taşı
echo "Dashboards dosyaları taşınıyor..."
sudo mv dashboards/*.json "$TARGET_DIR/dashboards/"

# sh dizinindeki script dosyalarını taşı
echo "Kurulum scriptleri taşınıyor..."
sudo mv sh/*.sh "$TARGET_DIR/sh/"

# prometheus.yml dosyasını taşı
echo "prometheus.yml dosyası taşınıyor..."
sudo mv prometheus.yml "$TARGET_DIR/"

# binpm'i PATH dizinine eklemek için /usr/local/bin dizinine kopyalayarak çalıştırılabilir hale getir
if [ -f ./binpm ]; then
    echo "binpm dosyası /usr/local/bin dizinine taşınıyor..."
    sudo cp ./binpm /usr/local/bin/binpm
    sudo chmod +x /usr/local/bin/binpm
else
    echo "binpm dosyası mevcut değil. Lütfen önce 'go build -o binpm main.go' komutuyla derleyin."
fi

# Projede artık gereksiz hale gelen dizinleri sil
echo "Geçici dizinleri temizleme..."
rm -rf dashboards sh

echo "Kurulum tamamlandı. binpm komutunu terminalde herhangi bir yerden çalıştırabilirsiniz."
