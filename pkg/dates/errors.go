package dates

import "github.com/sasalatart/batcoms/domain"

// ErrNotDate is used to communicate that a value is not able to be parsed as a date
const ErrNotDate = domain.Error("Value is not a date")
