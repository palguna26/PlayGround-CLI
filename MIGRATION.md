# Migration to PlayGround v2.0 (Local-First)

## Breaking Changes

PlayGround v2.0 is a **fundamental architectural pivot** from cloud-based LLMs to fully local, offline-first inference.

### What Changed

❌ **Removed:**
- OpenAI API integration
- Gemini API integration
- All API key configuration
- Cloud provider auto-detection
- Environment variables: `OPENAI_API_KEY`, `GEMINI_API_KEY`, `LLM_PROVIDER`

✅ **Added:**
- Local DeepSeek-Coder-7B-Instruct v1.5 inference
- llama.cpp integration
- Model management system
- Offline-first operation
- Enhanced privacy (100% local)

### Why This Change?

**Privacy**: Your code never leaves your machine  
**Cost**: No API fees, free forever  
**Speed**: No network latency after model loads  
**Reliability**: Works offline, no rate limits  
**Control**: You own the entire stack  

## Migration Steps

### 1. Update PlayGround

```bash
curl -fsSL https://raw.githubusercontent.com/palguna26/PlayGround-CLI/main/scripts/install.sh | sh
```

### 2. Install llama.cpp

**macOS:**
```bash
brew install llama.cpp
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt install llama.cpp

# Or build from source:
git clone https://github.com/ggerganov/llama.cpp
cd llama.cpp
make
sudo cp llama-cli /usr/local/bin/
```

**Windows:**
Download from [llama.cpp releases](https://github.com/ggerganov/llama.cpp/releases) and add to PATH.

### 3. Download Local Model

```bash
pg setup
```

This will:
- Prompt for model path
- Optionally download DeepSeek-Coder-7B-Instruct v1.5 (~4GB)
- Validate the model
- Save configuration

### 4. Remove Old Config (Optional)

```bash
rm ~/.playground/config.json
```

Then run `pg setup` again to create a fresh config.

## What Stays the Same

✅ **All safety guarantees remain:**
- Diff-only changes
- Explicit user approval required
- No auto-apply
- Deterministic behavior

✅ **Same commands:**
- `pg agent` - Interactive mode
- `pg start "goal"` - Start session
- `pg ask "question"` - One-off question
- `pg review` - Show diffs
- `pg apply` - Apply changes

✅ **Same workflow:**
- Agent reads files
- Proposes unified diffs
- User reviews
- User applies

## New Requirements

| Requirement | v1.x (Cloud) | v2.0 (Local) |
|-------------|--------------|--------------|
| Internet | Always | Only for setup |
| API Key | Required | Not needed |
| RAM | Any | 8GB minimum |
| Disk Space | Minimal | ~5GB for model |
| llama.cpp | Not needed | Required |

## Troubleshooting

### "No model configured"

Run `pg setup` to configure your local model path.

### "llama-cli: command not found"

Install llama.cpp (see step 2 above).

### Slow performance

- Ensure 8GB+ RAM available
- Close other applications
- Consider Q3_K_M quantization for faster inference

### Out of memory

- Use Q3_K_M instead of Q4_K_M quantization
- Reduce context size
- Close other applications

## FAQ

**Q: Can I still use OpenAI/Gemini?**  
A: No. v2.0 is local-only by design. This is a deliberate choice for privacy, cost, and reliability.

**Q: Is the local model as good as GPT-4?**  
A: DeepSeek-Coder-7B is optimized for coding tasks and performs very well for code generation, refactoring, and analysis. It's not as general-purpose as GPT-4, but it's excellent for its specific domain.

**Q: Can I use a different local model?**  
A: Currently, PlayGround is optimized for DeepSeek-Coder-7B-Instruct v1.5. Support for other models may be added in future versions.

**Q: What about my existing sessions?**  
A: Sessions from v1.x will continue to work. The session format hasn't changed.

**Q: Can I go back to v1.x?**  
A: Yes, you can install an older version from GitHub releases. However, we recommend embracing the local-first approach for better privacy and cost savings.

## Benefits of Local-First

### Privacy
- Code never leaves your machine
- No telemetry
- No cloud logging
- Complete data sovereignty

### Cost
- No API fees
- No rate limits
- No usage caps
- Free forever

### Reliability
- Works offline
- No network dependency
- No service outages
- Deterministic performance

### Speed
- No network latency (after model loads)
- Consistent response times
- No throttling

---

**Welcome to PlayGround v2.0 - Local, Safe, and Free.**
