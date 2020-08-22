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

type participantsIDsByWikiID map[int]uuid.UUID

func init() {
	config.Setup()
}

func main() {
	db := postgresql.Connect()
	defer db.Close()

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;`)
	schemas := []interface{}{&schema.BattleCommanderFaction{}, &schema.BattleFaction{}, &schema.BattleCommander{}, &schema.Faction{}, &schema.Commander{}, &schema.Battle{}}
	for _, s := range schemas {
		db.DropTableIfExists(s)
		db.AutoMigrate(s)
	}
	db.Model(&schema.BattleFaction{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleFaction{}).AddForeignKey("faction_id", "factions(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommander{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommander{}).AddForeignKey("commander_id", "commanders(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("battle_id", "battles(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("commander_id", "commanders(id)", "CASCADE", "CASCADE")
	db.Model(&schema.BattleCommanderFaction{}).AddForeignKey("faction_id", "factions(id)", "CASCADE", "CASCADE")

	importedData := &io.ImportedData{}
	battlesFileName := viper.GetString("SCRAPER_RESULTS.BATTLES")
	participantsFileName := viper.GetString("SCRAPER_RESULTS.PARTICIPANTS")
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

func seedCommanders(db *gorm.DB, importedData *io.ImportedData) participantsIDsByWikiID {
	store := postgresql.NewCommandersDataStore(db)
	cIDsByWikiID := make(participantsIDsByWikiID)
	counter := 0
	for _, sc := range importedData.SCommandersByID {
		fmt.Printf("\rSeeding commander %d/%d", counter, len(importedData.SCommandersByID))
		counter++
		input := domain.CreateCommanderInput{
			WikiID:  sc.ID,
			URL:     sc.URL,
			Name:    sc.Name,
			Summary: sc.Extract,
		}
		cID, err := store.CreateOne(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cIDsByWikiID[sc.ID] = cID
	}
	fmt.Println("\nFinished seeding commanders")
	return cIDsByWikiID
}

func seedFactions(db *gorm.DB, importedData *io.ImportedData) participantsIDsByWikiID {
	store := postgresql.NewFactionsDataStore(db)
	fIDsByWikiID := make(participantsIDsByWikiID)
	counter := 0
	for _, sf := range importedData.SFactionsByID {
		fmt.Printf("\rSeeding faction %d/%d", counter, len(importedData.SFactionsByID))
		counter++
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
	fmt.Println("\nFinished seeding factions")
	return fIDsByWikiID
}

func translateWikiIDs(from []int, idsMapper participantsIDsByWikiID) []uuid.UUID {
	result := []uuid.UUID{}
	for _, wikiID := range from {
		id, ok := idsMapper[wikiID]
		if !ok {
			fmt.Printf("\rID not found for WikiID %d\n", wikiID)
			continue
		}
		result = append(result, id)
	}
	return result
}

func seedBattles(db *gorm.DB, importedData *io.ImportedData, cIDsByWikiID, fIDsByWikiID participantsIDsByWikiID) {
	store := postgresql.NewBattlesDataStore(db)
	counter := 0
	for _, sb := range importedData.SBattlesByID {
		fmt.Printf("\rSeeding battle %d/%d", counter, len(importedData.SBattlesByID))
		counter++
		dates, err := parser.Date(sb.Date)
		if err != nil {
			fmt.Printf("\rUnable to parse date for %s (given date was %s)\n", sb.Name, sb.Date)
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
		input.FactionsBySide.A = translateWikiIDs(sb.Factions.A, fIDsByWikiID)
		input.FactionsBySide.B = translateWikiIDs(sb.Factions.B, fIDsByWikiID)
		input.CommandersBySide.A = translateWikiIDs(sb.Commanders.A, cIDsByWikiID)
		input.CommandersBySide.B = translateWikiIDs(sb.Commanders.B, cIDsByWikiID)
		for sfWikiID, scWikiIDs := range sb.CommandersByFaction {
			fID := fIDsByWikiID[sfWikiID]
			input.CommandersByFaction[fID] = translateWikiIDs(scWikiIDs, cIDsByWikiID)
		}
		if _, err := store.CreateOne(input); err != nil {
			fmt.Printf("\r%s\n", err)
			continue
		}
	}
	fmt.Println("\nFinished seeding battles")
}
