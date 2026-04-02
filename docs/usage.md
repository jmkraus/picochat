# Usage Guide

## Input workflow

PicoChat uses raw input mode for multiline prompts.

- Submit prompt: `Ctrl+D` (EOF)
- Cancel input: `Esc` (or `Ctrl+C`)
- Browse prompt history: Up/Down arrows

### Examples

Single line:

```text
>>> Tell me a joke! [Ctrl]+D
```

Multiline:

```text
>>> Hello, PicoChat! ↵
How are you today? ↵
↵
[Ctrl]+D
```

stdin pipe:

```bash
echo "Write a haiku about cheese" | picochat -quiet
```

Windows UTF-8 setup:

```text
chcp 65001
```

Piped command example:

```bash
echo "/models" | picochat -quiet
```

## Command-line arguments

| Argument   | Description                   |
| ---------- | ----------------------------- |
| `-config`  | Load a configuration file     |
| `-format`  | Path to JSON schema file      |
| `-history` | Load a specific session       |
| `-image`   | Path to image file            |
| `-model`   | Override configured model     |
| `-output`  | Response output format        |
| `-quiet`   | Suppress app messages         |
| `-version` | Show version and exit         |

## Commands

| Command        | Description                                       |
| -------------- | ------------------------------------------------- |
| `[Ctrl]+D`     | Submit multiline input                            |
| `[Esc]`        | Cancel multiline input and return to prompt       |
| `[Ctrl]+C`     | Cancel multiline input and return to prompt       |
| `[Up]/[Down]`  | Browse prompt history                             |
| `/copy`, `/c`  | Copy last answer to clipboard                     |
| `/paste`, `/v` | Paste clipboard content as user input and send    |
| `/info`        | Show system information                           |
| `/message`     | Show last message again                           |
| `/load`        | Load chat history from file                       |
| `/save`        | Save current chat history to file                 |
| `/models`      | List downloaded models (and switch models)        |
| `/clear`       | Clear chat history (retains system prompt)        |
| `/set`         | Set session variables (`key=value`)               |
| `/envs`        | Show environment variable status table            |
| `/image`       | Set image file path                               |
| `/retry`       | Resend chat history excluding last answer         |
| `/bye`         | Quit PicoChat                                     |
| `/help`, `/?`  | Show available commands                           |

### Command details

`/load <filename>` or `/load #<index>`:
- Without argument: lists history files and asks for selection.
- Empty input cancels loading.
- `.chat` suffix is optional.

`/save <filename>`:
- Without argument: uses a timestamp filename (for example `2025-05-11_20-26-32.chat`).

`/copy`:
- Default: copies latest assistant content.
- `/copy think`: includes reasoning section.
- `/copy code`: copies first fenced code block.
- `/copy user|system|assistant`: copies latest message by role.

`/models <index>`:
- Without argument: lists models.
- With index: switches model from cached model list.

`/set <key=value>`:
- Without argument: shows current configurable session values.
- With argument: changes runtime setting for current session only.

`/message <role>`, `/message #<index>`, `/message all`:
- No argument: shows latest message.
- Role: shows latest message for that role.
- Index: shows specific history item.
- `all`: shows full conversation with role formatting.
