# Installing PlayGround CLI

## One-Line Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

This automatically:
- Detects your OS (Linux, macOS, Windows)
- Detects your architecture (amd64, arm64)
- Downloads the correct prebuilt binary
- Verifies SHA256 checksum
- Installs to `/usr/local/bin` or `~/.local/bin`

## Platform-Specific Instructions

### Linux / macOS

```bash
# Install
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh

# Add to PATH if needed
export PATH="$PATH:$HOME/.local/bin"

# Verify
pg --version
```

### Windows

**Option 1: Git Bash / WSL**
```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

**Option 2: PowerShell (after curl install)**
```powershell
# Add to PATH for current session
$env:Path += ";$env:USERPROFILE\.local\bin"

# Add to PATH permanently
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$env:USERPROFILE\.local\bin", "User")

# Restart PowerShell, then verify
pg --version
```

## Manual Installation

### From GitHub Releases

1. Go to [Releases](https://github.com/palguna26/PlayGround-CLI/releases)
2. Download the appropriate file:
   - `pg_linux_amd64.tar.gz` - Linux x86_64
   - `pg_linux_arm64.tar.gz` - Linux ARM64
   - `pg_darwin_amd64.tar.gz` - macOS Intel
   - `pg_darwin_arm64.tar.gz` - macOS Apple Silicon
   - `pg_windows_amd64.zip` - Windows x86_64

3. Extract and install:

```bash
# Linux/macOS
tar -xzf pg_linux_amd64.tar.gz
sudo mv pg_linux_amd64/pg /usr/local/bin/

# Windows PowerShell
Expand-Archive pg_windows_amd64.zip -DestinationPath .
Move-Item pg_windows_amd64\pg.exe $env:USERPROFILE\bin\
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

## Post-Installation

### Configure API Key

```bash
pg setup
```

This guides you through setting up your Gemini or OpenAI API key.

### Verify Everything Works

```bash
pg --version      # Check version
pg --help         # See all commands
pg agent          # Start interactive mode
```

## Uninstall

```bash
# If installed to /usr/local/bin
sudo rm /usr/local/bin/pg

# If installed to ~/.local/bin
rm ~/.local/bin/pg

# Remove config (optional)
rm -rf ~/.playground
```

## Troubleshooting

### "pg: command not found"

Add the install directory to your PATH:

```bash
# Check where pg was installed
ls ~/.local/bin/pg

# Add to PATH
export PATH="$PATH:$HOME/.local/bin"

# Make permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
```

### Checksum Verification Failed

Network issue or corrupted download. Try again:

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

### Permission Denied

Run with sudo:

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sudo sh
```

---

*For more help, see the [User Guide](USER_GUIDE.md) or [open an issue](https://github.com/palguna26/PlayGround-CLI/issues).*
