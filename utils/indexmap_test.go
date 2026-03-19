package utils

import "testing"

func cloneIndexMap(src map[int]string) map[int]string {
	dst := make(map[int]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func TestSetHistoryListAndGetHistoryByIndex(t *testing.T) {
	prevHistory := cloneIndexMap(HistoryMap)
	t.Cleanup(func() {
		HistoryMap = prevHistory
	})

	setHistoryList([]string{"session-a.chat", "session-b.chat"})

	v1, ok := GetHistoryByIndex(1)
	if !ok || v1 != "session-a.chat" {
		t.Fatalf("GetHistoryByIndex(1) = (%q, %v), want (session-a.chat, true)", v1, ok)
	}

	v2, ok := GetHistoryByIndex(2)
	if !ok || v2 != "session-b.chat" {
		t.Fatalf("GetHistoryByIndex(2) = (%q, %v), want (session-b.chat, true)", v2, ok)
	}

	_, ok = GetHistoryByIndex(3)
	if ok {
		t.Fatalf("GetHistoryByIndex(3) should not exist")
	}
}

func TestSetModelsListAndGetModelsByIndex(t *testing.T) {
	prevModels := cloneIndexMap(ModelsMap)
	t.Cleanup(func() {
		ModelsMap = prevModels
	})

	setModelsList([]string{"model-a", "model-b"})

	v1, ok := GetModelsByIndex(1)
	if !ok || v1 != "model-a" {
		t.Fatalf("GetModelsByIndex(1) = (%q, %v), want (model-a, true)", v1, ok)
	}

	v2, ok := GetModelsByIndex(2)
	if !ok || v2 != "model-b" {
		t.Fatalf("GetModelsByIndex(2) = (%q, %v), want (model-b, true)", v2, ok)
	}

	_, ok = GetModelsByIndex(3)
	if ok {
		t.Fatalf("GetModelsByIndex(3) should not exist")
	}
}

func TestSetModelsList_ReplacesPreviousValues(t *testing.T) {
	prevModels := cloneIndexMap(ModelsMap)
	t.Cleanup(func() {
		ModelsMap = prevModels
	})

	setModelsList([]string{"old-a", "old-b", "old-c"})
	setModelsList([]string{"new-a"})

	v1, ok := GetModelsByIndex(1)
	if !ok || v1 != "new-a" {
		t.Fatalf("GetModelsByIndex(1) = (%q, %v), want (new-a, true)", v1, ok)
	}

	_, ok = GetModelsByIndex(2)
	if ok {
		t.Fatalf("GetModelsByIndex(2) should be cleared after replacing list")
	}
}
