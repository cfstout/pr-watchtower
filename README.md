# PR Watchtower

PR Watchtower is a terminal-based tool for monitoring GitHub Pull Requests. It provides a unified view of PRs that need your attention and PRs you've authored, with automation capabilities to streamline your workflow.

## Features

- **Unified Dashboard**: View "Needs Review" and "My PRs" lists in a single TUI.
- **Real-time Updates**: Automatically refreshes PR status at configurable intervals.
- **Change Tracking**: Highlights new and updated PRs since you last checked.
- **Automation**: Trigger GitHub Actions workflows directly from the TUI (e.g., for auto-fixing or deployment).
- **Configurable**: Customize search queries and refresh intervals via `config.yaml`.

## Installation

### Prerequisites

- Go 1.21+
- [GitHub CLI (`gh`)](https://cli.github.com/) installed and authenticated (`gh auth login`).

### Build from Source

```bash
git clone https://github.com/cfstout/pr-watchtower.git
cd pr-watchtower
go mod tidy
go build -o pr-watchtower ./cmd/watchtower
```

## Usage

### Quick Start (Makefile)

- **Run**: `make run` (or just `make`)
- **Build**: `make build`
- **Clean**: `make clean`

### Manual Commands

#### Development / Quick Run

You can run the application directly without building:

```bash
go run ./cmd/watchtower
```

### Run Binary

If you built the binary using the steps above:

```bash
./pr-watchtower
```

### Keybindings

- `Tab`: Switch between "Needs Review" and "My PRs" lists.
- `j` / `Down`: Move cursor down.
- `k` / `Up`: Move cursor up.
- `r`: Force refresh PRs.
- `a`: Trigger automation (runs `agent-fix.yml` workflow on selected PR).
- `q` / `Ctrl+C`: Quit.

## Configuration

The application looks for a `config.yaml` file in the current directory or `$HOME/.pr-watchtower`.

Default configuration:

```yaml
github:
  refresh_interval: 2m
  queries:
    needs_review: "review-requested:@me state:open"
    my_prs: "author:@me state:open"
```

## License

MIT License. See [LICENSE](LICENSE) for details.
