package utils

// Indexed list
var (
	ModelsMap  = make(map[int]string)
	HistoryMap = make(map[int]string)
)

// Filled by /list command
//
// Parameters:
//
//	list ([]string) - List of stored history sessions
//
// Returns:
//
//	none
func setHistoryList(list []string) {
	HistoryMap = make(map[int]string)
	for i, name := range list {
		HistoryMap[i+1] = name
	}
}

// GetHistoryByIndex retrieves a history session by its index
//
// Parameters:
//
//	i (int) - Index of the history session
//
// Returns:
//
//	string - session name
//	bool   - boolean indicating success
func GetHistoryByIndex(i int) (string, bool) {
	val, ok := HistoryMap[i]
	return val, ok
}

// Filled by /models command
//
// Parameters:
//
//	list ([]string) - List of available Models
//
// Returns:
//
//	none
func setModelsList(list []string) {
	ModelsMap = make(map[int]string)
	for i, name := range list {
		ModelsMap[i+1] = name
	}
}

// GetModelsByIndex retrieves a model by its index
//
// Parameters:
//
//	i (int) - Index of the model
//
// Returns:
//
//	string - model name
//	bool   - boolean indicating success
func GetModelsByIndex(i int) (string, bool) {
	val, ok := ModelsMap[i]
	return val, ok
}
