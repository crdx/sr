# sr

**sr** is a composable command-line search-and-replace tool that outputs unified diffs.

Pipe file paths via stdin, get a patch via stdout. No files are modified: pipe the output to `patch -p0` to apply.

Read [the introduction post](https://textplain.org/sr).

## Installation

```sh
go install crdx.org/sr@latest
```

## CLI

```
Usage:
    sr [options] <pattern> <replacement> [<pattern> <replacement>]...

Searches for patterns and outputs a unified diff that can be piped to
patch -p0. Pass in filenames via stdin (no xargs needed).

Examples:
    find | sr foo bar | patch -p0
    find | sr 'foo: (\d+)' 'bar: $1' | patch -p0

Use $$ for a literal $, and ${1} to disambiguate.
To put a newline in the replacement use $'\n'.

Options:
    -W, --whole   Match against whole file, not lines
    -F, --fixed   Use fixed strings instead of regex
```

Patterns use [Go's regex syntax](https://pkg.go.dev/regexp/syntax).

## Examples

Preview changes without applying them.

```sh
find -name '*.go' | sr oldFunc newFunc
```

Save to a file, review, edit, then apply.

```sh
find -name '*.go' | sr oldFunc newFunc > changes.patch
vim changes.patch
patch -p0 < changes.patch
```

Combine results from different scopes into a single patch.

```sh
{ find -name '*.go' | sr foo bar; find -name '*.ts' | sr baz qux; } | patch -p0
```

Reorder function arguments using capture groups.

```sh
find -name '*.go' | sr 'assertEqual\((.+), (.+)\)' 'assertEqual($2, $1)' | patch -p0
```

Collapse a multi-line import into a single line with `--whole`.

```sh
find -name '*.py' | sr -W 'import foo\nimport bar' 'import foo, bar' | patch -p0
```

Scope replacements using any tool that outputs file paths.

```sh
git ls-files '*.ts' | sr 'console\.log\(.*\)' '' | patch -p0
fd -e rb | sr 'require "foo"' 'require "bar"' -F | patch -p0
```

## Contributions

Open an [issue](https://github.com/crdx/sr/issues) or send a [pull request](https://github.com/crdx/sr/pulls).

## Licence

[GPLv3](LICENCE).
