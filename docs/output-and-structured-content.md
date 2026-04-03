# Output and Structured Content

## Output formats

Output formatting is post-processing after inference. It wraps plain text output into a selected format.

Example:

```bash
picochat -output json
```

Available formats:

| Format        | Description                  |
| ------------- | ---------------------------- |
| `plain`       | Default plain text           |
| `json`        | JSON response                |
| `json-pretty` | Pretty-printed JSON response |
| `yaml`        | YAML response                |

If an invalid format is provided, PicoChat falls back to `plain` and prints a warning.

## Structured content (`-format`)

Structured content is generated during inference from a JSON schema:
This feature is not available on all backends and works best with Ollama via the Ollama API, but it can also be used with the OpenAI-compatible completions API.

```bash
echo "Tell me about Canada" | picochat -format ./schema.json
```

Schema example:

```json
{
  "type": "object",
  "properties": {
    "name": { "type": "string" },
    "capital": { "type": "string" },
    "languages": {
      "type": "array",
      "items": { "type": "string" }
    }
  },
  "required": ["name", "capital", "languages"]
}
```

## Reasoning behavior

Backends can emit reasoning differently.

- Ollama typically separates reasoning and content.
- Other backends may embed reasoning markers in content.

PicoChat separates reasoning from content internally where possible.

Notes:

- Reasoning output depends on model and backend behavior.
- `/copy think` includes reasoning when available.
- Reasoning is not persisted by `/save`.
