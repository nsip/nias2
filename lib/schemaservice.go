// schemaservice.go
package nias2

import (
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	// "log"
)

// implementaiton of the json schema validation service
type SchemaService struct {
	schema     *gojsonschema.Schema
	schemaFile string
}

// use this method to return a schema validator initialised with the
// standard core json schema schemas/core.json
// typically returns an error if there are issues loading the core schema from file
func NewCoreSchemaService() (*SchemaService, error) {
	ss := SchemaService{}
	ss.schemaFile = "schemas/core.json"
	if err := ss.loadSchema(); err != nil {
		return &ss, err
	}
	return &ss, nil
}

// use this method to create a validator that will use the
// schema file given in filename (eg. local.json/sa_catholic.json)
// file is searched for in schemas/ directory
func NewCustomSchemaService(filename string) (*SchemaService, error) {
	ss := SchemaService{}
	ss.schemaFile = "schemas/" + filename
	if err := ss.loadSchema(); err != nil {
		return &ss, err
	}
	return &ss, nil
}

// loads the json-schema file for use in validation,
// will use the filename provided in SchemaService.schemaFile()
// and will look in directory /schemas for specified file.
func (ss *SchemaService) loadSchema() error {
	s, readerr := ioutil.ReadFile(ss.schemaFile)
	if readerr != nil {
		return readerr
	}
	schemaLoader := gojsonschema.NewStringLoader(string(s))
	schema, err := gojsonschema.NewSchema(schemaLoader)
	if err != nil {
		return err
	}
	// log.Println("loaded schema - schemas/" + ss.schemaFile)
	ss.schema = schema

	return nil
}

// implement the nias Service interface
func (ss *SchemaService) HandleMessage(req *NiasMessage) ([]NiasMessage, error) {

	responses := make([]NiasMessage, 0)

	// extract reg data from message as json
	data, err := json.Marshal(req.Body.(RegistrationRecord))
	if err != nil {
		return responses, err
	}

	// validate with schema
	payloadLoader := gojsonschema.NewStringLoader(string(data))
	result, err := ss.schema.Validate(payloadLoader)
	if err != nil {
		return responses, err
	}

	if !result.Valid() {

		for _, desc := range result.Errors() {

			// trap enum errors for large enums such as country code
			// and truncate the message to prevent unwieldy message
			if desc.Type() == "enum" {
				switch desc.Field() {
				case "CountryOfBirth":
					desc.SetDescription("Country Code is not one of SACC 1269.0 codeset")
				case "Parent1LOTE", "Parent2LOTE", "StudentLOTE":
					desc.SetDescription("Language Code is not one of ASCL 1267.0 codeset")
				case "VisaCode":
					desc.SetDescription("Visa Code is not one of known values from http://www.immi.gov.au")
				}

			}

			ve := ValidationError{
				Description:  desc.Description(),
				Field:        desc.Field(),
				OriginalLine: req.SeqNo,
				Vtype:        "content",
			}

			r := NiasMessage{}
			r.TxID = req.TxID
			r.SeqNo = req.SeqNo
			r.Target = VALIDATION_PREFIX
			r.Body = ve
			responses = append(responses, r)
		}
	}

	return responses, nil
}
