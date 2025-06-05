package parser

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type HelmRelease struct {
	ResourceName string
	Name         string
	Chart        string
	Version      string
	Repository   string
	Namespace    string
	Values       []string
	Sets         map[string]string
}

func ParseHelmRelease(resourceName string, block *hclsyntax.Block, src []byte, evalCtx *hcl.EvalContext) HelmRelease {
	helm := HelmRelease{
		ResourceName: resourceName,
		Sets:         make(map[string]string),
	}

	// Extract basic attributes
	for attrName, attr := range block.Body.Attributes {
		if attr.Expr != nil {
			val, diag := attr.Expr.Value(evalCtx)
			var value string
			if diag.HasErrors() {
				value = strings.Trim(string(attr.Expr.Range().SliceBytes(src)), `"`) // fallback
			} else if val.Type() == cty.String {
				value = val.AsString()
			} else {
				value = val.GoString()
			}
			switch attrName {
			case "name":
				helm.Name = value
			case "chart":
				helm.Chart = value
			case "version":
				helm.Version = value
			case "repository":
				helm.Repository = value
			case "namespace":
				helm.Namespace = value
			}
		}
	}

	// Extract nested blocks (values, set blocks)
	for _, nestedBlock := range block.Body.Blocks {
		switch nestedBlock.Type {
		case "set":
			// Handle set blocks for helm values
			var setName, setValue string
			for attrName, attr := range nestedBlock.Body.Attributes {
				if attr.Expr != nil {
					val, diag := attr.Expr.Value(evalCtx)
					var value string
					if diag.HasErrors() {
						value = strings.Trim(string(attr.Expr.Range().SliceBytes(src)), `"`)
					} else if val.Type() == cty.String {
						value = val.AsString()
					} else {
						value = val.GoString()
					}
					switch attrName {
					case "name":
						setName = value
					case "value":
						setValue = value
					}
				}
			}
			if setName != "" {
				helm.Sets[setName] = setValue
			}
		case "values":
			// Handle values blocks - simplified extraction
			for _, attr := range nestedBlock.Body.Attributes {
				if attr.Expr != nil {
					val, diag := attr.Expr.Value(evalCtx)
					if diag.HasErrors() {
						helm.Values = append(helm.Values, string(attr.Expr.Range().SliceBytes(src)))
					} else if val.Type() == cty.String {
						helm.Values = append(helm.Values, val.AsString())
					} else {
						helm.Values = append(helm.Values, val.GoString())
					}
				}
			}
		}
	}

	return helm
}
