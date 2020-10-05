package dates

// IsValid checks whether or not the input string is in "YYYY-MM-DD [BC]" format. For example, the
// following strings are all valid: "1769-08-15", "1769-8-15", "1769-8", "1769" and "1457-04-16 BC"
func IsValid(date string) bool {
	_, err := extract(date)
	return err == nil
}
