# Sessions

PicoChat supports multiple chat sessions during one interactive run. Each session has its own conversation history and context settings. The active session is shown by the numbered prompt.

## `/chat` command

| Command | Description |
| ------- | ----------- |
| `/chat` | List available sessions and mark the active session |
| `/chat new` | Create a new session using the current system prompt and context limit |
| `/chat <number>` | Switch to a session by its one-based session number |
| `/chat copy` | Copy the complete active session into a new session |
| `/chat copy <index>` | Copy messages from the active session through the given zero-based message index |

Session numbers shown by `/chat` start at `1`. Message indexes used by `/chat copy` start at `0` and include the system message.

### Examples

List sessions:

```text
[1]❯ /chat
* 1: 8 messages
  2: 3 messages
```

Create and switch between sessions:

```text
[2]❯ /chat new
New chat session 3 created.

[1]❯ /chat 2
Switched to chat session 2. 
```

Branch a session at a message index:

```text
[1]❯ /chat copy 7
Chat copied to session 2.
```

Creating, switching, or copying a session changes the active conversation used by subsequent prompts.
