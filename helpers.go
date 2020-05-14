package dmon

func stringIndex(strs []string, str string) int {
	for i, v := range strs {
		if v == str {
			return i
		}
	}
	return -1
}

func removeStringFromSlice(strs []string, str string) []string {
	id := stringIndex(strs, str)
	if id != -1 {
		return append(strs[:id], strs[id+1:]...)
	}
	return strs
}
