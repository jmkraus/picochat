# Debug and Known Issues

## Debug commands

`/hello`:
- Sends probe prompt: `Hello, are you there?`

`/test`:
- Generates an OS-aware script in the current directory.
- Iterates over detected models with a sample prompt.
- Exits PicoChat after script generation.

## Known issues

- Linux support is not broadly tested.
- Windows support is limited.
- On Windows, `Esc` may need to be pressed twice to cancel input (fallback: `Ctrl+C`).
- Command history may flicker on Windows with Up/Down navigation.
- Some Windows shells can have stdin/Unicode limitations.
- Soft-wrapped terminal lines are not fully handled by the input editor.
