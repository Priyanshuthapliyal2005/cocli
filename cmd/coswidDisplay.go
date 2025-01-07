package cmd

import (
    "encoding/json"
    "fmt"
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
    Short: "display one or more CBOR-encoded CoSWID(s) in human readable (JSON) format",
    Long: `display one or more CBOR-encoded CoSWID(s) in human readable (JSON) format.
You can supply individual CoSWID files or directories containing CoSWID files.

Display CoSWID in file s.cbor.

  cocli coswid display --file=s.cbor

Display CoSWIDs in files s1.cbor, s2.cbor and any cbor file in the coswids/ directory.

  cocli coswid display --file=s1.cbor --file=s2.cbor --dir=coswids
`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := checkCoswidDisplayArgs(); err != nil {
            return err
        }

        filesList := filesList([]string{coswidDisplayFile}, []string{coswidDisplayDir}, ".cbor")
        if len(filesList) == 0 {
            return fmt.Errorf("no files found")
        }

        for _, file := range filesList {
            if err := displayCoswid(file); err != nil {
                fmt.Printf(">> failed displaying %q: %v\n", file, err)
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

func displayCoswid(file string) error {
    var (
        coswidCBOR []byte
        s          swid.SoftwareIdentity
        err        error
    )

    if coswidCBOR, err = afero.ReadFile(fs, file); err != nil {
        return fmt.Errorf("error loading CoSWID from %s: %w", file, err)
    }

    if err = s.FromCBOR(coswidCBOR); err != nil {
        return fmt.Errorf("error decoding CoSWID from %s: %w", file, err)
    }

    coswidJSON, err := json.MarshalIndent(&s, "", "  ")
    if err != nil {
        return fmt.Errorf("error encoding CoSWID from %s: %w", file, err)
    }

    fmt.Printf(">> [%s]\n%s\n", file, string(coswidJSON))
    return nil
}

func init() {
    coswidCmd.AddCommand(coswidDisplayCmd)
    coswidDisplayCmd.Flags().StringVarP(&coswidDisplayFile, "file", "f", "", "a CoSWID file (in CBOR format)")
    coswidDisplayCmd.Flags().StringVarP(&coswidDisplayDir, "dir", "d", "", "a directory containing CoSWID files")
}