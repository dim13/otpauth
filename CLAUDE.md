# CLAUDE.md - otpauth codebase guidelines

## Build & Test Commands
- Build: `go build`
- Run tests: `go test ./...`
- Run specific test: `go test ./migration -run TestConvert`
- Run tests with verbose output: `go test -v ./...`
- Check code format: `gofmt -l .`
- Format code: `gofmt -w .`
- Run linter: `golangci-lint run`

## Code Style Guidelines
- Format: Standard Go style (gofmt)
- Imports: Group standard library first, then external packages
- Error handling: Check errors immediately with if err != nil pattern
- Naming: CamelCase for exported names, camelCase for unexported
- Comments: Package comments use // format, function comments explain purpose
- File organization: Each file has a specific focus (migrations, HTTP handlers)
- Error messages: Lowercase, no trailing punctuation
- Testing: Use table-driven tests where applicable

## Project Structure
- Main package at root with migration logic in separate migration package
- Protobuf definitions in migration.proto
- Web UI served from embedded static resources