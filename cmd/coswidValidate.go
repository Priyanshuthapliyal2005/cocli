package cmd

import (
    "encoding/json"
    "fmt"
    "encoding/base64"

    "github.com/spf13/afero"
    "github.com/spf13/cobra"
    "github.com/xeipuuv/gojsonschema"
    "github.com/fxamacker/cbor/v2" // Added CBOR package
)

var (
    coswidValidateFile   string
    coswidValidateSchema string
)

var coswidKeyMap = map[uint64]string{
    0: "schema-version", 
    1: "tag-id",
    2: "software-name",
    3: "tag-version",
    4: "patch-level",
    5: "version",
    6: "version-scheme",
    7: "lang",
    8: "directory",
    9: "file",
    10: "process",
    11: "resource",
    12: "size",
    13: "file-version",
    14: "entity",
    15: "evidence",
    16: "link",
    17: "payload",
    18: "hash",
    19: "hash-alg-id",
    20: "hash-value",
}

var coswidValidateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validate a CBOR-encoded CoSWID against the provided JSON schema",
    Long: `Validate a CBOR-encoded CoSWID against the provided JSON schema

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
        coswidJSON []byte
        err        error
    )

    if coswidCBOR, err = afero.ReadFile(fs, file); err != nil {
        return fmt.Errorf("error loading CoSWID from %s: %w", file, err)
    }

    // Decode CBOR with numeric key handling
    var data map[interface{}]interface{}
    if err = cbor.Unmarshal(coswidCBOR, &data); err != nil {
        return fmt.Errorf("error decoding CBOR from %s: %w", file, err)
    }

    // Convert map[interface{}]interface{} to map[string]interface{}
    stringMap := make(map[string]interface{})
    for key, value := range data {
        strKey := convertKeyToString(key)
        convertedValue := convertValue(value)
        stringMap[strKey] = convertedValue
    }

    // Debug: Iterate and print types
    for key, value := range stringMap {
        fmt.Printf("Field: %s, Type: %T, Value: %v\n", key, value, value)

        switch key {
        case "tag-id", "device-id", "location":
            if str, ok := value.(string); !ok {
                return fmt.Errorf("field %s is expected to be string, but got %T", key, value)
            } else {
                _ = str
            }

        case "software-name":
            // Accept either string or map
            switch v := value.(type) {
            case string:
                // OK
            case map[string]interface{}:
                // Handle it as a nested object if needed
                _ = v
            default:
                return fmt.Errorf("field %s has unexpected type %T", key, value)
            }

        case "tag-version", "hash-alg-id":
            switch v := value.(type) {
            case int, int32, int64:
            case float64:
                intValue := int(v)
                stringMap[key] = intValue
            default:
                return fmt.Errorf("field %s is expected to be integer, but got %T", key, value)
            }

        default:
            // Other fields
        }
    }

    // Marshal the decoded data to JSON
    if coswidJSON, err = json.Marshal(stringMap); err != nil {
        return fmt.Errorf("error marshaling CoSWID to JSON: %w", err)
    }

    schemaLoader := gojsonschema.NewReferenceLoader("file:///" + schema)
    documentLoader := gojsonschema.NewBytesLoader(coswidJSON)

    result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        return fmt.Errorf("error validating CoSWID from %s: %w", file, err)
    }

    if !result.Valid() {
        return fmt.Errorf("CoSWID from %s is invalid: %v", file, result.Errors())
    }

    return nil
}

func convertKeyToString(key interface{}) string {
    switch k := key.(type) {
    case string:
        return k
    case int:
        if mappedKey, ok := coswidKeyMap[uint64(k)]; ok {
            return mappedKey
        }
        return fmt.Sprintf("%d", k)
    case uint64:
        if mappedKey, ok := coswidKeyMap[k]; ok {
            return mappedKey
        }
        return fmt.Sprintf("%d", k)
    default:
        return fmt.Sprintf("%v", k)
    }
}

func convertValue(value interface{}) interface{} {
    switch v := value.(type) {
    case map[interface{}]interface{}:
        // Convert nested maps
        m := make(map[string]interface{})
        for k, val := range v {
            strKey := convertKeyToString(k)
            m[strKey] = convertValue(val)
        }
        return m
    case []interface{}:
        // Convert slice elements
        slice := make([]interface{}, len(v))
        for i, val := range v {
            slice[i] = convertValue(val)
        }
        return slice
    case []uint8:
        // Convert byte arrays to base64
        return base64.StdEncoding.EncodeToString(v)
    default:
        return v
    }
}

func init() {
    coswidCmd.AddCommand(coswidValidateCmd)
    coswidValidateCmd.Flags().StringVarP(&coswidValidateFile, "file", "f", "", "a CoSWID file (in CBOR format)")
    coswidValidateCmd.Flags().StringVarP(&coswidValidateSchema, "schema", "s", "", "a JSON schema file")
}
