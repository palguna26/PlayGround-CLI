# PlayGround CLI - User Guide

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Agent Mode](#agent-mode)
4. [Commands Reference](#commands-reference)
5. [Configuration](#configuration)
6. [Workflows](#workflows)
7. [Troubleshooting](#troubleshooting)

---

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

### Requirements

- **llama.cpp**: Install from [llama.cpp](https://github.com/ggerganov/llama.cpp)
- **Model**: DeepSeek-Coder-7B-Instruct v1.5 (~4GB)
- **RAM**: 8GB minimum, 16GB recommended
- **Disk**: ~5GB for model

See [SETUP.md](SETUP.md) for detailed installation instructions.

---

## Quick Start

### 1. Configure Local Model

```bash
pg setup
```

Follow the wizard to download or configure your model path.

### 2. Start Agent Mode

```bash
cd your-project
pg agent
```

### 3. Chat Naturally

```
You: Add a REST API endpoint for user registration

Agent: I'll create a user registration endpoint.
First, let me check your project structure...
```

---

## Agent Mode

Agent mode is the primary way to interact with PlayGround. It provides a conversational interface where you can:

- Describe what you want in natural language
- See the AI's plan before execution
- Review all changes as diffs
- Apply or reject modifications

### Starting Agent Mode

```bash
pg agent                    # New session
pg agent --resume pg-5      # Resume previous session
```

### In-Chat Commands

| Command | Action |
|---------|--------|
| `review` | Display all pending patches as diffs |
| `apply` | Apply pending patches (with confirmation) |
| `status` | Show session information |
| `help` | List available commands |
| `exit` | Save session and exit |

### Example Session

```
╔════════════════════════════════════════════════════════════╗
║           PlayGround Agent - Interactive Mode              ║
╚════════════════════════════════════════════════════════════╝

🤖 Model: DeepSeek-Coder-7B-Instruct-v1.5 (local)
📁 Path: ~/.playground/models/deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf

Session: pg-1
Goal: Add authentication

You: Create a config file parser for YAML

Agent: I'll create a YAML config parser. Let me check if you have 
any existing config handling...

[Executing: list_files, read_file]

I see you're using Go. I'll create a config package with:
1. Config struct with common fields
2. Load function for YAML parsing
3. Validation logic

I've prepared the changes. Type 'review' to see them.

You: review

═══ Patch 1/2 ═══
File: config/config.go
[diff output...]

═══ Patch 2/2 ═══
File: config/loader.go
[diff output...]

You: apply

✅ Successfully applied 2 patch(es)
```

---

## Commands Reference

### `pg setup`

Configure local model path.

```bash
pg setup
```

- Downloads or configures model
- Validates model file
- Saves to `~/.playground/config.json`

### `pg agent`

Start interactive agent mode.

```bash
pg agent                    # New session
pg agent --resume <id>      # Resume session
```

### `pg start`

Start a new session with a specific goal.

```bash
pg start "implement user authentication"
```

### `pg ask`

Ask a one-off question without starting a session.

```bash
pg ask "how do I parse JSON in Go?"
```

### `pg review`

Display all pending patches.

```bash
pg review
```

### `pg apply`

Apply all pending patches.

```bash
pg apply
```

### `pg status`

Show current session status.

```bash
pg status
```

### `pg resume`

Resume a previous session.

```bash
pg resume pg-3
```

---

## Configuration

### Config File Location

- **Linux/macOS**: `~/.playground/config.json`
- **Windows**: `%USERPROFILE%\.playground\config.json`

### Config Structure

```json
{
  "model_path": "/path/to/deepseek-coder-7b-instruct-v1.5.Q4_K_M.gguf"
}
```

### Changing Model

Run `pg setup` again to reconfigure:

```bash
pg setup
```

---

## Workflows

### Feature Development

```bash
# 1. Start session with clear goal
pg agent

# 2. Describe the feature
You: Add pagination to the user list API

# 3. Review changes
You: review

# 4. Apply if satisfied
You: apply

# 5. Continue iterating
You: Now add sorting support
```

### Code Review Assistance

```bash
pg agent

You: Review my changes in auth.go for security issues

Agent: I'll analyze auth.go for security vulnerabilities...
[Analysis and suggestions]
```

### Bug Fixing

```bash
pg agent

You: The login endpoint returns 500 when email is missing

Agent: Let me check the login handler...
[Proposes fix with validation]

You: apply
```

---

## Troubleshooting

### "llama-cli: command not found"

**Cause**: llama.cpp not installed or not in PATH.

**Solution**: Install llama.cpp:
```bash
# macOS
brew install llama.cpp

# Linux - see https://github.com/ggerganov/llama.cpp
```

### "No model configured"

**Solution**: Run the setup wizard.

```bash
pg setup
```

### Slow inference (> 10 seconds)

**Solutions**:
1. Close other applications to free RAM
2. Use Q3_K_M quantization for faster inference
3. Ensure 8GB+ RAM available

### Out of memory

**Solutions**:
- Close other applications
- Use Q3_K_M quantization (smaller model)
- Upgrade to 16GB RAM

### Session Not Found

**Solution**: List available sessions and resume the correct one.

```bash
pg status  # Shows current session
```

### Changes Not Applied

**Solution**: Ensure you run `apply` after reviewing.

```bash
pg review   # See pending changes
pg apply    # Apply them
```

---

## Best Practices

1. **Be Specific**: "Add JWT authentication with refresh tokens" > "Add auth"
2. **Review Before Applying**: Always check diffs before applying
3. **Iterate**: Start small, build incrementally
4. **Use Git**: While optional, Git provides better versioning
5. **Session Goals**: Set clear goals when starting sessions

---

## Privacy & Offline Operation

PlayGround runs **100% locally**:

- ✅ No data sent to cloud
- ✅ Works completely offline (after setup)
- ✅ No telemetry or tracking
- ✅ Your code stays on your machine

---

## Getting Help

- **In-agent help**: Type `help` in agent mode
- **Command help**: `pg --help` or `pg <command> --help`
- **Issues**: [GitHub Issues](https://github.com/palguna26/PlayGround-CLI/issues)

---

*PlayGround CLI - AI-assisted development, locally and safely.*
