package exports

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// JSON creates a file and saves the specified contents inside it in JSON format
func JSON(fileName string, d interface{}) error {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(d); err != nil {
		return errors.Wrapf(err, "Encoding %s as JSON", fileName)
	}
	if err := ioutil.WriteFile(fileName, buffer.Bytes(), 0644); err != nil {
		return errors.Wrapf(err, "Writing %s", fileName)
	}
	return nil
}
