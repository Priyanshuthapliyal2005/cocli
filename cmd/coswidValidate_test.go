package cmd

import (
    "testing"
    "github.com/spf13/afero"
    "github.com/spf13/cobra"
    "github.com/stretchr/testify/assert"
)

func TestCheckCoswidValidateArgs(t *testing.T) {
    // Setup
    coswidValidateFile = ""
    coswidValidateSchema = ""
    
    // Test cases
    err := checkCoswidValidateArgs()
    assert.Error(t, err, "no CoSWID file supplied")

    coswidValidateFile = "test.cbor"
    err = checkCoswidValidateArgs()
    assert.Error(t, err, "no schema supplied")

    coswidValidateSchema = "schema.json"
    err = checkCoswidValidateArgs()
    assert.NoError(t, err)
}

func TestValidateCoswid(t *testing.T) {
    // Setup
    fs = afero.NewMemMapFs()
    afero.WriteFile(fs, "test.cbor", []byte{0xA1, 0x01, 0x02}, 0644)
    afero.WriteFile(fs, "schema.json", []byte(`{"type": "object"}`), 0644)

    // Test case
    err := validateCoswid("test.cbor", "schema.json")
    assert.NoError(t, err)
}

func TestCoswidValidateCmd(t *testing.T) {
    // Setup
    cmd := &cobra.Command{}
    
    // Test case
    coswidValidateCmd.RunE(cmd, []string{})
    assert.NoError(t, cmd.Execute())
}
