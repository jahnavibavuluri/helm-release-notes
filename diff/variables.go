package diff

import (
	"fmt"
	"tf-diff/parser"
)

func DiffVariables(oldVars, newVars map[string]parser.TFVariable) {
	for name, oldVar := range oldVars {
		newVar, exists := newVars[name]
		if !exists {
			fmt.Printf("  ➖ Variable removed: %s\n", name)
			continue
		}

		if oldVar.Default != newVar.Default {
			fmt.Printf("  🔄 Variable changed: %s, Default: %s → %s\n", name, oldVar.Default, newVar.Default)
		}

		if oldVar.Type != newVar.Type {
			fmt.Printf("  🔄 Variable changed: %s, Type: %s → %s\n", name, oldVar.Type, newVar.Type)
		}
	}

	for name := range newVars {
		if _, exists := oldVars[name]; !exists {
			fmt.Printf("  ➕ Variable added: %s (default: %s)\n", name, newVars[name].Default)
		}
	}
}
