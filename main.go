package main

import (
	"bufio"
	"os"
	"regexp"

	"crdx.org/logger"

	"crdx.org/duckopt/v2"
)

func getUsage() string {
	return `
		Usage:
			$0 [options] <pattern> <replacement> [<pattern> <replacement>]...

		Outputs a unified diff that can be piped to patch -p0.
		Pass in filenames via stdin (no xargs needed).

		Examples:
			find | $0 foo bar | patch -p0
			find | $0 'foo: (\d+)' 'bar: $1' | patch -p0

		Use $$ for a literal $, and ${1} to disambiguate.
		To put a newline in the replacement use $'\n'.

		Options:
			-W, --whole   Match against whole file, not lines
			-F, --fixed   Use fixed strings instead of regex
	`
}

type Opts struct {
	Patterns     []string `docopt:"<pattern>"`
	Replacements []string `docopt:"<replacement>"`
	Fixed        bool     `docopt:"--fixed"`
	Whole        bool     `docopt:"--whole"`
}

func main() {
	logger.InitStderr()

	opts := duckopt.MustBind[Opts](getUsage(), "$0")

	if len(opts.Patterns) != len(opts.Replacements) {
		logger.Fatal("odd number of arguments")
	}

	var substitutions []substitution
	for i, pattern := range opts.Patterns {
		if opts.Fixed {
			substitutions = append(substitutions, substitution{
				fixedString: pattern,
				replacement: opts.Replacements[i],
				isFixed:     true,
			})
		} else {
			compiled, err := regexp.Compile(pattern)
			if err != nil {
				logger.Err("invalid regex pattern: %s", pattern)
				os.Exit(1)
			}

			substitutions = append(substitutions, substitution{
				pattern:     compiled,
				replacement: opts.Replacements[i],
			})
		}
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		process(opts.Whole, scanner.Text(), &substitutions)
	}
}
