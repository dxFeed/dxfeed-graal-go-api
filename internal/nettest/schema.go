// Copyright (c) 2026 Devexperts LLC.
// SPDX-License-Identifier: MPL-2.0

package nettest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

// ValidateJSON checks a JSON document against a JSON Schema with the keywords that nettest-summary.schema.json of the
// kit uses: type, enum, const, required, properties, additionalProperties false, items, minimum, maximum and
// minLength. The tests check the summary files with it.
func ValidateJSON(schema, document []byte) error {
	var s, d any
	if err := json.Unmarshal(schema, &s); err != nil {
		return fmt.Errorf("the schema: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(document))
	decoder.UseNumber()
	if err := decoder.Decode(&d); err != nil {
		return fmt.Errorf("the document: %w", err)
	}
	return validate(s.(map[string]any), d, "$")
}

func jsonType(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case bool:
		return "boolean"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	case json.Number:
		if f, err := v.Float64(); err == nil && f == math.Trunc(f) {
			return "integer"
		}
		return "number"
	}
	return "unknown"
}

// same compares a value of the document with one of the schema.
func same(document, schema any) bool {
	if n, ok := document.(json.Number); ok {
		f, err := n.Float64()
		s, isNumber := schema.(float64)
		return err == nil && isNumber && f == s
	}
	return fmt.Sprint(document) == fmt.Sprint(schema) && jsonType(document) != "number"
}

func validate(schema map[string]any, value any, path string) error {
	if types, ok := schema["type"]; ok {
		list, single := types.([]any)
		if !single {
			list = []any{types}
		}
		found := false
		for _, t := range list {
			actual := jsonType(value)
			found = found || t == actual || t == "number" && actual == "integer"
		}
		if !found {
			return fmt.Errorf("%s: %s is not of the type %v", path, jsonType(value), types)
		}
	}
	if c, ok := schema["const"]; ok && !same(value, c) {
		return fmt.Errorf("%s: %v is not %v", path, value, c)
	}
	if enum, ok := schema["enum"].([]any); ok {
		found := false
		for _, e := range enum {
			found = found || (e == nil && value == nil) || (e != nil && value != nil && same(value, e))
		}
		if !found {
			return fmt.Errorf("%s: %v is not one of %v", path, value, enum)
		}
	}
	switch v := value.(type) {
	case json.Number:
		f, _ := v.Float64()
		if minimum, ok := schema["minimum"].(float64); ok && f < minimum {
			return fmt.Errorf("%s: %v is below %v", path, v, minimum)
		}
		if maximum, ok := schema["maximum"].(float64); ok && f > maximum {
			return fmt.Errorf("%s: %v is above %v", path, v, maximum)
		}
	case string:
		if minLength, ok := schema["minLength"].(float64); ok && float64(len(v)) < minLength {
			return fmt.Errorf("%s: %q is shorter than %v", path, v, minLength)
		}
	case []any:
		if items, ok := schema["items"].(map[string]any); ok {
			for i, item := range v {
				if err := validate(items, item, fmt.Sprintf("%s[%d]", path, i)); err != nil {
					return err
				}
			}
		}
	case map[string]any:
		properties, _ := schema["properties"].(map[string]any)
		if required, ok := schema["required"].([]any); ok {
			for _, key := range required {
				if _, ok := v[key.(string)]; !ok {
					return fmt.Errorf("%s: %s is missing", path, key)
				}
			}
		}
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			property, ok := properties[key].(map[string]any)
			if !ok {
				if schema["additionalProperties"] == false {
					return fmt.Errorf("%s: %s is not allowed", path, key)
				}
				continue
			}
			if err := validate(property, v[key], path+"."+key); err != nil {
				return err
			}
		}
	}
	return nil
}
