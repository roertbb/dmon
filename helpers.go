package dmon

func stringIndex(strs []string, str string) int {
	for i, v := range strs {
		if v == str {
			return i
		}
	}
	return -1
}
