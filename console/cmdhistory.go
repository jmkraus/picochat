package console

var cmdHistory = &CommandHistory{}

func AddCommand(cmd string) {
	cmdHistory.add(cmd)
}

func PrevCommand() string {
	return cmdHistory.prev()
}

func NextCommand() string {
	return cmdHistory.next()
}

func (h *CommandHistory) add(cmd string) {
	if cmd == "" {
		return
	}
	// Optional: avoid duplicate entries
	if len(h.entries) == 0 || h.entries[len(h.entries)-1] != cmd {
		h.entries = append(h.entries, cmd)
	}
	h.index = len(h.entries)
}

func (h *CommandHistory) prev() string {
	if len(h.entries) == 0 {
		return ""
	}
	if h.index > 0 {
		h.index--
	}
	return h.entries[h.index]
}

func (h *CommandHistory) next() string {
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
