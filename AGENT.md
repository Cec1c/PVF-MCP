# PVF-MCP Agent Usage Guide

## Overview

PVF-MCP is an MCP server that wraps the pvfUtility HTTP API, providing programmatic access to DNF (Dungeon & Fighter) PVF files. With this server, you can read, search, modify, and save PVF content — all through the AI agent's tool-calling capability.

**Prerequisites:**
- pvfUtility must be running with a PVF file loaded
- pvfUtility HTTP service must be on port 27000 (default)
- The PVF-MCP server must be configured in your MCP client

## Quick Reference Card

| Task | Tool Chain |
|------|-----------|
| Verify connection | `get_version` |
| Find an item by name | `search_pvf` → `get_item_info` |
| Read item details | `get_file_data` (structured) or `get_file_content` (raw) |
| Find item code→path | `item_code_to_file_info` |
| List directory contents | `get_pvf_root_directory` → `get_file_list` |
| Full-text search | `search_pvf` with keyword |
| Modify a file | `get_file_data` → modify JSON → `serialize_file_data` → `import_file` → `save_pvf` |
| Modify raw text | `get_file_content` → modify text → `import_file` → `save_pvf` |
| Delete a file | `delete_file` → `save_pvf` |
| Check GUI state | `get_active_document`, `get_selected_files` |
| Browse lst index | `get_all_lst_file_list` → `get_lst_file_info` |

## All Tools Reference

### Connection & Status

#### `get_version`
Returns pvfUtility version string. Use this first to verify connectivity.
```
get_version()
→ "2023.3.9.1"
```

#### `get_loaded_pvf_path`
Returns the filesystem path of the currently loaded PVF. Empty string if none loaded.
```
get_loaded_pvf_path()
→ "C:\\Users\\...\\script.pvf"
```

### Directory & File Listing

#### `get_pvf_root_directory`
Returns top-level directory names inside the PVF.
```
get_pvf_root_directory()
→ ["equipment", "stackable", "skill", "npc", "monster", "creature", ...]
```

#### `get_file_list`
Lists files in a directory, optionally filtered by extension.
```
get_file_list(dir_name="equipment", file_type=".equ")
→ ["character/common/amulet/100300001.equ", ...]
```
- `dir_name`: required. Directory name (e.g. "equipment", "stackable")
- `file_type`: optional. File extension filter (e.g. ".equ", ".stk", ".shp")

#### `file_exists`
Check if a file path exists in the PVF.
```
file_exists(file_path="stackable/book_skill2.stk")
→ "true"
```

#### `folder_exists`
Check if a directory (or file) exists.
```
folder_exists(file_path="equipment/character")
→ "true"
```

### File Reading

#### `get_file_content`
Reads a file as raw PVF-format text. Use when you need the original format for text-based modifications.
```
get_file_content(file_path="stackable/book_skill2.stk")
→ "#PVF_File\r\n\r\n[name]\r\n\t`SP+20技能书`\r\n..."
```
- `file_path`: required
- `encoding_type`: optional. "TW", "CN", "KR", "JP", "UTF8", "Unicode" (for .nut/.str/.txt files)
- `use_compatible_decompiler`: optional boolean, default false

#### `get_file_data`
Reads a file as a structured JSON tree. Every `[section]` becomes a node with `Children`. **Preferred for programmatic access and modification.**
```
get_file_data(file_path="stackable/book_skill2.stk")
→ [
    {"SectionName":"[name]", "IsSection":true, "Children":[{"Value":"`SP+20技能书`", "DataType":7}]},
    {"SectionName":"[price]", "IsSection":true, "Children":[{"Value":"132234", "DataType":2}]},
    ...
  ]
```
Node structure:
- `SectionName`: section name or null for data lines
- `IsSection`: true if this is a `[section]` header
- `HasEndSection`: true if section has `[/section]` closing tag
- `DataType`: 2=int, 4=float, 5=section, 7=string
- `Value`: the data value as a JSON string (or null for sections)
- `Children`: sub-nodes (for sections)

#### `batch_get_file_contents`
Read multiple files at once. More efficient than repeated calls.
```
batch_get_file_contents(file_paths='["creature/a.cre","creature/b.cre"]')
```

### Search

#### `search_pvf`
Full-text search across PVF content. The primary way to find items, NPCs, skills, etc.
```
search_pvf(keyword="SP+20")
→ ["stackable/book_skill2.stk", "itemname.lst", "etc/iteminfo.dat"]

search_pvf(keyword="shop", search_folder="npc", file_types=".npc")
→ ["npc/tw_gooseout.npc", ...]
```
- `keyword`: required. Searches file content.
- `search_folder`: optional. Limit to a specific directory.
- `file_types`: optional. Extension filter.

### Item Information

#### `get_item_info`
Get an item's display name and item code from its file path.
```
get_item_info(file_path="stackable/book_skill2.stk")
→ {"ItemName": "SP+20技能书", "ItemCode": 1038}
```

#### `batch_get_item_infos`
Get info for multiple items at once.
```
batch_get_item_infos(file_paths='["path/a.stk","path/b.stk"]')
→ {"path/a.stk": {"ItemName":"...", "ItemCode":...}, ...}
```

#### `item_code_to_file_info`
Convert a numeric item code to its file path and name.
```
item_code_to_file_info(lst_names="stackable", item_code=1038)
→ {"FilePath": "stackable/book_skill2.stk", "ItemName": "SP+20技能书"}

item_code_to_file_info(lst_names="equipment,stackable", item_code=27098)
→ {"FilePath": "equipment/.../n_sswd_eleno.equ", "ItemName": "..."}
```
- `lst_names`: lst file names, comma-separated. Common: "equipment", "stackable", "skill/swordman"
- `item_code`: the numeric item ID

#### `batch_item_codes_to_file_infos`
Convert multiple item codes at once.
```
batch_item_codes_to_file_infos(lst_names='["equipment","stackable"]', item_codes='[1251,27098,1038]')
→ {"Infos": {"1251": {...}, "27098": {...}, "1038": {...}}}
```

### LST File Management

#### `get_all_lst_file_list`
Get all .lst index files. These are the master indexes mapping item codes to file paths.
```
get_all_lst_file_list()
→ ["equipment/equipment.lst", "stackable/stackable.lst", "skill/swordmanskill.lst", ...]
```

#### `get_lst_file_info`
Get all entries in a .lst file as a map of item codes to details.
```
get_lst_file_info(file_path="stackable/stackable.lst")
→ {"1038": {"ItemName":"SP+20技能书", "ItemCode":1038, "FullPath":"stackable/book_skill2.stk", ...}, ...}
```

#### `file_list_to_lst_rows`
Map file paths to their lst registrations. Given file paths, returns which lst registers each one.
```
file_list_to_lst_rows(file_paths='["stackable/book_skill2.stk"]')
→ {"stackable/stackable.lst": {"1038": "book_skill2.stk"}}
```

### Icons

#### `get_item_icon`
Get an item's icon as base64-encoded image. Rate limited: 1 call per 0.5 seconds.
```
get_item_icon(file_path="equipment/.../amulet/100300001.equ")
→ "data:image/jpeg;base64,..."
```

#### `batch_get_item_icons`
Get icons for multiple items.
```
batch_get_item_icons(file_paths='["path/a.equ","path/b.equ"]')
→ {"path/a.equ": "data:image/jpeg;base64,...", ...}
```

### GUI Interaction

#### `get_active_document`
Get the file path currently focused in pvfUtility's editor.
```
get_active_document()
→ "stackable/book_skill2.stk"
```

#### `get_selected_files`
Get files currently selected in pvfUtility's file explorer tree.
```
get_selected_files()
→ ["cashshop/arad_cashshop.shp", ...]
```

### String Table

#### `get_string_table`
Get the stringtable.bin contents as a string array (game text strings).
```
get_string_table()
→ ["...", "...", ...]
```

### Write Operations (with side effects)

**IMPORTANT:** Write operations modify the in-memory PVF. You MUST call `save_pvf` afterwards to persist changes to disk. The user should back up their PVF before modifications.

#### `import_file`
Create or overwrite a file in the PVF.
```
import_file(file_path="stackable/book_skill2.stk", content="#PVF_File\r\n...")
→ "File imported. IsError=false"
```

#### `delete_file`
Delete a file from the PVF.
```
delete_file(file_path="obsolete/old_item.stk")
→ "Deleted. IsError=false"
```

#### `save_pvf`
Save the modified PVF to disk. **Always call this after any write operation.**
```
save_pvf(output_path="C:\\Users\\...\\script.pvf")
→ "PVF saved to: C:\\Users\\...\\script.pvf, IsError=false"
```

### Utility

#### `serialize_file_data`
Convert a structured JSON tree (from `get_file_data`) back to PVF text format. Essential for the modify-workflow.
```
serialize_file_data(nodes='[{"SectionName":"[price]",...}]')
→ "#PVF_File\r\n\r\n[price]\r\n\t133322\r\n"
```

## Common Workflows

### Workflow 1: Find Item and Read Price
```
1. search_pvf(keyword="火之心项链")          → find file path
2. get_item_info(file_path="...")            → confirm item name + code
3. get_file_data(file_path="...")            → read structured data
   (find [price] section → first child Value = price)
```

### Workflow 2: Modify Item Price (Structured)
```
1. get_file_data(file_path="stackable/book_skill2.stk")
2. Modify the JSON: find [price] section → change Children[0].Value
3. serialize_file_data(nodes=<modified JSON>)
4. import_file(file_path="stackable/book_skill2.stk", content=<serialized>)
5. save_pvf(output_path="C:\...\script.pvf")
```

### Workflow 3: Modify Item Price (Raw Text)
```
1. get_file_content(file_path="stackable/book_skill2.stk")
2. Replace "132234" with "133322" in the text
3. import_file(file_path="stackable/book_skill2.stk", content=<modified text>)
4. save_pvf(output_path="C:\...\script.pvf")
```

### Workflow 4: Find Which NPCs Sell an Item
```
1. search_pvf(keyword="item_name")           → find the item file
2. get_item_info(file_path="...")            → get item code (e.g. 1038)
3. search_pvf(keyword="1038", file_types=".shp")  → find shop files
   Note: This PVF uses item-group based shops.
   All [etc shop] type shops sell items based on [item group index].
   Individual items are not explicitly listed in .shp files.
```

### Workflow 5: Batch Explore a Category
```
1. get_all_lst_file_list()                   → see all lst files
2. get_lst_file_info(file_path="equipment/equipment.lst")
   → get all equipment codes and paths
3. batch_get_item_infos(file_paths='[...]')  → get names for specific items
```

## Data Types Reference

From pvfUtility's `ScriptType` enum:
| DataType | Meaning | PVF Format Example |
|----------|---------|-------------------|
| 2 | Int | `500` |
| 4 | Float | `0.5` |
| 5 | Section | `[name]` ... `[/name]` |
| 7 | String | `` `value` `` |

## Error Handling

All tools return errors as text messages (not protocol errors). Check for:
- "文件不存在" — File not found in PVF
- Connection errors — pvfUtility not running or PVF not loaded
- Always call `get_version` or `get_loaded_pvf_path` first to verify connectivity

## Limitations

1. **No remote PVF loading** — PVF must be opened manually in pvfUtility GUI first
2. **Item group shops** — This PVF version uses item group indices, not explicit item lists in shops
3. **Icon rate limit** — 1 icon request per 0.5 seconds
4. **HTTP API port** — Changes if multiple pvfUtility instances are open
5. **Write-then-save required** — Modifications are in-memory until `save_pvf` is called
