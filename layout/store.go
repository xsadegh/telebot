package layout

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// store is a simple hierarchical key-value store that replaces spf13/viper.
// It wraps a map[string]interface{} and provides typed access via dot-separated keys.
type store struct {
	data map[string]interface{}
}

func newStore() *store {
	return &store{data: make(map[string]interface{})}
}

func storeFrom(m map[string]interface{}) *store {
	if m == nil {
		m = make(map[string]interface{})
	}
	return &store{data: m}
}

// get traverses nested maps using a dot-separated key path.
func (s *store) get(key string) interface{} {
	parts := strings.Split(key, ".")
	var cur interface{} = s.data
	for _, p := range parts {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil
		}
		cur, ok = m[p]
		if !ok {
			// Try case-insensitive lookup.
			found := false
			lp := strings.ToLower(p)
			for k, v := range m {
				if strings.ToLower(k) == lp {
					cur = v
					found = true
					break
				}
			}
			if !found {
				return nil
			}
		}
	}
	return cur
}

// sub returns a sub-store for the given key, or nil if not a map.
func (s *store) sub(key string) *store {
	v := s.get(key)
	if v == nil {
		return nil
	}
	m, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return &store{data: m}
}

// allKeys returns all leaf keys in flattened dot-notation.
func (s *store) allKeys() []string {
	var keys []string
	flatten("", s.data, &keys)
	return keys
}

func flatten(prefix string, m map[string]interface{}, keys *[]string) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		if sub, ok := v.(map[string]interface{}); ok {
			flatten(key, sub, keys)
		} else {
			*keys = append(*keys, key)
		}
	}
}

// getString returns the value at key as a string.
func (s *store) getString(key string) string {
	v := s.get(key)
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprint(val)
	}
}

// getInt returns the value at key as an int.
func (s *store) getInt(key string) int {
	v := s.get(key)
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int:
		return val
	case int64:
		return int(val)
	case float64:
		return int(val)
	case string:
		i, _ := strconv.Atoi(val)
		return i
	default:
		return 0
	}
}

// getInt64 returns the value at key as an int64.
func (s *store) getInt64(key string) int64 {
	v := s.get(key)
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case float64:
		return int64(val)
	case string:
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	default:
		return 0
	}
}

// getFloat64 returns the value at key as a float64.
func (s *store) getFloat64(key string) float64 {
	v := s.get(key)
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	default:
		return 0
	}
}

// getBool returns the value at key as a bool.
func (s *store) getBool(key string) bool {
	v := s.get(key)
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case string:
		b, _ := strconv.ParseBool(val)
		return b
	case int:
		return val != 0
	case float64:
		return val != 0
	default:
		return false
	}
}

// getDuration returns the value at key as a time.Duration.
func (s *store) getDuration(key string) time.Duration {
	v := s.get(key)
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case string:
		d, err := time.ParseDuration(val)
		if err != nil {
			// Try as integer nanoseconds
			if i, err := strconv.ParseInt(val, 10, 64); err == nil {
				return time.Duration(i)
			}
		}
		return d
	case int:
		return time.Duration(val)
	case int64:
		return time.Duration(val)
	case float64:
		return time.Duration(val)
	default:
		return 0
	}
}

// getStringSlice returns the value at key as a []string.
func (s *store) getStringSlice(key string) []string {
	v := s.get(key)
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []interface{}:
		ss := make([]string, len(val))
		for i, item := range val {
			ss[i] = fmt.Sprint(item)
		}
		return ss
	case []string:
		return val
	default:
		return nil
	}
}

// getIntSlice returns the value at key as a []int.
func (s *store) getIntSlice(key string) []int {
	v := s.get(key)
	if v == nil {
		return nil
	}
	switch val := v.(type) {
	case []interface{}:
		is := make([]int, len(val))
		for i, item := range val {
			switch n := item.(type) {
			case int:
				is[i] = n
			case int64:
				is[i] = int(n)
			case float64:
				is[i] = int(n)
			case string:
				is[i], _ = strconv.Atoi(n)
			}
		}
		return is
	default:
		return nil
	}
}

// unmarshal decodes the entire store into the provided struct via JSON round-trip.
func (s *store) unmarshal(out interface{}) error {
	b, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}

// unmarshalKey decodes a specific key's value into the provided struct via JSON round-trip.
func (s *store) unmarshalKey(key string, out interface{}) error {
	v := s.get(key)
	if v == nil {
		return fmt.Errorf("key %q not found", key)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
