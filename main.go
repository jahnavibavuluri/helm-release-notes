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
