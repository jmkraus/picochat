package utils

// Global cache for list selections (/list, /model)
var (
	ModelList   []string
	HistoryList []string
)

// Alternative mit Metadaten:
var (
	ModelMap   = make(map[int]string)
	HistoryMap = make(map[int]string)
)

// z.B. vom /list Command befüllbar
func SetHistoryList(list []string) {
	HistoryList = list
	HistoryMap = make(map[int]string)
	for i, name := range list {
		HistoryMap[i+1] = name
	}
}

func GetHistoryByIndex(i int) (string, bool) {
	val, ok := HistoryMap[i]
	return val, ok
}
