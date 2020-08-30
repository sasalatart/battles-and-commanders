package store

import (
	"github.com/sasalatart/batcoms/domain"
	uuid "github.com/satori/go.uuid"
)

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

// FactionsFinder is the interface through which factions may be found
type FactionsFinder interface {
	FindOne(query interface{}, args ...interface{}) (domain.Faction, error)
}

// FactionsCreator is the interface through which factions may be created
type FactionsCreator interface {
	CreateOne(data domain.CreateFactionInput) (uuid.UUID, error)
}

// CommandersFinder is the interface through which commanders may be found
type CommandersFinder interface {
	FindOne(query interface{}, args ...interface{}) (domain.Commander, error)
}

// CommandersCreator is the interface through which commanders may be created
type CommandersCreator interface {
	CreateOne(data domain.CreateCommanderInput) (uuid.UUID, error)
}

// BattlesFinder is the interface through which battles may be found
type BattlesFinder interface {
	FindOne(query interface{}, args ...interface{}) (domain.Battle, error)
}

// BattlesCreator is the interface through which battles may be created
type BattlesCreator interface {
	CreateOne(data domain.CreateBattleInput) (uuid.UUID, error)
}
