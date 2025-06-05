package parser

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func ParseLocalsTF(path string, evalCtx *hcl.EvalContext) (map[string]cty.Value, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", path, err)
	}

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(path)
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
						locals[name] = cty.StringVal(string(attr.Expr.Range().SliceBytes(src)))
					} else {
						locals[name] = val
					}
				}
			}
		}
	}
	return locals, nil
}
