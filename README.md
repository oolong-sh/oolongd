# Oolong Backend Daemon

Coming soon...

## Installation

TODO:

## Usage

TODO:

## Configuration

Oolong looks for a configuration file at `~/.config/oolong.toml`

| Option | Description | Recommended |
|--------|-------------|---------|
| `ngram_range` | Range of NGram sizes to use for keyword linking | `[1, 2, 3]` |
| `note_directories` | List of directories to use with Oolong | `["~/notes"]` |
| `ignored_directories` | Subdirectories to exclude from reading and linking | `[".git"]` |
| `allowed_extensions` | Whitelist of file extensions to use in linking | `[".md", ".txt", ".mdx", ".tex", ".typ"]` |
| `plugin_paths` | List of plugins (lua files) to load | `["./scripts/daily_note.lua"]` |

#### Example Configuration

```toml
# Range of NGram sizes to use for keyword linking
ngram_range = [ 1, 2, 3 ]

# List of directories to read into Oolong
note_directories = [
    "~/notes",
    "~/Documents/notes"
]

# Subdirectory patterns to exclude from the file watcher
ignored_directories = [
    ".git",
    ".templates",
]

# Whitelist of file extensions to use in linking
allowed_extensions = [
    ".md",
    ".mdx",
    ".tex",
    ".typ",
]

# List of plugins (lua files) to load
plugin_paths = [ "./scripts/daily_note.lua" ]
```
