# Pico Chat - a CLI Chat Client

## Purpose
Unlike similar tools like ollama, the Mac-only App "[Pico AI Homelab](https://picogpt.app/)" doesn't come with a dedicated CLI interface.

This tool fills the gap and has some additional tricks up its sleeve.

## Installation

## Usage

### Command line args

| ARG      | DESCRIPTION                  |
| -------- | ---------------------------- |
| -config  | Loads a configuration file   |
| -history | Loads the specific session   |
| -version | Shows version number and quits |


### Commands

| CMD    | DESCRIPTION |
| ------ | ------------------------------------------------- |
| /bye   | Exit the chat |
| /done  | Terminate the input |
| /save  | Save current chat history to a file |
| /load  | Load chat history from a file |
| /list  | List available saved history files |
| /show  | Show number of messages in history |
| /clear | Clear history and reinitialize with system prompt |
| /help  | Show available commands |

### Personas

Picochat allows basic persona handling: Store different configuration files in your config-path, e.g. "generic.toml" or "developer.toml" with specific system prompts.

Then load this configuration with a shortcut, e.g. `picochat -config @developer`. You can skip path (covered by "@") and suffix ".toml". Then picochat starts with the specified configuration file.
