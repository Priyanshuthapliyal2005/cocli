package cmd

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/xeipuuv/gojsonschema"
)

var (
	coswidValidateFile   string
	coswidValidateSchema string
)

var coswidValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate a CBOR-encoded CoSWID against the provided JSON schema",
	Long: `validate a CBOR-encoded CoSWID against the provided JSON schema

Validate the CoSWID in file s.cbor against the schema schema.json.

  cocli coswid validate --file=s.cbor --schema=schema.json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := checkCoswidValidateArgs(); err != nil {
			return err
		}

		if err := validateCoswid(coswidValidateFile, coswidValidateSchema); err != nil {
			return err
		}

		fmt.Printf(">> validated %q against %q\n", coswidValidateFile, coswidValidateSchema)
		return nil
	},
}

func checkCoswidValidateArgs() error {
	if coswidValidateFile == "" {
		return fmt.Errorf("no CoSWID file supplied")
	}
	if coswidValidateSchema == "" {
		return fmt.Errorf("no schema supplied")
	}
	return nil
}

func validateCoswid(file, schema string) error {
	var (
		coswidCBOR []byte
		err        error
	)

	if coswidCBOR, err = afero.ReadFile(fs, file); err != nil {
		return fmt.Errorf("error loading CoSWID from %s: %w", file, err)
	}

	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schema)
	documentLoader := gojsonschema.NewBytesLoader(coswidCBOR)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("error validating CoSWID from %s: %w", file, err)
	}

	if !result.Valid() {
		return fmt.Errorf("CoSWID from %s is invalid: %v", file, result.Errors())
	}

	return nil
}

func init() {
	coswidCmd.AddCommand(coswidValidateCmd)
	coswidValidateCmd.Flags().StringVarP(&coswidValidateFile, "file", "f", "", "a CoSWID file (in CBOR format)")
	coswidValidateCmd.Flags().StringVarP(&coswidValidateSchema, "schema", "s", "", "a JSON schema file")
}
