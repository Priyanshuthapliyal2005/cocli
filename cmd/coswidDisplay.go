package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "encoding/json"

    "github.com/spf13/afero"
    "github.com/spf13/cobra"
    "github.com/veraison/swid"
)

var (
    coswidDisplayFile string
    coswidDisplayDir  string
)

var coswidDisplayCmd = &cobra.Command{
    Use:   "display",
    Short: "Display one or more CBOR-encoded CoSWID(s) in human-readable (JSON) format",
    Long: `Display one or more CBOR-encoded CoSWID(s) in human-readable (JSON) format.
You can supply individual CoSWID files or directories containing CoSWID files.

Display CoSWID in file s.cbor.

  cocli coswid display --file=s.cbor

Display CoSWIDs in files s1.cbor, s2.cbor and any cbor file in the coswids/ directory.

  cocli coswid display --file=s1.cbor --file=s2.cbor --dir=coswids
`,
    RunE: func(cmd *cobra.Command, args []string) error {
        // Validate input arguments
        if err := checkCoswidDisplayArgs(); err != nil {
            return err
        }

        filesList := gatherFiles([]string{coswidDisplayFile}, []string{coswidDisplayDir}, ".cbor")
        if len(filesList) == 0 {
            return fmt.Errorf("no CoSWID files found")
        }

        for _, file := range filesList {
            if err := displayCoswid(file); err != nil {
                fmt.Printf("Error displaying %s: %v\n", file, err)
            }
        }

        return nil
    },
}

func checkCoswidDisplayArgs() error {
    if coswidDisplayFile == "" && coswidDisplayDir == "" {
        return fmt.Errorf("no CoSWID file or directory supplied")
    }
    return nil
}

func gatherFiles(files []string, dirs []string, ext string) []string {
    collectedMap := make(map[string]struct{})
    var collected []string

    // Collect files from specified file paths
    for _, file := range files {
        if filepath.Ext(file) == ext {
            collectedMap[file] = struct{}{}
        }
    }

    // Collect files from specified directories
    for _, dir := range dirs {
        if dir != "" {
            exists, err := afero.Exists(fs, dir)
            if err == nil && exists {
                afero.Walk(fs, dir, func(path string, info os.FileInfo, err error) error {
                    if err != nil {
                        fmt.Printf("Error accessing path %s: %v\n", path, err)
                        return nil
                    }
                    if !info.IsDir() && filepath.Ext(path) == ext {
                        collectedMap[path] = struct{}{}
                    }
                    return nil
                })
            }
        }
    }

    // Convert map keys to slice
    for file := range collectedMap {
        collected = append(collected, file)
    }

    return collected
}

func displayCoswid(file string) error {
    fmt.Printf("Processing file: %s\n", file)
    var (
        coswidCBOR []byte
        s          swid.SoftwareIdentity
        err        error
    )

    // Read the CBOR file
    if coswidCBOR, err = afero.ReadFile(fs, file); err != nil {
        return fmt.Errorf("error reading file %s: %w", file, err)
    }

    // Decode CBOR to SoftwareIdentity
    if err = s.FromCBOR(coswidCBOR); err != nil {
        return fmt.Errorf("error decoding CoSWID from %s: %w", file, err)
    }

    // Convert to JSON
    coswidJSON, err := json.MarshalIndent(&s, "", "  ")
    if err != nil {
        return fmt.Errorf("error marshaling CoSWID to JSON: %w", err)
    }

    fmt.Printf(">> [%s]\n%s\n", file, string(coswidJSON))
    return nil
}

func init() {
    coswidCmd.AddCommand(coswidDisplayCmd)
    coswidDisplayCmd.Flags().StringVarP(&coswidDisplayFile, "file", "f", "", "a CoSWID file (in CBOR format)")
    coswidDisplayCmd.Flags().StringVarP(&coswidDisplayDir, "dir", "d", "", "a directory containing CoSWID files")
}