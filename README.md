# a0ctl

The official command-line interface for [a0.run](https://a0.run) - a platform for deploying and managing applications.

## About

`a0ctl` is a Go-based CLI tool that provides a convenient way to interact with the a0.run platform. It allows you to authenticate, manage configurations, and perform various operations on your applications and deployments.

## Installation

### From Source

To build and install `a0ctl` from source, you need Go 1.24.3 or later:

```bash
git clone https://github.com/a0dotrun/a0ctl.git
cd a0ctl
go build -o a0ctl cmd/a0ctl/main.go
```

### Running Directly

You can also run the CLI directly without building:

```bash
go run cmd/a0ctl/main.go [command]
```

## Usage

### Getting Help

```bash
# Show general help
./a0ctl --help

# Show help for a specific command
./a0ctl [command] --help
```

### Available Commands

- **`auth`** - Manage authentication
  - `auth login` - Login to the platform
  - `auth whoami` - Show the current logged in user or token user
- **`config`** - Manage your CLI configuration
- **`version`** - Show version information for the a0ctl CLI
- **`completion`** - Generate the autocompletion script for the specified shell

### Examples

```bash
# Check version
./a0ctl version

# Login to a0.run
./a0ctl auth login

# Check current user
./a0ctl auth whoami

# Show help for configuration commands
./a0ctl config --help
```

## Authentication

Before using most commands, you'll need to authenticate with the a0.run platform:

```bash
./a0ctl auth login
```

This will open your browser and guide you through the authentication process.

## Configuration

The CLI stores configuration and authentication tokens in your home directory under `.a0/`. This directory is automatically created when needed.

## Development

### Prerequisites

- Go 1.24.3 or later
- Git

### Building

```bash
# Clone the repository
git clone https://github.com/a0dotrun/a0ctl.git
cd a0ctl

# Install dependencies
go mod download

# Build the binary
go build -o a0ctl cmd/a0ctl/main.go

# Run tests (if available)
go test ./...
```

### Project Structure

```
├── cmd/a0ctl/          # Main application entry point
├── internal/
│   ├── api/            # API client implementation
│   ├── cli/            # CLI utilities and helpers
│   ├── command/        # Command implementations
│   │   ├── auth/       # Authentication commands
│   │   ├── config/     # Configuration commands
│   │   ├── root/       # Root command setup
│   │   └── version/    # Version command
│   ├── flags/          # Command-line flag definitions
│   └── settings/       # Configuration and settings
├── examples/           # Example applications
└── go.mod             # Go module definition
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For support and documentation, visit [a0.run](https://a0.run) or check the help output of individual commands.