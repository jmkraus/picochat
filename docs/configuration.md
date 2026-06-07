# Configuration

## Config file discovery

A config file is optional. If none is found, PicoChat uses built-in defaults.

Search order:

1. Directory of executable
2. Path passed via `-config`
3. Path in `CONFIG_PATH`
4. `$XDG_CONFIG_HOME/picochat` (if set)
5. `~/.config/picochat`

History files are stored in the PicoChat config directory (for example `.config/picochat/history`).

## Config keys

| Key           | Type    | Description                                                         |
| ------------- | ------- | ------------------------------------------------------------------- |
| `Backend`     | string  | Backend flavor (`ollama`, `openai`, `responses`)                    |
| `URL`         | string  | Core API endpoint (default: `http://localhost:11434`)               |
| `APIKey`      | string  | API key for OpenAI-compatible backends (recommended via env var)    |
| `Model`       | string  | Model name (must be available on backend)                           |
| `Context`     | integer | Max messages in context (`3..100`)                                  |
| `Temperature` | float   | Model temperature (`0..2`)                                          |
| `Top_p`       | float   | Top-p sampling value (`0..1`)                                       |
| `Prompt`      | string  | System prompt/persona                                               |
| `Quiet`       | bool    | Suppress info/warn output                                           |
| `Reasoning`   | bool    | Enable or disable reasoning behavior                                |
| `Effort`      | string  | Tune the trace length of reasoning output (`low`, `medium`, `high`) |

## Environment variables

Load order: defaults -> config file -> environment variables.

Runtime CLI flags still take precedence for overlapping settings (for example `-model`, `-quiet`).

Supported variables:

- `PICOCHAT_BACKEND`
- `PICOCHAT_URL`
- `PICOCHAT_API_KEY`
- `PICOCHAT_MODEL`
- `PICOCHAT_CONTEXT`
- `PICOCHAT_TEMPERATURE`
- `PICOCHAT_TOP_P`
- `PICOCHAT_QUIET`
- `PICOCHAT_REASONING`
- `PICOCHAT_EFFORT`

`APIKey` can be set in `config.toml`, but this is not recommended for regular use because the key is then stored in plain text. A better approach is to fetch the key from your password manager in a shell script and export it as `PICOCHAT_API_KEY` before starting PicoChat. Here's an example for macOS:

```bash
# store api key (token) into Apple KeyChain
security add-generic-password -a "$USER" -s openai-personal-token -w "<token>"
```

where `<token>` is the actual token, and `openai-personal-token` is an arbitrary name. Be aware that the token still can be exposed via shell history. For full safety consider creating the entry directly in the Apple KeyChain.

Then it can be wrapped into a small shell script:

```bash
#!/bin/sh

# read api key (token) from KeyChain and start PicoChat.
export PICOCHAT_API_KEY="$(security find-generic-password -a "$USER" -s openai-personal-token -w)"
picochat
```


## Personas

You can maintain multiple config files (for example `generic.toml`, `developer.toml`) and load them with:

```bash
picochat -config @developer
```

`@name` implies config directory plus `.toml` suffix.

But it's also possible to enter a complete path:

```bash
picochat -config ./path/to/developer.toml
```
