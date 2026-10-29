# OPS — OpsPilot

A local-first DevOps terminal control center written in Go.

## Day 1

Cobra CLI, Bubble Tea v2 terminal dashboard, keyboard navigation, unit tests and GitHub Actions CI. Docker and CI/CD integration are planned for subsequent days.

### Prerequisites

Go 1.24+ and a modern terminal (Windows Terminal recommended).

### Run on Windows

```powershell
git clone https://github.com/j8y44/OPS.git
cd OPS
go mod tidy
go run ./cmd/opspilot
```

Navigate with Up/Down or j/k; Enter selects, Esc returns to Overview, Q quits.

### Commands

```powershell
go run ./cmd/opspilot about
go run ./cmd/opspilot --version
go test ./...
go build -o opspilot.exe ./cmd/opspilot
```

See [the 8-day roadmap](docs/ROADMAP.md).
