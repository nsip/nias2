// naprr_naplandata.go
package naprr

import (
	"github.com/nsip/nias2/xml"
	// "log"
)

// structs to hold naplan meta-data that is the same for all schools

type NAPLANData struct {
	Tests      map[string]xml.NAPTest
	Testlets   map[string]xml.NAPTestlet
	Items      map[string]xml.NAPTestItem
	Codeframes map[string]xml.NAPCodeFrame
}

func NewNAPLANData() *NAPLANData {
	nd := NAPLANData{
		Tests:      make(map[string]xml.NAPTest),
		Testlets:   make(map[string]xml.NAPTestlet),
		Items:      make(map[string]xml.NAPTestItem),
		Codeframes: make(map[string]xml.NAPCodeFrame),
	}
	return &nd
}
