package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Sabit dizin (örneğin /usr/local/share/binpm)
const defaultPackageDir = "/usr/local/share/binpm"

// Paket yapısı
type Package struct {
	Dependencies  []string `json:"dependencies"`
	InstallScript string   `json:"install_script"`
}

var packageList map[string]map[string]Package

// loadPackageList: Paket listesini sabit dizinden yükler
func loadPackageList() {
	// Çevresel değişkenden package dizinini al
	packageDir := os.Getenv("BINPM_PACKAGE_DIR")

	// Eğer çevresel değişken tanımlı değilse, varsayılan bir yol kullan
	if packageDir == "" {
		packageDir = "/usr/local/share/binpm"
	}

	// package_list.json dosyasının tam yolunu oluştur
	jsonPath := filepath.Join(packageDir, "package_list.json")

	// JSON dosyasını aç
	file, err := os.Open(jsonPath)
	if err != nil {
		fmt.Printf("Error opening package list: %v\n", err)
		return
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	// JSON'u çözümle
	err = json.Unmarshal(byteValue, &packageList)
	if err != nil {
		fmt.Printf("Error unmarshalling package list: %v\n", err)
	}
}

// InstallPackage: Paket yüklemesi ve bağımlılıkların kontrol edilmesi
func InstallPackage(pkgName string, osType string) {
	loadPackageList()

	// Uygun paketi bul
	pkg, exists := packageList[pkgName][osType]
	if !exists {
		fmt.Printf("Package not found for this OS: %s\n", pkgName)
		return
	}

	// Bağımlılıkları yükle
	installDependencies(pkg.Dependencies)

	// Yükleme script'ini çalıştır
	runInstallScript(pkg.InstallScript, osType)
}

// installDependencies: Bağımlılıkları yükleme fonksiyonu
func installDependencies(dependencies []string) {
	if len(dependencies) == 0 {
		fmt.Println("No dependencies to install.")
		return
	}

	fmt.Println("Installing dependencies...")
	for _, dep := range dependencies {
		fmt.Printf("Installing %s...\n", dep)
		cmd := exec.Command("sudo", "apt-get", "install", "-y", dep)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error installing dependency %s: %v\n", dep, err)
		}
	}
}

// runInstallScript: Yükleme script'ini çalıştırma fonksiyonu
func runInstallScript(scriptPath string, osType string) {
	// Sabit dizindeki script'in tam yolunu oluştur
	fullPath := filepath.Join(defaultPackageDir, scriptPath)

	fmt.Printf("Running install script: %s\n", fullPath)

	var cmd *exec.Cmd
	if osType == "windows" {
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", fullPath)
	} else {
		cmd = exec.Command("bash", fullPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running script %s: %v\n", fullPath, err)
	}
}
