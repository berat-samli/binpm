package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

func main() {
	binPath := "/usr/local/bin"

	var rootCmd = &cobra.Command{Use: "binpm"}

	var cmdInstall = &cobra.Command{
		Use:   "install [package]",
		Short: "Install a package",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			packageName := args[0]
			fmt.Printf("Installing package: %s\n", packageName)

			osType := runtime.GOOS

			var script string
			if osType == "linux" {
				script = filepath.Join(binPath, packageName+"_install_linux.sh")
			} else if osType == "darwin" {
				script = filepath.Join(binPath, packageName+"_install_macos.sh")
			} else {
				fmt.Println("Unsupported operating system")
				return
			}

			command := exec.Command("bash", script)
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr

			err := command.Run()
			if err != nil {
				fmt.Printf("Failed to install package: %s\n", err)
			} else {
				fmt.Println("Package installed successfully")
			}
		},
	}
	rootCmd.AddCommand(cmdInstall)
	rootCmd.Execute()
}
