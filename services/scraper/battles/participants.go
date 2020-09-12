package battles

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/parser"
	"github.com/sasalatart/batcoms/services/scraper/urls"
	"github.com/sasalatart/batcoms/services/summaries"

	"github.com/gocolly/colly"
)

type factionsMapper map[string]int
type commandersMapper map[string][]int
type onParticipantDone func(id int, flagURL string, err error)

func (s *Scraper) subscribeParticipants(ctx *battleCtx) {
	factionsByFlag := make(factionsMapper)
	commandersByFlag := make(commandersMapper)
	s.subscribeFactions(ctx, factionsByFlag)
	s.subscribeCommanders(ctx, commandersByFlag)
	s.subscribeGroupings(ctx, factionsByFlag, commandersByFlag)
}

func (s *Scraper) subscribeFactions(ctx *battleCtx, fm factionsMapper) {
	subscribeSetInfoBoxID(ctx, "belligerents", customFactionsID)

	handleFaction := func(id int, flag string, ids *[]int, err error) {
		if err != nil {
			s.logger.Error(err)
			return
		}
		if _, saved := fm[flag]; !saved && flag != "" {
			fm[flag] = id
		}
		*ids = append(*ids, id)
	}

	ctx.collector.OnHTML(infoBoxSelector, ctx.abortable(func(e *colly.HTMLElement) {
		s.participantsSide(ctx, e, domain.FactionKind, sideASelector, func(id int, flag string, err error) {
			handleFaction(id, flag, &ctx.battle.Factions.A, err)
		})
		s.participantsSide(ctx, e, domain.FactionKind, sideBSelector, func(id int, flag string, err error) {
			handleFaction(id, flag, &ctx.battle.Factions.B, err)
		})
	}))
}

func (s *Scraper) subscribeCommanders(ctx *battleCtx, cm commandersMapper) {
	subscribeSetInfoBoxID(ctx, "Commanders and leaders", customCommandersID)

	handleCommander := func(id int, flag string, ids *[]int, err error) {
		if err != nil {
			s.logger.Error(err)
			return
		}
		if flag != "" {
			cm[flag] = append(cm[flag], id)
		}
		*ids = append(*ids, id)
	}

	ctx.collector.OnHTML(infoBoxSelector, ctx.abortable(func(e *colly.HTMLElement) {
		s.participantsSide(ctx, e, domain.CommanderKind, sideASelector, func(id int, flag string, err error) {
			handleCommander(id, flag, &ctx.battle.Commanders.A, err)
		})
		s.participantsSide(ctx, e, domain.CommanderKind, sideBSelector, func(id int, flag string, err error) {
			handleCommander(id, flag, &ctx.battle.Commanders.B, err)
		})
	}))
}

func (s *Scraper) participantsSide(ctx *battleCtx, e *colly.HTMLElement, kind domain.ParticipantKind, sideSelector string, onDone onParticipantDone) {
	e.ForEach(participantsSelector(kind, sideSelector), func(_ int, side *colly.HTMLElement) {
		side.ForEach("a", func(_ int, node *colly.HTMLElement) {
			if ctx.err != nil {
				return
			}

			cleanText := parser.Clean(node.Text)
			if cleanText == "" {
				return
			}

			pURL := node.Attr("href")
			if !strings.Contains(pURL, "://") {
				pURL = "https://en.wikipedia.org" + pURL
			}
			if urls.ShouldSkip(pURL) {
				return
			}

			if p := s.sParticipantsStore.FindByURL(kind, pURL); p != nil {
				onDone(p.ID, flagURL(node), nil)
				return
			}

			summary, err := summaries.Fetch(pURL)
			if err != nil {
				s.logger.Error(errors.Wrapf(err, "Fetching summary for %s", pURL))
				return
			}

			name := summary.Title
			if urls.NotSpecific(pURL) {
				name = cleanText
			}

			flag := flagURL(node)
			participant := domain.SParticipant{
				Kind:        kind,
				ID:          int(summary.PageID),
				URL:         pURL,
				Flag:        flag,
				Name:        name,
				Description: summary.Description,
				Extract:     summary.Extract,
			}
			err = s.sParticipantsStore.Save(participant)
			onDone(participant.ID, flag, err)
		})
	})
}

func (s *Scraper) subscribeGroupings(ctx *battleCtx, factionsByFlag factionsMapper, commandersByFlag commandersMapper) {
	ctx.collector.OnScraped(func(_ *colly.Response) {
		ctx.battle.Factions.A = unique(ctx.battle.Factions.A)
		ctx.battle.Factions.B = unique(ctx.battle.Factions.B)
		ctx.battle.Commanders.A = unique(ctx.battle.Commanders.A)
		ctx.battle.Commanders.B = unique(ctx.battle.Commanders.B)

		for fFlag, fID := range factionsByFlag {
			for cFlag, cIDs := range commandersByFlag {
				if fFlag != cFlag {
					continue
				}

				ctx.battle.CommandersByFaction[fID] = cIDs
				break
			}
		}

		groupIfOneFaction := func(factions []int, commanders []int) {
			if len(factions) != 1 || len(commanders) == 0 {
				return
			}

			ctx.battle.CommandersByFaction[factions[0]] = commanders
		}
		groupIfOneFaction(ctx.battle.Factions.A, ctx.battle.Commanders.A)
		groupIfOneFaction(ctx.battle.Factions.B, ctx.battle.Commanders.B)
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

func unique(ids []int) []int {
	set := make(map[int]struct{})
	result := []int{}
	for _, id := range ids {
		if _, present := set[id]; !present {
			set[id] = struct{}{}
			result = append(result, id)
		}
	}
	return result
}
