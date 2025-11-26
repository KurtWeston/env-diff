# env-diff

Compare .env files and detect missing or extra variables to catch configuration issues before deployment

## Features

- Parse .env and .env.example files with support for comments and blank lines
- Detect missing variables (in .env.example but not in .env)
- Detect extra variables (in .env but not in .env.example)
- Detect value mismatches when both files have the same key
- Preserve and display inline comments from env files
- Color-coded diff output (red for missing, yellow for extra, green for matching)
- Exit with non-zero status code when differences are found (CI/CD friendly)
- Support comparing any two env files (not just .env and .env.example)
- Show line numbers for each variable in the source files
- Ignore commented-out variables (lines starting with #)
- Handle quoted values and special characters correctly
- Provide summary statistics (total variables, missing count, extra count)

## Installation

```bash
# Clone the repository
git clone https://github.com/KurtWeston/env-diff.git
cd env-diff

# Install dependencies
go build
```

## Usage

```bash
./main
```

## Built With

- go

## Dependencies

- `github.com/spf13/cobra`
- `github.com/fatih/color`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
