// Package util provides utility functions for FDIO
package util

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
)

// TomlTreeToMap converts a toml tree to an array of map[string]interface{}. It does so
// by introspecting the tree and looking for the items that match a specific key.
func TomlTreeToMap(tree *toml.Tree, key string) ([]map[string]interface{}, error) {
	// Get the correct key
	queryResult := tree.Get(key)
	if queryResult == nil {
		return nil, fmt.Errorf("No items found in the tree")
	}

	// Prepare the result
	resultArray := queryResult.([]*toml.Tree)
	datamap := make([]map[string]interface{}, len(resultArray))
	for idx, val := range resultArray {
		datamap[idx] = val.ToMap()
	}
	return datamap, nil
}
