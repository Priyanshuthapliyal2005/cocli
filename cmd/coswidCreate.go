package cmd

import (
    "encoding/json"
    "fmt"
    "path/filepath"

    "github.com/xeipuuv/gojsonschema"
    "github.com/fxamacker/cbor/v2"
    "github.com/spf13/afero"
    "github.com/spf13/cobra"
    "github.com/veraison/swid"
)

type CustomSoftwareIdentity struct {
    swid.SoftwareIdentity
    Evidence struct {
        Type  string `json:"type"`
        Value string `json:"value"`
    } `json:"evidence"`
}

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

        // Validate JSON against schema before processing
        schemaPath := "D:/opensource/cocli/data/coswid/templates/coswid-schema.json"
        err := validateJSON(coswidCreateTemplate, schemaPath)
        if err != nil {
            return fmt.Errorf("JSON validation failed: %v", err)
        }

        cborFile, err := coswidTemplateToCBOR(coswidCreateTemplate, coswidCreateOutputDir)
        if err != nil {
            return fmt.Errorf("error creating CBOR: %v", err)
        }
        fmt.Printf(">> created %q from %q\n", cborFile, coswidCreateTemplate)

        return nil
    },
}

func checkCoswidCreateArgs() error {
    if coswidCreateTemplate == "" {
        return fmt.Errorf("template file is required")
    }
    return nil
}

func coswidTemplateToCBOR(tmplFile, outputDir string) (string, error) {
    var (
        tmplData   []byte
        coswidCBOR []byte
        s          CustomSoftwareIdentity
        coswidFile string
        err        error
    )

    // Read the template file
    tmplData, err = afero.ReadFile(afero.NewOsFs(), tmplFile)
    if err != nil {
        return "", fmt.Errorf("unable to read template file: %v", err)
    }

    // Parse the JSON into the custom struct
    err = json.Unmarshal(tmplData, &s)
    if err != nil {
        return "", fmt.Errorf("error decoding template from %s: %v", tmplFile, err)
    }

    // Debugging: Print the parsed CustomSoftwareIdentity object
    fmt.Println("Decoded CustomSoftwareIdentity object:")
    fmt.Printf("%+v\n", s)

    // Encode the struct to CBOR using fxamacker/cbor
    coswidCBOR, err = cbor.Marshal(s)
    if err != nil {
        return "", fmt.Errorf("error encoding to CBOR: %v", err)
    }

    // Generate the output file name
    coswidFile = makeFileName(outputDir, tmplFile, ".cbor")

    // Write the CBOR data to the output file
    err = afero.WriteFile(afero.NewOsFs(), coswidFile, coswidCBOR, 0644)
    if err != nil {
        return "", fmt.Errorf("error writing CBOR file: %v", err)
    }

    return coswidFile, nil
}


// validateJSON validates the JSON template against the provided schema
func validateJSON(tmplFile, schemaFile string) error {
    schemaLoader := gojsonschema.NewReferenceLoader("file://" + filepath.ToSlash(schemaFile))
    documentLoader := gojsonschema.NewReferenceLoader("file://" + filepath.ToSlash(tmplFile))

    result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        return fmt.Errorf("error during JSON validation: %v", err)
    }

    if !result.Valid() {
        for _, desc := range result.Errors() {
            fmt.Printf("- %s\n", desc)
        }
        return fmt.Errorf("schema validation failed: JSON does not conform to schema")
    }

    fmt.Println("JSON validation successful.")
    return nil
}

func init() {
    coswidCmd.AddCommand(coswidCreateCmd)
    coswidCreateCmd.Flags().StringVarP(&coswidCreateTemplate, "template", "t", "", "a CoSWID template file (in JSON format)")
    coswidCreateCmd.Flags().StringVarP(&coswidCreateOutputDir, "output-dir", "o", ".", "output directory for CBOR file")
    coswidCreateCmd.MarkFlagRequired("template")
    coswidCreateCmd.MarkFlagRequired("output-dir")
}