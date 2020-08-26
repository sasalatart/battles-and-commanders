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

// FactionsFinder is the interface through which factions may be found
type FactionsFinder interface {
	FindOne(query interface{}, args ...interface{}) (domain.Faction, error)
}

// CommandersFinder is the interface through which commanders may be found
type CommandersFinder interface {
	FindOne(query interface{}, args ...interface{}) (domain.Commander, error)
}

// BattlesFinder is the interface through which battles may be found
type BattlesFinder interface {
	FindOne(query interface{}, args ...interface{}) (domain.Battle, error)
}
