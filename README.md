# PicoChat

A lightweight terminal chat client for local and OpenAI-compatible LLM backends.

## Introduction

PicoChat provides an interactive CLI workflow with multiline input, session history, structured output options, and image-capable prompts.

## Key Features

- Interactive multiline chat input (`Ctrl+D` submit, `Esc`/`Ctrl+C` cancel)
- Session history save/load
- Multiple backend protocols (`ollama`, `openai`, `responses`)
- Output formatting (`plain`, `json`, `json-pretty`, `yaml`)
- Structured content generation via JSON schema
- Image prompt support
- Clipboard helpers and runtime commands

## Installation

Build from source:

```bash
go mod download
go mod tidy
go test ./...
go build
```

On macOS, unsigned binaries may need quarantine removal:

```bash
sudo xattr -rd com.apple.quarantine ./picochat
```

## Quick Start

Interactive mode:

```bash
./picochat
```

Pipe mode:

```bash
echo "Write a Haiku about Cheese" | ./picochat -quiet
```

## Documentation

Detailed guides are in [`/docs`](docs):

- [Usage Guide](docs/usage.md)
- [Configuration](docs/configuration.md)
- [Output and Structured Content](docs/output-and-structured-content.md)
- [Image Processing](docs/image-processing.md)
- [Debug and Known Issues](docs/debug-and-known-issues.md)

## Acknowledgements

- [BurntSushi/toml](https://github.com/BurntSushi/toml)
- [atotto/clipboard](https://github.com/atotto/clipboard)
- [mattn/go-runewidth](https://github.com/mattn/go-runewidth)

## License

MIT
