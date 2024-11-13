#!/bin/bash

# Grafana versiyonu ve GPG anahtar URL'si
GRAFANA_GPG_KEY_URL="https://apt.grafana.com/gpg.key"
GRAFANA_REPO="deb [signed-by=/etc/apt/keyrings/grafana.gpg] https://apt.grafana.com stable main"

# Gereken dizini oluştur
sudo mkdir -p /etc/apt/keyrings

# GPG anahtarını yalnızca eğer dosya mevcut değilse ekle
if [ ! -f /etc/apt/keyrings/grafana.gpg ]; then
    echo "Grafana GPG anahtarı ekleniyor..."
    curl -fsSL $GRAFANA_GPG_KEY_URL | sudo tee /etc/apt/keyrings/grafana.gpg > /dev/null
fi

# Grafana deposunu yalnızca eğer ekli değilse ekle
if ! grep -q "^$GRAFANA_REPO" /etc/apt/sources.list.d/grafana.list 2>/dev/null; then
    echo "Grafana deposu ekleniyor..."
    echo "$GRAFANA_REPO" | sudo tee /etc/apt/sources.list.d/grafana.list > /dev/null
fi

# Paket listelerini güncelleyin
echo "Paket listeleri güncelleniyor..."
sudo apt-get update

# Grafana ve gerekli bağımlılıkları kurun
echo "Grafana kuruluyor..."
sudo apt-get install -y grafana --fix-missing

# Grafana servis dosyasını oluşturun eğer mevcut değilse
if [ ! -f /etc/systemd/system/grafana-server.service ]; then
    echo "Grafana servis dosyası oluşturuluyor..."
    sudo tee /etc/systemd/system/grafana-server.service > /dev/null << EOL
[Unit]
Description=Grafana instance
Documentation=http://docs.grafana.org
Wants=network-online.target
After=network-online.target

[Service]
EnvironmentFile=-/etc/default/grafana-server
User=grafana
ExecStart=/usr/sbin/grafana-server \
  --config=/etc/grafana/grafana.ini \
  --pidfile=/var/run/grafana-server.pid \
  --packaging=deb \
  --homepath=/usr/share/grafana \
  cfg:default.paths.logs=/var/log/grafana \
  cfg:default.paths.data=/var/lib/grafana \
  cfg:default.paths.plugins=/var/lib/grafana/plugins \
  cfg:default.paths.provisioning=/etc/grafana/provisioning
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOL
fi

# Grafana servisini başlatın ve etkinleştirin
echo "Grafana servisi başlatılıyor..."
sudo systemctl daemon-reload
sudo systemctl enable grafana-server
sudo systemctl start grafana-server

# Servisin durumu hakkında bilgi sağlayın
sudo systemctl status grafana-server --no-pager

# Grafana'nın başlatılması için bekleme (daha uzun süreyle)
echo "Grafana'nın başlatılması bekleniyor..."
for i in {1..20}; do
    if curl -s http://localhost:3000/api/health | grep -q '"database":"ok"'; then
        echo "Grafana başarıyla başlatıldı ve erişilebilir durumda."
        break
    fi
    echo "Grafana henüz erişilebilir değil, tekrar denenecek..."
    sleep 10
done

# Eğer Grafana hala erişilemiyorsa hata ver
if ! curl -s http://localhost:3000/api/health | grep -q '"database":"ok"'; then
    echo "Grafana sunucusuna bağlanılamıyor. Lütfen servisin çalıştığından emin olun."
    exit 1
fi
