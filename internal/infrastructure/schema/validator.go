package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Validator struct {
	compiler *jsonschema.Compiler
	schemas  map[string]*jsonschema.Schema
}

type fileLoader struct{}

func (fl *fileLoader) Load(s string) (any, error) {
	dir := "schemas"
	file := filepath.Join(dir, filepath.Base(s))
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(data), nil
}

func NewValidator() (*Validator, error) {
	v := &Validator{
		compiler: jsonschema.NewCompiler(),
		schemas:  make(map[string]*jsonschema.Schema),
	}

	v.compiler.UseLoader(&fileLoader{})

	schemaFiles := []string{
		"node.json",
		"device.json",
		"source.json",
		"flow.json",
		"sender.json",
		"receiver.json",
		"subscription.json",
		"resource_core.json",
		"clock_internal.json",
		"clock_ptp.json",
		"flow_video_raw.json",
		"flow_audio_raw.json",
		"source_audio.json",
		"receiver_video.json",
		"receiver_audio.json",
	}

	for _, file := range schemaFiles {
		data, err := os.ReadFile(filepath.Join("schemas", file))
		if err != nil {
			continue
		}
		var doc any
		if err := json.Unmarshal(data, &doc); err != nil {
			continue
		}
		v.compiler.AddResource("schemas/"+file, doc)
	}

	for _, file := range schemaFiles {
		name := strings.TrimSuffix(file, filepath.Ext(file))
		schema, err := v.compiler.Compile("schemas/" + file)
		if err != nil {
			continue
		}
		v.schemas[name] = schema
	}

	return v, nil
}

func (v *Validator) Validate(resourceType string, data interface{}) error {
	schema, ok := v.schemas[resourceType]
	if !ok {
		return fmt.Errorf("schema not found for type: %s", resourceType)
	}

	err := schema.Validate(data)
	if err != nil {
		return fmt.Errorf("validation failed for %s: %w", resourceType, err)
	}

	return nil
}

func (v *Validator) ValidateJSON(resourceType string, jsonData []byte) error {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	return v.Validate(resourceType, data)
}

func (v *Validator) ValidateAndSerialize(resourceType string, obj interface{}) ([]byte, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal object: %w", err)
	}

	if err := v.ValidateJSON(resourceType, data); err != nil {
		return nil, fmt.Errorf("schema validation failed: %w", err)
	}

	return data, nil
}

func ExtractValidationError(err error) string {
	if err == nil {
		return ""
	}
	errStr := err.Error()
	if strings.Contains(errStr, "validation failed") {
		parts := strings.Split(errStr, ":")
		if len(parts) > 2 {
			return strings.TrimSpace(parts[len(parts)-1])
		}
	}
	return errStr
}