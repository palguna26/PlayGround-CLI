# PlayGround CLI - Setup Guide

## Prerequisites

Before setting up PlayGround, ensure you have:

- **8GB RAM minimum** (16GB recommended)
- **~5GB free disk space** for the model
- **llama.cpp** installed and in PATH
- **Operating System**: Linux, macOS, or Windows

---

## Step 1: Install llama.cpp

### macOS

```bash
brew install llama.cpp
```

### Linux (Ubuntu/Debian)

```bash
# Option 1: Package manager (if available)
sudo apt install llama.cpp

# Option 2: Build from source
git clone https://github.com/ggerganov/llama.cpp
cd llama.cpp
make
sudo cp llama-cli /usr/local/bin/
```

### Windows

1. Download from [llama.cpp releases](https://github.com/ggerganov/llama.cpp/releases)
2. Extract to a directory (e.g., `C:\llama.cpp`)
3. Add to PATH:
   ```powershell
   $env:Path += ";C:\llama.cpp"
   # Make permanent:
   [Environment]::SetEnvironmentVariable("Path", "$env:Path;C:\llama.cpp", "User")
   ```

### Verify Installation

```bash
llama-cli --version
# Should output version information
```

---

## Step 2: Install PlayGround CLI

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

Or download from [GitHub Releases](https://github.com/palguna26/PlayGround-CLI/releases).

---

## Step 3: Download the Model

### Option 1: Automatic Download (Recommended)

```bash
pg setup
```

Follow the prompts:
1. Choose to download the model automatically
2. Wait for download (~4GB, may take 10-30 minutes)
3. Model will be saved to `~/.playground/models/`

### Option 2: Manual Download

1. Visit [TheBloke/deepseek-coder-7B-instruct-v1.5-GGUF](https://huggingface.co/TheBloke/deepseek-coder-7B-instruct-v1.5-GGUF)
2. Download `deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf` (~4GB)
3. Save to `~/.playground/models/` (create directory if needed)
4. Run `pg setup` and enter the path

### Model Variants

| Quantization | Size | RAM Required | Speed | Quality |
|--------------|------|--------------|-------|---------|
| Q3_K_M | ~3GB | 6GB | Faster | Good |
| **Q4_K_M** | **~4GB** | **8GB** | **Balanced** | **Recommended** |
| Q5_K_M | ~5GB | 10GB | Slower | Better |

**Recommended**: Q4_K_M for best balance of speed and quality.

---

## Step 4: Configure PlayGround

Run the setup wizard:

```bash
pg setup
```

You'll be prompted for:
- **Model path**: Path to your GGUF model file
- The wizard will validate the model and save the configuration

Configuration is saved to `~/.playground/config.json`.

---

## Step 5: Verify Setup

Test that everything works:

```bash
pg agent
```

You should see:
```
╔════════════════════════════════════════════════════════════╗
║           PlayGround Agent - Interactive Mode              ║
╚════════════════════════════════════════════════════════════╝

🤖 Model: DeepSeek-Coder-7B-Instruct-v1.5 (local)
📁 Path: ~/.playground/models/deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf

Session: pg-1
Goal: Interactive coding session

You: 
```

Type `exit` to quit.

---

## Troubleshooting

### "llama-cli: command not found"

llama.cpp is not installed or not in PATH.

**Solution**: Install llama.cpp (see Step 1) and ensure it's in your PATH.

### "No model configured"

You haven't run `pg setup` yet.

**Solution**: Run `pg setup` and configure your model path.

### "Model file does not exist"

The model path in your config is incorrect.

**Solution**: 
1. Check the path: `cat ~/.playground/config.json`
2. Verify the file exists: `ls -lh ~/.playground/models/`
3. Run `pg setup` again to reconfigure

### Slow inference (> 10 seconds per response)

Your system may not have enough RAM or CPU is slow.

**Solutions**:
- Close other applications to free up RAM
- Use Q3_K_M quantization for faster inference
- Ensure you have 8GB+ RAM available

### Out of memory errors

The model is too large for your available RAM.

**Solutions**:
- Close other applications
- Use Q3_K_M quantization (smaller, faster)
- Upgrade to 16GB RAM if possible

### Model download fails

Network issue or HuggingFace is down.

**Solutions**:
- Try again later
- Download manually from HuggingFace
- Use a different mirror if available

---

## Advanced Configuration

### Custom Model Path

Edit `~/.playground/config.json`:

```json
{
  "model_path": "/custom/path/to/model.gguf"
}
```

### Multiple Models

You can switch models by running `pg setup` again and entering a different path.

### Model Storage Location

Default: `~/.playground/models/`

To use a different location:
1. Move your models to the new location
2. Run `pg setup`
3. Enter the new path

---

## System Requirements Summary

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| RAM | 8GB | 16GB |
| Disk Space | 5GB | 10GB |
| CPU | Any modern CPU | Multi-core |
| Internet | Only for setup | Not required |
| llama.cpp | Latest | Latest |

---

## Next Steps

Once setup is complete:

1. **Start agent mode**: `pg agent`
2. **Read the user guide**: [USER_GUIDE.md](USER_GUIDE.md)
3. **Try a simple task**: Ask the agent to read a file or explain code

**You're ready to code with PlayGround!** 🚀
