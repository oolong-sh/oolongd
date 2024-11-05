# Oolong

Coming soon...

## Configuration

Oolong looks for a configuration file at `~/.oolong.json`

| Option | Description | Recommended |
|--------|-------------|---------|
| `ngramRange` | Range of NGram sizes to use for keyword linking | `[1, 2, 3]` |
| `noteDirectories` | List of directories to use with Oolong | `["~/notes"]` |
| `allowedExtensions` | Whitelist of file extensions to use in linking | `[".md", ".txt", ".mdx", ".tex", ".typ"]` |
| `pluginPaths` | List of plugins (lua files) to load | `["./scripts/daily_note.lua"]` |

#### Example Configuration

```json
{
    "ngramRange": [
        1,
        2,
        3
    ],
    "noteDirectories": [
        "~/notes"
    ],
    "allowedExtensions": [
        ".md",
        ".tex",
        ".typ"
    ],
    "pluginPaths": [
        "./scripts/daily_note.lua"
    ]
}
```
