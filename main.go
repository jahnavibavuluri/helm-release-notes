package main

import (
	"fmt"
	"os"
	"time"

	"tf-diff/diff"
	"tf-diff/fetch"
	"tf-diff/parser"

	"github.com/zclconf/go-cty/cty"
)

func main() {
	start := time.Now()

	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go <owner> <repo> <old_ref> <new_ref>")
		fmt.Println("Example: go run main.go hashicorp terraform-aws-modules main feature-branch")
		fmt.Println("Example: go run main.go jahnavibavuluri test-terraform-repo v1.0.0 v1.1.0")
		os.Exit(1)
	}

	owner := os.Args[1]
	repo := os.Args[2]
	oldRef := os.Args[3]
	newRef := os.Args[4]

	fmt.Printf("Comparing %s/%s: %s → %s\n\n", owner, repo, oldRef, newRef)

	// Download variables.tf files
	oldVariablesURL := fetch.BuildGitHubRawURL(owner, repo, oldRef, "variables.tf")
	newVariablesURL := fetch.BuildGitHubRawURL(owner, repo, newRef, "variables.tf")

	fmt.Printf("Downloading %s...\n", oldVariablesURL)
	oldVariablesContent, err := fetch.DownloadGitHubFile(oldVariablesURL)
	if err != nil {
		fmt.Printf("Error downloading old variables.tf: %v\n", err)
		return
	}

	fmt.Printf("Downloading %s...\n", newVariablesURL)
	newVariablesContent, err := fetch.DownloadGitHubFile(newVariablesURL)
	if err != nil {
		fmt.Printf("Error downloading new variables.tf: %v\n", err)
		return
	}

	// Parse variables
	oldVars, err := parser.ParseVariablesTF(oldVariablesContent, "variables.tf")
	if err != nil {
		fmt.Printf("Error parsing old variables.tf: %v\n", err)
		return
	}

	newVars, err := parser.ParseVariablesTF(newVariablesContent, "variables.tf")
	if err != nil {
		fmt.Printf("Error parsing new variables.tf: %v\n", err)
		return
	}

	//Download main.tf files
	oldMainURL := fetch.BuildGitHubRawURL(owner, repo, oldRef, "main.tf")
	newMainURL := fetch.BuildGitHubRawURL(owner, repo, newRef, "main.tf")

	fmt.Printf("\nDownloading %s...\n", oldMainURL)
	oldMainContent, err := fetch.DownloadGitHubFile(oldMainURL)
	if err != nil {
		fmt.Printf("Error downloading old main.tf: %v\n", err)
		return
	}

	fmt.Printf("Downloading %s...\n", newMainURL)
	newMainContent, err := fetch.DownloadGitHubFile(newMainURL)
	if err != nil {
		fmt.Printf("Error downloading new main.tf: %v\n", err)
		return
	}

	// Parse main.tf files with evaluation context
	// Step 1: Build a context with only variables
	evalCtxVars := parser.BuildEvalContext(oldVars, map[string]cty.Value{})

	// Step 2: Parse and evaluate locals using the context with variables
	oldLocals, _ := parser.ParseLocalsTF(oldMainContent, "main.tf", evalCtxVars)

	// Step 3: Build the final context with both variables and evaluated locals
	evalCtx := parser.BuildEvalContext(oldVars, oldLocals)

	oldResources, oldHelmReleases, err := parser.ParseMainTF(oldMainContent, "main.tf", evalCtx)
	if err != nil {
		fmt.Printf("Error parsing old main.tf: %v\n", err)
		return
	}

	evalCtxVars = parser.BuildEvalContext(newVars, map[string]cty.Value{})
	newLocals, _ := parser.ParseLocalsTF(newMainContent, "main.tf", evalCtxVars)
	evalCtx = parser.BuildEvalContext(newVars, newLocals)
	newResources, newHelmReleases, err := parser.ParseMainTF(newMainContent, "main.tf", evalCtx)
	if err != nil {
		fmt.Printf("Error parsing new main.tf: %v\n", err)
		return
	}

	fmt.Println("\nVariable Changes:")
	diff.DiffVariables(oldVars, newVars)

	fmt.Println("\nResource Changes:")
	diff.DiffResources(oldResources, newResources)

	fmt.Println("\nHelm Release Changes:")
	diff.DiffHelmReleases(oldHelmReleases, newHelmReleases)

	elapsed := time.Since(start)
	fmt.Printf("\nParsing completed in %s\n", elapsed)
}
