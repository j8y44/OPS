# OPS — OpsPilot

A local-first DevOps terminal control center written in Go.

## Progress: Day 2

The Cobra CLI and Bubble Tea v2 dashboard now include a responsive sidebar,
planned integration status cards, module information, and keyboard help.
Wide terminals show a sidebar and content area; smaller terminals switch to
compact navigation and stacked status rows. The layout updates when resized.

**This is a UI preview.** Docker, GitHub Actions, and diagnostics are planned;
no live service status is shown and no credentials are required.

### Prerequisites

Go 1.26+ (required by Bubble Tea v2.1.0) and a modern terminal.
Windows Terminal is recommended. Allow at least 40 columns by 16 rows;
120 by 32 or larger gives the full sidebar and card layout.

### Run on Windows

```powershell
git clone https://github.com/j8y44/OPS.git
cd OPS
go mod download
go run ./cmd/opspilot
```

While the Day 2 pull request is open, run `git switch feat/day-2-terminal-ui`
after cloning to try this version.

### Keyboard controls

| Key | Action |
|---|---|
| Arrows or h/j/k/l | Move the section cursor |
| Tab / Shift+Tab | Cycle forward / backward, wrapping at the ends |
| Home / End | Move to the first / last section |
| Enter | Open the highlighted section |
| 1–6 | Open a section directly |
| Esc | Return to Overview; close help when help is open |
| ? | Toggle keyboard help |
| Q / Ctrl+C | Quit |

The sidebar uses `>` for the cursor and `*` for the open section. Moving the
cursor leaves the current page open until Enter is pressed. Help pauses
navigation and preserves your place. Terminals below the minimum size show
a resize message; quit keys still work.

### Commands and checks

```powershell
go run ./cmd/opspilot about
go run ./cmd/opspilot --version
go test -mod=readonly ./...
go vet ./...
go build -o opspilot.exe ./cmd/opspilot
.\opspilot.exe
```

On Linux or macOS, build with `go build -o opspilot ./cmd/opspilot` and run
`./opspilot`. GitHub Actions runs tests, vet, and builds on Linux and Windows.

See [the 8-day roadmap](docs/ROADMAP.md). Next: Docker connection and container
discovery on Day 3.
