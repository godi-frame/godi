package internal

func KeyInStrings(key string, strings []string) bool {
	for _, str := range strings {
		if key == str {
			return true
		}
	}
	return false
}
