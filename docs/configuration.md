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
| `URL`         | string  | Core API endpoint (default: `http://localhost:11434/api`)          |
| `APIKey`      | string  | API key for OpenAI-compatible backends (recommended via env var)    |
| `Model`       | string  | Model name (must be available on backend)                          |
| `Context`     | integer | Max messages in context (`3..100`)                                 |
| `Temperature` | float   | Model temperature                                                   |
| `Top_p`       | float   | Top-p sampling value                                                |
| `Prompt`      | string  | System prompt/persona                                               |
| `Quiet`       | bool    | Suppress info/warn output                                           |
| `Reasoning`   | bool    | Enable or disable reasoning behavior                                |

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
- `PICOCHAT_REASONING`
- `PICOCHAT_QUIET`

`APIKey` can be set in `config.toml`, but this is not recommended for regular use because the key is then stored in plain text. A better approach is to fetch the key from your password manager in a shell script and export it as `PICOCHAT_API_KEY` before starting PicoChat.

## Personas

You can maintain multiple config files (for example `generic.toml`, `developer.toml`) and load them with:

```bash
picochat -config @developer
```

`@name` implies config directory plus `.toml` suffix.
