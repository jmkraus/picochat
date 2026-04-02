# Image Processing

PicoChat supports image prompts.

Set image via CLI argument:

```bash
picochat -image ./imgfile.jpg
```

Or via command:

```text
>>> /image ./imgfile.jpg
```

Behavior:

- Image path is consumed with the next user prompt.
- After sending, image path is discarded.
- Saved history stores image payloads as base64/data URL content.

Example with stdin pipe:

```bash
echo "What's on the image?" | picochat -model Qwen3-VL-8B-Instruct-4bit -image ./imgfile.jpg -quiet
```
