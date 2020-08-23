package store

import "github.com/sasalatart/batcoms/domain"

// SBattles is the interface through which SBattle may be saved, found and exported
type SBattles interface {
	Find(id int) *domain.SBattle
	Save(b domain.SBattle) error
	Export(fileName string) error
}

// SParticipants is the interface through which SParticipant may be saved, found and exported
type SParticipants interface {
	Find(kind domain.ParticipantKind, id int) *domain.SParticipant
	FindByURL(kind domain.ParticipantKind, url string) *domain.SParticipant
	Save(p domain.SParticipant) error
	Export(fileName string) error
}

// Factions is the interface through which Faction may be saved and found
type Factions interface {
	FindOne(query interface{}, args ...interface{}) (domain.Faction, error)
}
