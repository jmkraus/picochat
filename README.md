# PicoChat - a CLI Chat Client

## Purpose
The Mac‑only app [Pico AI Server](https://picogpt.app/), unlike tools such as Ollama, does not provide a dedicated CLI client.

This tool fills that gap and offers some additional features.

## Installation

The released binary is neither notarized nor signed. To use it on a Mac, there are two options:

 1. Build PicoChat from the sources
 2. Remove the quarantine flag from the binary

### Build PicoChat

 Add the required libraries to your local Go setup.

 ```text
  go get github.com/atotto/clipboard
  go get github.com/BurntSushi/toml
  go mod tidy
  go test ./...
  go build
 ```

### Remove the quarantine flag

Navigate to the folder containing the binary and enter the following:

`sudo xattr -rd com.apple.quarantine ./picochat`

Enter the administrator password to confirm. After that, the binary can be started without a warning.

## Usage

### Entering a prompt
PicoChat uses standard input, unlike Ollama’s raw input method. This approach was chosen for its simpler implementation and to avoid issues with multiline input enclosed in triple quotes.

Users can enter or paste as much text as needed for a prompt. Input ends when typing `/done` or `///` on a new line. Single-line input can also be ended by `///` on the same line.

Multiline input can be canceled at any time by entering `/cancel`, which returns the user to the prompt without sending the input. Here are a few examples (↵ indicates a RETURN):

#### Single line
```text
Hello are you there? ///↵
```

```text
Hello are you there?↵
/done↵
```

#### Multi line
```text
How can I improve the following GoLang code?↵
package main↵
import "fmt"↵
func main() {↵
    fmt.Println("Hello, World!")↵
}↵
/done↵
```

### Command line args

| ARG      | DESCRIPTION                    |
| -------- | ------------------------------ |
| -config  | Loads a configuration file     |
| -history | Loads the specific session     |
| -version | Shows version number and quits |

### Configuration files

PicoChat expects a configuration file. If no specific name is given, it looks for a file named `config.toml`

The lookup for the configuration file is in the following order:

 1. Full path is given with the `-config` argument.
 2. Environment variable `CONFIG_PATH` is set.
 3. Environment variable  `XDG_CONFIG_HOME` is set (searches for /picochat).
 4. User home directory (searches for .config/picochat).
 5. Same folder where the executable is placed.

History files (see below) are always stored in a subdirectory of the PicoChat config folder, e.g. `.config/picochat/history`.

PicoChat currently supports the following configurable values in the config file:

 * link to the core **API** endpoint (default usually `http://localhost:11434/api`).
 * model name (must be already downloaded!)
 * temperature
 * top_p
 * context size (must be a value between 5 and 100) - this is the number of messages!
 * system prompt ("Persona") where a specific skill or background can be specified.

Since Pico AI currently does not report token counts, it is difficult to calculate a proper context size. This may change in the future, but for now the context size is limited by the total number of messages, with the oldest ones being dropped when the limit is reached.

### Commands

| CMD        | DESCRIPTION |
| ---------- | ------------------------------------------------- |
| /done, /// | Terminate the input and send |
| /cancel    | Cancel multi-line input and return to prompt |
| /copy, /c  | Copy the last answer to clipboard |
| /paste, /v | Get clipboard content as user input and send |
| /info      | Show system information |
| /message   | Show last message again (e.g., after load) |
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

Without a filename, an input line appears, where the name can be entered. If the input is omitted (just pressing _ENTER_), the load process is canceled.

The filename alone is sufficient because the path is predefined (see above). The suffix can be omitted, as it defaults to `.chat`.

#### /save `<filename>`

Without a filename, the file is stored with a timestamp as filename, e.g. `2025-05-11_20-26-32.chat`.


#### /copy

This command copies the entire last answer to the clipboard but removes the `<think>` section from reasoning models. If the reasoning should be retained, then `/copy think` can be used instead.

If `/copy code` is entered, the first occurrence of a codeblock between ` ``` ` will be copied to the clipboard instead, skipping all descriptive text. If no codeblock is found, then the copy command is canceled.

#### /models `<index>`

Without an argument, this command lists the available models of the **LLM** server. If the list of models has been requested at least once, then it's possible to switch to another model by using the index of the list, e.g. `/models 3`.

#### /set `<key=value>`

Without an argument, the command lists the available parameters and show their current values.

With the optional argument the parameter values can be changed for the current session, e.g., `/set top_p=0.3`.

These values are not persistent and cannot be saved. For permanent changes, edit the entries in `config.toml`.

#### /message `<role>`

Without the argument, the last entry of the chat history (usually an assistant answer) will be shown.

With one of the possible roles (system, user, assistant), the specific last entry of the chat history can be chosen.

For example, `/message user` displays the last user question again.

### Personas

PicoChat supports basic persona handling: it is possible to store different configuration files in the config path, e.g., `generic.toml` or `developer.toml`, each with specific system prompts.

This configuration can be loaded using a shortcut, e.g., `picochat -config @developer`. The path and `.toml` suffix can be omitted because they are implied by the '@' symbol. PicoChat then starts with the specified configuration file.


## Acknowledgements

Special thanks to the developers of the libraries used in this project:

 * [BurntSushi's TOML library](https://github.com/BurntSushi/toml).
 * [Atotto's Go clipboard library](https://github.com/atotto/clipboard).


## License

This project is licensed under the MIT license.

JMK 2025
