# PVF-MCP

将 [pvfUtility](https://github.com/ledyxerago/pvfUtilityForEAssistant) 的 HTTP API 封装为 MCP (Model Context Protocol) 服务，让 AI 编程助手（Claude Code、Cursor、Codex 等）能直接操作 DNF PVF 文件。

## 功能

- **26 个 MCP 工具**：覆盖 PVF 文件的完整 CRUD 操作
- **单文件部署**：Go 编译的独立 .exe，无需 Python/Node.js 运行时
- **PVF 序列化器**：支持 get_file_data（结构化 JSON）→ 修改 → serialize_file_data → import_file 的闭环
- **多工具支持**：Claude Code / Cursor / Codex / 通用 MCP 客户端

## 快速开始

### 1. 编译

```powershell
cd PVF-MCP
go build -o pvf-mcp.exe .
```

### 2. 安装

```powershell
# 自动检测并安装到已安装的 AI 工具
.\install.ps1

# 或指定工具
.\install.ps1 -Target claude
.\install.ps1 -Target cursor
.\install.ps1 -Target codex

# 安装到所有检测到的工具
.\install.ps1 -Target all
```

### 3. 使用

1. 打开 pvfUtility，加载 PVF 文件
2. 重启你的 AI 编程助手
3. 对 AI 说：「搜索 SP+20 技能书的价格」

## 工具概览

| 类别 | 工具 |
|------|------|
| 连接 | `get_version`, `get_loaded_pvf_path` |
| 目录 | `get_pvf_root_directory`, `get_file_list`, `file_exists`, `folder_exists` |
| 读取 | `get_file_content`, `get_file_data`, `batch_get_file_contents` |
| 搜索 | `search_pvf` |
| 物品 | `get_item_info`, `batch_get_item_infos`, `item_code_to_file_info`, `batch_item_codes_to_file_infos` |
| 索引 | `get_all_lst_file_list`, `get_lst_file_info`, `file_list_to_lst_rows` |
| 图标 | `get_item_icon`, `batch_get_item_icons` |
| GUI | `get_active_document`, `get_selected_files` |
| 字符串 | `get_string_table` |
| 写入 | `import_file`, `delete_file`, `save_pvf` |
| 工具 | `serialize_file_data` |

## 系统要求

- **pvfUtility** 2022.9.30.2 或更高版本（HTTP API 支持）
- **Go** 1.21+（仅编译时需要）
- **PowerShell** 5.1+（Windows 10/11 自带，安装脚本用）

## 项目结构

```
PVF-MCP/
├── server.go              # 入口 + MCP 服务 + 26 个工具
├── client.go              # HTTP 客户端
├── serializer.go          # PVF 文本 ↔ JSON 序列化
├── types.go               # 类型定义
├── install.ps1            # 安装脚本 (PowerShell, Windows 自带)
├── test_mcp.py            # E2E 测试
├── go.mod / go.sum
├── AGENT.md               # Agent 详细使用文档
├── INSTALL.md             # 安装指南（给 Agent 看）
└── README.md              # 本文件
```

## 依赖

- [mcp-go](https://github.com/mark3labs/mcp-go) — Go MCP SDK
