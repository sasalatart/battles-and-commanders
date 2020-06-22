package store

import "github.com/sasalatart/batcoms/scraper/domain"

// ParticipantsStore is the interface through which participants may be saved, found and exported
type ParticipantsStore interface {
	Find(kind domain.ParticipantKind, id int) *domain.Participant
	FindByURL(kind domain.ParticipantKind, url string) *domain.Participant
	Save(p domain.Participant) error
	Export(fileName string) error
}
