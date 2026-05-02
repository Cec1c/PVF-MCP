package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	log.SetOutput(os.Stderr)

	pvf := NewPvfClient()

	s := server.NewMCPServer(
		"pvf-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	registerTools(s, pvf)

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func registerTools(s *server.MCPServer, pvf *PvfClient) {
	registerGetVersion(s, pvf)
	registerGetPvfRootDirectory(s, pvf)
	registerGetFileList(s, pvf)
	registerGetFileContent(s, pvf)
	registerGetFileData(s, pvf)
	registerBatchGetFileContents(s, pvf)
	registerSearchPvf(s, pvf)
	registerGetItemInfo(s, pvf)
	registerBatchGetItemInfos(s, pvf)
	registerItemCodeToFileInfo(s, pvf)
	registerBatchItemCodesToFileInfos(s, pvf)
	registerGetAllLstFileList(s, pvf)
	registerGetLstFileInfo(s, pvf)
	registerFileListToLstRows(s, pvf)
	registerGetItemIcon(s, pvf)
	registerBatchGetItemIcons(s, pvf)
	registerGetActiveDocument(s, pvf)
	registerGetSelectedFiles(s, pvf)
	registerGetStringTable(s, pvf)
	registerGetLoadedPvfPath(s, pvf)
	registerFileExists(s, pvf)
	registerFolderExists(s, pvf)
	registerImportFile(s, pvf)
	registerDeleteFile(s, pvf)
	registerSavePvf(s, pvf)
	registerSerializeFileData(s, pvf)
}

func jsonResult(data interface{}) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("JSON marshal error: %v", err)), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}

func jsonOrError(data interface{}, err error) (*mcp.CallToolResult, error) {
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return jsonResult(data)
}

// ============================================================
// Read tools
// ============================================================

func registerGetVersion(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_version",
		mcp.WithDescription("Get pvfUtility version number. Use to verify API accessibility."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		v, err := pvf.GetVersion()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(v), nil
	})
}

func registerGetPvfRootDirectory(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_pvf_root_directory",
		mcp.WithDescription("Get root directory listing of the currently loaded PVF. Returns categories like equipment, stackable, npc, skill, etc."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return jsonOrError(pvf.GetPvfRootDirectory())
	})
}

func registerGetFileList(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_file_list",
		mcp.WithDescription("List files in a PVF directory. Optionally filter by extension."),
		mcp.WithString("dir_name", mcp.Required(), mcp.Description("Directory inside PVF: equipment, stackable, npc, etc.")),
		mcp.WithString("file_type", mcp.Description("Optional extension filter: .equ, .stk, .npc, .shp")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		dirName, _ := req.RequireString("dir_name")
		fileType := req.GetString("file_type", "")
		return jsonOrError(pvf.GetFileList(dirName, fileType))
	})
}

func registerGetFileContent(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_file_content",
		mcp.WithDescription("Read raw PVF text content of a file. Use for viewing or modifying original format."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Full PVF path: 'stackable/book_skill2.stk'")),
		mcp.WithString("encoding_type", mcp.Description("Encoding: TW, CN, KR, JP, UTF8, Unicode")),
		mcp.WithBoolean("use_compatible_decompiler", mcp.Description("Use compatible decompiler")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		encodingType := req.GetString("encoding_type", "")
		useCompat := req.GetBool("use_compatible_decompiler", false)
		content, err := pvf.GetFileContent(filePath, encodingType, useCompat)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(content), nil
	})
}

func registerGetFileData(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_file_data",
		mcp.WithDescription("Read a PVF file as a structured JSON tree. Each [section] is a node with Children. Use this to programmatically navigate or modify content."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Full PVF path")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		return jsonOrError(pvf.GetFileData(filePath))
	})
}

func registerBatchGetFileContents(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("batch_get_file_contents",
		mcp.WithDescription("Read multiple files at once. More efficient than repeated get_file_content calls."),
		mcp.WithString("file_paths", mcp.Required(), mcp.Description("JSON array: '[\"path/a.cre\",\"path/b.atk\"]'")),
		mcp.WithBoolean("use_compatible_decompiler", mcp.Description("Use compatible decompiler")),
		mcp.WithString("encoding_type", mcp.Description("Encoding: TW, CN, KR, JP, UTF8, Unicode")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pathsStr, _ := req.RequireString("file_paths")
		var fileList []string
		if err := json.Unmarshal([]byte(pathsStr), &fileList); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("file_paths must be JSON array: %v", err)), nil
		}
		useCompat := req.GetBool("use_compatible_decompiler", false)
		encodingType := req.GetString("encoding_type", "")
		return jsonOrError(pvf.GetFileContents(fileList, useCompat, encodingType))
	})
}

func registerSearchPvf(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("search_pvf",
		mcp.WithDescription("Full-text search across PVF content. Primary way to find items, NPCs, skills by name or keyword."),
		mcp.WithString("keyword", mcp.Required(), mcp.Description("Search term. Searches file content, not just names.")),
		mcp.WithString("search_folder", mcp.Description("Limit to folder: equipment, npc, etc. Empty = all.")),
		mcp.WithString("file_types", mcp.Description("Extension filter: .stk, .npc, .shp")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		keyword, _ := req.RequireString("keyword")
		searchFolder := req.GetString("search_folder", "")
		fileTypes := req.GetString("file_types", "")
		return jsonOrError(pvf.SearchPvfSimple(keyword, searchFolder, fileTypes))
	})
}

func registerGetItemInfo(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_item_info",
		mcp.WithDescription("Get an item's display name and item code from its file path."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Path: 'stackable/book_skill2.stk'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		return jsonOrError(pvf.GetItemInfo(filePath))
	})
}

func registerBatchGetItemInfos(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("batch_get_item_infos",
		mcp.WithDescription("Get names and codes for multiple items at once."),
		mcp.WithString("file_paths", mcp.Required(), mcp.Description("JSON array of file paths")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pathsStr, _ := req.RequireString("file_paths")
		var filePaths []string
		if err := json.Unmarshal([]byte(pathsStr), &filePaths); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("file_paths must be JSON array: %v", err)), nil
		}
		return jsonOrError(pvf.GetItemInfos(filePaths))
	})
}

func registerItemCodeToFileInfo(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("item_code_to_file_info",
		mcp.WithDescription("Convert an item code (integer ID) to its file path and display name."),
		mcp.WithString("lst_names", mcp.Required(), mcp.Description("Lst name(s), comma-separated: 'equipment', 'stackable', 'equipment,stackable'")),
		mcp.WithNumber("item_code", mcp.Required(), mcp.Description("Numeric item code to look up")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lstNames, _ := req.RequireString("lst_names")
		itemCode, _ := req.RequireInt("item_code")
		return jsonOrError(pvf.ItemCodeToFileInfo(lstNames, itemCode))
	})
}

func registerBatchItemCodesToFileInfos(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("batch_item_codes_to_file_infos",
		mcp.WithDescription("Convert multiple item codes to file paths at once."),
		mcp.WithString("lst_names", mcp.Required(), mcp.Description("JSON array: '[\"equipment\",\"stackable\"]'")),
		mcp.WithString("item_codes", mcp.Required(), mcp.Description("JSON array: '[1251,27098,1038]'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		lstStr, _ := req.RequireString("lst_names")
		codesStr, _ := req.RequireString("item_codes")
		var lstNames []string
		var itemCodes []int
		if err := json.Unmarshal([]byte(lstStr), &lstNames); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("lst_names must be JSON array: %v", err)), nil
		}
		if err := json.Unmarshal([]byte(codesStr), &itemCodes); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("item_codes must be JSON array: %v", err)), nil
		}
		return jsonOrError(pvf.ItemCodesToFileInfos(lstNames, itemCodes))
	})
}

func registerGetAllLstFileList(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_all_lst_file_list",
		mcp.WithDescription("Get all .lst index files. These map item codes to file paths."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return jsonOrError(pvf.GetAllLstFileList())
	})
}

func registerGetLstFileInfo(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_lst_file_info",
		mcp.WithDescription("Get all entries in a .lst index file as a map of item codes to paths and names."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Path: 'equipment/equipment.lst'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		return jsonOrError(pvf.GetLstFileInfo(filePath))
	})
}

func registerFileListToLstRows(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("file_list_to_lst_rows",
		mcp.WithDescription("Map file paths to their lst registrations. Shows which lst file and item code each path belongs to."),
		mcp.WithString("file_paths", mcp.Required(), mcp.Description("JSON array of file paths")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pathsStr, _ := req.RequireString("file_paths")
		var filePaths []string
		if err := json.Unmarshal([]byte(pathsStr), &filePaths); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("file_paths must be JSON array: %v", err)), nil
		}
		return jsonOrError(pvf.FileListToLstRows(filePaths))
	})
}

func registerGetItemIcon(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_item_icon",
		mcp.WithDescription("Get an item's icon as a base64-encoded image. Rate limited: 1 per 0.5s."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Path to item file")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		icon, err := pvf.GetFileIcon(filePath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(icon), nil
	})
}

func registerBatchGetItemIcons(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("batch_get_item_icons",
		mcp.WithDescription("Get icons for multiple items at once."),
		mcp.WithString("file_paths", mcp.Required(), mcp.Description("JSON array of file paths")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		pathsStr, _ := req.RequireString("file_paths")
		var filePaths []string
		if err := json.Unmarshal([]byte(pathsStr), &filePaths); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("file_paths must be JSON array: %v", err)), nil
		}
		return jsonOrError(pvf.FilesToIconBase64(filePaths))
	})
}

func registerGetActiveDocument(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_active_document",
		mcp.WithDescription("Get the file path of the currently focused document in pvfUtility's editor."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		doc, err := pvf.GetActiveDocumentFilePath()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(doc), nil
	})
}

func registerGetSelectedFiles(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_selected_files",
		mcp.WithDescription("Get files currently selected in pvfUtility's file explorer tree."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return jsonOrError(pvf.GetTreeSelectedFiles())
	})
}

func registerGetStringTable(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_string_table",
		mcp.WithDescription("Get the stringtable.bin contents as string array."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return jsonOrError(pvf.GetStringTable())
	})
}

func registerGetLoadedPvfPath(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("get_loaded_pvf_path",
		mcp.WithDescription("Get the filesystem path of the currently loaded PVF. Use to verify a PVF is open."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		p, err := pvf.GetPvfPackFilePath()
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(p), nil
	})
}

func registerFileExists(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("file_exists",
		mcp.WithDescription("Check if a file exists in the PVF."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("File path to check")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		exists, err := pvf.FileIsExists(filePath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%v", exists)), nil
	})
}

func registerFolderExists(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("folder_exists",
		mcp.WithDescription("Check if a directory or file exists in the PVF."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Directory or file path to check")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		exists, err := pvf.FolderExists(filePath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("%v", exists)), nil
	})
}

// ============================================================
// Write tools
// ============================================================

func registerImportFile(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("import_file",
		mcp.WithDescription("Create or overwrite a file in the PVF. Provide full PVF-format content. WARNING: Modifies the loaded PVF. Use save_pvf afterwards."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("Target path: 'stackable/book_skill2.stk'")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Full file content in PVF text format")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		content, _ := req.RequireString("content")
		ar, err := pvf.ImportFile(filePath, content)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("File imported. IsError=%v", ar.IsError)), nil
	})
}

func registerDeleteFile(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("delete_file",
		mcp.WithDescription("Delete a file from the PVF. WARNING: Modifies the loaded PVF. Use save_pvf afterwards."),
		mcp.WithString("file_path", mcp.Required(), mcp.Description("File path to delete")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		filePath, _ := req.RequireString("file_path")
		ar, err := pvf.DeleteFile(filePath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("Deleted. IsError=%v", ar.IsError)), nil
	})
}

func registerSavePvf(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("save_pvf",
		mcp.WithDescription("Save the loaded PVF to disk. Use after import_file or delete_file to persist changes."),
		mcp.WithString("output_path", mcp.Required(), mcp.Description("Full filesystem path: 'C:\\output\\script.pvf'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		outputPath, _ := req.RequireString("output_path")
		ar, err := pvf.SaveAsPvfFile(outputPath)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("PVF saved to: %s, IsError=%v", outputPath, ar.IsError)), nil
	})
}

// ============================================================
// Utility
// ============================================================

func registerSerializeFileData(s *server.MCPServer, pvf *PvfClient) {
	s.AddTool(mcp.NewTool("serialize_file_data",
		mcp.WithDescription("Convert structured JSON tree (from get_file_data) back to PVF text format. Workflow: get_file_data → modify JSON → serialize_file_data → import_file."),
		mcp.WithString("nodes", mcp.Required(), mcp.Description("JSON array of WebApiFileData nodes (from get_file_data)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		nodesStr, _ := req.RequireString("nodes")
		var nodes []WebApiFileData
		if err := json.Unmarshal([]byte(nodesStr), &nodes); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("nodes must be JSON array: %v", err)), nil
		}
		text := SerializeWebApiFileData(nodes, "")
		return mcp.NewToolResultText(text), nil
	})
}
