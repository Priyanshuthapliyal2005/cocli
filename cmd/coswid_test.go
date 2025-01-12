package cmd

import (
    "testing"
    "github.com/spf13/cobra"
    "github.com/stretchr/testify/assert"
)

func TestCoswidCmd(t *testing.T) {
    // Setup
    cmd := &cobra.Command{}
    
    // Test case
    err := coswidCmd.RunE(cmd, []string{})
    assert.NoError(t, err)
    assert.NoError(t, cmd.Execute())
}
