package cmd

import (
    "testing"
    "github.com/spf13/afero"
    "github.com/stretchr/testify/assert"
)

func TestCheckCoswidDisplayArgs(t *testing.T) {
    // Setup
    coswidDisplayFile = ""
    coswidDisplayDir = ""
    
    // Test cases
    err := checkCoswidDisplayArgs()
    assert.Error(t, err, "no CoSWID file or directory supplied")

    coswidDisplayFile = "test.cbor"
    err = checkCoswidDisplayArgs()
    assert.NoError(t, err)
}

func TestGatherFiles(t *testing.T) {
    // Setup
    fs = afero.NewMemMapFs()
    afero.WriteFile(fs, "test.cbor", []byte{}, 0644)
    afero.WriteFile(fs, "dir/test.cbor", []byte{}, 0644)

    // Test case
    files := gatherFiles([]string{"test.cbor"}, []string{"dir"}, ".cbor")
    assert.Len(t, files, 2)
}

func TestDisplayCoswid(t *testing.T) {
    // Setup
    fs = afero.NewMemMapFs()
    afero.WriteFile(fs, "test.cbor", []byte{0xA1, 0x01, 0x02}, 0644)

    // Test case
    err := displayCoswid("test.cbor")
    assert.NoError(t, err)
}
