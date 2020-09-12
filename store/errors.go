package store

import "github.com/sasalatart/batcoms/domain"

// ErrNotFound is used to communicate that a record was not found in the database
const ErrNotFound = domain.Error("Record not Found")
