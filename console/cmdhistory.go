package console

var cmdHistory = &CommandHistory{}

// AddCommand adds a command to the history.
//
// Parameters:
//
//	cmd string - the command to add
//
// Returns:
//
//	none
func AddCommand(cmd string) {
	cmdHistory.add(cmd)
}

// PrevCommand returns the previous command in the history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the previous command
func PrevCommand() string {
	return cmdHistory.prev()
}

// NextCommand returns the next command in the history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the next command
func NextCommand() string {
	return cmdHistory.next()
}

// add adds a command to the history.
//
// Parameters:
//
//	cmd string - the command to add
//
// Returns:
//
//	none
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

// prev returns the previous command in the history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the previous command
func (h *CommandHistory) prev() string {
	if len(h.entries) == 0 {
		return ""
	}
	if h.index > 0 {
		h.index--
	}
	return h.entries[h.index]
}

// next returns the next command in the history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	string - the next command
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
