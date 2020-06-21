package exports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// JSON creates a file and saves the specified contents inside it in JSON format
func JSON(fileName string, d interface{}) error {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(d); err != nil {
		return fmt.Errorf("Failed encoding %s as JSON: %s", fileName, err)
	}
	if err := ioutil.WriteFile(fileName, buffer.Bytes(), 0644); err != nil {
		return fmt.Errorf("Failed to write %s: %s", fileName, err)
	}
	return nil
}
