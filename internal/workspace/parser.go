package workspace

import (
	"fmt"
	"regexp"
	"strings"
)

// ParseWorkspace parses a workspace file content string into a Workspace struct
func ParseWorkspace(content string) (*Workspace, error) {
	if content == "" {
		return nil, fmt.Errorf("workspace content cannot be empty")
	}

	workspace := NewWorkspace()
	lines := strings.Split(content, "\n")
	
	// State machine for parsing
	state := "start"
	var productLines []string
	var toolLines []string
	
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		switch state {
		case "start":
			if err := parseVersionLine(line, workspace); err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNum+1, err)
			}
			state = "header"
			
		case "header":
			if strings.HasPrefix(line, "organization ") {
				if err := parseOrganizationLine(line, workspace); err != nil {
					return nil, fmt.Errorf("line %d: %w", lineNum+1, err)
				}
			} else if strings.HasPrefix(line, "products ") {
				if err := parseProductsStartLine(line); err != nil {
					return nil, fmt.Errorf("line %d: %w", lineNum+1, err)
				}
				state = "products"
			} else if strings.HasPrefix(line, "tools ") {
				if err := parseToolsStartLine(line); err != nil {
					return nil, fmt.Errorf("line %d: %w", lineNum+1, err)
				}
				state = "tools"
			} else {
				return nil, fmt.Errorf("line %d: unexpected line in header section: %s", lineNum+1, line)
			}
			
		case "products":
			if line == ")" {
				state = "header"
			} else {
				productLines = append(productLines, line)
			}
			
		case "tools":
			if line == ")" {
				state = "header"
			} else {
				toolLines = append(toolLines, line)
			}
			
		case "end":
			return nil, fmt.Errorf("line %d: unexpected content after end: %s", lineNum+1, line)
		}
	}
	
	// Parse collected product lines
	if len(productLines) > 0 {
		if err := parseProductLines(productLines, workspace); err != nil {
			return nil, err
		}
	}
	
	// Parse collected tool lines
	if len(toolLines) > 0 {
		if err := parseToolLines(toolLines, workspace); err != nil {
			return nil, err
		}
	}
	
	// Validate final state
	if state == "products" {
		return nil, fmt.Errorf("products section not properly closed with ')'")
	}
	if state == "tools" {
		return nil, fmt.Errorf("tools section not properly closed with ')'")
	}
	
	return workspace, nil
}

// parseVersionLine parses the version line (e.g., "nimsforest 1.0")
func parseVersionLine(line string, workspace *Workspace) error {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return fmt.Errorf("invalid version line format, expected 'nimsforest <version>', got: %s", line)
	}
	
	if parts[0] != "nimsforest" {
		return fmt.Errorf("invalid version line, expected 'nimsforest', got: %s", parts[0])
	}
	
	// Validate version format (simple pattern: number.number)
	versionRegex := regexp.MustCompile(`^\d+\.\d+$`)
	if !versionRegex.MatchString(parts[1]) {
		return fmt.Errorf("invalid version format, expected format like '1.0', got: %s", parts[1])
	}
	
	workspace.Version = parts[1]
	return nil
}

// parseOrganizationLine parses the organization line (e.g., "organization ./acme-organization-workspace")
func parseOrganizationLine(line string, workspace *Workspace) error {
	parts := strings.Fields(line)
	if len(parts) != 2 {
		return fmt.Errorf("invalid organization line format, expected 'organization <path>', got: %s", line)
	}
	
	if parts[0] != "organization" {
		return fmt.Errorf("invalid organization line, expected 'organization', got: %s", parts[0])
	}
	
	workspace.Organization = parts[1]
	return nil
}

// parseProductsStartLine parses the products start line (e.g., "products (")
func parseProductsStartLine(line string) error {
	// Allow for flexible spacing around parentheses
	line = strings.TrimSpace(line)
	
	// Check if it's exactly "products (" or "products("
	if line == "products (" || line == "products(" {
		return nil
	}
	
	// Check if it starts with "products" and contains an opening parenthesis
	if strings.HasPrefix(line, "products") {
		remaining := strings.TrimSpace(line[8:]) // Remove "products"
		if remaining == "(" {
			return nil
		}
	}
	
	return fmt.Errorf("invalid products section start, expected 'products (', got: %s", line)
}

// parseProductLines parses the product lines within the products section
func parseProductLines(lines []string, workspace *Workspace) error {
	for i, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Remove leading whitespace that might be used for indentation
		productPath := strings.TrimSpace(line)
		
		// Validate that it's not empty after trimming
		if productPath == "" {
			return fmt.Errorf("empty product path at line %d in products section", i+1)
		}
		
		// Add the product to the workspace
		workspace.AddProduct(productPath)
	}
	
	return nil
}

// parseToolsStartLine parses the tools start line (e.g., "tools (")
func parseToolsStartLine(line string) error {
	// Allow for flexible spacing around parentheses
	line = strings.TrimSpace(line)
	
	// Check if it's exactly "tools (" or "tools("
	if line == "tools (" || line == "tools(" {
		return nil
	}
	
	// Check if it starts with "tools" and contains an opening parenthesis
	if strings.HasPrefix(line, "tools") {
		remaining := strings.TrimSpace(line[5:]) // Remove "tools"
		if remaining == "(" {
			return nil
		}
	}
	
	return fmt.Errorf("invalid tools section start, expected 'tools (', got: %s", line)
}

// parseToolLines parses the tool lines within the tools section
func parseToolLines(lines []string, workspace *Workspace) error {
	for i, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		// Parse tool line: "toolname mode path version"
		parts := strings.Fields(line)
		if len(parts) != 4 {
			return fmt.Errorf("invalid tool line format at line %d, expected 'name mode path version', got: %s", i+1, line)
		}
		
		tool := ToolEntry{
			Name:    parts[0],
			Mode:    parts[1],
			Path:    parts[2],
			Version: parts[3],
		}
		
		// Validate mode
		if tool.Mode != "binary" && tool.Mode != "clone" && tool.Mode != "submodule" {
			return fmt.Errorf("invalid tool mode '%s' at line %d, expected 'binary', 'clone', or 'submodule'", tool.Mode, i+1)
		}
		
		// Add the tool to the workspace
		workspace.AddTool(tool)
	}
	
	return nil
}

// ValidateWorkspaceFormat performs basic format validation on workspace content
func ValidateWorkspaceFormat(content string) error {
	if content == "" {
		return fmt.Errorf("workspace content cannot be empty")
	}
	
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("workspace content cannot be empty")
	}
	
	// Check for version line
	found := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		if strings.HasPrefix(line, "nimsforest ") {
			found = true
			break
		} else {
			return fmt.Errorf("first non-empty, non-comment line must be version line starting with 'nimsforest'")
		}
	}
	
	if !found {
		return fmt.Errorf("version line starting with 'nimsforest' not found")
	}
	
	return nil
}

// NormalizeWorkspaceContent normalizes workspace content by:
// - Removing extra whitespace
// - Ensuring consistent formatting
// - Preserving comments
func NormalizeWorkspaceContent(content string) string {
	lines := strings.Split(content, "\n")
	var normalized []string
	
	inProducts := false
	inTools := false
	
	for _, line := range lines {
		original := line
		line = strings.TrimSpace(line)
		
		// Skip empty lines but preserve them in output
		if line == "" {
			normalized = append(normalized, "")
			continue
		}
		
		// Preserve comments as-is
		if strings.HasPrefix(line, "#") {
			normalized = append(normalized, original)
			continue
		}
		
		// Handle different line types
		if strings.HasPrefix(line, "nimsforest ") {
			normalized = append(normalized, line)
		} else if strings.HasPrefix(line, "organization ") {
			normalized = append(normalized, line)
		} else if strings.HasPrefix(line, "products ") {
			normalized = append(normalized, "products (")
			inProducts = true
		} else if strings.HasPrefix(line, "tools ") {
			normalized = append(normalized, "tools (")
			inTools = true
		} else if line == ")" && (inProducts || inTools) {
			normalized = append(normalized, ")")
			inProducts = false
			inTools = false
		} else if inProducts {
			// Indent product lines
			normalized = append(normalized, fmt.Sprintf("    %s", line))
		} else if inTools {
			// Indent tool lines
			normalized = append(normalized, fmt.Sprintf("    %s", line))
		} else {
			normalized = append(normalized, line)
		}
	}
	
	return strings.Join(normalized, "\n")
}