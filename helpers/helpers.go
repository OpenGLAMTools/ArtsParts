package helpers

// StringInSlice checks if a string exists inside a slice of string
func StringInSlice(search string, sl []string) bool {
	for _, s := range sl {
		if s == search {
			return true
		}
	}
	return false
}
