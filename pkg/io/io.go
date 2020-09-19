package io

// Exporter interface defines behaviour for exporting data to a specified file
type Exporter interface {
	Export(fileName string, d interface{}) error
}

// ExporterFunc defines a function for exporting data to a specified file
type ExporterFunc func(fileName string, d interface{}) error
