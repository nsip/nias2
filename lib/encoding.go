// encoding.go
package lib

import (
	"bytes"
	"encoding/gob"
	"log"
	"strings"
)

// helper routines collection for encoding nias messages

// binary encding for messages going to internal q/store.
func EncodeNiasMessage(msg *NiasMessage) []byte {

	encBuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(encBuf)
	err := encoder.Encode(msg)
	if err != nil {
		log.Printf("Encoder unable to binary encode message for: %#v\n", msg)
	}
	return encBuf.Bytes()

}

// binary decoding for messages coming from internal q/store.
func DecodeNiasMessage(bytemsg []uint8) *NiasMessage {

	decBuf := bytes.NewBuffer(bytemsg)
	decoder := gob.NewDecoder(decBuf)
	var msgOut NiasMessage
	err := decoder.Decode(&msgOut)
	if err != nil {
		log.Println("Error decoding message from q/store(internal):", err)
	}
	return &msgOut
}

// truncate the record by removing items that have blank entries.
// this prevents the validation from throwing validation exceptions
// for fields that are not mandatory but included as empty in the
// dataset
func RemoveBlanks(m map[string]string) map[string]string {

	reducedmap := make(map[string]string)
	for key, val := range m {
		if val != "" {
			reducedmap[key] = strings.TrimSpace(val)
		}
	}
	return reducedmap
}
