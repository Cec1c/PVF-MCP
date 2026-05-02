# PVF-MCP

让 AI 直接操作 DNF PVF 文件。

---

## 这是什么？

[pvfUtility](https://github.com/ledyxerago/pvfUtilityForEAssistant) 为 DNF 的 PVF 文件提供了 HTTP API，可以读、写、搜索游戏数据。

**PVF-MCP** 把这个 API 包装成 MCP 服务 —— 从此你可以用自然语言对 AI 说：

> *"找到 SP+20 技能书，把价格改成 133322，然后保存"*

AI 会自动调用对应工具完成操作。

## 亮点

- **26 个工具** 覆盖 pvfUtility 全部 HTTP API——搜索物品、读写文件、查图标、翻索引
- **单文件部署** Go 静态编译的 `.exe`，11MB，不需要装 Python / Node / Java
- **结构化编辑** `get_file_data` → 改 JSON → `serialize_file_data` → `import_file` 完整闭环
- **零依赖安装** PowerShell 脚本一键搞定，Win10/11 开箱即用
- **多平台 Agent** Claude Code / Cursor / Codex / Windsurf / Continue.dev 全支持

## 快速开始

### 1. 编译

```powershell
git clone https://github.com/Cec1c/PVF-MCP.git
cd PVF-MCP
go build -o pvf-mcp.exe .
```

> 不想编译？直接去 [Releases](https://github.com/Cec1c/PVF-MCP/releases) 下载 `pvf-mcp.exe`。

### 2. 安装

```powershell
.\install.ps1                    # 自动检测你装了哪个 AI 工具
.\install.ps1 -Target claude      # 或手动指定
.\install.ps1 -Target cursor
.\install.ps1 -Target codex
```

### 3. 开用

1. 打开 pvfUtility，加载一个 PVF
2. 重启你的 AI 工具
3. 对它说：`用 get_version 检查连接`

搞定。

## 工具一览

| 分类 | 工具 | 一句话 |
|------|------|--------|
| 连接 | `get_version` `get_loaded_pvf_path` | 确认 pvfUtility 在线、PVF 已加载 |
| 浏览 | `get_pvf_root_directory` `get_file_list` `file_exists` `folder_exists` | 逛 PVF 目录树 |
| 阅读 | `get_file_content` `get_file_data` `batch_get_file_contents` | 读原始文本 / 结构化 JSON |
| 搜索 | `search_pvf` | 全文搜索，找物品/技能/NPC |
| 物品 | `get_item_info` `batch_get_item_infos` `item_code_to_file_info` `batch_item_codes_to_file_infos` | 物品名 ↔ 代码 ↔ 文件路径互查 |
| 索引 | `get_all_lst_file_list` `get_lst_file_info` `file_list_to_lst_rows` | lst 索引翻个底朝天 |
| 图标 | `get_item_icon` `batch_get_item_icons` | 拿物品图标 base64 |
| GUI | `get_active_document` `get_selected_files` | 跟 pvfUtility 窗口联动 |
| 修改 | `import_file` `delete_file` `save_pvf` | 改完记得保存 |
| 转换 | `serialize_file_data` | JSON 树 → PVF 文本 |

## 典型操作

### 改物品价格

```
get_file_data → 改 JSON 里的 price → serialize_file_data → import_file → save_pvf
```

### 物品代码查文件

```
item_code_to_file_info(lst_names="stackable", item_code=1038)
→ stackable/book_skill2.stk  /  SP+20技能书
```

### 搜某个 NPC 的商店

```
search_pvf(keyword="赛丽亚", search_folder="npc", file_types=".npc")
```

## 项目结构

```
PVF-MCP/
├── server.go          # 入口 + MCP 服务 + 26 个工具
├── client.go          # pvfUtility HTTP 客户端
├── serializer.go      # PVF 文本 ↔ JSON 树
├── types.go           # 数据类型
├── install.ps1        # 安装脚本
├── test_mcp.py        # E2E 测试
├── AGENT.md           # Agent 详细手册
├── INSTALL.md         # 安装指南（给 Agent 看）
└── README.md          # 你正在看
```

## 依赖

- [mcp-go](https://github.com/mark3labs/mcp-go) Go MCP SDK
- pvfUtility ≥ 2022.9.30.2
- Go 1.21+ (仅编译)
- PowerShell 5.1+ (Win10/11 自带)

## 反馈与交流
- QQ交流群 ： 908172336 （实际上是我开的公益群服，想玩也可以来看看）
