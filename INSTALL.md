# Installing PlayGround CLI

## One-Line Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

This automatically:
- Detects your OS (Linux, macOS, Windows via WSL)
- Detects your architecture (amd64, arm64)
- Downloads the correct prebuilt binary
- Verifies the checksum
- Installs to `/usr/local/bin` or `~/.local/bin`

## Manual Installation

### Download from GitHub Releases

1. Go to [Releases](https://github.com/palguna26/PlayGround-CLI/releases)
2. Download the appropriate archive for your system:
   - `pg_linux_amd64.tar.gz` - Linux x86_64
   - `pg_linux_arm64.tar.gz` - Linux ARM64
   - `pg_darwin_amd64.tar.gz` - macOS Intel
   - `pg_darwin_arm64.tar.gz` - macOS Apple Silicon
   - `pg_windows_amd64.zip` - Windows x86_64
3. Extract and move to your PATH:

```bash
# Linux/macOS
tar -xzf pg_linux_amd64.tar.gz
sudo mv pg_linux_amd64/pg /usr/local/bin/

# Windows (PowerShell)
Expand-Archive pg_windows_amd64.zip
Move-Item pg_windows_amd64\pg.exe C:\Windows\System32\
```

### Build from Source

Requires Go 1.19+:

```bash
git clone https://github.com/palguna26/PlayGround-CLI.git
cd PlayGround-CLI
go build -o pg ./cmd/pg
sudo mv pg /usr/local/bin/
```

Or using Make:

```bash
make install
```

## Verify Installation

```bash
pg --version
# Output: pg version v0.1.0
```

## Uninstall

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/pg

# If installed to ~/.local/bin
rm ~/.local/bin/pg
```

## Troubleshooting

### "pg: command not found"

Your PATH may not include the install directory:

```bash
# Add to ~/.bashrc or ~/.zshrc
export PATH="$PATH:$HOME/.local/bin"
```

### Permission denied

Run with sudo:

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sudo sh
```

### Checksum verification failed

Try downloading again or check network connectivity:

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```
