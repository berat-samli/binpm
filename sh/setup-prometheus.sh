sudo apt install prometheus prometheus-node-exporter prometheus-pushgateway prometheus-alertmanager -y

sudo systemctl enable prometheus
sudo systemctl start prometheus

sudo systemctl status prometheus | grep Active