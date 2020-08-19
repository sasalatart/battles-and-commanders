package io

import "github.com/sasalatart/batcoms/domain"

// Exporter interface defines behaviour for exporting data to a specified file
type Exporter interface {
	Export(fileName string, d interface{}) error
}

// ExporterFunc defines a function for exporting data to a specified file
type ExporterFunc func(fileName string, d interface{}) error

// ImportedData contains scraped battles and participants that have been read from a previously
// exported file. These have been indexed by their Wikipedia IDs
type ImportedData struct {
	SBattlesByID    map[string]domain.SBattle
	SFactionsByID   map[string]domain.SParticipant `json:"FactionsByID"`
	SCommandersByID map[string]domain.SParticipant `json:"CommandersByID"`
}
