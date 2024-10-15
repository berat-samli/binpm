
# binpm

`binpm` is a lightweight package manager built with Go, designed for easy installation and uninstallation of packages with shell script support. It provides simple command-line utilities to manage packages across different environments.

## Features
- **Cross-platform support**: Supports both macOS and Linux for package installation and uninstallation.
- **Shell script-based**: Easily configurable with shell scripts for package management.
- **Global command access**: Set up `binpm` to be accessible globally from anywhere on your system.
- **Modular design**: Easily extendable to handle various packages with installation and uninstallation scripts.

## Installation

### MacOS

To install `binpm` globally on your macOS system, follow the steps below:

1. Clone this repository:

   ```bash
   git clone https://github.com/bertt6/binpm.git
   ```

2. Run the setup script to build and install `binpm`:

   ```bash
   cd binpm
   chmod +x setup_macos.sh
   ./setup_macos.sh
   ```

This will move the `binpm` binary to `/usr/local/bin/` and set up the necessary shell scripts.

### Linux

For Linux systems, follow these steps to install `binpm` globally:

1. Clone this repository:

   ```bash
   git clone https://github.com/bertt6/binpm.git
   ```

2. Create .env file for SCRIPTS_PATH:
   ```
   SCRIPTS_PATH=/binpm/folder/path
   ```


3. Run the setup script to build and install `binpm`:

   ```bash
   cd binpm
   chmod +x setup.sh
   ./setup.sh
   ```

This will move the `binpm` binary to `/usr/local/bin/` and set up the necessary shell scripts.

## Usage

Once installed, `binpm` can be used globally to install or uninstall packages by running:

### Install a Package
```bash
binpm install <package-name>
```

For example, to install Docker:
```bash
binpm install docker
```

### Uninstall a Package
```bash
binpm uninstall <package-name>
```

For example, to uninstall Docker:
```bash
binpm uninstall docker
```

## Contributing

Feel free to submit issues and pull requests! Contributions are always welcome.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
