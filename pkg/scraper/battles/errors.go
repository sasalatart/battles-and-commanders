package battles

import "github.com/sasalatart/batcoms/domain"

// ErrMoreThanOneInfoBox is used to communicate that more than one info box was found
const ErrMoreThanOneInfoBox = domain.Error("More than one info box found")

// ErrNoInfoBox is used to communicate that no info box was found
const ErrNoInfoBox = domain.Error("No info box found. Is this a battle?")

// ErrNoSummaryExtract is used to communicate that no extract was found in the battle's summary
const ErrNoSummaryExtract = domain.Error("No extract found in summary")

// ErrNoDate is used to communicate that the info box did not have a date
const ErrNoDate = domain.Error("No date in info box")

// ErrNoResult is used to communicate that the info box did not specify the results
const ErrNoResult = domain.Error("No results in info box")

// ErrNoPlace is used to communicate that the info box did not specify a place
const ErrNoPlace = domain.Error("No place in info box")
