# Pico Chat - a CLI Chat Client

## Purpose
Unlike similar tools like ollama, the Mac-only app "[Pico AI Homelab](https://picogpt.app/)" doesn't come with a dedicated CLI client.

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

picochat expects a configuration file. If no specific name is given, it looks for a file `config.toml`

The lookup for the configuration file is in the following order:

 1. Full path is given with `-config` arg
 2. Environment variable `CONFIG_PATH` is set
 3. Environment variable  `XDG_CONFIG_HOME` is set (looking for /picochat)
 4. User Home Dir (looking for .config/picochat)
 5. Same folder where the executable is placed

The history files (see below) are invariably stored in a subfolder of the picochat config folder, e.g. `.config/picochat/history`.

Picochat currently supports four configurable values in the config file:

 * Link to the core api endpoint (default usually `http://localhost:11434/api`)
 * Model name (must be already downloaded!)
 * Context size (must be a value between 5 and 100) - this is the number of messages!
 * System prompt ("Persona") where a specific skill or background can be specified

Since Pico AI currently doesn't report token counts, it is difficult to calculate a proper context size. Maybe this changes in the future, but for now the context size is limited by the number of total messages, where the oldest ones are dropped when the limit is reached.

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
| /set     | Set session variables |
| /clear   | Clear history and reinitialize with system prompt |
| /help    | Show available commands |

Some commands can have an argument:

#### \load `<filename>`

Without a filename, an input line shows up, where the name can be entered. If the input is omitted (only _ENTER_), then the load process is cancelled.

Filename is sufficient since the path is invariably set (see above). Suffix can be omitted, it is always `.chat`.

#### \save `<filename>`

Without a filename, the file is stored with a timestamp as filename, e.g. `2025-05-11_20-26-32.chat`.


#### \copy

This command copies the full last answer into the clipboard. However, it removes the `<think>` section from reasoning models.

If `\copy code` is entered, the first occurrence of a codeblock between ` ``` ` will be copied to the clipboard instead, skipping all descriptive text.


### Multiline input
Unlike Ollama, Picochat uses standard input instead of raw input. Besides the simpler implementation, I was also uncomfortable with the approach of multiline input enclosed by """.

Sometimes I decided to enter more text but didn't start with """ so that I had to start from scratch. Therefore I considered a stop command as better solution for me. When entering a user prompt, as much text as desired can be entered or pasted. Press _ENTER_ for newline, then either enter `\done` or `\\\` followed by _ENTER_. Either will terminate the input and send it to the AI server.


### Personas

Picochat allows basic persona handling: Store different configuration files in your config-path, e.g. "generic.toml" or "developer.toml" with specific system prompts.

You can load this configuration using a shortcut, such as `picochat -config @developer`. The path and ".toml" suffix can be omitted since they are implied by the '@' symbol. Then picochat starts with the specified configuration file.


## Acknowledgements

Big shoutout to the providers of the libraries I used for this project:

 * [The TOML library by BurntSushi](https://github.com/BurntSushi/toml)
 * [Atotto's clipboard library for Go](https://github.com/atotto/clipboard)


## License

This project is distributed under the MIT license.

JMK 2025
