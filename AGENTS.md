# AGENTS.md

## Project Overview

**ginx** is a lightweight Go CLI tool that monitors remote Git repositories for changes and executes custom commands when updates are detected. Built for automating deployments, tasks, and workflows.

- **Module:** `github.com/didactiklabs/ginx`
- **Go version:** 1.23.3
- **License:** MIT
- **Current version:** v0.0.7

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go 1.23.3 |
| CLI framework | `github.com/spf13/cobra` v1.8.1 |
| Git operations | `github.com/go-git/go-git/v5` v5.13.2 |
| Logging | `go.uber.org/zap` v1.27.0 |
| Release/Build | GoReleaser v2 |
| CI/CD | GitHub Actions |
| Linting | gofumpt + golines (max 140 chars) |

## Directory Structure

```
.
├── main.go                       # Entry point
├── go.mod / go.sum               # Go module files
├── cmd/
│   └── root.go                   # Cobra CLI command definitions
├── internal/
│   └── utils/
│       ├── git.go                # Git clone, pull, commit checking
│       └── zap.go                # Logger initialization
├── scripts/
│   └── gen_doc.go                # CLI doc generator
├── docs/
│   └── ginx.md                   # Auto-generated CLI docs
├── public/
│   └── ginx.png                  # Project logo
└── .github/
    ├── workflows/
    │   ├── linting.yaml          # Lint on PRs
    │   ├── unittest.yaml         # Tests on push/PR to main
    │   ├── build.yaml            # GoReleaser dry run on PRs
    │   └── release.yaml          # GoReleaser release on tag push
    ├── dependabot.yml
    └── release.yml
```

## Commands

```bash
# Build
go build -o ginx ./

# Run
go run ./ --source <repo-url> -b <branch> -n <interval> -- <command>

# Test
go test -coverprofile=coverage.out ./...

# Lint (must pass CI)
gofumpt -d .
golines --max-len=140 . --dry-run

# Generate CLI docs
go run ./scripts/gen_doc.go

# GoReleaser dry run
goreleaser release --snapshot
```

## Code Conventions

- **Formatting:** `gofumpt` (stricter than gofmt), max line length 140 via `golines`
- **Package layout:** `main.go` at root, cobra commands in `cmd/`, shared utilities in `internal/utils/`
- **Naming:** Flag vars use camelCase with `Flag` suffix (e.g., `sourceFlag`). Exported types/funcs use PascalCase.
- **Error handling:** Fatal errors use `utils.Logger.Fatal(...)`. Non-fatal use `utils.Logger.Error(...)`. Clean up temp dirs with `os.RemoveAll(dir)` on error paths.
- **Logging:** Structured JSON via zap. Levels: debug, info, error. Set via `--log-level` flag.
- **No comments** in code unless explicitly requested.
- **No emojis** in code or docs unless explicitly requested.

## Branch & Commit Conventions

- **Branches:** `feat/<description>`, `fix/<description>`, `chore/<description>`
- **Commits:** Emoji-prefixed format, e.g., `feat : description`, `fix : description`, `chore : description`

## Architecture Notes

1. **Flow:** Create temp dir → clone remote repo → poll at interval → compare remote/local commits → pull + run command on change.
2. **`internal/utils/git.go`:** Pure Go git via `go-git`. Functions: `IsRepoCloned`, `CloneRepo`, `PullRepo`, `RunCommand`, `GetLatestRemoteCommit`, `GetLatestLocalCommit`.
3. **`internal/utils/zap.go`:** Exports a global `Logger` variable initialized in `cmd/root.go`'s `PersistentPreRun`.
4. **Version injection:** `cmd.version` set via ldflags at build time by GoReleaser.

## Known Gaps

- No test files (`*_test.go`) exist despite CI running `go test`.
- No `.gitignore` file.
- No Makefile or Dockerfile.
- `IsRepoCloned` only checks directory existence, not git validity.
- Temp dir cleanup lacks defer-based safety for unexpected process termination.
