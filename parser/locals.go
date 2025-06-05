package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"tf-diff/fetch"
)

func ParseLocalsTF(content []byte, path string, evalCtx *hcl.EvalContext) (map[string]cty.Value, error) {
	tmpFile, err := fetch.CreateTempFile(content, path)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(filepath.Dir(tmpFile))

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

	locals := make(map[string]cty.Value)
	for _, block := range body.Blocks {
		if block.Type == "locals" {
			for name, attr := range block.Body.Attributes {
				if attr.Expr != nil {
					val, diag := attr.Expr.Value(evalCtx)
					if diag.HasErrors() {
						locals[name] = cty.StringVal(string(attr.Expr.Range().SliceBytes(content)))
					} else {
						locals[name] = val
					}
				}
			}
		}
	}
	return locals, nil
}
