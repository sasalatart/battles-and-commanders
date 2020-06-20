package exports

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// JSON saves the specified contents into the specified file name
func JSON(fileName string, toExport interface{}) error {
	file, err := json.MarshalIndent(toExport, "", "")
	if err != nil {
		return fmt.Errorf("Failed to marshal data for %q: %s", fileName, err)
	}
	if err = ioutil.WriteFile(fileName, file, 0644); err != nil {
		return fmt.Errorf("Failed to write %q: %s", fileName, err)
	}
	return nil
}
