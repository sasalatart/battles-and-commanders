package battles

import (
	"fmt"
	"strings"

	"github.com/sasalatart/batcoms/domain"
)

const contentSelector = "#content"

const infoBoxSelector = ".infobox.vevent > tbody"

const partOfSelector = "tr:nth-child(2) > td"
const coordinatesSelector = "#coordinates"
const placeSelector = ".location " + coordinatesSelector

const sideASelector = "td:first-child"
const sideBSelector = "td:nth-child(2)"
const sideABSelector = "tr>td[colspan='2']"

const customFactionsID = "batcoms-factions"
const customCommandersID = "batcoms-commanders"

func customID(text string) string {
	return "batcoms-" + strings.ReplaceAll(strings.ToLower(text), " ", "")
}

func sideNumbersSelector(side, customID string) string {
	if side == sideABSelector {
		return fmt.Sprintf("#%s + %s", customID, side)
	}

	return fmt.Sprintf("#%s > %s", customID, side)
}

func participantsSelector(kind domain.ParticipantKind, sideSelector string) string {
	if kind == domain.FactionKind {
		return fmt.Sprintf("#%s > %s", customFactionsID, sideSelector)
	}

	return fmt.Sprintf("#%s > %s", customCommandersID, sideSelector)
}
