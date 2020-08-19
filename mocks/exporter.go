package mocks

// Exporter mocks a struct that implements the io.Exporter interface, defining behaviour for
// exporting scraped data indexed by ID into a file, such as battles or participants
type Exporter struct {
	CalledTimes   uint
	FileNamesUsed []string
}

// Export is the mock implementation of the Export func specified by the Exporter interface
func (e *Exporter) Export(fileName string, dataByID interface{}) error {
	e.CalledTimes++
	e.FileNamesUsed = append(e.FileNamesUsed, fileName)
	return nil
}
