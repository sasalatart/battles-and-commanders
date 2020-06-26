package urls

// BattlesLists returns the list of all lists of battles that should be scraped.
func BattlesLists() []string {
	var urlsParts = []string{
		"/Battles_of_the_Seven_Years%27_War",
		"/List_of_American_Civil_War_battles",
		"/List_of_American_Revolutionary_War_battles",
		"/List_of_battles_(alphabetical)",
		"/List_of_battles_(geographic)",
		"/List_of_battles_301-1300",
		"/List_of_battles_1301-1600",
		"/List_of_battles_1601-1800",
		"/List_of_battles_1801-1900",
		"/List_of_battles_1901-2000",
		"/List_of_battles_before_301",
		"/List_of_battles_since_2001",
		"/List_of_Hundred_Years%27_War_battles",
		"/List_of_military_engagements_of_World_War_I",
		"/List_of_military_engagements_of_World_War_II",
		"/List_of_Napoleonic_battles",
	}

	result := []string{}
	for _, p := range urlsParts {
		result = append(result, "https://en.wikipedia.org/wiki"+p)
	}
	return result
}
