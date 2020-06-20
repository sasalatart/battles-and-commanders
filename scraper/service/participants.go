package service

import (
	"fmt"
	"strings"

	"github.com/sasalatart/batcoms/scraper/domain"

	"github.com/gocolly/colly"
)

const factionsCSSID = "batcoms-factions"
const commandersCSSID = "batcoms-commanders"

type onParticipantsSideDone func(id int, flagURL string)

func (s *Scraper) subscribeFactions(c *colly.Collector, b *domain.Battle) {
	setInfoBoxID(c, "belligerents", factionsCSSID)
	c.OnHTML(infoBoxSelector, func(e *colly.HTMLElement) {
		s.participantsSide(e, domain.FactionKind, "td:first-child", func(id int, _ string) {
			*&b.Factions.A = append(*&b.Factions.A, id)
		})
		s.participantsSide(e, domain.FactionKind, "td:nth-child(2)", func(id int, _ string) {
			*&b.Factions.B = append(*&b.Factions.B, id)
		})
	})
}

func (s *Scraper) subscribeCommanders(c *colly.Collector, b *domain.Battle) {
	setInfoBoxID(c, "Commanders and leaders", commandersCSSID)

	glueCommander := func(id int, flag string) {
		if f := s.ParticipantsStore.FindFactionByFlag(flag); f != nil {
			b.CommandersByFaction[f.ID] = append(b.CommandersByFaction[f.ID], id)
		}
	}

	c.OnHTML(infoBoxSelector, func(e *colly.HTMLElement) {
		s.participantsSide(e, domain.CommanderKind, "td:first-child", func(id int, flag string) {
			*&b.Commanders.A = append(*&b.Commanders.A, id)
			glueCommander(id, flag)
		})
		s.participantsSide(e, domain.CommanderKind, "td:nth-child(2)", func(id int, flag string) {
			*&b.Commanders.B = append(*&b.Commanders.B, id)
			glueCommander(id, flag)
		})
	})
}

func participantSelector(kind domain.ParticipantKind, sideSelector string) string {
	if kind == domain.FactionKind {
		return fmt.Sprintf("#%s > %s", factionsCSSID, sideSelector)
	}

	return fmt.Sprintf("#%s > %s", commandersCSSID, sideSelector)
}

func (s *Scraper) participantsSide(e *colly.HTMLElement, kind domain.ParticipantKind, sideSelector string, onDone onParticipantsSideDone) {
	fullSelector := participantSelector(kind, sideSelector)
	e.ForEach(fullSelector, func(_ int, side *colly.HTMLElement) {
		side.ForEach("a", func(_ int, node *colly.HTMLElement) {
			if node.Text == "" {
				return
			}

			participantURL := "https://en.wikipedia.org" + node.Attr("href")
			if p := s.ParticipantsStore.FindByURL(kind, participantURL); p != nil {
				onDone(p.ID, flagURL(node))
				return
			}

			summary, err := PageSummary(participantURL)
			if err != nil {
				message := fmt.Sprintf("Failed to fetch summary URL %s: %s", participantURL, err.Error())
				s.logger.Write([]byte(message))
				return
			}

			flag := flagURL(node)
			participant := domain.Participant{
				Kind:        kind,
				ID:          int(summary.PageID),
				URL:         participantURL,
				Flag:        flag,
				Name:        summary.DisplayTitle,
				Description: summary.Description,
				Extract:     summary.Extract,
			}
			// TODO: handle error case
			s.ParticipantsStore.Save(participant)
			onDone(participant.ID, flag)
		})
	})
}

func flagURL(participantNode *colly.HTMLElement) string {
	prevNode := participantNode.DOM.Prev()
	if !prevNode.HasClass("flagicon") {
		return ""
	}

	flagSRC, flagExists := prevNode.Find("img").Attr("src")
	if !flagExists {
		return ""
	}
	return strings.ReplaceAll(flagSRC, "//upload.wikimedia.org/wikipedia/commons", "")
}
