package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/sasalatart/batcoms/config"
	"github.com/sasalatart/batcoms/domain"
	"github.com/sasalatart/batcoms/services/io"
	"github.com/sasalatart/batcoms/services/io/json"
	"github.com/sasalatart/batcoms/services/parser"
	"github.com/sasalatart/batcoms/store/postgresql"
	"github.com/sasalatart/batcoms/store/postgresql/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
)

func main() {
	config.Setup()
	vpr := viper.GetViper()

	db := postgresql.Connect()
	defer db.Close()

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;`)
	schemas := []interface{}{&schema.BattleCommanderFaction{}, &schema.BattleCommander{}, &schema.BattleFaction{}, &schema.Battle{}, &schema.Commander{}, &schema.Faction{}}
	for _, s := range schemas {
		db.DropTableIfExists(s)
		db.AutoMigrate(s)
	}
	db.Model(&schema.BattleCommander{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommander{}).AddForeignKey("commander_id", "commanders(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleFaction{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleFaction{}).AddForeignKey("faction_id", "factions(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("commander_id", "commanders(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("faction_id", "factions(id)", "CASCADE", "CASCADE")

	importedData := &io.ImportedData{}
	battlesFileName := vpr.GetString("SCRAPER_RESULTS.BATTLES")
	participantsFileName := vpr.GetString("SCRAPER_RESULTS.PARTICIPANTS")
	if err := json.Import(battlesFileName, &importedData.SBattlesByID); err != nil {
		fmt.Println(err)
	}
	if err := json.Import(participantsFileName, &importedData); err != nil {
		fmt.Println(err)
	}
	fIDsByWikiID := seedFactions(db, importedData)
	cIDsByWikiID := seedCommanders(db, importedData)
	seedBattles(db, importedData, cIDsByWikiID, fIDsByWikiID)
}

func seedCommanders(db *gorm.DB, importedData *io.ImportedData) map[int]uuid.UUID {
	store := postgresql.NewCommandersDataStore(db)
	cIDsByWikiID := make(map[int]uuid.UUID)
	for _, sc := range importedData.SCommandersByID {
		input := domain.CreateCommanderInput{
			WikiID:  sc.ID,
			URL:     sc.URL,
			Name:    sc.Name,
			Summary: sc.Extract,
		}
		id, err := store.CreateOne(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cIDsByWikiID[sc.ID] = id
	}
	return cIDsByWikiID
}

func seedFactions(db *gorm.DB, importedData *io.ImportedData) map[int]uuid.UUID {
	store := postgresql.NewFactionsDataStore(db)
	fIDsByWikiID := make(map[int]uuid.UUID)
	for _, sf := range importedData.SFactionsByID {
		input := domain.CreateFactionInput{
			WikiID:  sf.ID,
			URL:     sf.URL,
			Name:    sf.Name,
			Summary: sf.Extract,
		}
		fID, err := store.CreateOne(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fIDsByWikiID[sf.ID] = fID
	}
	return fIDsByWikiID
}

func collectGroupingIDs(from []int, to *[]uuid.UUID, idsMapper map[int]uuid.UUID) {
	for _, wikiID := range from {
		id, ok := idsMapper[wikiID]
		if !ok {
			fmt.Printf("ID not found for WikiID %d\n", wikiID)
			continue
		}
		*to = append(*to, id)
	}
}

func seedBattles(db *gorm.DB, importedData *io.ImportedData, cIDsByWikiID, fIDsByWikiID map[int]uuid.UUID) {
	store := postgresql.NewBattlesDataStore(db)
	for _, sb := range importedData.SBattlesByID {
		dates, err := parser.Date(sb.Date)
		if err != nil {
			fmt.Printf("Unable to parse date for %s (%s)\n", sb.Name, sb.Date)
			continue
		}
		input := domain.CreateBattleInput{
			WikiID:              sb.ID,
			URL:                 sb.URL,
			Name:                sb.Name,
			PartOf:              sb.PartOf,
			Summary:             sb.Extract,
			StartDate:           dates[0],
			EndDate:             dates[len(dates)-1],
			Location:            sb.Location,
			Result:              sb.Result,
			TerritorialChanges:  sb.TerritorialChanges,
			Strength:            sb.Strength,
			Casualties:          sb.Casualties,
			CommandersByFaction: make(domain.CommandersByFaction),
		}
		collectGroupingIDs(sb.Factions.A, &input.FactionsBySide.A, fIDsByWikiID)
		collectGroupingIDs(sb.Factions.B, &input.FactionsBySide.B, fIDsByWikiID)
		collectGroupingIDs(sb.Commanders.A, &input.CommandersBySide.A, cIDsByWikiID)
		collectGroupingIDs(sb.Commanders.B, &input.CommandersBySide.B, cIDsByWikiID)
		for sfWikiID, scWikiIDs := range sb.CommandersByFaction {
			fID := fIDsByWikiID[sfWikiID]
			var cIDS []uuid.UUID
			for _, cWikiID := range scWikiIDs {
				cIDS = append(cIDS, cIDsByWikiID[cWikiID])
			}
			input.CommandersByFaction[fID] = cIDS
		}
		if _, err := store.CreateOne(input); err != nil {
			fmt.Println(err)
			continue
		}
	}
}
