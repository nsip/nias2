// encoding.go

package naprrql

import (
	"encoding/json"
	"log"

	"github.com/nsip/nias2/xml"
)

func jsonEncode(v interface{}) ([]byte, error) {

	switch t := v.(type) {
	case xml.NAPTestScoreSummary:
		tss, ok := v.(xml.NAPTestScoreSummary)
		if !ok {
			log.Println("Error asserting type for NAPTestScoreSummary")
		}
		result, err := json.Marshal(tss)
		// log.Printf("\n\n%s\n\n", string(result))
		return result, err
	case xml.NAPResponseSet:
		rs, ok := v.(xml.NAPResponseSet)
		if !ok {
			log.Println("Error asserting type for NAPResponseSet")
		}
		result, err := json.Marshal(rs)
		// log.Printf("\n\n%s\n\n", string(result))
		return result, err
	default:
		// unknown type, ignore...
		_ = t
	}

	return []byte(""), nil

}

func jsonDecode(data []byte, v interface{}) error {

	err := json.Unmarshal(data, v)
	// log.Printf("\n\n%s\n\n", string(m))
	return err

}
