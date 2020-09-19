package seeder_test

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/sasalatart/batcoms/db/seeder"
	"github.com/sasalatart/batcoms/domain/wikiactors"
	"github.com/sasalatart/batcoms/domain/wikibattles"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/pkg/logger"
)

func TestSeeder(t *testing.T) {
	importedData := seeder.ImportedData{
		WikiBattlesByID: map[string]wikibattles.Battle{
			strconv.Itoa(mocks.WikiBattle().ID): mocks.WikiBattle(),
		},
		WikiFactionsByID: map[string]wikiactors.Actor{
			strconv.Itoa(mocks.WikiFaction().ID):  mocks.WikiFaction(),
			strconv.Itoa(mocks.WikiFaction2().ID): mocks.WikiFaction2(),
			strconv.Itoa(mocks.WikiFaction3().ID): mocks.WikiFaction3(),
		},
		WikiCommandersByID: map[string]wikiactors.Actor{
			strconv.Itoa(mocks.WikiCommander().ID):  mocks.WikiCommander(),
			strconv.Itoa(mocks.WikiCommander2().ID): mocks.WikiCommander2(),
			strconv.Itoa(mocks.WikiCommander3().ID): mocks.WikiCommander3(),
			strconv.Itoa(mocks.WikiCommander4().ID): mocks.WikiCommander4(),
			strconv.Itoa(mocks.WikiCommander5().ID): mocks.WikiCommander5(),
		},
	}

	fr := new(mocks.FactionsRepository)
	fr.On("CreateOne", mocks.FactionCreationInput()).Return(mocks.Faction().ID, nil)
	fr.On("CreateOne", mocks.FactionCreationInput2()).Return(mocks.Faction2().ID, nil)
	fr.On("CreateOne", mocks.FactionCreationInput3()).Return(mocks.Faction3().ID, nil)

	cr := new(mocks.CommandersRepository)
	cr.On("CreateOne", mocks.CommanderCreationInput()).Return(mocks.Commander().ID, nil)
	cr.On("CreateOne", mocks.CommanderCreationInput2()).Return(mocks.Commander2().ID, nil)
	cr.On("CreateOne", mocks.CommanderCreationInput3()).Return(mocks.Commander3().ID, nil)
	cr.On("CreateOne", mocks.CommanderCreationInput4()).Return(mocks.Commander4().ID, nil)
	cr.On("CreateOne", mocks.CommanderCreationInput5()).Return(mocks.Commander5().ID, nil)

	br := new(mocks.BattlesRepository)
	br.On("CreateOne", mocks.BattleCreationInput()).Return(mocks.Battle().ID, nil)

	seeder.Seed(&importedData, fr, cr, br, logger.New(ioutil.Discard, ioutil.Discard))
	fr.AssertExpectations(t)
	cr.AssertExpectations(t)
	br.AssertExpectations(t)
}
