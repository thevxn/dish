package socket

import "regexp"

// IsFilePath checks whether input is a file path or URL.
func IsFilePath(source string) bool {
	matched, _ := regexp.MatchString("^(http|https)://", source)
	return !matched
}
