sudo apt-get remove grafana
sudo apt-get remove --auto-remove grafana
sudo systemctl stop grafana-agent-flow
sudo apt-get remove grafana-agent-flow
sudo rm -i /etc/apt/sources.list.d/grafana.list