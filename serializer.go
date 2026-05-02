package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SerializeWebApiFileData converts a structured WebApiFileData tree back to PVF text format.
// This enables: get_file_data → modify JSON → serialize → import_file workflow.
// The output includes the standard #PVF_File header.
func SerializeWebApiFileData(nodes []WebApiFileData, indent string) string {
	var sb strings.Builder
	sb.WriteString("#PVF_File\r\n")
	for _, node := range nodes {
		writeNode(&sb, node, indent)
	}
	return sb.String()
}

// stripBrackets removes surrounding [ ] from a section name if present.
// pvfUtility's API returns SectionName as "[name]" — we need just "name".
func stripBrackets(name string) string {
	name = strings.TrimPrefix(name, "[")
	name = strings.TrimSuffix(name, "]")
	return name
}

func writeNode(sb *strings.Builder, node WebApiFileData, indent string) {
	if node.IsSection {
		name := stripBrackets(*node.SectionName)
		sb.WriteString(indent)
		sb.WriteString(fmt.Sprintf("[%s]\r\n", name))
		for _, child := range node.Children {
			writeNode(sb, child, indent+"\t")
		}
		if node.HasEndSection {
			sb.WriteString(indent)
			sb.WriteString(fmt.Sprintf("[/%s]\r\n", name))
		}
		return
	}

	// Data value
	sb.WriteString(indent)
	if node.Value != nil {
		var raw string
		if err := json.Unmarshal(node.Value, &raw); err == nil {
			raw = strings.TrimSpace(raw)
		}
		switch node.DataType {
		case ScriptTypeString:
			sb.WriteString(fmt.Sprintf("`%s`\r\n", raw))
		case ScriptTypeInt, ScriptTypeIntEx:
			sb.WriteString(fmt.Sprintf("%s\r\n", raw))
		case ScriptTypeFloat:
			sb.WriteString(fmt.Sprintf("%s\r\n", raw))
		case ScriptTypeCommand:
			sb.WriteString(fmt.Sprintf("`%s`\r\n", raw))
		default:
			if raw != "" {
				sb.WriteString(fmt.Sprintf("%s\r\n", raw))
			}
		}
	}

	// If node has children but is not a section, write them without extra indent
	for _, child := range node.Children {
		writeNode(sb, child, indent)
	}
}

// ParseRawValue extracts the raw string value from a WebApiFileData node's Value field.
func ParseRawValue(node WebApiFileData) string {
	if node.Value == nil {
		return ""
	}
	var s string
	if err := json.Unmarshal(node.Value, &s); err != nil {
		return string(node.Value)
	}
	return strings.TrimSpace(s)
}

// FindNodeByPath searches a WebApiFileData tree for a section by dot-separated path.
// E.g. "shop.goods" finds [shop] → [goods].
func FindNodeByPath(nodes []WebApiFileData, path string) *WebApiFileData {
	parts := strings.Split(path, ".")
	var current []WebApiFileData = nodes
	var lastMatch *WebApiFileData
	for _, part := range parts {
		found := false
		for i := range current {
			if current[i].IsSection && current[i].SectionName != nil && *current[i].SectionName == part {
				lastMatch = &current[i]
				current = current[i].Children
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return lastMatch
}

// FindNodeBySection finds the first child node with the given section name.
func FindNodeBySection(nodes []WebApiFileData, sectionName string) *WebApiFileData {
	for i := range nodes {
		if nodes[i].IsSection && nodes[i].SectionName != nil && *nodes[i].SectionName == sectionName {
			return &nodes[i]
		}
	}
	return nil
}

// UpdateNodeValue updates the first child value in a section node.
func UpdateNodeValue(node *WebApiFileData, newValue string) {
	for i := range node.Children {
		if !node.Children[i].IsSection {
			raw, _ := json.Marshal(newValue)
			node.Children[i].Value = raw
			return
		}
	}
	// No existing child value — append one
	raw, _ := json.Marshal(newValue)
	node.Children = append(node.Children, WebApiFileData{
		DataType: ScriptTypeString,
		Value:    raw,
	})
}
