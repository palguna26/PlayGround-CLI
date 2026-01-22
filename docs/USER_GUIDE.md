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

### One-Line Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

### Windows (PowerShell)

```powershell
# After running the curl installer above, add to PATH:
$env:Path += ";$env:USERPROFILE\.local\bin"

# Make permanent:
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$env:USERPROFILE\.local\bin", "User")
```

### Verify Installation

```bash
pg --version
# Output: pg version 0.1.0
```

---

## Quick Start

### 1. Configure API Key

```bash
pg setup
```

Follow the wizard to enter your Gemini or OpenAI API key.

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

| Command | Description |
|---------|-------------|
| `review` | Display all pending patches as diffs |
| `apply` | Apply pending patches (with confirmation) |
| `status` | Show session information |
| `help` | List available commands |
| `exit` | Save session and exit |

### Example Session

```
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

Interactive configuration wizard.

```bash
pg setup
```

- Configures API keys for Gemini/OpenAI
- Sets preferred provider
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
  "gemini_api_key": "AIza...",
  "openai_api_key": "sk-...",
  "llm_provider": "gemini"
}
```

### Environment Variables

Environment variables override config file settings:

```bash
export GEMINI_API_KEY="your-key"
export OPENAI_API_KEY="your-key"
export LLM_PROVIDER="gemini"  # or "openai"
```

### Provider Priority

1. `LLM_PROVIDER` environment variable (if set)
2. Config file `llm_provider` setting
3. Auto-detect: Gemini preferred if both keys present

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

### Rate Limit Errors (429)

**Cause**: API rate limit exceeded.

**Solutions**:
1. Wait 1-2 minutes and retry
2. Switch to a different provider
3. Upgrade to a paid API tier

```bash
# Switch provider temporarily
export LLM_PROVIDER="openai"
pg agent
```

### "No LLM API key found"

**Solution**: Run the setup wizard.

```bash
pg setup
```

### "pg: command not found"

**Solution**: Add install directory to PATH.

```bash
# Linux/macOS
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
source ~/.bashrc

# Windows PowerShell
[Environment]::SetEnvironmentVariable("Path", "$env:Path;$env:USERPROFILE\.local\bin", "User")
# Restart terminal
```

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

## Getting Help

- **In-agent help**: Type `help` in agent mode
- **Command help**: `pg --help` or `pg <command> --help`
- **Issues**: [GitHub Issues](https://github.com/palguna26/PlayGround-CLI/issues)

---

*PlayGround CLI - AI-assisted development, safely.*
