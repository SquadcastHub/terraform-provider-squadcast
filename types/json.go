// Package types implements all the types needed for this
// service
package types

import (
	"bytes"
	"encoding/json"
	"io"
)

// JSON maps go type to a JSON struct
type JSON map[string]interface{}

// ToReader returns a io.Reader after encoding the map into JSON byte array
func (j JSON) ToReader() io.Reader {
	b := bytes.NewBuffer(nil)
	json.NewEncoder(b).Encode(j)
	return b
}

// Array corresponds to an array of elements with any type
type Array interface{}
