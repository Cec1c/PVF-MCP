package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const pvfBaseURL = "http://localhost:27000/Api/PvfUtiltiy"

// PvfClient wraps the pvfUtility HTTP API.
type PvfClient struct {
	baseURL string
	client  *http.Client
}

// NewPvfClient creates a new client.
func NewPvfClient() *PvfClient {
	return &PvfClient{
		baseURL: pvfBaseURL,
		client:  &http.Client{},
	}
}

func (c *PvfClient) get(path string, query url.Values) (*ApiResponse, error) {
	u := c.baseURL + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}
	resp, err := c.client.Get(u)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", path, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	var ar ApiResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w\nbody: %s", path, err, string(body))
	}
	if ar.IsError && ar.Msg != nil {
		return &ar, fmt.Errorf("%s", *ar.Msg)
	}
	return &ar, nil
}

func (c *PvfClient) post(path string, body interface{}) (*ApiResponse, error) {
	return c.postRaw(path, body, nil)
}

func (c *PvfClient) postRaw(path string, body interface{}, rawBody []byte) (*ApiResponse, error) {
	u := c.baseURL + path
	var reqBody io.Reader
	if rawBody != nil {
		reqBody = bytes.NewReader(rawBody)
	} else if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}
	resp, err := c.client.Post(u, "application/json", reqBody)
	if err != nil {
		return nil, fmt.Errorf("POST %s: %w", path, err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	var ar ApiResponse
	if err := json.Unmarshal(respBody, &ar); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %w\nbody: %s", path, err, string(respBody))
	}
	if ar.IsError && ar.Msg != nil {
		return &ar, fmt.Errorf("%s", *ar.Msg)
	}
	return &ar, nil
}

// ------ Read APIs ------

func (c *PvfClient) GetVersion() (string, error) {
	ar, err := c.get("/getVersion", nil)
	if err != nil {
		return "", err
	}
	var s string
	if err := json.Unmarshal(ar.Data, &s); err != nil {
		return "", err
	}
	return s, nil
}

func (c *PvfClient) GetPvfRootDirectory() ([]string, error) {
	ar, err := c.get("/getPvfRootDirectory", nil)
	if err != nil {
		return nil, err
	}
	var dirs []string
	if err := json.Unmarshal(ar.Data, &dirs); err != nil {
		return nil, err
	}
	return dirs, nil
}

func (c *PvfClient) GetFileList(dirName, fileType string) ([]string, error) {
	q := url.Values{}
	q.Set("dirName", dirName)
	q.Set("returnType", "0")
	if fileType != "" {
		q.Set("fileType", fileType)
	}
	ar, err := c.get("/GetFileList", q)
	if err != nil {
		return nil, err
	}
	var files []string
	if err := json.Unmarshal(ar.Data, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (c *PvfClient) GetFileContent(filePath, encodingType string, useCompatibleDecompiler bool) (string, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	if useCompatibleDecompiler {
		q.Set("useCompatibleDecompiler", "true")
	} else {
		q.Set("useCompatibleDecompiler", "false")
	}
	if encodingType != "" {
		q.Set("encodingType", encodingType)
	}
	ar, err := c.get("/GetFileContent", q)
	if err != nil {
		return "", err
	}
	// Data is either a string or a structured object
	var s string
	if err := json.Unmarshal(ar.Data, &s); err != nil {
		// Try structured
		return string(ar.Data), nil
	}
	return s, nil
}

func (c *PvfClient) GetFileData(filePath string) ([]WebApiFileData, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	ar, err := c.get("/getFileData", q)
	if err != nil {
		return nil, err
	}
	var nodes []WebApiFileData
	if err := json.Unmarshal(ar.Data, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func (c *PvfClient) GetFileContents(fileList []string, useCompatibleDecompiler bool, encodingType string) (*ApiResponse, error) {
	req := GetFileContentsRequest{
		FileList:                fileList,
		UseCompatibleDecompiler: useCompatibleDecompiler,
		EncodingType:            nil,
	}
	if encodingType != "" {
		req.EncodingType = encodingType
	}
	return c.post("/GetFileContents", req)
}

func (c *PvfClient) SearchPvf(params SearchPvfRequest) ([]string, error) {
	ar, err := c.post("/SearchPvf", params)
	if err != nil {
		return nil, err
	}
	var files []string
	if err := json.Unmarshal(ar.Data, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (c *PvfClient) GetItemInfo(filePath string) (*ItemInfo, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	ar, err := c.get("/GetItemInfo", q)
	if err != nil {
		return nil, err
	}
	var info ItemInfo
	if err := json.Unmarshal(ar.Data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *PvfClient) GetItemInfos(filePaths []string) (map[string]ItemInfo, error) {
	ar, err := c.post("/GetItemInfos", filePaths)
	if err != nil {
		return nil, err
	}
	var infos map[string]ItemInfo
	if err := json.Unmarshal(ar.Data, &infos); err != nil {
		return nil, err
	}
	return infos, nil
}

func (c *PvfClient) ItemCodeToFileInfo(lstNames string, itemCode int) (*ItemCodeFileInfo, error) {
	q := url.Values{}
	q.Set("lstNames", lstNames)
	q.Set("itemCode", fmt.Sprintf("%d", itemCode))
	ar, err := c.get("/ItemCodeToFileInfo", q)
	if err != nil {
		return nil, err
	}
	var info ItemCodeFileInfo
	if err := json.Unmarshal(ar.Data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *PvfClient) ItemCodesToFileInfos(lstNames []string, itemCodes []int) (*ItemCodesResponse, error) {
	req := ItemCodesRequest{
		LstNames:  lstNames,
		ItemCodes: itemCodes,
	}
	ar, err := c.post("/ItemCodesToFileInfos", req)
	if err != nil {
		return nil, err
	}
	var resp ItemCodesResponse
	if err := json.Unmarshal(ar.Data, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *PvfClient) GetAllLstFileList() ([]string, error) {
	ar, err := c.get("/GetAllLstFileList", nil)
	if err != nil {
		return nil, err
	}
	var lsts []string
	if err := json.Unmarshal(ar.Data, &lsts); err != nil {
		return nil, err
	}
	return lsts, nil
}

func (c *PvfClient) GetLstFileInfo(filePath string) (map[string]LstFileInfo, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	ar, err := c.get("/getLstFileInfo", q)
	if err != nil {
		return nil, err
	}
	var infos map[string]LstFileInfo
	if err := json.Unmarshal(ar.Data, &infos); err != nil {
		return nil, err
	}
	return infos, nil
}

func (c *PvfClient) FileListToLstRows(filePaths []string) (FileListToLstResponse, error) {
	ar, err := c.post("/FileListToLstRows", filePaths)
	if err != nil {
		return nil, err
	}
	var resp FileListToLstResponse
	if err := json.Unmarshal(ar.Data, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *PvfClient) GetFileIcon(filePath string) (string, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	ar, err := c.get("/getFileIcon", q)
	if err != nil {
		return "", err
	}
	var s string
	if err := json.Unmarshal(ar.Data, &s); err != nil {
		return "", err
	}
	return s, nil
}

func (c *PvfClient) FilesToIconBase64(filePaths []string) (map[string]string, error) {
	ar, err := c.post("/filesToIconBase64", filePaths)
	if err != nil {
		return nil, err
	}
	var icons map[string]string
	if err := json.Unmarshal(ar.Data, &icons); err != nil {
		return nil, err
	}
	return icons, nil
}

func (c *PvfClient) GetActiveDocumentFilePath() (string, error) {
	ar, err := c.get("/GetActiveDocumentFilePath", nil)
	if err != nil {
		return "", err
	}
	var s string
	if err := json.Unmarshal(ar.Data, &s); err != nil {
		return "", err
	}
	return s, nil
}

func (c *PvfClient) GetTreeSelectedFiles() ([]string, error) {
	ar, err := c.get("/GetTreeSelectedFiles", nil)
	if err != nil {
		return nil, err
	}
	var files []string
	if err := json.Unmarshal(ar.Data, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (c *PvfClient) GetSearchPanelSelectedFiles() ([]string, error) {
	ar, err := c.get("/GetSearchPanelSelectedFiles", nil)
	if err != nil {
		return nil, err
	}
	var files []string
	if err := json.Unmarshal(ar.Data, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func (c *PvfClient) GetStringTable() ([]string, error) {
	ar, err := c.get("/getStringTable", nil)
	if err != nil {
		return nil, err
	}
	var tbl []string
	if err := json.Unmarshal(ar.Data, &tbl); err != nil {
		return nil, err
	}
	return tbl, nil
}

func (c *PvfClient) GetPvfPackFilePath() (string, error) {
	ar, err := c.get("/GetPvfPackFilePath", nil)
	if err != nil {
		return "", err
	}
	var s string
	if err := json.Unmarshal(ar.Data, &s); err != nil {
		return "", err
	}
	return s, nil
}

func (c *PvfClient) FileIsExists(filePath string) (bool, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	ar, err := c.get("/FileIsExists", q)
	if err != nil {
		return false, err
	}
	return !ar.IsError, nil
}

func (c *PvfClient) FolderExists(filePath string) (bool, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	ar, err := c.get("/folderExists", q)
	if err != nil {
		return false, err
	}
	return !ar.IsError, nil
}

func (c *PvfClient) GetTreeListFocusedFilePath() (string, error) {
	ar, err := c.get("/GetTreeListFocusedFilePath", nil)
	if err != nil {
		return "", err
	}
	var s string
	if err := json.Unmarshal(ar.Data, &s); err != nil {
		return "", err
	}
	return s, nil
}

func (c *PvfClient) GetSearchPanelTreeListFocusedFilePath() (string, error) {
	ar, err := c.get("/GetSearchPanelTreeListFocusedFilePath", nil)
	if err != nil {
		return "", err
	}
	var s string
	if err := json.Unmarshal(ar.Data, &s); err != nil {
		return "", err
	}
	return s, nil
}

// ------ Write APIs ------

func (c *PvfClient) ImportFile(filePath string, content string) (*ApiResponse, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	path := "/ImportFile?" + q.Encode()
	return c.postRaw(path, nil, []byte(content))
}

func (c *PvfClient) ImportFiles(entries []FileContentEntry) (*ApiResponse, error) {
	return c.post("/ImportFiles", entries)
}

func (c *PvfClient) DeleteFile(filePath string) (*ApiResponse, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	return c.get("/DeleteFile", q)
}

func (c *PvfClient) DeleteFiles(filePaths []string) (*ApiResponse, error) {
	return c.post("/DeleteFiles", filePaths)
}

func (c *PvfClient) SaveAsPvfFile(filePath string) (*ApiResponse, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	return c.get("/SaveAsPvfFile", q)
}

// GoToTreeListNode navigates GUI to a file and optionally opens it.
func (c *PvfClient) GoToTreeListNode(filePath string, openDocument bool) (*ApiResponse, error) {
	q := url.Values{}
	q.Set("filePath", filePath)
	if openDocument {
		q.Set("openTextDocument", "1")
	} else {
		q.Set("openTextDocument", "0")
	}
	return c.get("/goToTreeListNode", q)
}

// SearchPvfSimple builds a simple keyword search from defaults.
func (c *PvfClient) SearchPvfSimple(keyword, searchFolder, fileTypesString string) ([]string, error) {
	params := SearchPvfRequest{
		SearchFolder:            searchFolder,
		Keyword:                 keyword,
		Type:                    1,
		SourceType:              0,
		NormalUsing:             1,
		IsStartMatch:            false,
		SearchResult:            nil,
		ScriptContentSearchMode: 1,
		IsUseLikeSearchPath:     false,
		Trait:                   false,
		UseRegularExpression:    false,
		WholeWordMatch:          false,
		RemoveOrKeep:            1,
		FileTypesString:         nil,
		ScriptContent:           "",
		ScriptContentStart:      "",
		ScriptContentStop:       "",
	}
	if fileTypesString != "" {
		params.FileTypesString = fileTypesString
	}
	return c.SearchPvf(params)
}

// urlEncode encodes a string for use in a URL path segment.
func urlEncode(s string) string {
	// Use url.PathEscape for path segments, but keep '/' intact
	parts := strings.Split(s, "/")
	for i, p := range parts {
		parts[i] = url.PathEscape(p)
	}
	return strings.Join(parts, "/")
}
