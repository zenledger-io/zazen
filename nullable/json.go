package nullable

import (
	"bytes"
	"encoding/json"
)

var (
	nullJSON = []byte("null")
)

func marshalJSON(i any, valid bool) ([]byte, error) {
	if !valid {
		return nullJSON, nil
	}

	return json.Marshal(i)
}

func unmarshalJSON(b []byte, ptr any, valid *bool) error {
	var v bool
	defer func() {
		*valid = v
	}()

	if bytes.Equal(nullJSON, b) {
		return nil
	}

	if err := json.Unmarshal(b, ptr); err != nil {
		return err
	}

	v = true
	return nil
}
