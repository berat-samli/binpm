package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

func main() {
	var privateKeyPath string
	var user string

	var rootCmd = &cobra.Command{
		Use:   "binpm",
		Short: "Binpm - Sunucu izleme (monitoring) yönetim aracı",
		Long:  `Binpm, bir veya daha fazla sunucuya monitoring bileşenlerini (Prometheus ve Grafana) kurmak için kullanılan bir komut satırı aracıdır.`,
	}

	var setupCmd = &cobra.Command{
		Use:   "setup monitor",
		Short: "Config dosyasındaki sunuculara Prometheus ve Grafana kurulumunu yapar",
		Long:  `Setup monitor komutu, config dosyasında tanımlanan her sunucuya SSH ile bağlanarak Prometheus ve Grafana gibi monitoring bileşenlerinin kurulumunu yapar.`,
		Run: func(cmd *cobra.Command, args []string) {
			setupMonitor(privateKeyPath, user)
		},
	}

	setupCmd.Flags().StringVarP(&privateKeyPath, "key", "k", "", "SSH private key yolu (zorunlu)")
	setupCmd.Flags().StringVarP(&user, "user", "u", "root", "SSH kullanıcı adı (varsayılan: root)")
	setupCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(setupCmd)
	rootCmd.Execute()
}
func setupMonitor(privateKeyPath, user string) {
	configFile := "/tmp/binpm/config.ini"

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		fmt.Println("Config dosyası bulunamadı. Lütfen /tmp/binpm/config.ini dosyasını oluşturun.")
		return
	}

	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("Config dosyası açılamadı: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var prometheusHost, grafanaHost string
	var setupDashboard bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "[") {
			prometheusHost, grafanaHost, setupDashboard = "", "", false
		} else if strings.HasPrefix(line, "prometheus-server-url") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				prometheusHost = strings.TrimSpace(parts[1])

				// Prometheus adresi yerelse doğrudan çalıştır
				if isLocalHost(prometheusHost) {
					runLocalScript("/tmp/binpm/sh/setup-prometheus.sh")
				} else {
					runRemoteScript(prometheusHost, privateKeyPath, user, "/tmp/binpm/sh/setup-prometheus.sh")
				}
			}
		} else if strings.HasPrefix(line, "grafana-server-url") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				grafanaHost = strings.TrimSpace(parts[1])

				// Grafana adresi yerelse doğrudan çalıştır
				if isLocalHost(grafanaHost) {
					runLocalScript("/tmp/binpm/sh/setup-grafana.sh")
					if setupDashboard {
						configureGrafana("localhost", prometheusHost)
					}
				} else {
					runRemoteScript(grafanaHost, privateKeyPath, user, "/tmp/binpm/sh/setup-grafana.sh")
					if setupDashboard {
						configureGrafana(grafanaHost, prometheusHost)
					}
				}
			}
		} else if strings.HasPrefix(line, "type") && strings.Contains(line, "node_exporter") {
			if isLocalHost(prometheusHost) {
				runLocalScript("/tmp/binpm/sh/setup-node-exporter.sh")
				updatePrometheusConfig("localhost", privateKeyPath, user)
			} else {
				runRemoteScript(prometheusHost, privateKeyPath, user, "/tmp/binpm/sh/setup-node-exporter.sh")
				updatePrometheusConfig(prometheusHost, privateKeyPath, user)
			}
			setupDashboard = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Config dosyasını okurken hata oluştu: %v", err)
	}
}

// Yerel kurulum için script çalıştırma fonksiyonu
func runLocalScript(scriptFile string) {
	fmt.Printf("Yerel makinede script çalıştırılıyor: %s\n", scriptFile)
	cmd := exec.Command("bash", "-c", fmt.Sprintf("sudo chmod +x %s && sudo bash %s", scriptFile, scriptFile))
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Yerel makinede script çalıştırılamadı: %v, Çıktı: %s, Hata Çıktısı: %s", err, stdout.String(), stderr.String())
	} else {
		fmt.Printf("Script yerel makinede başarıyla çalıştırıldı: %s\n", scriptFile)
	}
}

// IP'nin localhost veya yerel adreslerden biri olup olmadığını kontrol eden fonksiyon
func isLocalHost(host string) bool {
	localHosts := []string{"localhost", "127.0.0.1", "0.0.0.0"}
	for _, lh := range localHosts {
		if host == lh {
			return true
		}
	}
	return false
}


func configureGrafana(grafanaHost, prometheusHost string) {
	grafanaURL := fmt.Sprintf("http://%s:3000", grafanaHost)

	addDataSource(grafanaURL, prometheusHost)
	addDashboard(grafanaURL)
}

func addDataSource(grafanaURL, prometheusHost string) {
	url := fmt.Sprintf("%s/api/datasources", grafanaURL)
	dataSource := map[string]interface{}{
		"name":      "Prometheus",
		"type":      "prometheus",
		"access":    "proxy",
		"url":       fmt.Sprintf("http://%s:9090", prometheusHost),
		"isDefault": true,
	}
	body, _ := json.Marshal(dataSource)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Data Source isteği oluşturulamadı: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Grafana Data Source isteği başarısız: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Grafana Data Source oluşturulamadı, durum kodu: %d", resp.StatusCode)
	}

	fmt.Println("Grafana'da Prometheus data source başarıyla oluşturuldu.")
}

func addDashboard(grafanaURL string) {
	url := fmt.Sprintf("%s/api/dashboards/db", grafanaURL)
	dashboardFile := "/tmp/binpm/dashboards/node_exporter.json"

	dashboardContent, err := ioutil.ReadFile(dashboardFile)
	if err != nil {
		log.Fatalf("Dashboard dosyası okunamadı: %v", err)
	}

	var dashboard map[string]interface{}
	if err := json.Unmarshal(dashboardContent, &dashboard); err != nil {
		log.Fatalf("Dashboard JSON dosyası çözümlenemedi: %v", err)
	}

	payload := map[string]interface{}{
		"dashboard": dashboard,
		"overwrite": true,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Dashboard isteği oluşturulamadı: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "admin")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Grafana Dashboard isteği başarısız: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Grafana Dashboard oluşturulamadı, durum kodu: %d", resp.StatusCode)
	}

	fmt.Println("Grafana'da node_exporter dashboard başarıyla oluşturuldu.")
}

func runRemoteScript(host, privateKeyPath, user, scriptFile string) {
	// Hedef sunucuda dizin yapısını oluştur
	conn, err := createSSHClient(host, user, privateKeyPath)
	if err != nil {
		log.Fatalf("SSH bağlantısı sağlanamadı: %v", err)
	}
	defer conn.Close()

	// Hedef sunucuda dizinleri oluşturmak için SSH oturumu başlatıyoruz
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("SSH oturumu başlatılamadı (dizin oluşturma için): %v", err)
	}
	defer session.Close()

	// Hedef sunucuda /tmp/binpm ve /tmp/binpm/sh dizinlerini oluşturma komutu
	makeDirCommand := "sudo mkdir -p /tmp/binpm/sh && sudo mkdir -p /tmp/binpm/dashboards && sudo chmod -R 777 /tmp/binpm"
	if err := session.Run(makeDirCommand); err != nil {
		log.Fatalf("Hedef sunucuda dizinler oluşturulamadı: %v", err)
	} else {
		fmt.Printf("Hedef sunucuda dizinler oluşturuldu: /tmp/binpm/sh ve /tmp/binpm/dashboards\n")
	}

	// `scp` komutunu çalıştırarak script dosyasını gönderme
	scpCommand := fmt.Sprintf("scp -i %s %s %s@%s:/tmp/binpm/sh/", privateKeyPath, scriptFile, user, host)
	fmt.Printf("Sunucuya gönderiliyor: %s, Hedef IP: %s\n", scriptFile, host)
	output, err := exec.Command("bash", "-c", scpCommand).CombinedOutput()

	if err != nil {
		log.Fatalf("Script dosyası %s sunucusuna gönderilemedi. Hedef IP: %s, Dosya: %s, Hata: %v, Çıktı: %s", host, host, scriptFile, err, string(output))
	} else {
		fmt.Printf("Script dosyası başarıyla gönderildi: %s, Hedef IP: %s\n", scriptFile, host)
	}

	// SSH üzerinden chmod + script çalıştırmayı tek bir komutta çalıştırıyoruz
	session, err = conn.NewSession()
	if err != nil {
		log.Fatalf("SSH oturumu başlatılamadı (script çalıştırma için): %v", err)
	}
	defer session.Close()

	var stdout, stderr strings.Builder
	session.Stdout = &stdout
	session.Stderr = &stderr
	// Doğru dosya yolunu belirtelim
	//scriptPath := fmt.Sprintf("/tmp/binpm/sh/%s", scriptFile)
	scriptPath := fmt.Sprintf("/%s", scriptFile)
	runCommand := fmt.Sprintf("sudo chmod +x %s && sudo bash %s", scriptPath, scriptPath)
	if err := session.Run(runCommand); err != nil {
		log.Fatalf("Sunucuda script çalıştırılamadı: %v, Çıktı: %s, Hata Çıktısı: %s", err, stdout.String(), stderr.String())
	} else {
		fmt.Printf("Script dosyası başarıyla çalıştırıldı: %s\n", scriptFile)
	}
}

func updatePrometheusConfig(prometheusHost, privateKeyPath, user string) {
	conn, err := createSSHClient(prometheusHost, user, privateKeyPath)
	if err != nil {
		log.Fatalf("SSH bağlantısı sağlanamadı: %v", err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("SSH oturumu başlatılamadı (prometheus.yml güncelleme için): %v", err)
	}
	defer session.Close()

	prometheusConfig := `
  - job_name: 'node_exporter'
    static_configs:
      - targets: ['localhost:9100']
`

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	updateCommand := fmt.Sprintf(`echo "%s" | sudo tee -a /tmp/binpm/prometheus.yml > /dev/null`, prometheusConfig)
	if err := session.Run(updateCommand); err != nil {
		log.Fatalf("prometheus.yml güncellenemedi: %v, Çıktı: %s, Hata Çıktısı: %s", err, stdout.String(), stderr.String())
	}

	fmt.Printf("Prometheus konfigürasyon dosyasına node_exporter eklendi. Hedef IP: %s\n", prometheusHost)

	session, err = conn.NewSession()
	if err != nil {
		log.Fatalf("SSH oturumu başlatılamadı (Prometheus yeniden başlatma için): %v", err)
	}
	defer session.Close()

	if err := session.Run("sudo systemctl restart prometheus"); err != nil {
		log.Fatalf("Prometheus yeniden başlatılamadı: %v", err)
	}

	fmt.Println("Prometheus yeniden başlatıldı ve node_exporter konfigürasyonu aktif edildi.")
}

func createSSHClient(host, user, privateKeyPath string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("SSH anahtarı okunamadı: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("SSH anahtarı hatalı: %v", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, fmt.Errorf("%s sunucusuna bağlantı sağlanamadı: %v", host, err)
	}

	return client, nil
}