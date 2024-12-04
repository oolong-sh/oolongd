# Oolong Backend Daemon

The Oolong backend service handles note parsing, aggregation, NGram extraction, and cloud synchronization.

## Installation

Oolongd can be installed with go install:
```sh
go install github.com/oolong-sh/oolongd
```

It can also be built from source:
```sh
git clone https://github.com/oolong-sh/oolongd.git
cd oolongd
go build
```

TODO: Docker image, Nix and Homebrew packages?

## Usage

After defining your configuration in `~/.config/oolong.toml` [See configuration](#configuration), oolongd can be run as follows:

```sh
oolong
```

The service will run in the background, and the API will be accessible on port 11975.

To view the constructed graph and use the web editor see [oolong-web](https://github.com/oolong-sh/oolong-web).

(Serving up a copy of the graph from the backend is a WIP)

## Configuration

Oolong looks for a configuration file at `~/.config/oolong.toml`


**Core Settings** (required):

| Option | Description | Recommended |
|--------|-------------|---------|
| `ngram_range` | Range of NGram sizes to use for keyword linking | `[1, 2, 3]` |
| `note_directories` | List of directories to use with Oolong | `["~/notes"]` |
| `ignored_directories` | Subdirectories to exclude from reading and linking | `[".git"]` |
| `allowed_extensions` | Whitelist of file extensions to use in linking | `[".md", ".txt", ".mdx", ".tex", ".typ"]` |
| `stop_words` | Additional stop words to exclude from keyword extraction | `[]` |


**Graph Settings** (required):

| Option | Description | Recommended |
|--------|-------------|-------------|
| min_node_weight | Minimum NGram weight to allow to show up in the graph | `2.0` (Increase to a larger number for large note directories) |
| max_node_weight | Maximum size of a node in the graph (larger values are clamped to this size) | `10.0` |
| min_link_weight | The minimum allowed link strength between a note and NGram | `0.1` (Increase to a larger number (0.2-0.3) for larger note directories) |


**Cloud Synchronization Settings** (optional):

| Option | Description |
|--------|-------------|
| host | Server IP address or hostname |
| user | Server user account name |
| port | Server SSH port |
| private_key_path | Path to your SSH private key file |

**Plugin Settings** (optional -- not recommended):
| Option | Description | Recommended |
|--------|-------------|-------------|
| plugin_paths | List of scripts to load | `[]` |


### Example Configuration

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
    ".txt"
    ".tex",
    ".typ",
]

stop_words = [
    "hello",
]


# graph settings (required)
[graph]
min_node_weight = 8.0
max_node_weight = 12.0
min_link_weight = 0.2

# optional plugins section (not currently recommended)
[plugins]
# List of plugins (lua files) to load
plugin_paths = [ "./scripts/daily_note.lua" ]

#
# NOTE: do not include the following section if you do not want to use cloud sync
#
[sync]
host = "127.0.0.1" # replace with your server hostname/ip
user = "your_username" # replace with your server username
port = 22 # server ssh port
private_key_path = "/home/<your_username>/.ssh/<your_ssh_key_path>" # replace with your private key path
```
