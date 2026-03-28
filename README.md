# clickup-cli (nick-preda fork)

A fast, opinionated CLI for managing ClickUp from the terminal. Built for teams that live in the terminal and work with AI agents like Claude Code.

Forked from [triptechtravel/clickup-cli](https://github.com/triptechtravel/clickup-cli) because the upstream search was too basic, and I kept losing tasks in the void.

## Why this fork?

The upstream CLI had a few pain points that drove me crazy:

### Search that actually finds things

**The problem:** You search for "Miramonti" and get... nothing. The task exists, but it's assigned to a colleague, so the search stops at "your tasks" and never looks further. Or you search for "Posta Vecia" and it's invisible because the list has 150+ tasks and the search only checks the first page.

**The fix:** Search now drills through all levels (server-side, your tasks, your space, all spaces) and only stops early when it finds an **exact** substring match — not a random fuzzy hit. Lists with 100+ tasks are now paginated properly.

```sh
# Before (upstream): nothing, or irrelevant fuzzy matches
clickup task search "Miramonti"
# => dare più info sulle mail di ricordarsi... (fuzzy garbage)

# After (this fork): finds it, first result
clickup task search "Miramonti"
# => 86c8w6nk4  Miramonti Ristorante Pizzeria B&B  primo contatto  Michela  name
```

### No more hunting for list IDs

```sh
# See all your lists with their IDs
clickup list ls

# Create tasks by list name — no copy-pasting IDs
clickup task create --list-name "Issues" --name "[Bug] Fix login timeout"
```

### Chat integration

```sh
# Post to ClickUp Chat channels directly
clickup chat send khpgh-10335 "Deploy v2.3 done, all green"
```

## What this fork adds

| Feature | Command | The problem it solves |
|---------|---------|----------------------|
| **Smart search** | `clickup task search` | Finds tasks assigned to *anyone*, not just you |
| **Paginated search** | (automatic) | Finds tasks in lists with 100+ items |
| **List all lists** | `clickup list ls` | See list IDs without digging through the UI |
| **Chat messages** | `clickup chat send` | Post reports and alerts to Chat channels |
| **Create by list name** | `--list-name "Issues"` | No need to look up numeric IDs |
| **Assignee filter** | `--assignee me` | Filter search to your tasks only |

## Install

```sh
# From source
git clone https://github.com/nick-preda/clickup-cli.git
cd clickup-cli
make install

# Or directly with Go
go install github.com/nick-preda/clickup-cli/cmd/clickup@latest
```

## Quick start

```sh
clickup auth login         # authenticate with your API token
clickup space select       # pick a default space
clickup list ls            # see all lists and their IDs
clickup task search "bug"  # find tasks across the whole workspace
```

## Real-world examples

These are actual things I do every day:

```sh
# Find a restaurant task assigned to a colleague
clickup task search "Miramonti"

# Create a task in Issues without memorizing list IDs
clickup task create --list-name "Issues" \
  --name "[Bug] Fix login timeout" --priority 2

# Search only your own tasks
clickup task search "deploy" --assignee me

# Only show exact matches (no fuzzy noise)
clickup task search "commissione" --exact

# Add a comment and @mention someone
clickup comment add 86abc123 "@Michela this is ready for review"

# View task activity and comments
clickup task activity 86abc123

# Send an update to a Chat channel
clickup chat send khpgh-10335 "Daily report: all systems green"

# Quick status changes
clickup status set online 86abc123

# JSON output for scripts and AI agents
clickup task search "deploy" --json | jq '.[0].id'
```

## How search works

Search uses a progressive drill-down strategy, from fastest to most thorough:

1. **Server-side search** — ClickUp's own search API (fast, but limited)
2. **Sprint tasks** — if you have a sprint configured
3. **Your assigned tasks** — tasks assigned to you
4. **Default space** — paginated search in your configured space
5. **All spaces** — parallel scan of every list in every space (8 concurrent)
6. **Workspace fallback** — last resort paginated search

The key insight: search **only stops early on exact substring matches**. Fuzzy-only results are accumulated and the search continues deeper. This means you'll always find what you're looking for, while exact matches still return instantly.

## All commands

| Area | Commands |
|------|----------|
| **Tasks** | `task view`, `task create`, `task edit`, `task search`, `task list`, `task recent`, `task delete` |
| **Lists** | `list ls` |
| **Chat** | `chat send` |
| **Comments** | `comment add`, `comment list`, `comment edit`, `comment delete` |
| **Status** | `status set`, `status list`, `status add` |
| **Git** | `link pr`, `link sync`, `link branch`, `link commit` |
| **Sprints** | `sprint current`, `sprint list` |
| **Workspace** | `space list`, `space select`, `member list`, `inbox`, `tag list`, `field list` |

## Using with AI agents

This CLI is designed to pair well with Claude Code and other AI agents:

```sh
# JSON output for programmatic use
clickup task search "deploy" --json
clickup list ls --json

# Create tasks without knowing list IDs
clickup task create --list-name "Issues" --name "task name"

# Post automated reports to Chat
clickup chat send <channel-id> "automated report here"
```

## Configuration

Config lives in `~/.config/clickup/config.yml`:

```yaml
workspace: "20503057"        # your team/workspace ID
space: "90060297766"         # default space for list ls, task create --list-name
```

Set per-directory defaults with `directory_defaults` in the config file.

## Upstream docs

For features inherited from upstream, see the [original documentation](https://triptechtravel.github.io/clickup-cli/).

## License

[MIT](LICENSE)
