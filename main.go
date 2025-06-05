package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"tf-diff/diff"
	"tf-diff/parser"

	"github.com/zclconf/go-cty/cty"
)

// type TFVariable struct {
// 	Name    string
// 	Default string
// 	Type    string
// }

// type TFResource struct {
// 	Type       string
// 	Name       string
// 	Config     map[string]string
// 	SourceLine int
// }

// type HelmRelease struct {
// 	ResourceName string
// 	Name         string
// 	Chart        string
// 	Version      string
// 	Repository   string
// 	Namespace    string
// 	Values       []string
// 	Sets         map[string]string
// }

// This is used to build the evaluation context for HCL expressions.
// Evalutes expressions in the context of variables and locals.
// func buildEvalContext(vars map[string]TFVariable, locals map[string]cty.Value) *hcl.EvalContext {
// 	varMap := map[string]cty.Value{}
// 	for name, variable := range vars {
// 		val := strings.Trim(variable.Default, `"`)
// 		varMap[name] = cty.StringVal(val)
// 	}

// 	return &hcl.EvalContext{
// 		Variables: map[string]cty.Value{
// 			"var":   cty.ObjectVal(varMap),
// 			"local": cty.ObjectVal(locals),
// 		},
// 	}
// }

// func parseVariablesTF(path string) (map[string]TFVariable, error) {

// 	// Read the file content, later used for extracting default and type values
// 	src, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read file %s: %v", path, err)
// 	}

// 	parser := hclparse.NewParser()
// 	file, diags := parser.ParseHCLFile(path)

// 	if diags.HasErrors() {
// 		return nil, fmt.Errorf("failed to parse file %s: %v", path, diags)
// 	}

// 	body, ok := file.Body.(*hclsyntax.Body)
// 	if !ok {
// 		return nil, fmt.Errorf("file %s does not contain a valid HCL body", path)
// 	}

// 	vars := make(map[string]TFVariable)
// 	for _, block := range body.Blocks {
// 		if block.Type == "variable" && len(block.Labels) == 1 {
// 			varName := block.Labels[0]
// 			variable := TFVariable{Name: varName}

// 			// Iterate over attributes to find default and type
// 			for _, attr := range block.Body.Attributes {
// 				switch attr.Name {
// 				case "default":
// 					if attr.Expr != nil {
// 						variable.Default = string(attr.Expr.Range().SliceBytes(src))
// 					}
// 				case "type":
// 					if attr.Expr != nil {
// 						variable.Type = string(attr.Expr.Range().SliceBytes(src))
// 					}
// 				}
// 			}

// 			vars[varName] = variable
// 		}
// 	}

// 	// for _, v := range vars {
// 	// 	fmt.Printf("Variable: %s, Type: %s, Default: %s\n", v.Name, v.Type, v.Default)
// 	// }

// 	return vars, nil
// }

// func parseMainTF(path string, evalCtx *hcl.EvalContext) (map[string]TFResource, []HelmRelease, error) {
// 	src, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("failed to read file %s: %v", path, err)
// 	}

// 	parser := hclparse.NewParser()
// 	file, diags := parser.ParseHCLFile(path)

// 	if diags.HasErrors() {
// 		return nil, nil, fmt.Errorf("failed to parse file %s: %v", path, diags)
// 	}

// 	body, ok := file.Body.(*hclsyntax.Body)
// 	if !ok {
// 		return nil, nil, fmt.Errorf("file %s does not contain a valid HCL body", path)
// 	}

// 	resources := make(map[string]TFResource)
// 	var helmReleases []HelmRelease

// 	for _, block := range body.Blocks {
// 		if block.Type == "resource" && len(block.Labels) == 2 {
// 			resourceType := block.Labels[0]
// 			resourceName := block.Labels[1]
// 			resourceKey := fmt.Sprintf("%s.%s", resourceType, resourceName)

// 			resource := TFResource{
// 				Type:       resourceType,
// 				Name:       resourceName,
// 				Config:     make(map[string]string),
// 				SourceLine: block.DefRange().Start.Line,
// 			}

// 			// Extract key attributes
// 			for attrName, attr := range block.Body.Attributes {
// 				if attr.Expr != nil {
// 					// Evaluate the expression in the context of evalCtx
// 					val, diag := attr.Expr.Value(evalCtx)
// 					if diag.HasErrors() {
// 						resource.Config[attrName] = string(attr.Expr.Range().SliceBytes(src)) // fallback
// 					} else if val.Type() == cty.String {
// 						resource.Config[attrName] = val.AsString()
// 					} else {
// 						resource.Config[attrName] = val.GoString() // evaluated value
// 					}

// 				}
// 			}

// 			resources[resourceKey] = resource

// 			// Special handling for helm_release resources
// 			if resourceType == "helm_release" {
// 				helmRelease := parseHelmRelease(resourceName, block, src, evalCtx)
// 				helmReleases = append(helmReleases, helmRelease)
// 			}
// 		}
// 	}

// 	// for _, vars := range resources {
// 	// 	fmt.Println("Resource:", vars.Type, vars.Name)
// 	// 	for key, value := range vars.Config {
// 	// 		fmt.Printf("  %s: %s\n", key, value)
// 	// 	}
// 	// }

// 	return resources, helmReleases, nil
// }

// func parseHelmRelease(resourceName string, block *hclsyntax.Block, src []byte, evalCtx *hcl.EvalContext) HelmRelease {
// 	helm := HelmRelease{
// 		ResourceName: resourceName,
// 		Sets:         make(map[string]string),
// 	}

// 	// Extract basic attributes
// 	for attrName, attr := range block.Body.Attributes {
// 		if attr.Expr != nil {
// 			val, diag := attr.Expr.Value(evalCtx)
// 			var value string
// 			if diag.HasErrors() {
// 				value = strings.Trim(string(attr.Expr.Range().SliceBytes(src)), `"`) // fallback
// 			} else if val.Type() == cty.String {
// 				value = val.AsString()
// 			} else {
// 				value = val.GoString()
// 			}
// 			switch attrName {
// 			case "name":
// 				helm.Name = value
// 			case "chart":
// 				helm.Chart = value
// 			case "version":
// 				helm.Version = value
// 			case "repository":
// 				helm.Repository = value
// 			case "namespace":
// 				helm.Namespace = value
// 			}
// 		}
// 	}

// 	// Extract nested blocks (values, set blocks)
// 	for _, nestedBlock := range block.Body.Blocks {
// 		switch nestedBlock.Type {
// 		case "set":
// 			// Handle set blocks for helm values
// 			var setName, setValue string
// 			for attrName, attr := range nestedBlock.Body.Attributes {
// 				if attr.Expr != nil {
// 					val, diag := attr.Expr.Value(evalCtx)
// 					var value string
// 					if diag.HasErrors() {
// 						value = strings.Trim(string(attr.Expr.Range().SliceBytes(src)), `"`)
// 					} else if val.Type() == cty.String {
// 						value = val.AsString()
// 					} else {
// 						value = val.GoString()
// 					}
// 					switch attrName {
// 					case "name":
// 						setName = value
// 					case "value":
// 						setValue = value
// 					}
// 				}
// 			}
// 			if setName != "" {
// 				helm.Sets[setName] = setValue
// 			}
// 		case "values":
// 			// Handle values blocks - simplified extraction
// 			for _, attr := range nestedBlock.Body.Attributes {
// 				if attr.Expr != nil {
// 					val, diag := attr.Expr.Value(evalCtx)
// 					if diag.HasErrors() {
// 						helm.Values = append(helm.Values, string(attr.Expr.Range().SliceBytes(src)))
// 					} else if val.Type() == cty.String {
// 						helm.Values = append(helm.Values, val.AsString())
// 					} else {
// 						helm.Values = append(helm.Values, val.GoString())
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return helm
// }

// func diffVariables(oldVars, newVars map[string]TFVariable) {
// 	for name, oldVar := range oldVars {
// 		newVar, exists := newVars[name]
// 		if !exists {
// 			fmt.Printf("  ➖ Variable removed: %s\n", name)
// 			continue
// 		}

// 		if oldVar.Default != newVar.Default {
// 			fmt.Printf("  🔄 Variable changed: %s, Default: %s → %s\n", name, oldVar.Default, newVar.Default)
// 		}

// 		if oldVar.Type != newVar.Type {
// 			fmt.Printf("  🔄 Variable changed: %s, Type: %s → %s\n", name, oldVar.Type, newVar.Type)
// 		}
// 	}

// 	for name := range newVars {
// 		if _, exists := oldVars[name]; !exists {
// 			fmt.Printf("  ➕ Variable added: %s (default: %s)\n", name, newVars[name].Default)
// 		}
// 	}
// }

// func diffResources(oldResources, newResources map[string]TFResource) {
// 	for key, oldRes := range oldResources {
// 		if oldRes.Type == "helm_release" {
// 			continue // Skip helm_release resources, handled in diffHelmReleases
// 		}
// 		newRes, exists := newResources[key]
// 		if !exists {
// 			fmt.Printf("  ➖ Resource removed: %s\n", key)
// 			continue
// 		}

// 		if oldRes.Type != newRes.Type {
// 			fmt.Printf("  🔄 Resource type changed: %s, Type: %s → %s\n", key, oldRes.Type, newRes.Type)
// 		}

// 		if oldRes.Name != newRes.Name {
// 			fmt.Printf("  🔄 Resource name changed: %s, Name: %s → %s\n", key, oldRes.Name, newRes.Name)
// 		}

// 		for attrName, oldValue := range oldRes.Config {
// 			newValue, exists := newRes.Config[attrName]
// 			if !exists {
// 				fmt.Printf("  ➖ Attribute removed from resource %s: %s\n", key, attrName)
// 				continue
// 			}
// 			if oldValue != newValue {
// 				fmt.Printf("  🔄 Attribute changed in resource %s: %s, Value: %s → %s\n", key, attrName, oldValue, newValue)
// 			}
// 		}

// 		for attrName := range newRes.Config {
// 			if _, exists := oldRes.Config[attrName]; !exists {
// 				fmt.Printf("  ➕ Attribute added to resource %s: %s (value: %s)\n", key, attrName, newRes.Config[attrName])
// 			}
// 		}
// 	}

// 	for key, newRes := range newResources {
// 		if newRes.Type == "helm_release" {
// 			continue // Skip helm_release resources, handled in diffHelmReleases
// 		}
// 		if _, exists := oldResources[key]; !exists {
// 			fmt.Printf("  ➕ Resource added: %s\n", key)
// 		}
// 	}
// }

// func helmReleasesEqual(oldRelease, newRelease HelmRelease) bool {
// 	if oldRelease.ResourceName != newRelease.ResourceName ||
// 		oldRelease.Name != newRelease.Name ||
// 		oldRelease.Chart != newRelease.Chart ||
// 		oldRelease.Version != newRelease.Version ||
// 		oldRelease.Repository != newRelease.Repository ||
// 		oldRelease.Namespace != newRelease.Namespace {
// 		return false
// 	}

// 	if len(oldRelease.Values) != len(newRelease.Values) {
// 		return false
// 	}
// 	for i, value := range oldRelease.Values {
// 		if value != newRelease.Values[i] {
// 			return false
// 		}
// 	}

// 	if len(oldRelease.Sets) != len(newRelease.Sets) {
// 		return false
// 	}
// 	for key, value := range oldRelease.Sets {
// 		if newValue, exists := newRelease.Sets[key]; !exists || value != newValue {
// 			return false
// 		}
// 	}

// 	return true
// }

// func diffHelmReleases(oldReleases, newReleases []HelmRelease) {
// 	oldMap := make(map[string]HelmRelease)
// 	newMap := make(map[string]HelmRelease)

// 	for _, h := range oldReleases {
// 		oldMap[h.ResourceName] = h
// 	}
// 	for _, h := range newReleases {
// 		newMap[h.ResourceName] = h
// 	}

// 	for name, oldRelease := range oldMap {
// 		newRelease, exists := newMap[name]
// 		if !exists {
// 			fmt.Printf("  ➖ Removed: %s\n", oldRelease.Name)
// 			continue
// 		}
// 		if !helmReleasesEqual(oldRelease, newRelease) {
// 			fmt.Printf("  🔄 Changed: %s\n", newRelease.Name)
// 			if oldRelease.Chart != newRelease.Chart {
// 				fmt.Printf("    Chart: %s → %s\n", oldRelease.Chart, newRelease.Chart)
// 			}
// 			if oldRelease.Version != newRelease.Version {
// 				fmt.Printf("    Version: %s → %s\n", oldRelease.Version, newRelease.Version)
// 			}
// 			if oldRelease.Repository != newRelease.Repository {
// 				fmt.Printf("    Repository: %s → %s\n", oldRelease.Repository, newRelease.Repository)
// 			}
// 			if oldRelease.Namespace != newRelease.Namespace {
// 				fmt.Printf("    Namespace: %s → %s\n", oldRelease.Namespace, newRelease.Namespace)
// 			}
// 			// Print changes in Sets
// 			for key, oldVal := range oldRelease.Sets {
// 				newVal, exists := newRelease.Sets[key]
// 				if !exists {
// 					fmt.Printf("    ➖ Set removed: %s (was: %s)\n", key, oldVal)
// 				} else if oldVal != newVal {
// 					fmt.Printf("    🔄 Set changed: %s, Value: %s → %s\n", key, oldVal, newVal)
// 				}
// 			}
// 			for key, newVal := range newRelease.Sets {
// 				if _, exists := oldRelease.Sets[key]; !exists {
// 					fmt.Printf("    ➕ Set added: %s (value: %s)\n", key, newVal)
// 				}
// 			}
// 			// Print changes in Values
// 			maxLen := len(oldRelease.Values)
// 			if len(newRelease.Values) > maxLen {
// 				maxLen = len(newRelease.Values)
// 			}
// 			for i := 0; i < maxLen; i++ {
// 				var oldVal, newVal string
// 				if i < len(oldRelease.Values) {
// 					oldVal = oldRelease.Values[i]
// 				}
// 				if i < len(newRelease.Values) {
// 					newVal = newRelease.Values[i]
// 				}
// 				if oldVal != newVal {
// 					if oldVal == "" {
// 						fmt.Printf("    ➕ Value added: %s\n", newVal)
// 					} else if newVal == "" {
// 						fmt.Printf("    ➖ Value removed: %s\n", oldVal)
// 					} else {
// 						fmt.Printf("    🔄 Value changed: %s → %s\n", oldVal, newVal)
// 					}
// 				}
// 			}
// 		}
// 	}

// 	for name, newRelease := range newMap {
// 		if _, exists := oldMap[name]; !exists {
// 			fmt.Printf("  ➕ Added: %s (Chart: %s, Version: %s)\n",
// 				newRelease.Name, newRelease.Chart, newRelease.Version)
// 		}
// 	}
// }

// func parseLocalsTF(path string, evalCtx *hcl.EvalContext) (map[string]cty.Value, error) {
// 	src, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read file %s: %v", path, err)
// 	}

// 	parser := hclparse.NewParser()
// 	file, diags := parser.ParseHCLFile(path)
// 	if diags.HasErrors() {
// 		return nil, fmt.Errorf("failed to parse file %s: %v", path, diags)
// 	}

// 	body, ok := file.Body.(*hclsyntax.Body)
// 	if !ok {
// 		return nil, fmt.Errorf("file %s does not contain a valid HCL body", path)
// 	}

// 	locals := make(map[string]cty.Value)
// 	for _, block := range body.Blocks {
// 		if block.Type == "locals" {
// 			for name, attr := range block.Body.Attributes {
// 				if attr.Expr != nil {
// 					val, diag := attr.Expr.Value(evalCtx)
// 					if diag.HasErrors() {
// 						locals[name] = cty.StringVal(string(attr.Expr.Range().SliceBytes(src)))
// 					} else {
// 						locals[name] = val
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return locals, nil
// }

func main() {
	start := time.Now()

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <old_version_path> <new_version_path>")
		// Example: go run main.go ./advanced-test-env/mock-istio-v1.11/ ./advanced-test-env/mock-istio-v1.12/
		os.Exit(1)
	}

	oldPath := filepath.Join(os.Args[1], "variables.tf")
	newPath := filepath.Join(os.Args[2], "variables.tf")

	oldVars, err := parser.ParseVariablesTF(oldPath)
	if err != nil {
		fmt.Printf("Error parsing old variables.tf: %v\n", err)
		return
	}

	newVars, err := parser.ParseVariablesTF(newPath)
	if err != nil {
		fmt.Printf("Error parsing new variables,tf: %v\n", err)
		return
	}

	fmt.Println("Variable Changes:")
	diff.DiffVariables(oldVars, newVars)

	// Step 1: Build a context with only variables
	evalCtxVars := parser.BuildEvalContext(oldVars, map[string]cty.Value{})

	// Step 2: Parse and evaluate locals using the context with variables
	oldLocals, _ := parser.ParseLocalsTF(filepath.Join(os.Args[1], "main.tf"), evalCtxVars)

	// Step 3: Build the final context with both variables and evaluated locals
	evalCtx := parser.BuildEvalContext(oldVars, oldLocals)

	oldResources, oldHelmReleases, err := parser.ParseMainTF(filepath.Join(os.Args[1], "main.tf"), evalCtx)
	if err != nil {
		fmt.Printf("Error parsing old main.tf: %v\n", err)
		return
	}

	evalCtxVars = parser.BuildEvalContext(newVars, map[string]cty.Value{})
	newLocals, _ := parser.ParseLocalsTF(filepath.Join(os.Args[2], "main.tf"), evalCtxVars)
	evalCtx = parser.BuildEvalContext(newVars, newLocals)
	newResources, newHelmReleases, err := parser.ParseMainTF(filepath.Join(os.Args[2], "main.tf"), evalCtx)
	if err != nil {
		fmt.Printf("Error parsing new main.tf: %v\n", err)
		return
	}

	fmt.Println("Resource Changes:")
	diff.DiffResources(oldResources, newResources)

	fmt.Println("Helm Release Changes:")
	diff.DiffHelmReleases(oldHelmReleases, newHelmReleases)

	elapsed := time.Since(start)
	fmt.Printf("Parsing completed in %s\n", elapsed)
}
