package api_test

import (
	"fmt"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/check.v1"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
)

func validateSchema(c *check.C, json []byte, schema string) {
	absPath, err := filepath.Abs(schema)
	c.Assert(err, check.IsNil)

	jsonSchema := gojsonschema.NewReferenceLoader("file://" + absPath)
	loader := gojsonschema.NewStringLoader(string(json))
	result, err := gojsonschema.Validate(jsonSchema, loader)
	c.Assert(err, check.IsNil)

	if !result.Valid() {
		jsonValErr := ""
		for _, desc := range result.Errors() {
			jsonValErr = fmt.Sprintf("%s- %s\n", jsonValErr, desc)
		}
		c.Error(fmt.Errorf("(%s)\n%s", schema, jsonValErr))
	}
}
