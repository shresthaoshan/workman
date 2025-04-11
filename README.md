# Workman: A CLI Tool to Manage EC2 Instances


[![Release](https://github.com/shresthaoshan/workman/actions/workflows/release.yml/badge.svg)](https://github.com/shresthaoshan/workman/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue)](https://golang.org/) 
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**Workman** is a lightweight, open-source CLI tool designed to simplify the management of AWS EC2 instances. With Workman, you can start, stop, SSH into, configure, list, and remove EC2 instances from the command line. It is built using Go and integrates seamlessly with AWS SDK for Go.

---

## Table of Contents

1. [Features](#features)
2. [Installation](#installation)
3. [Configuration](#configuration)
4. [Usage](#usage)
   - [Start an Instance](#start-an-instance)
   - [Stop an Instance](#stop-an-instance)
   - [SSH into an Instance](#ssh-into-an-instance)
   - [List Configured Instances](#list-configured-instances)
   - [Configure a New Instance](#configure-a-new-instance)
   - [Remove an Instance Configuration](#remove-an-instance-configuration)
5. [Contributing](#contributing)
6. [License](#license)

---

## Features

- **Start/Stop EC2 Instances**: Start or stop instances with a single command.
- **SSH Integration**: Seamlessly SSH into instances using their PEM files.
- **Interactive Configuration**: Add new instance configurations interactively.
- **List Instances**: View all configured instances along with metadata like `LastAccessed`.
- **Remove Configurations**: Safely remove unused instance configurations.
- **Cross-Platform Support**: Build binaries for Linux, macOS (including M-series CPUs), and Windows.
- **AWS Profile Support**: Use different AWS profiles for managing multiple accounts.

---

## Installation

1. Ensure you have Go (>= 1.20) installed on your system.
2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/workman.git
   cd workman
   ```
3. Build the binary:
   ```bash
   go build -o workman
   ```
4. Move the binary to a directory in your PATH:
   ```bash
   mv workman /usr/local/bin/
   ```

---

## Configuration

1. Ensure you have AWS credentials configured. You can use the AWS CLI to set them up:
   ```bash
   aws configure
   ```
2. Run the `workman` CLI to initialize the configuration:
   ```bash
   workman configure
   ```
3. Follow the interactive prompts to add your EC2 instance details.

---

## Usage

### Start an Instance
```bash
workman start <instance-id>
```

### Stop an Instance
```bash
workman stop <instance-id>
```

### SSH into an Instance
```bash
workman ssh <instance-id>
```

### List Configured Instances
```bash
workman list
```

### Configure a New Instance
```bash
workman configure
```

### Remove an Instance Configuration
```bash
workman remove <instance-id>
```

---

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository.
2. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b feature-name
   ```
3. Commit your changes:
   ```bash
   git commit -m "Add feature-name"
   ```
4. Push to your fork:
   ```bash
   git push origin feature-name
   ```
5. Open a pull request on the main repository.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

