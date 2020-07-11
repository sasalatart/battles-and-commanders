package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/sasalatart/batcoms/parser"
	"github.com/sasalatart/batcoms/scraper/domain"
	"github.com/sasalatart/batcoms/scraper/urls"

	"github.com/gocolly/colly"
)

const factionsCSSID = "batcoms-factions"
const commandersCSSID = "batcoms-commanders"

type factionsMapper map[string]int
type commandersMapper map[string][]int
type onParticipantDone func(id int, flagURL string, err error)

func (s *Scraper) subscribeParticipants(c *colly.Collector, b *domain.Battle) {
	factionsByFlag := make(map[string]int)
	commandersByFlag := make(map[string][]int)
	s.subscribeFactions(c, b, factionsByFlag)
	s.subscribeCommanders(c, b, commandersByFlag)
	subscribeGroupings(c, b, factionsByFlag, commandersByFlag)
}

func (s *Scraper) subscribeFactions(c *colly.Collector, b *domain.Battle, fm factionsMapper) {
	subscribeInfoBoxID(c, "belligerents", factionsCSSID)

	handleFaction := func(id int, flag string, ids *[]int, err error) {
		if err != nil {
			log.Println(err)
			return
		}
		if _, saved := fm[flag]; !saved && flag != "" {
			fm[flag] = id
		}
		*ids = append(*ids, id)
	}

	c.OnHTML(infoBoxSelector, func(e *colly.HTMLElement) {
		s.participantsSide(e, domain.FactionKind, "td:first-child", func(id int, flag string, err error) {
			handleFaction(id, flag, &b.Factions.A, err)
		})
		s.participantsSide(e, domain.FactionKind, "td:nth-child(2)", func(id int, flag string, err error) {
			handleFaction(id, flag, &b.Factions.B, err)
		})
	})
}

func (s *Scraper) subscribeCommanders(c *colly.Collector, b *domain.Battle, cm commandersMapper) {
	subscribeInfoBoxID(c, "Commanders and leaders", commandersCSSID)

	handleCommander := func(id int, flag string, ids *[]int, err error) {
		if err != nil {
			log.Println(err)
			return
		}
		if flag != "" {
			cm[flag] = append(cm[flag], id)
		}
		*ids = append(*ids, id)
	}

	c.OnHTML(infoBoxSelector, func(e *colly.HTMLElement) {
		s.participantsSide(e, domain.CommanderKind, "td:first-child", func(id int, flag string, err error) {
			handleCommander(id, flag, &b.Commanders.A, err)
		})
		s.participantsSide(e, domain.CommanderKind, "td:nth-child(2)", func(id int, flag string, err error) {
			handleCommander(id, flag, &b.Commanders.B, err)
		})
	})
}

func participantSelector(kind domain.ParticipantKind, sideSelector string) string {
	if kind == domain.FactionKind {
		return fmt.Sprintf("#%s > %s", factionsCSSID, sideSelector)
	}

	return fmt.Sprintf("#%s > %s", commandersCSSID, sideSelector)
}

func (s *Scraper) participantsSide(e *colly.HTMLElement, kind domain.ParticipantKind, sideSelector string, onDone onParticipantDone) {
	fullSelector := participantSelector(kind, sideSelector)
	e.ForEach(fullSelector, func(_ int, side *colly.HTMLElement) {
		side.ForEach("a", func(_ int, node *colly.HTMLElement) {
			cleanText := parser.Clean(node.Text)
			if cleanText == "" {
				return
			}

			pURL := node.Attr("href")
			if !strings.Contains(pURL, "://") {
				pURL = "https://en.wikipedia.org" + pURL
			}

			if p := s.ParticipantsStore.FindByURL(kind, pURL); p != nil {
				onDone(p.ID, flagURL(node), nil)
				return
			}
			if urls.ShouldSkip(pURL) {
				return
			}

			summary, err := PageSummary(pURL)
			if err != nil {
				log.Printf("Failed to fetch summary for %s: %s", pURL, err)
				return
			}

			name := summary.DisplayTitle
			if urls.NotSpecific(pURL) {
				name = cleanText
			}

			flag := flagURL(node)
			participant := domain.Participant{
				Kind:        kind,
				ID:          int(summary.PageID),
				URL:         pURL,
				Flag:        flag,
				Name:        name,
				Description: summary.Description,
				Extract:     summary.Extract,
			}
			err = s.ParticipantsStore.Save(participant)
			onDone(participant.ID, flag, err)
		})
	})
}

func subscribeGroupings(c *colly.Collector, b *domain.Battle, factionsByFlag map[string]int, commandersByFlag map[string][]int) {
	handler := func(_ *colly.Response) {
		for fFlag, fID := range factionsByFlag {
			for cFlag, cIDs := range commandersByFlag {
				if fFlag != cFlag {
					continue
				}

				b.CommandersByFaction[fID] = cIDs
				break
			}
		}

		groupIfOneFaction := func(factions []int, commanders []int) {
			if len(factions) != 1 || len(commanders) == 0 {
				return
			}

			b.CommandersByFaction[factions[0]] = commanders
		}
		groupIfOneFaction(b.Factions.A, b.Commanders.A)
		groupIfOneFaction(b.Factions.B, b.Commanders.B)
	}
	c.OnScraped(handler)
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
