package store

import "github.com/sasalatart/batcoms/scraper/domain"

// BattlesStore is the interface through which battles may be saved, found and exported
type BattlesStore interface {
	Find(id int) *domain.Battle
	Save(b domain.Battle) error
	Export(fileName string) error
}
