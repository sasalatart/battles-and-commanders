package store

// Error represents an abstraction of database operations errors for the datastore level
type Error string

func (e Error) Error() string {
	return string(e)
}

// ErrNotFound is used to communicate that a record was not found in the database
const ErrNotFound = Error("Record not Found")
