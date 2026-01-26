# PlayGround CLI

**Local-First AI Coding Assistant with Safety Guarantees**

PlayGround is a CLI tool that brings AI-assisted development to your terminal with a critical difference: **fully local and offline**. No API keys, no cloud dependencies, just you and your code.

## Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

Then set up your local model:

```bash
pg setup
```

## Features

- 🤖 **Fully Local** - Runs DeepSeek-Coder-7B-Instruct v1.5 on your machine
- 🔒 **Safe by Design** - All changes shown as diffs, never auto-applied
- 📴 **Offline-First** - Works without internet after initial setup
- 💬 **Interactive Agent Mode** - Chat naturally with AI (like Claude Code)
- 📂 **Git Optional** - Works with or without Git (uses snapshots)
- ⚡ **Streaming Responses** - See AI thinking in real-time
- 🔁 **Session Resumption** - Pick up where you left off

## Quick Start

```bash
# 1. Install PlayGround
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh

# 2. Setup (downloads model automatically)
pg setup

# 3. Start interactive agent mode
pg agent

# 4. Chat naturally!
You: Add a login endpoint with JWT authentication
Agent: I'll create a JWT-based login system...
```

## Requirements

- **RAM**: 8GB minimum, 16GB recommended
- **Disk**: ~5GB for model
- **CPU**: Any modern CPU (GPU not required)
- **OS**: Linux, macOS, Windows
- **llama.cpp**: Installed and in PATH (see [Setup Guide](docs/SETUP.md))

## Commands

| Command | Description |
|---------|-------------|
| `pg setup` | Configure local model path |
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

🤖 Model: DeepSeek-Coder-7B-Instruct-v1.5 (local)
📁 Path: ~/.playground/models/deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf

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

## Why Local-First?

| Feature | PlayGround (Local) | Cloud AI Tools |
|---------|-------------------|----------------|
| API keys required | ❌ Never | ✅ Always |
| Internet required | ❌ Only for setup | ✅ Always |
| Data privacy | ✅ 100% local | ❌ Sent to cloud |
| Cost | ✅ Free forever | ❌ Pay per token |
| Speed (after load) | ✅ Fast | ⚠️ Network dependent |
| Auto-apply changes | ❌ Never | ⚠️ Sometimes |

## Safety Guarantees

PlayGround is built with safety as the core principle:

- ❌ **Never** writes files directly
- ✅ **Always** shows diffs before applying
- ✅ **Always** requires explicit user approval
- ✅ **Always** deterministic and reviewable
- ✅ **Full** rollback support
- ✅ **Works** with or without Git

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
│  Local LLM   │        │    Tools     │
│  DeepSeek    │        │  • read_file │
│  Coder 7B    │        │  • list_files│
│  (llama.cpp) │        │  • git_*     │
└──────────────┘        │  • patches   │
                        └──────────────┘
                              │
┌─────────────────────────────▼───────────┐
│            Workspace Layer              │
│   • Git Mode (commits, stash)           │
│   • Snapshot Mode (SHA-based)           │
└─────────────────────────────────────────┘
```

## Building from Source

```bash
git clone https://github.com/palguna26/PlayGround-CLI.git
cd PlayGround-CLI
go build -o pg ./cmd/pg
sudo mv pg /usr/local/bin/
```

## Troubleshooting

### "No model configured"

Run the setup wizard:
```bash
pg setup
```

### "llama-cli: command not found"

Install llama.cpp:
```bash
# macOS
brew install llama.cpp

# Linux
# See: https://github.com/ggerganov/llama.cpp

# Windows
# Download from: https://github.com/ggerganov/llama.cpp/releases
```

### Slow inference

- Ensure you have 8GB+ RAM
- Close other applications
- Consider using Q3_K_M quantization for faster inference

### Out of memory

- Close other applications
- Use a smaller quantization (Q3_K_M instead of Q4_K_M)
- Reduce context size in config

---

## Migration from v1.x

See [MIGRATION.md](MIGRATION.md) for upgrading from cloud-based versions.

## License

MIT License - see [LICENSE](LICENSE)

## Contributing

Contributions welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

---

**PlayGround CLI** - AI-assisted development, locally and safely.

*For more help, see the [User Guide](docs/USER_GUIDE.md) or [open an issue](https://github.com/palguna26/PlayGround-CLI/issues).*
