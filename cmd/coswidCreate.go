package cmd

import (
    "fmt"

    "github.com/spf13/afero"
    "github.com/spf13/cobra"
    "github.com/veraison/swid"
)

var (
    coswidCreateTemplate  string
    coswidCreateOutputDir string
)

var coswidCreateCmd = &cobra.Command{
    Use:   "create",
    Short: "Create a CBOR-encoded CoSWID from the supplied JSON template",
    Long: `Create a CBOR-encoded CoSWID from the supplied JSON template.

Create a CoSWID from template t1.json and save it to the current directory.

  cocli coswid create --template=t1.json

Create a CoSWID from template t1.json and save it to the specified directory.
`,
    RunE: func(cmd *cobra.Command, args []string) error {
        if err := checkCoswidCreateArgs(); err != nil {
            return err
        }

        cborFile, err := coswidTemplateToCBOR(coswidCreateTemplate, coswidCreateOutputDir)
        if err != nil {
            return err
        }
        fmt.Printf(">> created %q from %q\n", cborFile, coswidCreateTemplate)

        return nil
    },
}

func checkCoswidCreateArgs() error {
    if coswidCreateTemplate == "" {
        return fmt.Errorf("no CoSWID template supplied")
    }
    return nil
}

func coswidTemplateToCBOR(tmplFile, outputDir string) (string, error) {
    var (
        tmplData, coswidCBOR []byte
        s                    swid.SoftwareIdentity
        coswidFile           string
        err                  error
    )

    // Read the template file
    if tmplData, err = afero.ReadFile(fs, tmplFile); err != nil {
        return "", fmt.Errorf("error loading template from %s: %w", tmplFile, err)
    }

    // Parse the JSON into the struct
    if err = s.FromJSON(tmplData); err != nil {
        return "", fmt.Errorf("error decoding template from %s: %w", tmplFile, err)
    }

    // Debugging: Print the parsed SoftwareIdentity object
    fmt.Println("Decoded SoftwareIdentity object:")
    fmt.Printf("%+v\n", s)

    // Encode the struct to CBOR
    if coswidCBOR, err = s.ToCBOR(); err != nil {
        fmt.Printf("SoftwareIdentity object before CBOR encoding: %+v\n", s)
        return "", fmt.Errorf("error encoding CoSWID to CBOR: %w", err)
    }

    // Generate the output file name
    coswidFile = makeFileName(outputDir, tmplFile, ".cbor")

    // Write the CBOR data to the output file
    if err = afero.WriteFile(fs, coswidFile, coswidCBOR, 0644); err != nil {
        return "", fmt.Errorf("error saving CoSWID to file %s: %w", coswidFile, err)
    }

    return coswidFile, nil
}

func init() {
    coswidCmd.AddCommand(coswidCreateCmd)
    coswidCreateCmd.Flags().StringVarP(&coswidCreateTemplate, "template", "t", "", "a CoSWID template file (in JSON format)")
    coswidCreateCmd.Flags().StringVarP(&coswidCreateOutputDir, "output-dir", "o", ".", "directory where the created files are stored")
}