package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"tf-diff/fetch"
)

type TFVariable struct {
	Name    string
	Default string
	Type    string
}

func ParseVariablesTF(content []byte, path string) (map[string]TFVariable, error) {

	// Create temp file to parse the HCL content
	tmpFile, err := fetch.CreateTempFile(content, path)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(filepath.Dir(tmpFile))

	// // Read the file content, later used for extracting default and type values
	// src, err := os.ReadFile(path)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to read file %s: %v", path, err)
	// }

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(tmpFile)

	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse file %s: %v", path, diags)
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("file %s does not contain a valid HCL body", path)
	}

	vars := make(map[string]TFVariable)
	for _, block := range body.Blocks {
		if block.Type == "variable" && len(block.Labels) == 1 {
			varName := block.Labels[0]
			variable := TFVariable{Name: varName}

			// Iterate over attributes to find default and type
			for _, attr := range block.Body.Attributes {
				switch attr.Name {
				case "default":
					if attr.Expr != nil {
						variable.Default = string(attr.Expr.Range().SliceBytes(content))
					}
				case "type":
					if attr.Expr != nil {
						variable.Type = string(attr.Expr.Range().SliceBytes(content))
					}
				}
			}

			vars[varName] = variable
		}
	}

	// for _, v := range vars {
	// 	fmt.Printf("Variable: %s, Type: %s, Default: %s\n", v.Name, v.Type, v.Default)
	// }

	return vars, nil
}
