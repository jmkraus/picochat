package console

func (h *CommandHistory) Add(cmd string) {
	if cmd == "" {
		return
	}
	// Optional: avoid duplicate entries
	if len(h.entries) == 0 || h.entries[len(h.entries)-1] != cmd {
		h.entries = append(h.entries, cmd)
	}
	h.index = len(h.entries)
}

func (h *CommandHistory) Prev() string {
	if len(h.entries) == 0 {
		return ""
	}
	if h.index > 0 {
		h.index--
	}
	return h.entries[h.index]
}

func (h *CommandHistory) Next() string {
	if len(h.entries) == 0 {
		return ""
	}
	if h.index < len(h.entries)-1 {
		h.index++
		return h.entries[h.index]
	}
	h.index = len(h.entries)
	return ""
}
