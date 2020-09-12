package domain

// Error represents error messages originated from business logic, and also provides an easier
// mechanism to tell one error from another
type Error string

// Error is like errors.Error, but for domain.Error
func (e Error) Error() string {
	return string(e)
}
