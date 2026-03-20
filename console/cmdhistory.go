package console

var cmdHistory = &commandHistory{}

// AddCommand adds a command to the history.
//
// Parameters:
//
//	cmd (string) - the command to add
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
//	cmd (string) - the command to add
//
// Returns:
//
//	none
func (h *commandHistory) add(cmd string) {
	if cmd == "" {
		return
	}
	// Optional: avoid duplicate entries
	if h.len() == 0 || h.entries[h.len()-1] != cmd {
		h.entries = append(h.entries, cmd)
	}
	h.index = h.len()
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
func (h *commandHistory) prev() string {
	if h.len() == 0 {
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
func (h *commandHistory) next() string {
	if h.len() == 0 {
		return ""
	}
	if h.index < h.len()-1 {
		h.index++
		return h.entries[h.index]
	}
	h.index = h.len()
	return ""
}

// len returns the length of the command history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	int - length
func (h *commandHistory) len() int {
	return len(h.entries)
}
