package main

import "encoding/json"

// ApiResponse is the generic wrapper returned by pvfUtility HTTP API.
type ApiResponse struct {
	Data    json.RawMessage `json:"Data"`
	IsError bool            `json:"IsError"`
	Msg     *string         `json:"Msg"`
	ErrorId int             `json:"ErrorId"`
}

// ItemInfo from GetItemInfo/GetItemInfos.
type ItemInfo struct {
	ItemName string `json:"ItemName"`
	ItemCode int    `json:"ItemCode"`
}

// ItemCodeFileInfo from ItemCodeToFileInfo/ItemCodesToFileInfos.
type ItemCodeFileInfo struct {
	FilePath string `json:"FilePath"`
	ItemName string `json:"ItemName"`
}

// LstFileInfo from getLstFileInfo.
type LstFileInfo struct {
	PathHeader string `json:"PathHeader"`
	ItemPath   string `json:"ItemPath"`
	FullPath   string `json:"FullPath"`
	ItemName   string `json:"ItemName"`
	ItemCode   int    `json:"ItemCode"`
}

// WebApiFileData is the structured JSON tree node returned by getFileData.
type WebApiFileData struct {
	SectionName   *string          `json:"SectionName"`
	IsSection     bool             `json:"IsSection"`
	HasEndSection bool             `json:"HasEndSection"`
	DataType      int              `json:"DataType"`
	Value         json.RawMessage  `json:"Value"`
	Children      []WebApiFileData `json:"Children"`
}

// ScriptType enum matching pvfUtility's ScriptType.
const (
	ScriptTypeInt             = 2
	ScriptTypeIntEx           = 3
	ScriptTypeFloat           = 4
	ScriptTypeSection         = 5
	ScriptTypeCommand         = 6
	ScriptTypeString          = 7
	ScriptTypeCommandSeparator = 8
	ScriptTypeStringLinkIndex  = 9
	ScriptTypeStringLink       = 10
)

// SearchPvfRequest payload for the SearchPvf endpoint.
type SearchPvfRequest struct {
	SearchFolder             string      `json:"SearchFolder"`
	Keyword                  string      `json:"Keyword"`
	Type                     int         `json:"Type"`
	SourceType               int         `json:"SourceType"`
	NormalUsing              int         `json:"NormalUsing"`
	IsStartMatch             bool        `json:"IsStartMatch"`
	SearchResult             interface{} `json:"SearchResult"`
	ScriptContentSearchMode  int         `json:"ScriptContentSearchMode"`
	IsUseLikeSearchPath      bool        `json:"IsUseLikeSearchPath"`
	Trait                    bool        `json:"Trait"`
	UseRegularExpression     bool        `json:"UseRegularExpression"`
	WholeWordMatch            bool        `json:"WholeWordMatch"`
	RemoveOrKeep             int         `json:"RemoveOrKeep"`
	FileTypesString          interface{} `json:"FileTypesString"`
	ScriptContent            string      `json:"ScriptContent"`
	ScriptContentStart       string      `json:"ScriptContentStart"`
	ScriptContentStop        string      `json:"ScriptContentStop"`
}

// FileContentEntry for batch file content request/response.
type FileContentEntry struct {
	FilePath    string `json:"FilePath"`
	FileContent string `json:"FileContent"`
}

// GetFileContentsRequest payload.
type GetFileContentsRequest struct {
	FileList                 []string    `json:"FileList"`
	UseCompatibleDecompiler  bool        `json:"UseCompatibleDecompiler"`
	EncodingType             interface{} `json:"EncodingType"`
}

// ItemCodesRequest payload for ItemCodesToFileInfos.
type ItemCodesRequest struct {
	LstNames  []string `json:"lstNames"`
	ItemCodes []int    `json:"ItemCodes"`
}

// ItemCodesResponse from ItemCodesToFileInfos.
type ItemCodesResponse struct {
	Infos map[string]ItemCodeFileInfo `json:"Infos"`
}

// FileListToLstResponse maps lst file paths to their entries.
type FileListToLstResponse map[string]map[string]string
