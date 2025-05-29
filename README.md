# Pico Chat - a CLI Chat Client

## Purpose
Unlike similar tools like ollama, the Mac-only App "[Pico AI Homelab](https://picogpt.app/)" doesn't come with a dedicated CLI interface.

This tool fills the gap and has some additional tricks up its sleeve.

## Installation

## Usage

### Command line args

| ARG      | DESCRIPTION                    |
| -------- | ------------------------------ |
| -config  | Loads a configuration file     |
| -history | Loads the specific session     |
| -version | Shows version number and quits |

### Configuration files

picochat expects a configuration file. If no dedicated name is specified, it looks for a file `config.toml`

The lookup for the configuration file is in the following order:

 1. Full path is given with `-config` arg
 2. Environment variable `CONFIG_PATH` is set
 3. Environment variable  `XDG_CONFIG_HOME` is set (looking for /picochat)
 4. User Home Dir (looking for .config/picochat)
 5. Same folder where the executable is placed

The history files (see below) are invariably stored in a subfolder of the picochat config folder, e.g. `.config/picochat/history`.

Currently picochat supports four values in the config file:

 * Link to the core api endpoint (default usually `http://localhost:11434/api`)
 * Model name (must be already downloaded!)
 * Context size (must be a value between 5 and 100) - this is the number of messages!
 * System prompt ("Persona") where a specific skill or background can be specified


### Commands

| CMD      | DESCRIPTION |
| -------- | ------------------------------------------------- |
| /bye     | Exit the chat |
| /done    | Terminate the input |
| /save    | Save current chat history to a file |
| /load    | Load chat history from a file |
| /list    | List available saved history files |
| /models  | List downloaded models |
| /show    | Show number of messages in history |
| /clear   | Clear history and reinitialize with system prompt |
| /help    | Show available commands |

Some commands can have an argument:

#### \load `<filename>`

Without a filename, an input line shows up, where the name can be entered. If the input is omitted (only __ENTER__), then the load process is cancelled.

Filename is sufficient since the path is invariably set (see above). Suffix can be omitted, it is always `.chat`.

#### \save `<filename>`

Without a filename, the file is stored with a timestamp as filename, e.g. `2025-05-11_20-26-32.chat`.


#### \copy

This command copies the full last answer into the clipboard. However, it removes the `<think>` section from reasoning models.

If `\copy code` is entered, the first occurrence of a codeblock between ` ``` ` will be copied to the clipboard instead, skipping all descriptive text.


### Personas

Picochat allows basic persona handling: Store different configuration files in your config-path, e.g. "generic.toml" or "developer.toml" with specific system prompts.

Then load this configuration with a shortcut, e.g. `picochat -config @developer`. You can skip path (covered by "@") and suffix ".toml". Then picochat starts with the specified configuration file.
