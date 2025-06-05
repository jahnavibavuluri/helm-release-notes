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

type TFResource struct {
	Type       string
	Name       string
	Config     map[string]string
	SourceLine int
}

func ParseMainTF(content []byte, path string, evalCtx *hcl.EvalContext) (map[string]TFResource, []HelmRelease, error) {
	tmpFile, err := fetch.CreateTempFile(content, path)
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(filepath.Dir(tmpFile))

	// src, err := os.ReadFile(path)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to read file %s: %v", path, err)
	// }

	parser := hclparse.NewParser()
	file, diags := parser.ParseHCLFile(tmpFile)

	if diags.HasErrors() {
		return nil, nil, fmt.Errorf("failed to parse file %s: %v", path, diags)
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, nil, fmt.Errorf("file %s does not contain a valid HCL body", path)
	}

	resources := make(map[string]TFResource)
	var helmReleases []HelmRelease

	for _, block := range body.Blocks {
		if block.Type == "resource" && len(block.Labels) == 2 {
			resourceType := block.Labels[0]
			resourceName := block.Labels[1]
			resourceKey := fmt.Sprintf("%s.%s", resourceType, resourceName)

			resource := TFResource{
				Type:       resourceType,
				Name:       resourceName,
				Config:     make(map[string]string),
				SourceLine: block.DefRange().Start.Line,
			}

			// Extract key attributes
			for attrName, attr := range block.Body.Attributes {
				if attr.Expr != nil {
					// Evaluate the expression in the context of evalCtx
					val, diag := attr.Expr.Value(evalCtx)
					if diag.HasErrors() {
						resource.Config[attrName] = string(attr.Expr.Range().SliceBytes(content)) // fallback
					} else if val.Type() == cty.String {
						resource.Config[attrName] = val.AsString()
					} else {
						resource.Config[attrName] = val.GoString() // evaluated value
					}

				}
			}

			resources[resourceKey] = resource

			// Special handling for helm_release resources
			if resourceType == "helm_release" {
				helmRelease := ParseHelmRelease(resourceName, block, content, evalCtx)
				helmReleases = append(helmReleases, helmRelease)
			}
		}
	}

	// for _, vars := range resources {
	// 	fmt.Println("Resource:", vars.Type, vars.Name)
	// 	for key, value := range vars.Config {
	// 		fmt.Printf("  %s: %s\n", key, value)
	// 	}
	// }

	return resources, helmReleases, nil
}
