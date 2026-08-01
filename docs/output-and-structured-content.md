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

## Structured content (`-schema`)

Structured content is generated during inference from a JSON schema. Support depends on the configured backend and model. Ollama supports schemas through its native `format` option, while OpenAI-compatible Responses API backends may also support schemas through their structured-output format.

```bash
echo "Tell me about Canada" | picochat -schema ./country.json
```

Schema example for `country.json`:

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
  "required": ["name", "capital", "languages"],
  "additionalProperties": false
}
```

For `ollama`, PicoChat additionally performs local post-processing: content streaming is disabled, the complete output is collected and validated against the provided schema, then reformatted into pretty JSON and stored in the chat history. It is rendered as one complete output block after inference finishes.

This local validation, pretty-printing, history formatting, and single-block rendering is currently specific to Ollama. Other backends, including OpenAI-compatible Responses API backends, may enforce structured output themselves, but their behavior depends on the endpoint and model.

If reasoning is enabled, reasoning output may still be displayed during inference. The structured JSON content itself is not displayed until processing is complete.

If the Ollama response is invalid JSON or does not satisfy the schema, PicoChat reports an error and does not add the assistant response to the chat history.

When a schema is supplied, PicoChat uses plain rendering for the structured Ollama result. This displays the pretty-printed JSON as one complete block instead of wrapping it in the normal JSON or YAML `ChatResult` output.

LLM output based on the example schema:

```json
{
  "capital": "Ottawa",
  "languages": [
    "English",
    "French"
  ],
  "name": "Canada"
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
