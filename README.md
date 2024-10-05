# binpm - binpm is not package manager (still in dev)

`binpm` is a lightweight package manager written in Go that allows you to install tools like Docker, Terraform, and more on different operating systems, including Linux, MacOS, and Windows. It automatically detects your operating system and architecture, installs the necessary dependencies, and runs the correct installation scripts.

## Supported Platforms
- Linux (Ubuntu, Debian, CentOS)
- MacOS
- Windows

## Features
- Install tools like Docker, Terraform, and more with a simple command.
- Automatically installs the necessary dependencies for each tool.
- Platform-specific installation scripts.
- Supports Linux, MacOS, and Windows.

---

## Installation

### 1. Clone the repository

To start using `binpm`, first clone the repository:

```bash
git clone https://github.com/bertt6/binpm.git
cd binpm
```

### 2. Build the project

To build the binary for your system:

#### Linux / MacOS:

```bash
go build -o binpm cmd/main.go
```

#### Windows:

```bash
go build -o binpm.exe cmd/main.go
```

### 3. Environment Variable Setup

`binpm` uses environment variables to find the necessary `package_list.json` file and installation scripts. You will need to set an environment variable `BINPM_PACKAGE_DIR` to specify the directory where your scripts and `package_list.json` are located.

You can do this by creating a `.env` file in your project directory and specifying the path to your package directory.

1. Create a `.env` file in your project directory:

```bash
touch .env
```

2. Add the following line to your `.env` file:

```bash
BINPM_PACKAGE_DIR=/path/to/your/packages
```

Replace `/path/to/your/packages` with the actual path where your `package_list.json` and scripts are located.

3. Load the `.env` file into your environment:

For **Linux/MacOS**:

```bash
export $(grep -v '^#' .env | xargs)
```

For **Windows (PowerShell)**:

```powershell
$env:BINPM_PACKAGE_DIR = "C:\path\to\your\packages"
```

After setting up the environment variable, `binpm` will automatically look for the `package_list.json` and installation scripts in the specified directory.

### 4. Verify Installation

To verify that `binpm` is correctly installed, run:

```bash
binpm --help
```

This should display the usage instructions.

---

## Usage

### Install a tool

To install a tool, use the following command:

```bash
binpm install <tool_name>
```

For example, to install Docker:

```bash
binpm install docker
```

This will automatically detect your operating system and architecture, install the required dependencies, and run the appropriate installation script.

### Uninstall binpm

For now, uninstalling `binpm` is a manual process. You just need to remove the binary from your system.

#### Linux / MacOS:

To uninstall `binpm`, simply remove the binary from `/usr/local/bin`:

```bash
sudo rm /usr/local/bin/binpm
```

#### Windows:

To uninstall `binpm` on Windows, delete the `binpm.exe` file from the directory where you placed it (e.g., `C:\Windows\System32` or another PATH directory):

```powershell
Remove-Item C:\Windows\System32\binpm.exe
```

This will remove `binpm` from your system.

---

## Platform-Specific Installation Details

### Linux

On Linux, `binpm` uses `apt-get` (for Ubuntu/Debian) or `yum` (for CentOS/RHEL) to install dependencies and tools. Make sure you have `sudo` permissions.

For example, to install Docker on Linux:

```bash
binpm install docker
```

This will:
- Install dependencies like `curl`, `apt-transport-https`, etc.
- Download and install Docker using the official Docker repository.

### MacOS

On MacOS, `binpm` relies on `brew` (Homebrew) to install packages. Ensure that Homebrew is installed on your system.

To install Docker on MacOS:

```bash
binpm install docker
```

This will:
- Use Homebrew to install Docker.

### Windows

On Windows, `binpm` uses PowerShell scripts to download and install tools. Make sure to run `binpm` in a PowerShell session with administrator rights.

To install Docker on Windows:

```bash
binpm install docker
```

This will:
- Download the Docker installer and run it.

---

## Development

If you want to contribute or modify `binpm`, you can run the project locally and make changes.

### Run Locally

To run `binpm` locally:

1. Clone the repository:
   ```bash
   git clone https://github.com/bertt6/binpm.git
   cd binpm
   ```

2. Make your changes.

3. Build the project:
   - **Linux / MacOS**:
     ```bash
     go build -o binpm cmd/main.go
     ```
   - **Windows**:
     ```bash
     go build -o binpm.exe cmd/main.go
     ```

### Testing

You can test the package manager by running it locally:

```bash
./binpm install docker
```

---

## Roadmap

- [ ] Add uninstall functionality for tools
- [ ] Add versioning support for package installations
- [ ] Improve error handling and logging
- [ ] Add more tools (Node.js, Python, etc.)

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.