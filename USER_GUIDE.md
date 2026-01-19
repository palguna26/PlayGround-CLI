# PlayGround CLI - User Guide

**Version 0.1.0**

Welcome to PlayGround (`pg`) - your AI-powered coding assistant that keeps you in control. This guide will help you get started and master the tool.

---

## Table of Contents

1. [What is PlayGround?](#what-is-playground)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Commands Reference](#commands-reference)
5. [Workflows & Examples](#workflows--examples)
6. [Best Practices](#best-practices)
7. [Troubleshooting](#troubleshooting)
8. [FAQ](#faq)

---

## What is PlayGround?

PlayGround is a **session-based CLI tool** that lets AI agents help you code while maintaining complete safety and control:

✅ **You're always in control** - AI proposes changes, you review and approve  
✅ **No surprises** - All changes are shown as diffs before applying  
✅ **Safe by design** - AI can't write files directly  
✅ **Session-based** - Track your work with clear goals  
✅ **Local-first** - Everything stays on your machine

### Core Concepts

- **Session**: A work session with a specific goal (e.g., "add JWT auth")
- **Patch**: A proposed code change shown as a unified diff
- **Tool**: Actions the AI can take (read files, check Git status, propose changes)
- **Provider**: The AI service (OpenAI or Gemini)

---

## Installation

### Step 1: Prerequisites

Make sure you have:
- **Go 1.19+** installed ([download](https://go.dev/dl/))
- **Git** installed
- An API key from **OpenAI** or **Google Gemini**

### Step 2: Build PlayGround

```bash
# Clone the repository
git clone https://github.com/yourusername/playground.git
cd playground

# Build the binary
go build -o pg cmd/pg/main.go

# Optional: Install globally
# Linux/macOS:
sudo mv pg /usr/local/bin/
# Windows: Move pg.exe to a folder in your PATH
```

### Step 3: Set Up Your API Key

Choose your preferred AI provider:

**Option A: Use Gemini (Recommended - Fast & Free Tier)**
```bash
export GEMINI_API_KEY="your-gemini-api-key-here"
```

**Option B: Use OpenAI**
```bash
export OPENAI_API_KEY="your-openai-api-key-here"
```

**Option C: Both (Gemini will be preferred)**
```bash
export GEMINI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key"
# To force OpenAI:
export LLM_PROVIDER="openai"
```

### Step 4: Verify Installation

```bash
pg --version
# Output: pg version 0.1.0
```

---

## Quick Start

Let's walk through your first PlayGround session:

### 1. Navigate to Your Project

```bash
cd /path/to/your/project
# Must be a Git repository
```

### 2. Start a Session

```bash
pg start "add input validation to user registration"
```

**Output:**
```
✓ Started new session: pg-1
  Goal: add input validation to user registration
  Repo: /path/to/your/project
```

### 3. Ask the AI for Help

```bash
pg ask "show me the current user registration code"
```

**Output:**
```
Using LLM provider: Gemini
🤖 Agent working...

Agent: I've read the registration code in auth/register.go. 
The function currently doesn't validate email format or password strength.
Would you like me to add validation?
```

### 4. Request Changes

```bash
pg ask "add email validation and password strength checks"
```

**Output:**
```
🤖 Agent working...

Agent: I've proposed changes to add email validation using regex 
and password strength requirements (min 8 chars, 1 uppercase, 1 number).
Review with 'pg review' and apply with 'pg apply'.
```

### 5. Review Proposed Changes

```bash
pg review
```

**Output:**
```
Session: pg-1
Pending patches: 1

═══ Patch 1/1 ═══
File: auth/register.go
Created: 2026-01-19 16:15:30

--- auth/register.go
+++ auth/register.go
@@ -5,6 +5,7 @@
 
 import (
     "database/sql"
+    "regexp"
 )
 
@@ -12,6 +13,25 @@
+func validateEmail(email string) bool {
+    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
+    return emailRegex.MatchString(email)
+}
+
+func validatePassword(password string) error {
+    if len(password) < 8 {
+        return fmt.Errorf("password must be at least 8 characters")
+    }
+    // ... (full validation code)
+}
...
```

### 6. Apply the Changes

```bash
pg apply
```

**Output:**
```
About to apply 1 patch(es) to the repository.
Apply all patches? [y/N]: y
Applying patch 1/1: auth/register.go... ✓

✓ Successfully applied 1 patch(es)
```

### 7. Check Session Status

```bash
pg status
```

**Output:**
```
Session: pg-1
Goal: add input validation to user registration
Repository: /path/to/your/project
Created: 2026-01-19 16:10:00

Pending Patches: 0
Tool History: 4 calls

Recent Tool Calls:
  ✓ read_file - 16:15:20
  ✓ read_file - 16:15:25
  ✓ propose_patch - 16:15:30
```

**Congratulations! You've completed your first PlayGround session.** 🎉

---

## Commands Reference

### `pg start "<goal>"`

**Creates a new coding session**

```bash
pg start "refactor database queries to use prepared statements"
```

- Generates a unique session ID (pg-1, pg-2, etc.)
- Sets the session goal
- Makes it the active session
- Must be run inside a Git repository

### `pg ask "<question or task>"`

**Interact with the AI agent**

```bash
pg ask "what files handle authentication?"
pg ask "add error logging to the API handlers"
pg ask "explain how the cache invalidation works"
```

- Requires an active session
- AI can read files, check Git status, and propose changes
- Shows which LLM provider is being used
- Saves progress automatically

**AI Capabilities:**
- Read any file in the repository
- List directory contents
- Check Git status and diffs
- Propose code changes as patches
- Run commands (requires your approval)

### `pg status`

**Display current session information**

```bash
pg status
```

Shows:
- Session ID and goal
- Repository path
- Creation time
- Number of pending patches
- Recent tool executions
- Context summary (if AI has added one)

### `pg review`

**Review all pending patches before applying**

```bash
pg review
```

- Shows each patch as a unified diff
- Displays file path and creation time
- Lets you inspect exactly what will change
- No patches? Shows "No pending patches"

### `pg apply`

**Apply all pending patches to your code**

```bash
pg apply
```

- Asks for confirmation before applying
- Validates each patch before applying
- Applies atomically (all or nothing per patch)
- Creates backups and rolls back on failure
- Clears pending patches on success

**Interactive prompt:**
```
About to apply 2 patch(es) to the repository.
Apply all patches? [y/N]:
```

### `pg resume <session-id>`

**Resume a previous session**

```bash
pg resume pg-3
```

- Loads session from `.pg/sessions/`
- Validates repository matches
- Sets as active session
- Restores all context and pending patches

To see available sessions, check `.pg/sessions/` in your repo.

---

## Workflows & Examples

### Workflow 1: Adding a New Feature

**Goal:** Add rate limiting to an API

```bash
# Start session
pg start "add rate limiting to API endpoints"

# Explore codebase
pg ask "show me how requests are currently handled"

# Ask for implementation
pg ask "add rate limiting middleware with 100 requests per minute limit"

# Review and iterate
pg review
pg ask "add Redis-based rate limiting instead of in-memory"

# Apply when satisfied
pg review
pg apply
```

### Workflow 2: Bug Investigation & Fix

**Goal:** Fix a memory leak

```bash
# Start session
pg start "investigate and fix memory leak in worker pool"

# Investigation
pg ask "which files manage the worker pool?"
pg ask "show me the worker lifecycle code"
pg ask "are workers being properly closed?"

# Fix proposal
pg ask "add defer close() calls and proper cleanup"

# Review and apply
pg review
pg apply
```

### Workflow 3: Refactoring

**Goal:** Extract common logic into a helper

```bash
pg start "extract database connection code into helper"

# Analysis
pg ask "find all places where we create database connections"

# Refactor
pg ask "create a db package with a NewConnection helper and update all callers"

# Review large changes carefully
pg review  # Check all files being modified

# Apply
pg apply
```

### Workflow 4: Code Review Assistance

**Goal:** Understand unfamiliar code

```bash
pg start "understand authentication flow"

# Exploration (read-only, no patches)
pg ask "explain the authentication flow step by step"
pg ask "what security measures are in place?"
pg ask "trace what happens when login fails"

# Check status (should show 0 pending patches)
pg status
```

### Workflow 5: Testing & Validation

**Goal:** Add test coverage

```bash
pg start "add unit tests for auth package"

# Generate tests
pg ask "create unit tests for the login function"
pg ask "add tests for edge cases like empty email or wrong password"

# Review tests
pg review

# Apply
pg apply

# Run tests (with approval)
pg ask "run go test ./auth/..."
# Prompts: Allow this command? [y/N]:
```

---

## Best Practices

### ✅ DO

1. **Start with clear goals**
   ```bash
   pg start "migrate from MySQL to PostgreSQL"  # ✅ Clear
   # NOT: pg start "database stuff"              # ❌ Vague
   ```

2. **Review before applying**
   - Always run `pg review` before `pg apply`
   - Read the diffs carefully
   - Verify the AI understood your intent

3. **Use descriptive questions**
   ```bash
   pg ask "add JWT authentication to the login endpoint"  # ✅ Specific
   # NOT: pg ask "make it better"                         # ❌ Vague
   ```

4. **Commit after applying patches**
   ```bash
   pg apply
   git add .
   git commit -m "Add JWT auth (via PlayGround pg-5)"
   ```

5. **Keep sessions focused**
   - One goal per session
   - Start a new session for unrelated work

6. **Check status regularly**
   ```bash
   pg status  # See what's been done
   ```

### ❌ DON'T

1. **Don't apply without reviewing**
   ```bash
   pg apply  # Without pg review first - risky!
   ```

2. **Don't ignore context mismatches**
   - If `pg apply` fails due to context mismatch, the file has changed
   - Regenerate patches with a fresh `pg ask`

3. **Don't use for binary files**
   - Patches work on text files only
   - AI can't see or modify binary files

4. **Don't expect perfection**
   - AI-generated code should be reviewed
   - Test changes before committing

5. **Don't work across repos**
   - One session = one repository
   - Start fresh when switching repos

---

## Troubleshooting

### "not a git repository"

**Problem:** Command fails with this error

**Solution:** 
```bash
cd /path/to/your/git/repo
git status  # Verify it's a Git repo
pg start "your goal"
```

### "no active session"

**Problem:** Running `pg ask`, `pg status`, or `pg review` with no session

**Solution:**
```bash
pg start "describe your goal here"
```

Or resume an existing session:
```bash
pg resume pg-3
```

### "no LLM API key found"

**Problem:** No API key configured

**Solution:**
```bash
export GEMINI_API_KEY="your-key-here"
# OR
export OPENAI_API_KEY="your-key-here"
```

Add to your shell profile (~/.bashrc, ~/.zshrc) to persist:
```bash
echo 'export GEMINI_API_KEY="your-key"' >> ~/.zshrc
```

### Patch Application Failed

**Problem:** `pg apply` fails with "context mismatch"

**Cause:** The file has changed since the patch was created

**Solution:**
```bash
# Check what changed
git status
git diff

# Regenerate patches
pg ask "rebase the previous changes on the current code"
pg review
pg apply
```

### Provider Not Working

**Problem:** Want to switch from OpenAI to Gemini (or vice versa)

**Solution:**
```bash
export LLM_PROVIDER="gemini"  # or "openai"
export GEMINI_API_KEY="your-gemini-key"
```

### Session Files Corrupted

**Problem:** Session won't load

**Solution:**
```bash
# Check session files
ls .pg/sessions/

# View session JSON
cat .pg/sessions/pg-1.json

# Start fresh if needed
pg start "new goal"
```

### Agent Not Responding

**Problem:** `pg ask` hangs or times out

**Check:**
1. API key is valid
2. Internet connection is working
3. API service is not experiencing downtime

**Try:**
```bash
# Use the other provider
export LLM_PROVIDER="openai"  # if using Gemini
# OR
export LLM_PROVIDER="gemini"  # if using OpenAI
```

---

## FAQ

### **Q: Is my code sent to OpenAI/Gemini?**

A: Yes, when you use `pg ask`, the AI needs to see your code to help. Only the files the agent reads are sent. If privacy is critical, use a local-only LLM (not supported in v0.1).

### **Q: Can the AI delete files?**

A: No. The AI can only propose changes as patches. Deletions would appear in the diff for your review. You control what gets applied.

### **Q: What happens if my session crashes?**

A: Sessions are saved after every agent action. Resume with `pg resume <session-id>` and your progress is restored.

### **Q: Can I use this for non-code files?**

A: Yes! It works with any text files - markdown, config files, documentation, etc.

### **Q: How do I see all my sessions?**

A: Check the `.pg/sessions/` directory:
```bash
ls .pg/sessions/
# Output: pg-1.json  pg-2.json  pg-3.json
```

### **Q: Can I edit sessions manually?**

A: Yes, they're JSON files. But be careful - invalid JSON will break the session.

### **Q: Does this replace Git?**

A: No! PlayGround works *with* Git. You still commit, push, and manage your repository as usual. PlayGround just helps you generate changes.

### **Q: What's the difference between `pg ask` and just using ChatGPT?**

A: PlayGround:
- Maintains session context
- Can read your actual files
- Proposes changes as reviewable patches
- Enforces safety (no direct writes)
- Integrates with your Git workflow

### **Q: Can multiple people use the same session?**

A: Not recommended. Sessions are local to your machine. Use Git to collaborate.

### **Q: How much does this cost?**

A: That depends on your LLM provider:
- **Gemini**: Generous free tier, Gemini 2.0 Flash is very affordable
- **OpenAI**: Pay-per-token pricing

PlayGround itself is free and open-source.

### **Q: Can I customize the AI's behavior?**

A: Not in v0.1. The system prompt is fixed. Future versions may support customization.

---

## Tips & Tricks

### Tip 1: Use Git Branches for Experiments

```bash
git checkout -b ai-experiment
pg start "try refactoring with AI"
# ... work with pg ...
# If you like it: git merge
# If not: git checkout main && git branch -D ai-experiment
```

### Tip 2: Chain Questions for Context

```bash
pg ask "show me the login function"
pg ask "now add 2FA support to it"  # AI remembers the previous context
```

### Tip 3: Ask for Explanations

```bash
pg ask "explain how the caching works"  # Read-only, no patches
```

### Tip 4: Use Descriptive Session Goals

Your future self will thank you:
```bash
pg start "migrate database schema to v2"  # ✅
# NOT: pg start "fix stuff"                # ❌
```

### Tip 5: Keep a Session Log

```bash
pg status > session-log.txt
git add session-log.txt
git commit -m "Session pg-5 completed"
```

---

## Next Steps

Now that you know the basics, try:

1. **Start a real session** in one of your projects
2. **Experiment with different questions** to see what the AI can do
3. **Review the README** for architecture details
4. **Check the walkthrough** to understand how it works internally

Happy coding with PlayGround! 🚀

---

## Getting Help

- **Issues**: File on GitHub (when available)
- **Questions**: Check this guide first, then the README
- **Contributing**: See CONTRIBUTING.md (follow the safety-first philosophy)

---

**Remember: PlayGround is infrastructure, not magic. You're in control. Always review, always verify, always test.**
