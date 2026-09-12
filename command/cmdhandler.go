package command

import (
	"fmt"
	"picochat/clipb"
	"picochat/messages"
	"picochat/output"
	"regexp"
	"strings"
)

// handleChatCommand processes the /chat command based on the given arguments.
//
// Parameters:
//
//	args ([]string) - arguments passed to the /chat command
//	history (*messages.ChatHistory) - chat history used as source for message lookup
//	sessions (*messages.SessionManager) - the instance of the chat history session manager
//
// Returns:
//
//	CommandResult - a struct containing status flags and messages of the command processing
func handleChatCommand(args []string, history *messages.ChatHistory, sessions *messages.SessionManager) CommandResult {
	if sessions == nil {
		return CommandResult{Error: fmt.Errorf("session manager is unavailable")}
	}

	if len(args) == 0 || args[0] == "" {
		return listChats(sessions)
	}

	switch strings.ToLower(args[0]) {
	case "new":
		if len(args) != 1 {
			return CommandResult{Error: fmt.Errorf("/chat new does not accept arguments")}
		}
		if history == nil {
			return CommandResult{Error: fmt.Errorf("current history is unavailable")}
		}

		system, found := history.GetLastRole(messages.RoleSystem)
		if !found {
			return CommandResult{Error: fmt.Errorf("current history has no system prompt")}
		}
		index, _, err := sessions.Create(system.Content, history.MaxCtx())
		if err != nil {
			return CommandResult{Error: fmt.Errorf("create chat failed: %w", err)}
		}
		return CommandResult{
			Info:           fmt.Sprintf("New chat session %d created.", index+1),
			SessionChanged: true,
		}

	case "copy":
		if len(args) > 2 {
			return CommandResult{Error: fmt.Errorf("usage: /chat copy [message-index]")}
		}

		var messageIndex *int
		if len(args) == 2 {
			index, err := parseIndex(args[1])
			if err != nil {
				return CommandResult{Error: fmt.Errorf("copy chat failed: %w", err)}
			}
			messageIndex = &index
		}

		index, _, err := sessions.Copy(messageIndex)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("copy chat failed: %w", err)}
		}
		return CommandResult{
			Info:           fmt.Sprintf("Chat copied to session %d.", index+1),
			SessionChanged: true,
		}

	default:
		if len(args) != 1 {
			return CommandResult{Error: fmt.Errorf("usage: /chat, /chat new, /chat <number>, or /chat copy [message-index]")}
		}

		index, err := parseIndex(args[0])
		if err != nil || index < 1 {
			return CommandResult{Error: fmt.Errorf("chat number must be a positive integer")}
		}
		if err := sessions.Switch(index - 1); err != nil {
			return CommandResult{Error: fmt.Errorf("switch chat session failed: %w", err)}
		}
		return CommandResult{
			Info:           fmt.Sprintf("Chat switched to session %d.", index),
			SessionChanged: true,
		}
	}
}

func listChats(sessions *messages.SessionManager) CommandResult {
	chats := sessions.List()
	if len(chats) == 0 {
		return CommandResult{Warn: "No chats available."}
	}

	var lines []string
	for i, history := range chats {
		marker := " "
		if i == sessions.ActiveIndex() {
			marker = "*"
		}
		lines = append(lines, fmt.Sprintf("%s %d: %d messages", marker, i+1, history.Len()))
	}
	return CommandResult{Output: strings.Join(lines, "\n")}
}

// handleCopyCommand processes the /copy command based on its argument.
//
// Parameters:
//
//	args (string) - argument passed to the /copy command
//	history (*messages.ChatHistory) - chat history used as source for message lookup
//
// Returns:
//
//	CommandResult - a struct containing the outcome of the command
func handleCopyCommand(args string, history *messages.ChatHistory) CommandResult {
	nothing := "Nothing to copy."
	if indexStr, ok := strings.CutPrefix(args, "#"); ok {
		index, err := parseIndex(indexStr)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("copy message failed: %w", err)}
		}
		msg, err := history.GetByIndex(index)
		if err != nil {
			return CommandResult{Error: fmt.Errorf("copy message failed: %w", err)}
		}
		return copyTextToClipboard(msg.Content, fmt.Sprintf("Message #%d copied to clipboard.", index))
	}

	if args == "" {
		args = messages.RoleAssistant
	}

	switch args {
	case messages.RoleAssistant, messages.RoleUser, messages.RoleSystem:
		lastMessage, found := history.GetLastRole(args)
		if !found || lastMessage.Content == "" {
			return CommandResult{Warn: nothing}
		}
		return copyTextToClipboard(lastMessage.Content, fmt.Sprintf("Last %s prompt copied to clipboard.", args))

	case "all":
		conversation := output.FormatConversation(history.Get(), false)
		return copyTextToClipboard(conversation, "Full conversation copied to clipboard.")
	case "think":
		lastMessage, found := history.GetLastRole(messages.RoleAssistant)
		if !found || (lastMessage.Content == "" && lastMessage.Reasoning == "") {
			return CommandResult{Warn: nothing}
		}
		return copyTextToClipboard(encloseThinkingTags(lastMessage.Reasoning)+lastMessage.Content, "Last assistant prompt (with thinking) copied to clipboard.")

	case "code":
		lastMessage, found := history.GetLastRole(messages.RoleAssistant)
		if !found || lastMessage.Content == "" {
			return CommandResult{Warn: nothing}
		}
		codeBlock, found := extractCodeBlock(lastMessage.Content)
		if !found {
			return CommandResult{Warn: nothing}
		}
		return copyTextToClipboard(codeBlock, "First code block copied to clipboard.")

	default:
		return CommandResult{Error: fmt.Errorf("unknown copy argument")}
	}
}

var writeClipboard = clipb.WriteClipboard

func copyTextToClipboard(text, info string) CommandResult {
	if err := writeClipboard(text); err != nil {
		return CommandResult{Error: err}
	}
	return CommandResult{Info: info}
}

// extractCodeBlock extracts the first code block from a string
// formatted with triple backticks.
//
// Parameters:
//
//	s (string) - the input string containing code blocks.
//
// Returns:
//
//	string - the extracted code block content.
//	bool   - true if a code block was found, false otherwise.
func extractCodeBlock(s string) (string, bool) {
	re := regexp.MustCompile("(?s)```\\w*\\n(.*?)```")
	match := re.FindStringSubmatch(s)
	if len(match) >= 2 {
		return match[1], true
	}
	return "", false
}

// encloseThinkingTags adds tags around a given string to
// identify it as the reasoning part of the text
//
// Parameters:
//
//	s (string) - the string to be tagged
//
// Returns:
//
//	string - the tagged text
func encloseThinkingTags(s string) string {
	return fmt.Sprintf("<think>\n%s\n</think>\n\n", s)
}

// handleMessageCommand processes the /message command based on the given argument.
//
// Parameters:
//
//	args ([]string) - arguments passed to the /message command
//	history (*messages.ChatHistory) - chat history used as source for message lookup
//
// Returns:
//
//	CommandResult - a struct containing the outcome of the command
func handleMessageCommand(args []string, history *messages.ChatHistory) CommandResult {
	if idxArg, ok := strings.CutPrefix(args[0], "#"); ok {
		msg, err := getMessageByIndex(idxArg, history)
		if err != nil {
			return CommandResult{Error: err}
		}
		return CommandResult{Output: msg}
	}

	switch args[0] {
	case "":
		msg := history.GetLast().Content
		return CommandResult{Output: msg}
	case "all":
		conversation := output.FormatConversation(history.Get(), true)
		return CommandResult{Output: conversation}
	case messages.RoleAssistant, messages.RoleUser, messages.RoleSystem:
		msg, found := history.GetLastRole(args[0])
		if found {
			return CommandResult{Output: msg.Content}
		}
		return CommandResult{Warn: fmt.Sprintf("No element for role %q found.", args)}
	default:
		return CommandResult{Error: fmt.Errorf("unknown argument")}
	}
}
