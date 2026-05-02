# PVF-MCP Installation Guide for Agents

This document is written for AI agents tasked with installing PVF-MCP. It covers all supported platforms and fallback strategies.

## Architecture

```
AI Agent Tool (Claude Code/Cursor/Codex/Other)
    │ MCP stdio protocol
    ▼
pvf-mcp.exe (Go binary, 11MB, zero dependencies)
    │ HTTP REST
    ▼
pvfUtility (localhost:27000)
```

The `pvf-mcp` binary is an MCP stdio server. The AI tool spawns it as a subprocess and communicates via stdin/stdout JSON-RPC.

## Step 1: Verify Prerequisites

Before installing PVF-MCP, check:
1. **pvfUtility is installed** — the HTTP API must be accessible at `http://localhost:27000/Api/PvfUtiltiy/getVersion`
2. **A PVF file is loaded** — pvfUtility must have a PVF open for meaningful operations
3. **The binary exists** — `pvf-mcp.exe` (Windows) or `pvf-mcp` (Linux/macOS) must be available

## Step 2: Get the Binary

### Option A: Build from source
```bash
cd PVF-MCP/
go build -o pvf-mcp.exe .   # Windows
go build -o pvf-mcp .       # Linux/macOS
```

### Option B: Pre-built binary
Download from the release page and place it in a permanent location:
- Windows: `%APPDATA%\pvf-mcp\pvf-mcp.exe`
- Linux: `~/.local/share/pvf-mcp/pvf-mcp`
- macOS: `~/Library/Application Support/pvf-mcp/pvf-mcp`

## Step 3: Configure the MCP Client

MCP servers are configured in JSON config files. Each tool looks in different locations.

### Claude Code

**File location:** `~/.claude/mcp.json` (user-level) or `<project>/.mcp.json` (project-level)

```json
{
  "mcpServers": {
    "pvf-mcp": {
      "command": "C:\\Users\\<username>\\AppData\\Roaming\\pvf-mcp\\pvf-mcp.exe"
    }
  }
}
```

For project-level install (shared with team):
```bash
echo '{"mcpServers":{"pvf-mcp":{"command":"<full-path>/pvf-mcp"}}}' > .mcp.json
```

Restart Claude Code after configuration.

### Cursor

**File location:** `<project>/.cursor/mcp.json`

```json
{
  "mcpServers": {
    "pvf-mcp": {
      "command": "C:\\Users\\<username>\\AppData\\Roaming\\pvf-mcp\\pvf-mcp.exe"
    }
  }
}
```

Restart Cursor after configuration.

### Codex (OpenAI)

**File location:** `~/.codex/mcp.json`

```json
{
  "mcpServers": {
    "pvf-mcp": {
      "command": "C:\\Users\\<username>\\AppData\\Roaming\\pvf-mcp\\pvf-mcp.exe"
    }
  }
}
```

### Windsurf

**File location:** `~/.windsurf/mcp.json`

Same format as above.

### GitHub Copilot (VS Code)

**File location:** `.vscode/mcp.json` (project-level) or VS Code settings

```json
{
  "servers": {
    "pvf-mcp": {
      "type": "stdio",
      "command": "C:\\Users\\<username>\\AppData\\Roaming\\pvf-mcp\\pvf-mcp.exe"
    }
  }
}
```

Note: Copilot uses `"servers"` key instead of `"mcpServers"` and requires `"type": "stdio"`.

### Continue.dev

**File location:** `~/.continue/config.json`

```json
{
  "experimental": {
    "mcpServers": {
      "pvf-mcp": {
        "command": "C:\\Users\\<username>\\AppData\\Roaming\\pvf-mcp\\pvf-mcp.exe"
      }
    }
  }
}
```

### Generic MCP Client / Other Tools

**Universal format** (write to the tool's MCP config file):

```json
{
  "mcpServers": {
    "pvf-mcp": {
      "command": "<absolute-path-to-pvf-mcp-binary>"
    }
  }
}
```

Key rules for generic MCP configuration:
- Use **absolute paths** for the binary (not relative)
- No `args` needed — the binary is self-contained
- The MCP transport is **stdio** (no SSE/HTTP server needed)
- The binary has no external dependencies (Go static compilation)

If the tool uses a different config schema, adapt accordingly:
- Some tools use `"servers"` instead of `"mcpServers"`
- Some require `"transport": "stdio"` explicitly
- Some use TOML or YAML instead of JSON — convert the structure

## Step 4: Verify Installation

After restarting the AI tool, verify the MCP server is loaded:

1. Ask the agent to run `get_version`
2. Expected response: `"2023.3.9.1"` (or the installed pvfUtility version)
3. If the tool is not found, check:
   - Is the binary path absolute and correct?
   - Does the binary have execute permission? (`chmod +x pvf-mcp` on Linux/macOS)
   - Is pvfUtility running with a PVF loaded?

## Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|-------------|-----|
| Tool not appearing | Binary path wrong | Use full absolute path |
| Tool not appearing | No execute permission | `chmod +x` (Linux/macOS) |
| All tools return errors | pvfUtility not running | Start pvfUtility, load PVF |
| Empty responses | No PVF loaded | Load a PVF in pvfUtility GUI |
| Port connection refused | Wrong port | Check pvfUtility log for actual port |
| `get_version` returns empty | API path typo | Ensure pvfUtility ≥ 2022.9.30.2 |

## Automated Install (Windows)

```powershell
.\install.ps1                          # Auto-detect
.\install.ps1 -Target claude           # Claude Code only
.\install.ps1 -Target cursor           # Cursor only
.\install.ps1 -Target all              # All detected
.\install.ps1 -Target generic -ConfigPath "C:\path\mcp.json"
```

The script copies `pvf-mcp.exe` to `%APPDATA%\pvf-mcp\` and writes the correct MCP configuration.

For Linux/macOS, follow the manual steps in Step 3.
