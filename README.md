# PlayGround CLI

**AI-Powered Coding Assistant with Safety Guarantees**

PlayGround is a CLI tool that brings AI-assisted development to your terminal with a critical difference: **you stay in control**. Unlike other AI coding tools, PlayGround never applies changes automatically. Every modification is shown as a diff and requires your explicit approval.

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

Or download from [Releases](https://github.com/palguna26/PlayGround-CLI/releases).

## Features

- 🤖 **Interactive Agent Mode** - Chat naturally with AI (like Claude Code)
- 🔒 **Safe by Design** - All changes shown as diffs, never auto-applied
- 🔄 **Multi-Provider** - Supports OpenAI GPT-4 and Google Gemini
- 📂 **Git Optional** - Works with or without Git (uses snapshots)
- ⚡ **Streaming Responses** - See AI thinking in real-time
- 🔁 **Session Resumption** - Pick up where you left off

## Quick Start

```bash
# 1. Configure your API key
pg setup

# 2. Start interactive agent mode
pg agent

# 3. Chat naturally!
You: Add a login endpoint with JWT authentication
Agent: I'll create a JWT-based login system...
```

## Commands

| Command | Description |
|---------|-------------|
| `pg setup` | Interactive API key configuration |
| `pg agent` | Start interactive chat mode |
| `pg start "goal"` | Start a new session with a goal |
| `pg ask "question"` | Ask a one-off question |
| `pg review` | Show pending changes as diffs |
| `pg apply` | Apply approved changes |
| `pg status` | Show current session status |
| `pg resume <id>` | Resume a previous session |

## Agent Mode

The flagship feature - an interactive chat interface:

```bash
pg agent
```

```
╔════════════════════════════════════════════════════════════╗
║           PlayGround Agent - Interactive Mode              ║
╚════════════════════════════════════════════════════════════╝

Session: pg-1
Goal: Add authentication

You: Create a user model with email and password

Agent: I'll create a user model with bcrypt password hashing.
Let me first check your project structure...

[Reading files...]

I've prepared the changes. Type 'review' to see them.

You: review

═══ Patch 1/1 ═══
File: models/user.go

--- /dev/null
+++ models/user.go
@@ -0,0 +1,25 @@
+package models
+
+import "golang.org/x/crypto/bcrypt"
+
+type User struct {
+    ID       int
+    Email    string
+    Password string
+}
...

You: apply

✅ Successfully applied 1 patch(es)
```

### In-Chat Commands

| Command | Action |
|---------|--------|
| `review` | Show pending patches |
| `apply` | Apply patches |
| `status` | Session info |
| `help` | Show commands |
| `exit` | Save and quit |

## Configuration

### Setup Wizard (Recommended)

```bash
pg setup
```

This guides you through:
- Choosing your LLM provider (Gemini or OpenAI)
- Entering your API key
- Saving configuration securely

Config is stored in `~/.playground/config.json`.

### Environment Variables (Alternative)

```bash
# Gemini (recommended - generous free tier)
export GEMINI_API_KEY="your-key"

# OpenAI
export OPENAI_API_KEY="your-key"

# Force specific provider
export LLM_PROVIDER="gemini"  # or "openai"
```

## Safety Guarantees

PlayGround is built with safety as the core principle:

| Feature | PlayGround | Other AI Tools |
|---------|------------|----------------|
| Auto-apply changes | ❌ Never | ✅ Yes |
| Show diffs before apply | ✅ Always | ❌ Sometimes |
| Require approval | ✅ Always | ❌ No |
| Rollback support | ✅ Full | ❌ Limited |
| Git optional | ✅ Yes | ❌ Usually required |

## Architecture

```
┌─────────────────────────────────────────┐
│              User Interface             │
│   (pg agent / pg ask / pg start)        │
└───────────────────┬─────────────────────┘
                    │
┌───────────────────▼─────────────────────┐
│              Agent Core                  │
│   • System Prompts                       │
│   • Tool Definitions                     │
│   • Session Management                   │
└───────────────────┬─────────────────────┘
                    │
        ┌───────────┴───────────┐
        ▼                       ▼
┌──────────────┐        ┌──────────────┐
│  LLM Layer   │        │    Tools     │
│  • OpenAI    │        │  • read_file │
│  • Gemini    │        │  • list_files│
│  • Streaming │        │  • git_*     │
└──────────────┘        │  • patches   │
                        └──────────────┘
                              │
┌─────────────────────────────▼───────────┐
│            Workspace Layer              │
│   • Git Mode (commits, stash)           │
│   • Snapshot Mode (SHA-based)           │
└─────────────────────────────────────────┘
```

## Requirements

- **API Key**: Gemini (free) or OpenAI (paid)
- **OS**: Linux, macOS, Windows
- **Git**: Optional (uses snapshots if not available)

## Building from Source

```bash
git clone https://github.com/palguna26/PlayGround-CLI.git
cd PlayGround-CLI
go build -o pg ./cmd/pg
sudo mv pg /usr/local/bin/
```

## Troubleshooting

### "Rate limit error (429)"

You've hit the API rate limit. Options:
- Wait 1-2 minutes and retry
- Switch providers: `pg setup`
- Use a paid API tier

### "pg: command not found"

Add the install directory to your PATH:

```bash
# Linux/macOS
export PATH="$PATH:$HOME/.local/bin"

# Windows PowerShell
$env:Path += ";$env:USERPROFILE\.local\bin"
```

### "No LLM API key found"

Run the setup wizard:
```bash
pg setup
```

## License

MIT License - see [LICENSE](LICENSE)

## Contributing

Contributions welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

**PlayGround CLI** - AI-assisted development, safely.

*For more help, see the [User Guide](docs/USER_GUIDE.md) or [open an issue](https://github.com/palguna26/PlayGround-CLI/issues).*
