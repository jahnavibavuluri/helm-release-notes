package diff

import (
	"fmt"

	"tf-diff/parser"
)

func DiffResources(oldResources, newResources map[string]parser.TFResource) {
	for key, oldRes := range oldResources {
		if oldRes.Type == "helm_release" {
			continue // Skip helm_release resources, handled in diffHelmReleases
		}
		newRes, exists := newResources[key]
		if !exists {
			fmt.Printf("  ➖ Resource removed: %s\n", key)
			continue
		}

		if oldRes.Type != newRes.Type {
			fmt.Printf("  🔄 Resource type changed: %s, Type: %s → %s\n", key, oldRes.Type, newRes.Type)
		}

		if oldRes.Name != newRes.Name {
			fmt.Printf("  🔄 Resource name changed: %s, Name: %s → %s\n", key, oldRes.Name, newRes.Name)
		}

		for attrName, oldValue := range oldRes.Config {
			newValue, exists := newRes.Config[attrName]
			if !exists {
				fmt.Printf("  ➖ Attribute removed from resource %s: %s\n", key, attrName)
				continue
			}
			if oldValue != newValue {
				fmt.Printf("  🔄 Attribute changed in resource %s: %s, Value: %s → %s\n", key, attrName, oldValue, newValue)
			}
		}

		for attrName := range newRes.Config {
			if _, exists := oldRes.Config[attrName]; !exists {
				fmt.Printf("  ➕ Attribute added to resource %s: %s (value: %s)\n", key, attrName, newRes.Config[attrName])
			}
		}
	}

	for key, newRes := range newResources {
		if newRes.Type == "helm_release" {
			continue // Skip helm_release resources, handled in diffHelmReleases
		}
		if _, exists := oldResources[key]; !exists {
			fmt.Printf("  ➕ Resource added: %s\n", key)
		}
	}
}
