package json

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Export creates a file and saves the specified contents inside it in JSON format
func Export(fileName string, d interface{}) error {
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

// Import reads the specified JSON file's contents into memory
func Import(fileName string, d interface{}) error {
	battlesFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrapf(err, "Reading file %s", fileName)
	}
	if err = json.Unmarshal(battlesFile, d); err != nil {
		return errors.Wrapf(err, "Unmarshalling file contents in %s", fileName)
	}
	return nil
}
