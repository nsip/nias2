package napval

import (
	"encoding/gob"
)

// ensures transmissable types are registered for binary encoding
func init() {
	// make gob encoder aware of local types
	gob.Register(ValidationError{})
}

// struct to handle reporting of validation errors found in
// naplan registration files
type ValidationError struct {
	Field        string `json:"errField"`     // the field that has an error
	Description  string `json:"description"`  // error description
	OriginalLine string `json:"originalLine"` // input file record line that has the error
	Vtype        string `json:"validationType"`
	Severity     string `json:"severity"` // warning, error
}

// helper method for writing out csv encoding of error reports
func (ve *ValidationError) ToSlice() []string {

	return []string{ve.OriginalLine, ve.Vtype, ve.Field, ve.Description, ve.Severity}
}
