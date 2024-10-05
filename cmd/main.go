package main

import (
	"binpm/internal"
	"fmt"
	"os"
	"runtime"
)

func main() {
	// En az 3 argüman bekliyoruz
	if len(os.Args) < 3 {
		fmt.Println("Usage: binpm <command> <package> [--packagedir=<path_to_package_dir>]")
		return
	}

	command := os.Args[1]
	pkgName := os.Args[2]

	switch command {
	case "install":
		fmt.Printf("Installing package: %s\n", pkgName)

		osType := runtime.GOOS     // "linux", "darwin", "windows" gibi döner
		archType := runtime.GOARCH // "amd64", "arm64" gibi döner

		fmt.Printf("Detected OS: %s, Architecture: %s\n", osType, archType)

		internal.InstallPackage(pkgName, osType)
	default:
		fmt.Println("Unknown command:", command)
	}
}
