package diff

import (
	"fmt"
	"tf-diff/parser"
)

func helmReleasesEqual(oldRelease, newRelease parser.HelmRelease) bool {
	if oldRelease.ResourceName != newRelease.ResourceName ||
		oldRelease.Name != newRelease.Name ||
		oldRelease.Chart != newRelease.Chart ||
		oldRelease.Version != newRelease.Version ||
		oldRelease.Repository != newRelease.Repository ||
		oldRelease.Namespace != newRelease.Namespace {
		return false
	}

	if len(oldRelease.Values) != len(newRelease.Values) {
		return false
	}
	for i, value := range oldRelease.Values {
		if value != newRelease.Values[i] {
			return false
		}
	}

	if len(oldRelease.Sets) != len(newRelease.Sets) {
		return false
	}
	for key, value := range oldRelease.Sets {
		if newValue, exists := newRelease.Sets[key]; !exists || value != newValue {
			return false
		}
	}

	return true
}

func DiffHelmReleases(oldReleases, newReleases []parser.HelmRelease) {
	oldMap := make(map[string]parser.HelmRelease)
	newMap := make(map[string]parser.HelmRelease)

	for _, h := range oldReleases {
		oldMap[h.ResourceName] = h
	}
	for _, h := range newReleases {
		newMap[h.ResourceName] = h
	}

	for name, oldRelease := range oldMap {
		newRelease, exists := newMap[name]
		if !exists {
			fmt.Printf("  ➖ Removed: %s\n", oldRelease.Name)
			continue
		}
		if !helmReleasesEqual(oldRelease, newRelease) {
			fmt.Printf("  🔄 Changed: %s\n", newRelease.Name)
			if oldRelease.Chart != newRelease.Chart {
				fmt.Printf("    Chart: %s → %s\n", oldRelease.Chart, newRelease.Chart)
			}
			if oldRelease.Version != newRelease.Version {
				fmt.Printf("    Version: %s → %s\n", oldRelease.Version, newRelease.Version)
			}
			if oldRelease.Repository != newRelease.Repository {
				fmt.Printf("    Repository: %s → %s\n", oldRelease.Repository, newRelease.Repository)
			}
			if oldRelease.Namespace != newRelease.Namespace {
				fmt.Printf("    Namespace: %s → %s\n", oldRelease.Namespace, newRelease.Namespace)
			}
			// Print changes in Sets
			for key, oldVal := range oldRelease.Sets {
				newVal, exists := newRelease.Sets[key]
				if !exists {
					fmt.Printf("    ➖ Set removed: %s (was: %s)\n", key, oldVal)
				} else if oldVal != newVal {
					fmt.Printf("    🔄 Set changed: %s, Value: %s → %s\n", key, oldVal, newVal)
				}
			}
			for key, newVal := range newRelease.Sets {
				if _, exists := oldRelease.Sets[key]; !exists {
					fmt.Printf("    ➕ Set added: %s (value: %s)\n", key, newVal)
				}
			}
			// Print changes in Values
			maxLen := len(oldRelease.Values)
			if len(newRelease.Values) > maxLen {
				maxLen = len(newRelease.Values)
			}
			for i := 0; i < maxLen; i++ {
				var oldVal, newVal string
				if i < len(oldRelease.Values) {
					oldVal = oldRelease.Values[i]
				}
				if i < len(newRelease.Values) {
					newVal = newRelease.Values[i]
				}
				if oldVal != newVal {
					if oldVal == "" {
						fmt.Printf("    ➕ Value added: %s\n", newVal)
					} else if newVal == "" {
						fmt.Printf("    ➖ Value removed: %s\n", oldVal)
					} else {
						fmt.Printf("    🔄 Value changed: %s → %s\n", oldVal, newVal)
					}
				}
			}
		}
	}

	for name, newRelease := range newMap {
		if _, exists := oldMap[name]; !exists {
			fmt.Printf("  ➕ Added: %s (Chart: %s, Version: %s)\n",
				newRelease.Name, newRelease.Chart, newRelease.Version)
		}
	}
}
