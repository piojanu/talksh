# talksh

**talksh** lets you describe your intent in natural language and returns the corresponding shell command.
This tool makes it easy to discover and remember shell commands without searching manually.

## Usage

```bash
$ talksh ask "find all .txt files modified today"
> find . -type f -name '*.txt' -mtime 0
```

## Installation

Make sure you have Go installed (>= 1.18).

```bash
go install github.com/piojanu/talksh/cmd/talksh@latest
```

This will install the `talksh` binary in your `$GOPATH/bin`.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
