package cmd

import (
    "testing"
    "github.com/spf13/afero"
    "github.com/stretchr/testify/assert"
)

func TestCheckCoswidCreateArgs(t *testing.T) {
    // Setup
    coswidCreateTemplate = ""
    
    // Test cases
    err := checkCoswidCreateArgs()
    assert.Error(t, err, "template file is required")

    coswidCreateTemplate = "template.json"
    err = checkCoswidCreateArgs()
    assert.NoError(t, err)
}

func TestCoswidTemplateToCBOR(t *testing.T) {
    // Setup
    template := `{"evidence": {"type": "test", "value": "test"}}`
    afero.WriteFile(afero.NewOsFs(), "template.json", []byte(template), 0644)

    // Test case
    output, err := coswidTemplateToCBOR("template.json", ".")
    assert.NoError(t, err)
    assert.NotEmpty(t, output)
}

func TestValidateJSON(t *testing.T) {
    // Setup
    template := `{"evidence": {"type": "test", "value": "test"}}`
    schema := `{"type": "object"}`
    afero.WriteFile(afero.NewOsFs(), "template.json", []byte(template), 0644)
    afero.WriteFile(afero.NewOsFs(), "schema.json", []byte(schema), 0644)

    // Test case
    err := validateJSON("template.json", "schema.json")
    assert.NoError(t, err)
}
