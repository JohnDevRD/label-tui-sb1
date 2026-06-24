# Contributing

Thanks for considering contributing to label-tui-sb1.

## How to contribute

1. Fork the repo
2. Create a branch: `git checkout -b feature/my-feature`
3. Commit your changes: `git commit -am 'feat: add my feature'`
4. Push: `git push origin feature/my-feature`
5. Open a Pull Request

## Guidelines

- Run `go fmt` before committing
- Keep commits focused and use conventional commit messages
- Add tests when adding functionality
- Keep the TUI responsive — avoid blocking the event loop

## Development

```bash
go run ./cmd/label-tui-sb1
```

The app expects templates in `~/.label-tui/templates/` and settings in `~/.label-tui/settings.json`.
