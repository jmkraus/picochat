# Pico Chat - a CLI Chat Client

## Purpose
The Mac‑only app [Pico AI Server](https://picogpt.app/), unlike tools such as Ollama, does not include a dedicated CLI client.

This tool fills the gap and has some additional tricks up its sleeve.

## Installation

The released binary is unsigned. To use it on a Mac, there are two options:

 1. Build picochat from the source code
 2. Remove the quarantine tag from the binary

### Build picochat

 Add the required libraries to your local Go setup.

 ```
  go get github.com/atotto/clipboard
  go get github.com/BurntSushi/toml
  go mod tidy
  go test ./...
  go build
 ```

### Remove the quarantine flag

Navigate to the target folder where the binary was placed and enter the following:

`sudo xattr -rd com.apple.quarantine ./picochat`

Enter the administrator password to confirm. Then the binary works.

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

 1. Full path is given with the `-config` argument.
 2. Environment variable `CONFIG_PATH` is set.
 3. Environment variable  `XDG_CONFIG_HOME` is set (searches for /picochat).
 4. User home directory (searches for .config/picochat).
 5. Same folder where the executable is placed.

The history files (see below) are invariably stored in a subfolder of the picochat config folder, e.g. `.config/picochat/history`.

Picochat currently supports four configurable values in the config file:

 * Link to the core **API** endpoint (default usually `http://localhost:11434/api`).
 * Model name (must be already downloaded!)
 * Context size (must be a value between 5 and 100) - this is the number of messages!
 * System prompt ("Persona") where a specific skill or background can be specified.

Since Pico AI currently doesn't report token counts, it is difficult to calculate a proper context size. Maybe this changes in the future, but for now the context size is limited by the number of total messages, where the oldest ones are dropped when the limit is reached.

### Commands

| CMD        | DESCRIPTION |
| ---------- | ------------------------------------------------- |
| /done, /// | Terminate the input and send |
| /cancel    | Cancel multi-line input and return to prompt |
| /copy, /c  | Copy the last answer to clipboard |
| /paste, /v | Get clipboard content as user input and send |
| /info      | Show system information |
| /message   | Output last message again (e.g., after load) |
| /load      | Load chat history from a file |
| /save      | Save current chat history to a file |
| /list      | List available saved history files |
| /models    | List (and switch) downloaded models |
| /clear     | Clear history and reinitialize with system prompt |
| /set       | Set session variables (key=value) |
| /retry     | Sends chat history again, but without last answer |
| /bye       | Exit the chat |
| /help, /?  | Show available commands |

Some commands can have an argument:

#### /load `<filename>`

Without a filename, an input line appears, where the name can be entered. If the input is omitted (only _ENTER_), then the load process is cancelled.

The filename is sufficient because the path is invariably set (see above). Suffix can be omitted, it is always `.chat`.

#### /save `<filename>`

Without a filename, the file is stored with a timestamp as filename, e.g. `2025-05-11_20-26-32.chat`.


#### /copy

This command copies the full last answer into the clipboard. However, it removes the `<think>` section from reasoning models. If the reasoning should be retained, then `/copy think` can be used instead.

If `/copy code` is entered, the first occurrence of a codeblock between ` ``` ` will be copied to the clipboard instead, skipping all descriptive text.

#### /models `<index>`

Without an argument, this command lists the available models of the **LLM** server. If the list of models has been requested at least once, then it's possible to switch to another model by using the index of the list, e.g. `/models 3`.

### Multiline input
Picochat utilizes standard input, unlike Ollama’s raw input method. This approach was preferred for its simpler implementation and to avoid issues with multiline input enclosed in triple quotes.

Users can enter or paste as much text as needed for a prompt. Input is terminated by entering `/done` or `///` on a new line. Single lines can also be terminated by `///` in the same line.

Multiline input can be cancelled at any time by entering `/cancel`, returning the user to the prompt without sending the input.

### Personas

Picochat allows basic persona handling: Store different configuration files in the config path, e.g., `generic.toml` or `developer.toml`, each with specific system prompts.

This configuration can be loaded using a shortcut, e.g., `picochat -config @developer`. The path and `.toml` suffix can be omitted because they are implied by the '@' symbol. Picochat then starts with the specified configuration file.


## Acknowledgements

Big shoutout to the providers of the libraries used in this project:

 * [The TOML library by BurntSushi](https://github.com/BurntSushi/toml).
 * [Atotto's clipboard library for Go](https://github.com/atotto/clipboard).


## License

This project is distributed under the MIT license.

JMK 2025
