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
  go mod download
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

PicoChat uses a raw input mode for entering prompts.
This provides a more natural multiline editing experience than the triple-quote (""") approach used by other applications, while still keeping the implementation simple and terminal-friendly.

Users can type or paste as much text as needed for a prompt.
Input is submitted by pressing Ctrl+D (End-of-File).
This method ensures that pasted text — including blank lines or code blocks — is handled safely without triggering unintended input termination.

Input can be canceled at any time by pressing Esc, which immediately returns the user to the main prompt without sending the text.

Here are a few examples (↵ indicates pressing Enter):

#### Single line
```text
>>> Tell me a joke! [Ctrl]+D
```

#### Multi line
```text
>>> Hello, PicoChat! ↵
How are you today? ↵
↵
[Ctrl]+D        ← Send prompt
```

#### Via stdin pipe
Use PicoChat in scripts via Pipe, e.g.:
```
echo "Write a Haiku about Cheese" | picochat -quiet
```

The `-quiet` argument is optional and suppresses all app messages, so that only the LLM response is displayed.

For Windows turn the command shell into UTF-8 mode first for proper handling of non-Western characters:
```
chcp 65001
```

It's also possible to run a PicoChat command via stdin pipe. This is experimental and has not yet been thoroughly tested and might have unexpected side effects. 
```
echo "/models" | picochat -quiet
```


### Command line args

| Arg      | Description                             |
| -------- | --------------------------------------- |
| -config  | Loads a configuration file              |
| -history | Loads the specific session              |
| -model   | Overrides config setting with new model |
| -image   | Sets a path for an image file           |
| -format  | Defines output format of the response   |
| -quiet   | Suppresses all app messages             |
| -version | Shows version number and quits          |


### Output formats

PicoChat provides output formats other than plain text. This simplifies the output processing in pipelines etc.

Example usage: `picochat -format json`

The following formats are available:

| Format      | Description                                        |
| ----------- | -------------------------------------------------- |
| plain       | (default) Plain text output, can be omitted.       |
| json        | Response formatted as json output.                 |
| json-pretty | Response formatted as json output in pretty print. |
| yaml        | Response formatted as yaml output.                 |

If a wrong format is given, a warning shows up and PicoChat uses plain text as fallback.

**Example output: plain**

`echo "Write a Haiku about Cheese" | picochat -quiet`

```
Golden, sharp, and bold,
Melts upon your waiting tongue,
Dairy bliss unfolds.
```

**Example output: json**

`echo "Write a Haiku about Cheese" | picochat -quiet -format json`

```
{"output":"Golden, sharp, and bold,\nMelts upon your waiting tongue,\nDairy bliss unfolds.","elapsed":"00:02","tokens_per_sec":7.8}
```

**Example output: json-pretty**

`echo "Write a Haiku about Cheese" | picochat -quiet -format json-pretty`

```
{
  "output": "Golden, sharp, and bold,\nMelts upon your waiting tongue,\nDairy bliss unfolds."
  "elapsed": "00:02",
  "tokens_per_sec": 7.8
}
```

**Example output: yaml**

`echo "Write a Haiku about Cheese" | picochat -quiet -format yaml`

```
output: "Golden, sharp, and bold,\nMelts upon your waiting tongue,\nDairy bliss unfolds."
elapsed: "00:02",
tokens_per_sec: 7.8
```

Using formatted output enables pipelines like this:

```
echo "Write a Haiku about Cheese" | picochat -quiet -format json | jq -r '.elapsed'
```


### Reasoning

Pico AI Server and Ollama handle reasoning differently. While Ollama separates reasoning and content output, Pico AI's content also includes reasoning, which is embedded in `<think>` tags.

PicoChat can handle both forms directly. However, there are differences in how they are displayed during output. The output in Ollama is in gray color to distinguish it from the content. Since some models in Pico AI are faulty and do not output an opening `<think>` tag, the reasoning part cannot be correctly recognized. Therefore, in this case, the reasoning remains part of the content and is not colorized.

Internally, reasoning is separated from content for both, Pico AI and Ollama. This means that even if the answer is queried later using `/message`, no reasoning is displayed in Pico AI. However, the behavior of `/copy` is now consistent: only the content is copied in Pico AI as well. If reasoning is also to be included, `/copy think` can be used.

Care should also be taken to ensure that the reasoning flag is set when using reasoning models. Furthermore, reasoning behavior is highly dependent on the model used and can occasionally lead to unexpected behavior and freezes, or simply be ignored by the model.

It should be noted that reasoning is not saved when using `/save`, as it is not considered an official part of the conversation. 


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

| Key           | Type    | Value                                                                              |
| ------------- | ------- | ---------------------------------------------------------------------------------- |
| `URL`         | string  | link to the core **API** endpoint (default usually `http://localhost:11434/api`).  |
| `Model`       | string  | Model name (must be already downloaded)                                            |
| `Context`     | integer | Context size (must be a value between 3 and 100) - this is the number of messages! |
| `Temperature` | float   | Model temperature                                                                  |
| `TopP`        | float   | Model top_p                                                                        |
| `Prompt`      | string  | System prompt ("Persona") where a specific skill or background can be specified.   |
| `Quiet`       | bool    | Suppresses all messages (except for errors)                                        |
| `Reasoning`   | bool    | Enables or disables reasoning                                                      |

Since Pico AI currently does not report token counts, it is difficult to calculate a proper context size. This may change in the future, but for now the context size is limited by the total number of messages, with the oldest ones being dropped when the limit is reached.

### Commands

| CMD             | DESCRIPTION |
| --------------- | ------------------------------------------------- |
| [Ctrl]+D        | Submit multiline input (EOF)                      |
| [Esc], [Ctrl]+C | Cancel multiline input and return to prompt       |
| [Up]/[Down]     | Browse prompt command history                     |
| /copy, /c       | Copy the last answer to clipboard                 |
| /paste, /v      | Get clipboard content as user input and send      |
| /info           | Show system information                           |
| /message        | Show last message again (e.g., after load)        |
| /load           | Load chat history from a file                     |
| /save           | Save current chat history to a file               |
| /list           | List available saved history files                |
| /models         | List (and switch) downloaded models               |
| /clear          | Clear session context                             |
| /set            | Set session variables (key=value)                 |
| /image          | Set image file path                               |
| /retry          | Sends chat history again, but without last answer |
| /bye            | Quit PicoChat                                     |
| /help, /?       | Show available commands                           |

Some commands can have an argument:

#### /load `<filename>`

Without a filename, an input line appears, where the name can be entered. If the input is omitted (just pressing _ENTER_), the load process is canceled.

The filename alone is sufficient because the path is predefined (see above). The suffix can be omitted, as it defaults to `.chat`.

If the command `/list` has been executed before, it is also possible to load a session by index, e.g.: `/load #3`. The hash mark indicates that an index is given rather than a filename.

#### /save `<filename>`

Without a filename, the file is stored with a timestamp as filename, e.g. `2025-05-11_20-26-32.chat`.


#### /copy

This command copies the entire last answer to the clipboard but removes the `<think>` section from reasoning models. If the reasoning should be retained, then `/copy think` can be used instead.

If `/copy code` is entered, the first occurrence of a codeblock between ` ``` ` will be copied to the clipboard instead, skipping all descriptive text. If no codeblock is found, then the copy command is canceled.

With the argument `/copy user` the last user prompt is put into the clipboard.

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


### Image processing

PicoChat allows for simple use of images. An image can either be passed as an argument:

```
picochat -image ./imgfile.jpg
```

or by using the `/image` command:

```
>>> /image ./imgfile.jpg
```

It's possible to use the tilde "~" as abbreviation for the user home directory. With the next user prompt this image will be processed and then discarded. When passing an image path to PicoChat, the existence of that file will be checked (and a warning message is shown if it doesn't exist), but not if it is a valid image.

In combination with stdin pipe and `-model` argument a simple commandline image analytics is possible:

```
echo "What's on the image?" | picochat -model Qwen3-VL-8B-Instruct-4bit -image ./imgfile.jpg -quiet
```

When saving the chat history, the user prompt contains the image as base64 encoded data.


### Personas

PicoChat supports basic persona handling: it is possible to store different configuration files in the config path, e.g., `generic.toml` or `developer.toml`, each with specific system prompts.

This configuration can be loaded using a shortcut, e.g., `picochat -config @developer`. The path and `.toml` suffix can be omitted because they are implied by the '@' symbol. PicoChat then starts with the specified configuration file.


### Known issues

While everything should work as expected on macOS, PicoChat is completely untested on Linux and only minimally tested on Windows.

* On Windows, the Esc key must be pressed twice to cancel a prompt input. Alternatively, Ctrl + C can be used.
* The command history flickers on Windows when switching between entries using Up + Down keys.
* Since Pico AI is not available on Windows, other tools (e.g., ollama) have to be used instead. These may differ in their behavior in some details.
* Specific ollama features (e.g., new "thinking" output) are not supported.
* Text input cannot deal with soft wrap of lines in the Terminal.
* Stdin Pipe doesn't work with non-Western characters in Windows Powershell. Workaround: Use Command Prompt (cmd.exe) instead.


## Acknowledgements

Special thanks to the developers of the libraries used in this project:

 * [BurntSushi's TOML library](https://github.com/BurntSushi/toml).
 * [Atotto's Go clipboard library](https://github.com/atotto/clipboard).
 * [mattn's go-runewidth](https://github.com/mattn/go-runewidth). 

## License

This project is licensed under the MIT license.

JMK 2026
