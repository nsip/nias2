// encoding.go

package naprrql

import (
	"bytes"
	"encoding/gob"

	"github.com/nsip/nias2/xml"
)

// ensure all data types are registered with the encoder
func init() {
	gob.Register(xml.NAPEvent{})
	gob.Register(xml.NAPCodeFrame{})
	gob.Register(xml.NAPResponseSet{})
	gob.Register(xml.NAPTest{})
	gob.Register(xml.NAPTestItem{})
	gob.Register(xml.NAPTestlet{})
	gob.Register(xml.NAPTestScoreSummary{})
	gob.Register(xml.RegistrationRecord{})
	gob.Register(xml.SchoolInfo{})
}

// GobEncoder is a Go specific GOB Encoder implementation for EncodedConn.
// This encoder will use the builtin encoding/gob to Marshal
// and Unmarshal most types, including structs.
type GobEncoder struct {
	// Empty
}

// Encode
// note: encoding the pointer is deliberate as forces
// encoded types to be interface{}
func (ge *GobEncoder) Encode(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	if err := enc.Encode(&v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Decode
func (ge *GobEncoder) Decode(data []byte, vPtr interface{}) (err error) {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	err = dec.Decode(vPtr)
	return
}
