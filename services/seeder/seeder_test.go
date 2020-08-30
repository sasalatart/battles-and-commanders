package seeder_test

import (
	"strconv"
	"testing"

	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/mocks"
	"github.com/sasalatart/batcoms/services/io"
	"github.com/sasalatart/batcoms/services/seeder"
)

func TestSeeder(t *testing.T) {
	importedData := io.ImportedData{
		SBattlesByID: map[string]domain.SBattle{
			strconv.Itoa(mocks.SBattle().ID): mocks.SBattle(),
		},
		SFactionsByID: map[string]domain.SParticipant{
			strconv.Itoa(mocks.SFaction().ID):  mocks.SFaction(),
			strconv.Itoa(mocks.SFaction2().ID): mocks.SFaction2(),
			strconv.Itoa(mocks.SFaction3().ID): mocks.SFaction3(),
		},
		SCommandersByID: map[string]domain.SParticipant{
			strconv.Itoa(mocks.SCommander().ID):  mocks.SCommander(),
			strconv.Itoa(mocks.SCommander2().ID): mocks.SCommander2(),
			strconv.Itoa(mocks.SCommander3().ID): mocks.SCommander3(),
			strconv.Itoa(mocks.SCommander4().ID): mocks.SCommander4(),
			strconv.Itoa(mocks.SCommander5().ID): mocks.SCommander5(),
		},
	}

	fc := new(mocks.FactionsDataStore)
	fc.On("CreateOne", mocks.CreateFactionInput()).Return(mocks.Faction().ID, nil)
	fc.On("CreateOne", mocks.CreateFactionInput2()).Return(mocks.Faction2().ID, nil)
	fc.On("CreateOne", mocks.CreateFactionInput3()).Return(mocks.Faction3().ID, nil)

	cc := new(mocks.CommandersDataStore)
	cc.On("CreateOne", mocks.CreateCommanderInput()).Return(mocks.Commander().ID, nil)
	cc.On("CreateOne", mocks.CreateCommanderInput2()).Return(mocks.Commander2().ID, nil)
	cc.On("CreateOne", mocks.CreateCommanderInput3()).Return(mocks.Commander3().ID, nil)
	cc.On("CreateOne", mocks.CreateCommanderInput4()).Return(mocks.Commander4().ID, nil)
	cc.On("CreateOne", mocks.CreateCommanderInput5()).Return(mocks.Commander5().ID, nil)

	bc := new(mocks.BattlesDataStore)
	bc.On("CreateOne", mocks.CreateBattleInput()).Return(mocks.Battle().ID, nil)

	service := seeder.Service{
		FactionsCreator:   fc,
		CommandersCreator: cc,
		BattlesCreator:    bc,
		ImportedData:      &importedData,
		Logger:            new(mocks.Logger),
	}
	service.Seed()
	fc.AssertExpectations(t)
	cc.AssertExpectations(t)
	bc.AssertExpectations(t)
}
