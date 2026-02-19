package layout

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// MapItem represents a single key-value pair in an ordered map.
type MapItem struct {
	Key   interface{}
	Value interface{}
}

// MapSlice is an ordered slice of MapItems, preserving YAML map key ordering.
type MapSlice []MapItem

// UnmarshalYAML implements the yaml.Unmarshaler interface for gopkg.in/yaml.v3.
func (ms *MapSlice) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.MappingNode {
		return fmt.Errorf("expected a mapping node, got %v", value.Kind)
	}

	for i := 0; i < len(value.Content)-1; i += 2 {
		var key, val interface{}
		if err := value.Content[i].Decode(&key); err != nil {
			return err
		}
		if err := value.Content[i+1].Decode(&val); err != nil {
			return err
		}
		*ms = append(*ms, MapItem{Key: key, Value: val})
	}
	return nil
}
