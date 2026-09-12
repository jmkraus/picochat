package messages

import (
	"fmt"
	"sync"
)

// MaxSessions is the maximum number of in-memory chat sessions.
const MaxSessions = 8

// SessionManager manages ordered chat histories and the currently active chat.
type SessionManager struct {
	mu        sync.RWMutex
	histories []*ChatHistory
	active    int
}

// NewSessionManager creates an empty session manager.
//
// Parameters:
//
//	none
//
// Returns:
//
//	*SessionManager - an empty session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{active: -1}
}

// Add registers an existing chat history and makes it active.
//
// Parameters:
//
//	history (*ChatHistory) - chat history to register
//
// Returns:
//
//	int   - the history's new zero-based slot index
//	error - error if the history is nil or the session limit is reached
func (sm *SessionManager) Add(history *ChatHistory) (int, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if history == nil {
		return -1, fmt.Errorf("cannot add nil session")
	}
	if len(sm.histories) >= MaxSessions {
		return -1, fmt.Errorf("maximum of %d sessions reached", MaxSessions)
	}

	sm.histories = append(sm.histories, history)
	sm.active = len(sm.histories) - 1

	return sm.active, nil
}

// Create adds a new chat history and makes it active.
//
// Parameters:
//
//	systemPrompt (string) - system prompt for the new history
//	maxContext  (int)    - maximum context size for the new history
//
// Returns:
//
//	int          - the new chat's zero-based slot index
//	*ChatHistory - the new chat history
//	error        - error if the session limit is reached
func (sm *SessionManager) Create(systemPrompt string, maxContext int) (int, *ChatHistory, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if len(sm.histories) >= MaxSessions {
		return -1, nil, fmt.Errorf("maximum of %d sessions reached", MaxSessions)
	}

	history := NewHistory(systemPrompt, maxContext)
	sm.histories = append(sm.histories, history)
	sm.active = len(sm.histories) - 1

	return sm.active, history, nil
}

// Active returns the currently active chat history.
//
// Parameters:
//
//	none
//
// Returns:
//
//	*ChatHistory - the active chat history
//	error        - error if there is no active session
func (sm *SessionManager) Active() (*ChatHistory, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if sm.active < 0 || sm.active >= len(sm.histories) {
		return nil, fmt.Errorf("no active session")
	}

	return sm.histories[sm.active], nil
}

// ActiveIndex returns the zero-based index of the active chat.
//
// Parameters:
//
//	none
//
// Returns:
//
//	int - the active chat's zero-based index, or -1 if no chat is active
func (sm *SessionManager) ActiveIndex() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.active
}

// Get returns the chat history at the given zero-based slot index.
//
// Parameters:
//
//	index (int) - zero-based slot index
//
// Returns:
//
//	*ChatHistory - the chat history at the requested index
//	error        - error if the index is out of bounds
func (sm *SessionManager) Get(index int) (*ChatHistory, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	if index < 0 || index >= len(sm.histories) {
		return nil, fmt.Errorf("session index %d out of bounds", index+1)
	}

	return sm.histories[index], nil
}

// List returns the chat histories in stable slot order.
//
// Parameters:
//
//	none
//
// Returns:
//
//	[]*ChatHistory - a copy of the chat history slots
func (sm *SessionManager) List() []*ChatHistory {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return append([]*ChatHistory(nil), sm.histories...)
}

// Switch makes the chat at the given zero-based slot index active.
//
// Parameters:
//
//	index (int) - zero-based slot index
//
// Returns:
//
//	error - error if the index is out of bounds
func (sm *SessionManager) Switch(index int) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if index < 0 || index >= len(sm.histories) {
		return fmt.Errorf("session index %d out of bounds", index+1)
	}

	sm.active = index
	return nil
}

// Copy creates a new active chat from the current chat.
// If messageIndex is nil, the complete current history is copied. Otherwise,
// messages from index 0 through messageIndex are copied, inclusive.
//
// Parameters:
//
//	messageIndex (*int) - optional inclusive upper message index to copy
//
// Returns:
//
//	int          - the new chat's zero-based slot index
//	*ChatHistory - the copied chat history
//	error        - error if there is no active session, the session limit is reached, or the message index is invalid
func (sm *SessionManager) Copy(messageIndex *int) (int, *ChatHistory, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.active < 0 || sm.active >= len(sm.histories) {
		return -1, nil, fmt.Errorf("no active session")
	}
	if len(sm.histories) >= MaxSessions {
		return -1, nil, fmt.Errorf("maximum of %d sessions reached", MaxSessions)
	}

	source := sm.histories[sm.active]
	end := len(source.Messages)
	if messageIndex != nil {
		if *messageIndex < 0 || *messageIndex >= end {
			return -1, nil, fmt.Errorf("message index %d out of bounds", *messageIndex)
		}
		end = *messageIndex + 1
	}

	copied := &ChatHistory{
		Messages:          cloneMessages(source.Messages[:end]),
		MaxContext:        source.MaxContext,
		MaxContextReached: end >= source.MaxContext,
	}
	sm.histories = append(sm.histories, copied)
	sm.active = len(sm.histories) - 1

	return sm.active, copied, nil
}

// cloneMessages creates a copy of the supplied messages and their image lists.
//
// Parameters:
//
//	source ([]Message) - messages to clone
//
// Returns:
//
//	[]Message - cloned messages
func cloneMessages(source []Message) []Message {
	cloned := make([]Message, len(source))
	copy(cloned, source)

	for i := range cloned {
		cloned[i].Images = append([]string(nil), source[i].Images...)
	}

	return cloned
}
