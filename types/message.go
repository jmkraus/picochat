package types

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatHistory struct {
	Messages []Message
}

func NewHistory(systemPrompt string) *ChatHistory {
	return &ChatHistory{
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
		},
	}
}

func (h *ChatHistory) Add(role, content string) {
	h.Messages = append(h.Messages, Message{Role: role, Content: content})
}

func (h *ChatHistory) Get() []Message {
	return h.Messages
}

func (h *ChatHistory) Replace(newMessages []Message) {
	h.Messages = newMessages
}
