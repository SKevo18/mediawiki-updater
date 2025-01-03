# MediaWiki Updater & Downloader

This Go application automates the process of updating or downloading and setting up MediaWiki core, extensions, and skins. It's designed to simplify the installation process for MediaWiki, making it easier to get a new wiki up and running quickly.

## Features

- Downloads MediaWiki core (version 1.43 by default);
- Downloads specified extensions and skins;
- Uses a temporary directory for downloads to prevent overwriting existing files;
- Ignores specified paths when copying files to the target directory;

## Prerequisites

- Go 1.16 or higher (latest version is recommended)
- Git (for cloning the repository)

> [!CAUTION]
> Backup your existing installation (this tool comes with no warranty, and bugs may happen - please, ALWAYS make backups)!

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/SKevo18/mediawiki-updater.git
   cd mediawiki-updater
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

1. Create two text files in the same directory as the main.go file:
   - `extensions.txt`: List the names of desired extensions, one per line.
   - `skins.txt`: List the names of desired skins, one per line.

   If these files do not exist, nothing happens as the tool will ignore that and move on.

2. Run the application:

   ```bash
   go run main.go [target_directory]
   ```

   If no target directory is specified, the current directory will be used.

## Configuration

You can modify the following constants in the `main.go` file to customize the download sources:

- `versionTag`: The MediaWiki version tag (default: "REL1_42")
- `mediaWikiURL`: The URL for the MediaWiki core tarball
- `extensionURL`: The base URL for extensions
- `skinURL`: The base URL for skins

## Ignored Paths

The following paths are ignored when copying files to the target directory:

- `LocalSettings.php`
- `.htaccess`
- `images`

This ensures that existing configurations and user-uploaded content are not overwritten during the update process.

## License

MIT

## Disclaimer

This tool is not officially associated with MediaWiki or the Wikimedia Foundation. Use at your own risk.
