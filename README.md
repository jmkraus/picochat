# PicoChat - a CLI Chat Client

## Purpose
The Mac‑only app [Pico AI Server](https://picogpt.app/), unlike tools such as Ollama, doesn't provide a dedicated CLI client.

This tool fills that gap and adds a few extra features.

## Installation

The released binary is neither notarized nor code-signed. To use it on macOS, there are two options:

 1. Build PicoChat from source
 2. Remove the quarantine flag from the binary

### Build PicoChat

 Install the dependencies and build.

 ```text
  go mod download
  go mod tidy
  go test ./...
  go build
 ```

### Remove the quarantine flag

Navigate to the folder containing the binary and enter the following:

`sudo xattr -rd com.apple.quarantine ./picochat`

Enter your administrator password to confirm. After that, you can run the binary without warnings.

## Usage

### Entering a prompt

PicoChat uses raw input mode to enter prompts.
This provides a more natural multi-line editing experience than the triple-quote (""") approach used by other applications, while still keeping the implementation simple and terminal-friendly.

Users can type or paste as much text as needed for a prompt.
Submit input by pressing Ctrl+D (end-of-file).
This method ensures that pasted text — including blank lines or code blocks — is handled safely without triggering accidental termination.

Input can be canceled at any time by pressing Esc, which immediately returns the user to the main prompt without sending the text.

Here are a few examples (↵ indicates pressing Enter):

#### Single line
```text
>>> Tell me a joke! [Ctrl]+D
```

#### Multi-line
```text
>>> Hello, PicoChat! ↵
How are you today? ↵
↵
[Ctrl]+D        ← Send prompt
```

#### Via stdin pipe
Use PicoChat in scripts via pipe, e.g.:
```
echo "Write a Haiku about Cheese" | picochat -quiet
```

The `-quiet` argument is optional and suppresses all app messages, so that only the LLM response is displayed.

On Windows, switch the shell to UTF-8 first for proper handling of non-Western characters:
```
chcp 65001
```

You can also pipe a PicoChat command via stdin. This is experimental, not thoroughly tested, and may have unexpected side effects. 
```
echo "/models" | picochat -quiet
```


### Command-line arguments

| Arg      | Description                             |
| -------- | --------------------------------------- |
| -config  | Loads a configuration file              |
| -format  | Sets the path to a JSON schema file     |
| -history | Loads a specific session                |
| -image   | Sets a path for an image file           |
| -model   | Overrides the configured model          |
| -output  | Sets the response output format         |
| -quiet   | Suppresses all app messages             |
| -version | Shows the version and exits             |


### Output formats and structured content
 
#### Output format

PicoChat provides output formats other than plain text. This simplifies output processing in pipelines, etc.
This step happens *after* inference and wraps the plain-text output into a structure while keeping
the output itself as plain text. This works with any server.

Example usage: `picochat -output json`

The following formats are available:

| Format      | Description                                  |
| ----------- | -------------------------------------------- |
| plain       | (default) Plain text output, can be omitted. |
| json        | Response formatted as JSON.                  |
| json-pretty | Response formatted as pretty-printed JSON.   |
| yaml        | Response formatted as YAML.                  |

If an invalid format is provided, PicoChat prints a warning and falls back to plain text.

**Example output: plain**

`echo "Write a Haiku about Cheese" | picochat -quiet`

```
Golden, sharp, and bold,
Melts upon your waiting tongue,
Dairy bliss unfolds.
```

**Example output: json**

`echo "Write a Haiku about Cheese" | picochat -quiet -output json`

```json
{"output":"Golden, sharp, and bold,\nMelts upon your waiting tongue,\nDairy bliss unfolds.","elapsed":"00:02","tokens_per_sec":7.8}
```

**Example output: json-pretty**

`echo "Write a Haiku about Cheese" | picochat -quiet -output json-pretty`

```json
{
  "output": "Golden, sharp, and bold,\nMelts upon your waiting tongue,\nDairy bliss unfolds."
  "elapsed": "00:02",
  "tokens_per_sec": 7.8
}
```

**Example output: yaml**

`echo "Write a Haiku about Cheese" | picochat -quiet -output yaml`

```yaml
output: "Golden, sharp, and bold,\nMelts upon your waiting tongue,\nDairy bliss unfolds."
elapsed: "00:02",
tokens_per_sec: 7.8
```

Using formatted output enables pipelines like this:

`echo "Write a Haiku about Cheese" | picochat -quiet -output json | jq -r '.elapsed'`


#### Structured content

Unlike the output formats above, structured content is generated *during* inference.
It currently works only with Ollama (untested with other local LLM servers) and works as follows:

`echo "Tell me about Canada" | picochat -format ./schema.json`

where the provided `schema.json` file looks like this:

```json
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    },
    "capital": {
      "type": "string"
    },
    "languages": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  },
  "required": [
    "name",
    "capital", 
    "languages"
  ]
}
```

In this case, generating the resulting JSON structure is an integral part of the response process
resulting in well-formed, structured content taylored to the specific needs:

```json
{
  "capital": "Ottawa",
  "languages": ["English", "French"],
  "name": "Canada"
}
```

If the exact structure doesn't matter (as long as the result is valid json), then it's also possible to
create a dummy file that contains only the word `json`. This will be accepted as well. All other schema files
are validated. In most cases this works, but the output quality depends on the model chosen
and may not be satisfactory.


### Reasoning

Pico AI Server and Ollama handle reasoning differently. While Ollama separates reasoning and content output,
Pico AI's output also includes reasoning, which is embedded in `<think>` tags.

PicoChat can handle both forms directly. However, they are displayed differently.
In Ollama, the reasoning is shown in gray color to distinguish it from the content.
Because some Pico AI models sometimes omit the opening `<think>` tag,
the reasoning part cannot be correctly recognized.
Therefore, in that case, the reasoning remains part of the content and is not colorized.

Internally, reasoning is separated from content for both Pico AI and Ollama. This means that even if
you later query the answer using `/message`, no reasoning is displayed in Pico AI. The 
`/copy` behavior is consistent: only the content is copied in Pico AI as well. To include reasoning
`/copy think` can be used.

Make sure the reasoning flag is enabled when using reasoning models.
Furthermore, reasoning behavior is highly dependent on the model used and
can occasionally cause freezes, behave unexpectedly, or be ignored by the model.

Note that reasoning isn't saved when using `/save`, as it is not considered an official
part of the conversation. 


### Configuration files

PicoChat does not require a config file: if none is found, it falls back to sensible built-in defaults.  
However, depending on the AI server you’re using, you may need to explicitly select a model via `-model`,
since the default model name may not exist (or may differ) across the wide range of available models.

If no filename is provided, PicoChat looks for a file named config.toml and searches for it in the following order:

1. The directory of the executable.
2. The full path provided via `-config`.
3. The path specified by the `CONFIG_PATH` environment variable.
4. `$XDG_CONFIG_HOME/picochat` (if `XDG_CONFIG_HOME` is set).
5. `~/.config/picochat`.

History files (see below) are always stored in the PicoChat config directory, e.g. `.config/picochat/history`.

PicoChat currently supports the following configurable values in the config file:

| Key           | Type    | Value                                                                           |
| ------------- | ------- | ------------------------------------------------------------------------------- |
| `URL`         | string  | URL of the core **API** endpoint (default: `http://localhost:11434/api`).       |
| `Model`       | string  | Model name (must already be downloaded)                                         |
| `Context`     | integer | Context size (a value between 3 and 100) i.e., the maximum number of messages!  |
| `Temperature` | float   | Model temperature                                                               |
| `TopP`        | float   | Model top-p value                                                               |
| `Prompt`      | string  | System prompt ("persona") used to specify a skill or background.                |
| `Quiet`       | bool    | Suppresses all messages (except for errors)                                     |
| `Reasoning`   | bool    | Enables or disables reasoning                                                   |

Since Pico AI currently doesn't report token counts, it is difficult to calculate an accurate context size.
This may change in the future, but for now the context size is limited by the total number of messages,
with the oldest ones being dropped when the limit is reached.


### Commands

| CMD             | DESCRIPTION                                        |
| --------------- | -------------------------------------------------- |
| [Ctrl]+D        | Submit multiline input (EOF)                       |
| [Esc], [Ctrl]+C | Cancel multiline input and return to prompt        |
| [Up]/[Down]     | Browse prompt history                              |
| /copy, /c       | Copy the last answer to clipboard                  |
| /paste, /v      | Paste clipboard contents as user input and send    |
| /info           | Show system information                            |
| /message        | Show the last message again (e.g., after load)     |
| /load           | Load chat history from a file                      |
| /save           | Save current chat history to a file                |
| /list           | List available saved history files                 |
| /models         | List downloaded models (and switch models)         |
| /clear          | Clear session context                              |
| /set            | Set session variables (key=value)                  |
| /image          | Set image file path                                |
| /retry          | Resend the chat history, excluding the last answer |
| /bye            | Quit PicoChat                                      |
| /help, /?       | Show available commands                            |

Some commands accept an argument:

#### /load `<filename>`

Without a filename, you'll be prompted to enter a name. If you leave it blank
(just pressing *ENTER*), the load process is canceled.

The filename alone is sufficient because the path is predefined (see above). The suffix can be omitted,
as it defaults to `.chat`.

If you've run `/list` before, it is also possible to load a session by index,
e.g.: `/load #3`. The hash mark indicates that an index is given rather than a filename.

#### /save `<filename>`

Without a filename, the file is saved with a timestamp as the filename, e.g. `2025-05-11_20-26-32.chat`.


#### /copy

This command copies the entire last response to the clipboard and removes the `<think>` section
for reasoning models. To keep the reasoning `/copy think` can be used instead.

If `/copy code` is entered, the first occurrence of a code block between triple backticks (` ``` `)
will be copied to the clipboard instead, ommitting surrounding explanatory text. If no code block is found,
then the command is canceled.

With a role argument (system or user; e.g., `/copy user`), the most recent prompt for that role is copied
to the clipboard.

#### /models `<index>`

Without an argument, this command lists the available models available on the **LLM** server.
If the list of models has been fetched at least once, then it's possible to switch to another model
by using the list index, e.g., `/models 3`.

#### /set `<key=value>`

Without an argument, the command lists the available parameters and shows their current values.

With an argument the parameter values can be changed for the current session, e.g., `/set top_p=0.3`.

These changes aren't persisted and cannot be saved. For permanent changes, edit the values in `config.toml`.

#### /message `<role>`

Without an argument, the last entry of the chat history (usually an assistant answer) will be shown.

With one of the possible roles (system, user, assistant), you can choose the most recent entry for a specific role.

For example, `/message user` displays the last user question again.


### Image processing

PicoChat supports basic image input. An image can either be passed as a command-line argument:

```
picochat -image ./imgfile.jpg
```

or by using the `/image` command:

```
>>> /image ./imgfile.jpg
```

You can use ~ as a shortcut for the user home directory. With the next user prompt, the image is
processed and then discarded. When passing an image path to PicoChat, the existence of that
file will be checked (and a warning is shown if it doesn't exist), but not if it is a valid image.

In combination with stdin pipe and `-model` argument a simple command-line image analysis is possible:

```
echo "What's on the image?" | picochat -model Qwen3-VL-8B-Instruct-4bit -image ./imgfile.jpg -quiet
```

When you save chat history, the user prompt contains the image as base64 encoded data.


### Personas

PicoChat supports basic personas: you can store different configuration files in the config directory,
 e.g., `generic.toml` or `developer.toml`, each with specific system prompts.

This configuration can be loaded using a shortcut, e.g., `picochat -config @developer`.
The path and `.toml` suffix can be omitted because '@' implies them.
PicoChat then starts with the specified configuration file.


### Known issues

While everything should work as expected on macOS, PicoChat is untested on Linux and only minimally
tested on Windows.

* On Windows, the Esc key must be pressed twice to cancel a prompt input. Alternatively, Ctrl + C can be used.
* The command history flickers on Windows when switching between entries using the Up/Down keys.
* Since Pico AI is not available on Windows, other tools (e.g., ollama) must be used instead. These may differ in their behavior in some ways.
* Text input doesn't handle soft-wrapped lines in the Terminal.
* Piping via stdin doesn't work with non-Western characters in Windows PowerShell. Workaround: Use Command Prompt (cmd.exe) instead.


## Acknowledgements

Special thanks to the developers of the libraries used in this project:

 * [BurntSushi's TOML library](https://github.com/BurntSushi/toml).
 * [Atotto's Go clipboard library](https://github.com/atotto/clipboard).
 * [mattn's go-runewidth](https://github.com/mattn/go-runewidth). 

## License

This project is licensed under the MIT License.

JMK 2026
