package helpers

import "os"

// StringInSlice checks if a string exists inside a slice of string
func StringInSlice(search string, sl []string) bool {
	for _, s := range sl {
		if s == search {
			return true
		}
	}
	return false
}

// FileExists returns true if a file exists
func FileExists(fpath string) bool {
	_, err := os.Stat(fpath)
	if err == nil {
		return true
	}
	return false
}
