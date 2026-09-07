<div align="center">

# PVF-MCP

[![语言：Go](https://img.shields.io/static/v1?label=%E8%AF%AD%E8%A8%80&message=Go&color=00ADD8&style=flat-square&logo=go&logoColor=white)](https://go.dev/)
[![协议：MCP stdio](https://img.shields.io/static/v1?label=%E5%8D%8F%E8%AE%AE&message=MCP%20stdio&color=8B5CF6&style=flat-square)](https://modelcontextprotocol.io/)
[![平台：Windows 10 / 11](https://img.shields.io/static/v1?label=%E5%B9%B3%E5%8F%B0&message=Windows%2010%20%2F%2011&color=0078D4&style=flat-square)](#快速开始)
[![发行版](https://img.shields.io/github/v/release/Cec1c/PVF-MCP?sort=date&display_name=tag&style=flat-square&label=%E5%8F%91%E8%A1%8C%E7%89%88&color=2EA44F)](https://github.com/Cec1c/PVF-MCP/releases/latest)
[![Stars](https://img.shields.io/github/stars/Cec1c/PVF-MCP?style=flat-square&label=Stars&color=E3B341)](https://github.com/Cec1c/PVF-MCP/stargazers)

让 AI 通过 pvfUtility 搜索、读取和编辑 DNF PVF 游戏数据。

基于 Go 的 MCP（Model Context Protocol）服务，提供 26 个工具；下载单个可执行文件，即可接入支持 stdio 的 AI 客户端。

[English](README.en.md) ｜ [下载](https://github.com/Cec1c/PVF-MCP/releases/latest) ｜ [快速开始](#快速开始) ｜ [工具手册](AGENT.md)

</div>

> [!NOTE]
> PVF-MCP 当前连接同一台电脑的 `http://localhost:27000`。AI 客户端通过 **stdio** 启动 PVF-MCP，无需把这个 HTTP 地址填作 MCP 服务地址。

## 能做什么？

- **找到需要的数据。** 搜索物品、技能和 NPC，按物品代码查文件，读取 LST 索引与图标。
- **按结构修改内容。** 将 PVF 解析为 JSON 树，修改字段后转回原始文本，导入并保存。
- **配合编辑器工作。** 获取 pvfUtility 当前打开的文档和选中文件，让 AI 接着你正在看的内容处理。

例如，你可以这样提问：

> 找到 SP+20 技能书，列出文件路径、物品代码和当前价格，先不要修改。

## 快速开始

### 1. 准备 pvfUtility

在 Windows 10 / 11 上打开 [pvfUtility](https://github.com/ledyxerago/pvfUtilityForEAssistant)（≥ 2022.9.30.2），加载一个 PVF 文件，并确认 HTTP API 可用。

### 2. 下载并接入 AI 客户端

从 [最新 Release](https://github.com/Cec1c/PVF-MCP/releases/latest) 下载 `pvf-mcp.exe`，放到固定目录，例如 `C:\Tools\pvf-mcp\pvf-mcp.exe`。使用发行包无需安装 Go、Python 或 Node.js。

Codex CLI 用户可运行以下命令，将路径替换为实际存放位置：

```powershell
codex mcp add pvf-mcp -- "C:\Tools\pvf-mcp\pvf-mcp.exe"
```

其他客户端选择 **stdio**，将启动命令设为可执行文件的绝对路径，无需额外参数。

<details>
<summary>JSON 配置示例（支持 mcpServers 格式的客户端）</summary>

将此条目合并到客户端现有 MCP 配置中，并按该客户端的说明选择配置文件位置：

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

### 3. 确认连接

重启或重新加载 AI 客户端，让它依次调用 `get_version` 和 `get_loaded_pvf_path`。能返回 pvfUtility 版本和当前 PVF 路径，就可以开始使用。

连接失败时，先检查 pvfUtility 是否已启动、PVF 是否已加载，以及 HTTP API 是否使用端口 `27000`；当前服务地址固定在 [`client.go`](client.go) 中。

## 常用操作

| 想做的事 | 使用的工具 |
| --- | --- |
| 按关键词找物品、技能或 NPC | `search_pvf` |
| 按物品代码查文件与名称 | `item_code_to_file_info` |
| 读取原文或结构化数据 | `get_file_content` / `get_file_data` |
| 批量读取文件或物品信息 | `batch_get_file_contents` / `batch_get_item_infos` |
| 查看索引和图标 | `get_lst_file_info` / `get_item_icon` |
| 获取编辑器当前文档与选中项 | `get_active_document` / `get_selected_files` |
| 导入、删除或保存文件 | `import_file` / `delete_file` / `save_pvf` |

> [!IMPORTANT]
> 写操作先改变已加载的 PVF，调用 `save_pvf` 才会写入磁盘。修改前保留备份，并通过 `output_path` 指定另存路径。

结构化修改的完整流程：

```text
get_file_data → 修改 JSON → serialize_file_data → import_file → save_pvf
```

全部工具、参数和示例见 [Agent 工具手册](AGENT.md)。

## 从源码构建

当前源码要求 **Go 1.26.2 或更高版本**，具体以 [`go.mod`](go.mod) 为准：

```powershell
git clone https://github.com/Cec1c/PVF-MCP.git
cd PVF-MCP
go build -o pvf-mcp.exe .
```

构建后按上方步骤接入客户端。仓库也提供 [`install.ps1`](install.ps1)（PowerShell 5.1+），但脚本及旧版 [`INSTALL.md`](INSTALL.md) 中包含历史客户端配置路径；Codex 请使用本页的 `codex mcp add` 命令。

## 文档与来源

| 内容 | 入口 |
| --- | --- |
| 工具参数、数据类型与操作示例 | [AGENT.md](AGENT.md) |
| MCP 工具注册与服务入口 | [server.go](server.go) |
| PVF HTTP 客户端与文本序列化 | [client.go](client.go) · [serializer.go](serializer.go) |

基于 [pvfUtility](https://github.com/ledyxerago/pvfUtilityForEAssistant) 的 HTTP API，使用 [mcp-go](https://github.com/mark3labs/mcp-go) 实现 MCP 服务。

## 反馈与交流

- 问题与建议：[GitHub Issues](https://github.com/Cec1c/PVF-MCP/issues)。
- QQ交流群 ： 908172336 （实际上是我开的公益群服，想玩也可以来看看）
