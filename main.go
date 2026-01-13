package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/monaco-io/cmd/ascii_art"
	"github.com/monaco-io/cmd/fanyi"
	"github.com/monaco-io/cmd/timestamp"
)

func init() {
	{
		subcommands.Register(subcommands.HelpCommand(), "default")
		subcommands.Register(subcommands.FlagsCommand(), "default")
		subcommands.Register(subcommands.CommandsCommand(), "default")
	}

	{
		subcommands.Register(ascii_art.New(), "udf")
		subcommands.Register(fanyi.New(), "udf")
		subcommands.Register(timestamp.New(), "udf")
	}
	{
		subcommands.Register(fanyi.New(), "")
	}

}

func main() {
	flag.Parse()
	os.Exit(int(subcommands.Execute(context.Background())))
}
