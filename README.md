# vtdict

Terminal dictionary with translation. Look up English words, see definitions and examples, auto-translate to Vietnamese (or any language), and track your lookup history.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/NBVTien/vtdict/main/install.sh | bash
```

No Go required. Downloads the correct binary for your OS and architecture.

**Self-update:**
```bash
vtdict update
```

**Manual:** grab a binary from [releases](https://github.com/NBVTien/vtdict/releases/latest).

## Usage

```bash
vtdict <word>               # look up a word
vtdict -s <word>            # force-lookup (use for words that clash with commands)
vtdict help                 # show help
```

### Examples

```bash
vtdict hello
vtdict ephemeral
vtdict ephemeral --lang fr        # translate to French instead
vtdict ephemeral --no-translate   # skip translation
vtdict -s log                     # look up the word "log"
vtdict -s config                  # look up the word "config"
```

## Commands

### Other

```bash
vtdict version   # show current version
vtdict update    # update to latest release
```

### `vtdict log` — history

```bash
vtdict log                        # show last 20 looked-up words
vtdict log -n 50                  # show last 50
vtdict log --clear                # interactive picker — select words to remove
vtdict log --clear --all          # wipe entire history
vtdict log --clear --before 7d    # remove entries older than 7 days
vtdict log --clear --before 2w    # remove entries older than 2 weeks
vtdict log --clear --before 1m    # remove entries older than 1 month
vtdict log --clear --word hello   # remove one specific word
```

### `vtdict config` — configuration

```bash
vtdict config                     # interactive TUI editor
vtdict config list                # show all values + config file path
vtdict config set lang fr         # set a value
vtdict config get lang            # read a value
vtdict config reset               # restore defaults
```

**Config keys:**

| Key | Default | Description |
|-----|---------|-------------|
| `lang` | `vi` | Translation target language (ISO 639-1 code) |
| `translate` | `true` | Auto-translate on lookup |
| `phonetic` | `true` | Show phonetic transcription |
| `examples` | `true` | Show usage examples |
| `pos` | `true` | Show part of speech |

## Config file

Auto-created on first run at `~/.config/vtdict/config.toml`. Edit directly:

```toml
# vtdict configuration
lang         = "vi"
translate    = true
phonetic     = true
examples     = true
pos          = true
```

Changes take effect on the next `vtdict` invocation. No restart needed.

## Global flags

These work on every command:

| Flag | Description |
|------|-------------|
| `-s, --search <word>` | Force word lookup (bypasses subcommand names) |
| `-l, --lang <code>` | Override translation language for this call |
| `--no-translate` | Skip translation for this call |

## Data

| File | Path |
|------|------|
| Config | `~/.config/vtdict/config.toml` |
| History DB | `~/Library/Application Support/vtdict/history.db` (macOS) |

## APIs used

- [Free Dictionary API](https://dictionaryapi.dev/) — definitions, no key needed
- [MyMemory](https://mymemory.translated.net/) — translation, no key needed, 1000 req/day free
