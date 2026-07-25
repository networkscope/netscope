# Contributing to NetScope

NetScope is open source and accepts contributions under the terms of the GNU General Public License v3.0.

## Code of conduct

Be respectful. This project is maintained by volunteers. Criticism, discussion, and feedback are welcome, but harassment or personal attacks are not.

## How to contribute

1. Fork the repository and create a branch from `main`
2. Make your change
3. Add tests that cover the new behavior
4. Run `go test ./...` and `go vet ./...`
5. Open a pull request with a clear description

## Development setup

```bash
git clone https://github.com/networkscope/netscope.git
cd netscope
go build ./...
go test ./...
```

## Code style

- Idiomatic Go
- Small functions, clear names
- Explicit error handling
- Avoid global mutable state
- Keep comments sparse and focused on non-obvious reasoning

## License

By contributing, you agree that your work will be licensed under the project's GPL-3.0 license.
