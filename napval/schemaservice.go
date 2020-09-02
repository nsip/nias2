// schemaservice.go
package napval

import (
	"encoding/json"
	"github.com/nsip/nias2/lib"
	"github.com/nsip/nias2/xml"
	"github.com/tidwall/sjson"
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
	"log"
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
func (ss *SchemaService) HandleMessage(req *lib.NiasMessage) ([]lib.NiasMessage, error) {

	responses := make([]lib.NiasMessage, 0)

	// extract reg data from message as json
	r := req.Body.(xml.RegistrationRecord)
	var r_ptr = &r
	r_ptr.Flatten()
	data, err := json.Marshal(*r_ptr)
	if err != nil {
		return responses, err
	}

	// remove extraneous elements
	data, _ = sjson.DeleteBytes(data, "XMLName")
	data, _ = sjson.DeleteBytes(data, "OtherIdList")
	data, _ = sjson.DeleteBytes(data, "RefId")
	data, _ = sjson.DeleteBytes(data, "OtherEnrollmentSchoolACARAId")
	// validate with schema
	payloadLoader := gojsonschema.NewStringLoader(string(data))
	result, err := ss.schema.Validate(payloadLoader)
	if err != nil {
		return responses, err
	}

	if !result.Valid() {
		seen := make(map[string]bool)
		for _, desc := range result.Errors() {
			field := desc.Field()
			field1 := desc.Details()["field"].(string)
			field2 := desc.Details()["property"]
			//log.Printf("ERRRORRRRR: %+v: [%s] [%s] [%s] %s\n", desc, field, field1, field2, req.SeqNo)
			if field == "(root)" && field1 != "" {
				field = field1
			}
			if field == "(root)" && field2 != nil && field2.(string) != "" {
				field = field2.(string)
			}

			// Handle dependencies  errors
			if desc.Field() == "(root)" {
				if desc.Description() == "Has a dependency on Parent2NonSchoolEducation" {
					field = "Parent2NonSchoolEducation"
					desc.SetDescription("Must be present if other Parent2 fields are present")
				}
				if desc.Description() == "Has a dependency on Parent2SchoolEducation" {
					field = "Parent2SchoolEducation"
					desc.SetDescription("Must be present if other Parent2 fields are present")
				}
				if desc.Description() == "Has a dependency on Parent2Occupation" {
					field = "Parent2Occupation"
					desc.SetDescription("Must be present if other Parent2 fields are present")
				}
				if desc.Description() == "Has a dependency on Parent2LOTE" {
					field = "Parent2LOTE"
					desc.SetDescription("Must be present if other Parent2 fields are present")
				}
				// report any given field only once per line
				if seen[string(req.SeqNo)+field] {
					continue
				}
				seen[string(req.SeqNo)+field] = true
			}

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
			if desc.Type() == "pattern" {
				switch desc.Field() {
				case "PlatformId", "PreviousPlatformId":
					desc.SetDescription("PlatformId is not in correct format")
				}
			}

			ve := ValidationError{
				Description:  desc.Description(),
				Field:        field,
				OriginalLine: req.SeqNo,
				Vtype:        "content",
				Severity:     "error",
			}

			r := lib.NiasMessage{}
			r.TxID = req.TxID
			r.SeqNo = req.SeqNo
			// r.Target = VALIDATION_PREFIX
			r.Body = ve
			responses = append(responses, r)
		}
	}

	return responses, nil
}
