<div align="center">

# PVF-MCP

[![Language: Go](https://img.shields.io/static/v1?label=Language&message=Go&color=00ADD8&style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![Protocol: MCP stdio](https://img.shields.io/static/v1?label=Protocol&message=MCP%20stdio&color=8B5CF6&style=flat-square)](https://modelcontextprotocol.io/)
[![Platform: Windows 10 / 11](https://img.shields.io/static/v1?label=Platform&message=Windows%2010%20%2F%2011&color=0078D4&style=flat-square)](#quick-start)
[![Release](https://img.shields.io/github/v/release/Cec1c/PVF-MCP?sort=date&display_name=tag&style=flat-square&label=Release&color=2EA44F)](https://github.com/Cec1c/PVF-MCP/releases/latest)
[![Stars](https://img.shields.io/github/stars/Cec1c/PVF-MCP?style=flat-square&label=Stars&color=E3B341)](https://github.com/Cec1c/PVF-MCP/stargazers)

Let AI search, read, and edit Dungeon & Fighter PVF game data through pvfUtility.

A Go-based Model Context Protocol (MCP) server with 26 tools. Download one executable and connect it to an AI client that supports stdio.

[简体中文](README.md) ｜ [Download](https://github.com/Cec1c/PVF-MCP/releases/latest) ｜ [Quick Start](#quick-start) ｜ [Tool Reference](AGENT.md)

</div>

> [!NOTE]
> PVF-MCP currently connects to `http://localhost:27000` on the same computer. Your AI client launches PVF-MCP over **stdio**; this HTTP address is not an MCP server URL.

## What Can It Do?

- **Find the data you need.** Search items, skills, and NPCs; resolve item codes to files; read LST indexes and icons.
- **Edit structured content.** Read PVF data as a JSON tree, change fields, serialize it back to text, then import and save.
- **Work with your editor.** Read the active document and selected files in pvfUtility so the AI can work on what you are viewing.

For example, ask your AI client:

> Find the SP+20 skill book and list its file path, item code, and current price. Do not change anything yet.

## Quick Start

### 1. Prepare pvfUtility

On Windows 10 / 11, open [pvfUtility](https://github.com/ledyxerago/pvfUtilityForEAssistant) (≥ 2022.9.30.2), load a PVF file, and make sure its HTTP API is available.

### 2. Download and Connect

Download `pvf-mcp.exe` from the [latest Release](https://github.com/Cec1c/PVF-MCP/releases/latest) and keep it at a permanent path, such as `C:\Tools\pvf-mcp\pvf-mcp.exe`. The release binary does not require Go, Python, or Node.js.

For Codex CLI, run the following command with your actual executable path:

```powershell
codex mcp add pvf-mcp -- "C:\Tools\pvf-mcp\pvf-mcp.exe"
```

For other clients, select **stdio** and set the launch command to the executable's absolute path. No additional arguments are needed.

<details>
<summary>JSON example for clients that use the mcpServers format</summary>

Merge this entry into the client's existing MCP configuration. Follow that client's documentation for the configuration file location:

```json
{
  "mcpServers": {
    "pvf-mcp": {
      "command": "C:\\Tools\\pvf-mcp\\pvf-mcp.exe"
    }
  }
}
```

</details>

### 3. Verify the Connection

Restart or reload your AI client, then ask it to call `get_version` and `get_loaded_pvf_path`. Once they return the pvfUtility version and the loaded PVF path, you are ready to use the tools.

If the connection fails, check that pvfUtility is running, a PVF is loaded, and its HTTP API uses port `27000`. The current server address is fixed in [`client.go`](client.go).

## Common Tasks

| Task | Tools |
| --- | --- |
| Search items, skills, or NPCs by keyword | `search_pvf` |
| Resolve an item code to a file and name | `item_code_to_file_info` |
| Read raw text or structured data | `get_file_content` / `get_file_data` |
| Read multiple files or item records | `batch_get_file_contents` / `batch_get_item_infos` |
| Inspect indexes and icons | `get_lst_file_info` / `get_item_icon` |
| Read the active editor document and selection | `get_active_document` / `get_selected_files` |
| Import, delete, or save files | `import_file` / `delete_file` / `save_pvf` |

> [!IMPORTANT]
> Write operations change the loaded PVF in memory. Call `save_pvf` to persist them to disk. Keep a backup before editing and use `output_path` to save to a separate file.

A complete structured-edit workflow:

```text
get_file_data → Edit JSON → serialize_file_data → import_file → save_pvf
```

See the [Agent Tool Reference](AGENT.md) for all tools, parameters, and examples.

## Build from Source

The current source requires **Go 1.26.2 or later**. See [`go.mod`](go.mod) for the declared version:

```powershell
git clone https://github.com/Cec1c/PVF-MCP.git
cd PVF-MCP
go build -o pvf-mcp.exe .
```

Connect the built executable using the steps above. The repository also includes [`install.ps1`](install.ps1) (PowerShell 5.1+), but the script and older [`INSTALL.md`](INSTALL.md) contain legacy client configuration paths. For Codex, use the `codex mcp add` command on this page.

## Documentation and Credits

| Topic | Reference |
| --- | --- |
| Tool parameters, data types, and workflows | [AGENT.md](AGENT.md) |
| MCP tool registration and server entry point | [server.go](server.go) |
| PVF HTTP client and text serialization | [client.go](client.go) · [serializer.go](serializer.go) |

Built on the [pvfUtility](https://github.com/ledyxerago/pvfUtilityForEAssistant) HTTP API, using [mcp-go](https://github.com/mark3labs/mcp-go) to implement the MCP server.

## Feedback and Community

- Bugs and suggestions: [GitHub Issues](https://github.com/Cec1c/PVF-MCP/issues).
- QQ group: **908172336** — also the author's community game-server group; players are welcome.
