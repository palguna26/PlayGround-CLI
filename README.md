# PlayGround CLI (`pg`)

**A session-based, local-first CLI that allows AI agents to safely modify code through diff-only, reviewable patches.**

PlayGround is developer infrastructure, not a product UI. It's an open runtime for coding agents that prioritizes **correctness**, **safety**, and **long-term maintainability**.

## Core Principles

✅ **Diff-only writes** — The agent can NEVER write files directly  
✅ **Session-based** — Users explicitly start sessions with clear goals  
✅ **Reviewable patches** — All changes must be reviewed before applying  
✅ **Local-first** — No telemetry, no background daemons  
✅ **CLI-first** — No TUI, each command does one thing and exits

## Installation

### Prerequisites

- Go 1.19 or later
- Git
- OpenAI API key (for AI features)

### Build from source

```bash
git clone https://github.com/yourusername/playground.git
cd playground
go build -o pg cmd/pg/main.go
```

### Install globally

```bash
# Linux/macOS
sudo mv pg /usr/local/bin/

# Windows
# Move pg.exe to a directory in your PATH
```

## Setup

Set your OpenAI API key:

```bash
export OPENAI_API_KEY="sk-abcdef1234567890abcdef1234567890abcdef12"
```

## Usage

### 1. Start a session

```bash
pg start "add jwt auth to api"
```

This creates a new session with a unique ID and stores it in `.pg/sessions/`.

### 2. Ask the agent questions or give it tasks

```bash
pg ask "what authentication do we currently use?"
pg ask "add jwt middleware to the auth package"
```

The agent will:
- Read files using `read_file` and `list_files` tools
- Check repository state with `git_status` and `git_diff`
- Propose code changes as unified diffs via `propose_patch`

### 3. Check session status

```bash
pg status
```

Shows:
- Current session ID and goal
- Number of pending patches
- Tool execution history

### 4. Review proposed patches

```bash
pg review
```

Displays all pending patches as unified diffs for manual review.

### 5. Apply patches

```bash
pg apply
```

Applies all validated patches after user confirmation. Each patch is:
- Validated before application
- Applied atomically with backup/rollback
- Guaranteed to match current file state

### 6. Resume a previous session

```bash
pg resume pg-12
```

## Architecture

```
playground/
├── cmd/pg/           # CLI entry point
├── internal/
│   ├── cli/          # Command implementations
│   ├── session/      # Session persistence
│   ├── tools/        # Agent tools (read, list, git, run, propose_patch)
│   ├── patch/        # Patch validation & application
│   ├── llm/          # LLM provider abstraction (OpenAI)
│   └── agent/        # Agent runtime & loop
└── go.mod
```

## Safety Guarantees

### 1. No Direct File Writes
The agent can only propose changes as unified diffs. The patch engine is the ONLY code that modifies files.

### 2. Hard Stop Conditions
The agent loop terminates after:
- Max iterations (default: 10)
- LLM signals completion
- No tool calls in response
- Error encountered

### 3. Atomic Operations
Patch application is atomic:
- Backup created before modification
- Validation performed upfront
- Rollback on any failure

### 4. Session Persistence
Sessions are saved after every agent iteration, ensuring no data loss even if the process crashes.

## Agent Tools

| Tool | Description | Safety |
|------|-------------|--------|
| `read_file(path)` | Read file contents | Path validation, repo bounds check |
| `list_files(path)` | List directory contents | Path validation, hidden file filtering |
| `git_status()` | Get Git status | Read-only |
| `git_diff()` | Get uncommitted changes | Read-only |
| `run_command(cmd)` | Execute shell command | **Requires user approval** |
| `propose_patch(file_path, unified_diff)` | Propose code change | Stored only, not applied |

## Example Workflow

```bash
# Start a new session
$ pg start "refactor authentication to use JWT"
✓ Started new session: pg-1
  Goal: refactor authentication to use JWT
  Repo: /home/user/myproject

# Ask agent to analyze current auth
$ pg ask "what authentication method do we use now?"
🤖 Agent working...

Agent: I've analyzed the codebase. Currently using session-based 
authentication in auth/session.go with cookies.

# Ask agent to propose JWT implementation
$ pg ask "add JWT authentication"
🤖 Agent working...

Agent: I've proposed changes to implement JWT authentication.
Review with 'pg review' and apply with 'pg apply'.

# Review proposed changes
$ pg review
Session: pg-1
Pending patches: 2

═══ Patch 1/2 ═══
File: auth/jwt.go
Created: 2026-01-18 17:45:32

--- /dev/null
+++ auth/jwt.go
...

═══ Patch 2/2 ═══
File: main.go
Created: 2026-01-18 17:45:35

--- main.go
+++ main.go
...

# Apply patches
$ pg apply
About to apply 2 patch(es) to the repository.
Apply all patches? [y/N]: y
Applying patch 1/2: auth/jwt.go... ✓
Applying patch 2/2: main.go... ✓

✓ Successfully applied 2 patch(es)
```

## Session Storage

Sessions are stored in `.pg/sessions/` as JSON files:

```json
{
  "id": "pg-1",
  "repo": "/absolute/path/to/repo",
  "goal": "add jwt auth to api",
  "context_summary": "",
  "pending_patches": [],
  "tool_history": [],
  "created_at": "2026-01-18T17:30:00Z"
}
```

The active session is tracked in `.pg/active`.

## Limitations (v0.1)

- Single agent only
- Supports OpenAI and Google Gemini providers
- No plugins or MCP integration
- No LSP or IDE integration
- No auto-apply (by design)
- English only (agent prompts)

## Troubleshooting

**"not a git repository"**  
→ Run `pg` commands from within a Git repository.

**"no active session"**  
→ Start a session with `pg start "<goal>"`.

**"OPENAI_API_KEY environment variable not set"**  
→ Set your API key: `export OPENAI_API_KEY="sk-..."`

**Patch application failed**  
→ File may have changed since patch was proposed. Check `git status` and regenerate patches.

## Contributing

PlayGround is infrastructure. Contributions should prioritize:
1. **Correctness** over features
2. **Safety** over convenience  
3. **Simplicity** over cleverness

## License

MIT

## Philosophy

> "Treat PlayGround as infrastructure, not a product UI."

PlayGround is designed to be:
- **Trustworthy** — Never modifies files without review
- **Transparent** — Every action is logged and auditable
- **Predictable** — Deterministic behavior, no surprises
- **Respectful** — You control when and how code changes

It is NOT designed to be:
- Flashy
- Automatic
- A replacement for your IDE
- A background service
