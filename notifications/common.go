package notifications

// Contains returns true if the array a contains the string x, and false otherwise.
func Contains(a []string, x string) bool {
	for _, s := range a {
		if s == x {
			return true
		}
	}
	return false
}
