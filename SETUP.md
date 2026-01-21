# Quick Setup Guide

## For Windows Users

Great news! You don't need to mess with environment variables anymore. Just use the setup wizard:

```powershell
# Run the interactive setup
pg setup
```

This will guide you through:
1. Choosing your LLM provider (Gemini or OpenAI)
2. Entering your API key
3. Saving it securely

The setup stores your configuration in `C:\Users\YourName\.playground\config.json` so it persists across sessions.

## First Time Setup (from scratch)

If you just installed PlayGround:

```powershell
# 1. Navigate to your project
cd E:\Projects\ExpenseTracker

# 2. Run setup wizard
pg setup
# Follow the prompts to enter your Gemini or OpenAI API key

# 3. Start coding!
pg start "your goal here"
pg ask "your question"
```

## Manual Setup (Alternative)

If you prefer environment variables:

### PowerShell (Windows)
```powershell
# Temporary  (current session only)
$env:GEMINI_API_KEY = "your-api-key-here"

# Permanent (persists across sessions)
[System.Environment]::SetEnvironmentVariable('GEMINI_API_KEY', 'your-api-key-here', 'User')
```

### Bash/Zsh (Linux/Mac)
```bash
# Add to ~/.bashrc or ~/.zshrc
export GEMINI_API_KEY="your-api-key-here"
```

## What `pg setup` Does

- ✅ Guides you through provider selection
- ✅ Securely stores API keys in config file
- ✅ Sets provider preference
- ✅ Shows current configuration
- ✅ Works cross-platform (Windows, Mac, Linux)

## Config File Location

Your configuration is stored in:
- **Windows**: `C:\Users\YourName\.playground\config.json`
- **Mac/Linux**: `~/.playground/config.json`

The file is created with restricted permissions (user-only access) for security.

## Reconfiguring

Already ran setup but want to change settings?

```powershell
pg setup  # Run again to update configuration
```

It will show your current settings and let you update them.

## Troubleshooting

**"pg: command not found"**
- Make sure `pg.exe` is in your PATH (see main README)
- Or use `.\pg.exe` from the PlayGround-CLI directory

**"No LLM API key found"**
- Run `pg setup` to configure
- Or set environment variable manually

**Want to use different provider temporarily?**
```powershell
# Override with environment variable (doesn't change config)
$env:LLM_PROVIDER = "openai"
pg ask "your question"
```
