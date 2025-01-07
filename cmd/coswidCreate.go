package cmd

import (
    "fmt"
    "github.com/spf13/afero"
    "github.com/spf13/cobra"
    "github.com/veraison/swid"
)

var (
    coswidCreateTemplate string
    coswidCreateOutputDir string
)

var coswidCreateCmd = &cobra.Command{
    Use:   "create",
    Short: "create a CBOR-encoded CoSWID from the supplied JSON template",
    Long: `create a CBOR-encoded CoSWID from the supplied JSON template

Create a CoSWID from template t1.json and save it to the current directory.

  cocli coswid create --template=t1.json

Create a CoSWID from template t1.json and save it to the specified directory.

  cocli coswid create --template=t1.json --output-dir=/tmp
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

    if tmplData, err = afero.ReadFile(fs, tmplFile); err != nil {
        return "", fmt.Errorf("error loading template from %s: %w", tmplFile, err)
    }

    if err = s.FromJSON(tmplData); err != nil {
        return "", fmt.Errorf("error decoding template from %s: %w", tmplFile, err)
    }

    if coswidCBOR, err = s.ToCBOR(); err != nil {
        return "", fmt.Errorf("error encoding CoSWID to CBOR: %w", err)
    }

    coswidFile = makeFileName(outputDir, tmplFile, ".cbor")

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