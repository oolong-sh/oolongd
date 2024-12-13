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

<!-- TODO: Docker image, Nix and Homebrew packages? -->

## Usage

After defining your configuration in `~/.config/oolong.toml` [See configuration](#configuration), oolongd can be run as follows:

```sh
oolong
```

The service will run in the background, and the API will be accessible on port 11975.
The graph can be opened by navigating to [http://localhost:11975](http://localhost:11975), and clicking on note nodes will open the file in an editor of your choice (as defined in config).


To use the Oolong web editor see [oolong-web](https://github.com/oolong-sh/oolong-web). (The web editor also includes a graph integration)

## Configuration

Oolong looks for a configuration file at `~/.config/oolong.toml`


**Core Settings** (required):

| Option | Description | Recommended |
|--------|-------------|---------|
| `note_directories` | List of directories to use with Oolong | `["~/notes"]` |
| `ignored_directories` | Subdirectories to exclude from reading and linking | `[".git"]` |
| `allowed_extensions` | Whitelist of file extensions to use in linking | `[".md", ".txt", ".mdx", ".tex", ".typ"]` |
| `open_command` | Command to run when clicking a graph node | `["code"]` (See below for more details) |
| `pinning_enabled` | **Optional** boolean indicating if note pinning should be enabled | `false` |

The `open_command` option is used by the graph to allow you to open a clicked note in an editor of your choice.

For example, to open a note in VSCode use `open_command = ["code"]`

To use your system default editor:
- Linux: `open_command = ["xdg-open"]`
- MacOS: `open_command = ["open"]`
- Windows: `open_command = ["start"]`

For more situations where you want to run a more complex command, separate consecutive arguments:
- `open_command = ["tmux", "neww", "-c", "shell", "nvim"]` (opens Neovim in a new tmux window in the active session)

<!-- TODO: example using a script -->


**Linker Config** (optional -- falls back to defaults)

| Option | Description | Default |
|--------|-------------|-------------|
| `ngram_range` | Range of NGram sizes to use for keyword linking | `[1, 2, 3]` |
| `stop_words` | Additional stop words to exclude from keyword extraction | `[]` |

**Graph Settings** (optional -- falls back to defaults):

| Option | Description | Default |
|--------|-------------|-------------|
| min_node_weight | Minimum NGram weight to allow to show up in the graph | `2.0` (Increase to a larger number for large note directories) |
| max_node_weight | Maximum size of a node in the graph (larger values are clamped to this size) | `10.0` |
| min_link_weight | The minimum allowed link strength between a note and NGram | `0.1` (Increase to a larger number (0.2-0.3) for larger note directories) |
| default_mode | Default graph mode (2d/3d) | `"2d"` |


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
    ".txt",
    ".tex",
    ".typ",
]

# Enables note pinning
pinning_enabled = true

# Command to run when open endpoint it called (a note node is clicked on the graph)
# Note: All arguments MUST be separated into separate strings (see config for more details)
open_command = [ "code" ]

[linker]
# Range of NGram sizes to use for keyword linking
ngram_range = [ 1, 2, 3 ]

# Extra stop words to exclude from NGram generation
stop_words = [ "hello" ]

# graph settings (required)
[graph]
min_node_weight = 2.0
max_node_weight = 12.0
min_link_weight = 0.1
default_mode = "3d"

#
# NOTE: do not include the following section if you do not want to use cloud sync
#
[sync]
host = "127.0.0.1" # replace with your server hostname/ip
user = "your_username" # replace with your server username
port = 22 # server ssh port
private_key_path = "/home/<your_username>/.ssh/<your_ssh_key_path>" # replace with your private key path


# optional plugins section (not currently recommended)
[plugins]
# List of plugins (lua files) to load
plugin_paths = [ "./scripts/daily_note.lua" ]
```
