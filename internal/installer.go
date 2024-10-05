package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

type Package struct {
	Dependencies  []string `json:"dependencies"`
	InstallScript string   `json:"install_script"`
}

var packageList map[string]map[string]Package

func loadPackageList() {
	file, err := os.Open("packages/package_list.json")
	if err != nil {
		fmt.Println("Error opening package list:", err)
		return
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	json.Unmarshal(byteValue, &packageList)
}

func InstallPackage(pkgName string) {
	loadPackageList()

	osType := runtime.GOOS // "linux", "darwin", "windows" gibi değerler döner

	pkg, exists := packageList[pkgName][osType]
	if !exists {
		fmt.Println("Package not found for this OS:", pkgName)
		return
	}

	installDependencies(pkg.Dependencies)

	runInstallScript(pkg.InstallScript)
}

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

func runInstallScript(scriptPath string) {
	fmt.Printf("Running install script: %s\n", scriptPath)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", scriptPath)
	} else {
		cmd = exec.Command("bash", scriptPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running script %s: %v\n", scriptPath, err)
	}
}
