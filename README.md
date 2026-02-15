# MediaWiki Updater & Downloader

A modern, feature-rich CLI tool for automating MediaWiki installation and updates. Built with Go and Cobra, this tool downloads MediaWiki core along with configured extensions and skins from multiple sources.

## ✨ Features

- **📦 MediaWiki Core**: Downloads any version from official releases with automatic URL parsing
- **🧩 Extensions & Skins**: Support for both ExtDist and Git repositories
- **🔧 Flexible Configuration**: INI-based configuration with version-specific downloads
- **🏗️ Modular Architecture**: Clean, maintainable codebase with separated concerns
- **🛡️ Safe Operations**: Preserves important files during updates (LocalSettings.php, images, etc.)
- **📋 Discovery Tools**: List available versions, extensions, and skins
- **⚠️ Graceful Handling**: Continues operation even if individual components fail to download

## 🚀 Installation

### Option 1: Download Binary

Download the latest binary from the [releases page](https://github.com/SKevo18/mediawiki-updater/releases).

### Option 2: Build from Source

```bash
git clone https://github.com/SKevo18/mediawiki-updater.git
cd mediawiki-updater
go build -o mediawiki-updater
```

## 📖 Usage

### Basic Update

```bash
./mediawiki-updater --config config.ini --target /path/to/mediawiki
```

### Available Commands

```bash
# Show help
./mediawiki-updater --help

# List available MediaWiki versions
./mediawiki-updater list versions

# List available extensions
./mediawiki-updater list extensions

# List available skins
./mediawiki-updater list skins

# Update with verbose output
./mediawiki-updater --verbose --config my-config.ini --target /var/www/mediawiki
```

## ⚙️ Configuration

Create an INI configuration file (default: `config.ini`) with the following sections:

```ini
[mediawiki]
version=1.43.1

[skins]
; From ExtDist (official distribution)
extdist1=Vector
extdist2=Timeless

; From Git repository
git1=https://github.com/StarCitizenWiki/mediawiki-skins-Citizen.git|main

[extensions]
; From ExtDist
extdist1=Cite
extdist2=VisualEditor
extdist3=WikiEditor

; With specific version
extdist4=Math|REL1_43

; From Git repository
git1=https://github.com/wikimedia/mediawiki-extensions-MobileFrontend.git|REL1_43
```

### Configuration Sections

#### `[mediawiki]`

- `version`: MediaWiki version to download (e.g., "1.43.1")

#### `[extensions]` and `[skins]`

- `extdist=<name>`: Download from ExtDist using the MediaWiki version
- `extdist=<name>:<version>`: Download specific version from ExtDist
- `git=<repo-url>:<branch>`: Clone from Git repository (branch defaults to "master")

## 🗂️ Project Structure

```plaintext
mediawiki-updater/
├── cmd/                  # Cobra CLI commands
│   ├── root.go            # Main command
│   └── list.go            # List subcommands
├── internal/             # Internal packages
│   ├── config/             # Configuration parsing
│   ├── downloader/        # Download management
│   ├── extractor/         # Archive extraction
│   ├── mediawiki/         # MediaWiki-specific logic
│   └── updater/           # Main update orchestration
├── config.ini        # Default configuration
├── extensions-sample.ini # Example configuration
└── main.go               # Application entry point
```

## 🔧 CLI Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--config` | `-c` | `config.ini` | Path to configuration file |
| `--target` | `-t` | `.` | Target directory for installation |
| `--verbose` | `-v` | `false` | Enable verbose output |

## 🛡️ Preserved Files

The following files/directories are preserved during updates:

- `LocalSettings.php`
- `.htaccess`
- `images/`

## 🏗️ Architecture

The application follows clean architecture principles:

- **Config**: Handles INI file parsing and validation
- **Downloader**: Manages downloads from ExtDist and Git
- **Extractor**: Handles archive extraction and file operations
- **MediaWiki**: Parses official release pages for download URLs
- **Updater**: Orchestrates the entire update process

## 🤝 Contributing

The usual, e. g.:

1. Fork the repository
2. Create a new branch
3. Commit your changes
4. Push the branch, to your forked repository
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## ⚠️ Disclaimer

> [!CAUTION]
> Always backup your existing MediaWiki installation before running this tool. This software comes with no warranty and bugs may occur.

This tool is not officially associated with MediaWiki or the Wikimedia Foundation. Use at your own risk.
