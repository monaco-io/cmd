package ascii_art

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"path"

	"github.com/google/subcommands"
	"github.com/monaco-io/cmd/ascii_art/font"
)

//go:embed logo/*
var content embed.FS

type asciiArtCmd struct {
	name    string
	face    string
	list    bool
	viewAll bool
}

func New() subcommands.Command {
	return &asciiArtCmd{}
}

func (*asciiArtCmd) Name() string     { return "ascii-art" }
func (*asciiArtCmd) Synopsis() string { return "Generate ASCII art from string." }
func (*asciiArtCmd) Usage() string {
	return `ascii-art [-list] [-view-all] [-name string] [-face string]:
	Generate ASCII art for a string.`
}

func (c *asciiArtCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.list, "list", false, "list all fonts")
	f.BoolVar(&c.viewAll, "view-all", false, "view all fonts")
	f.StringVar(&c.name, "name", "bilibili.com", "content of the ascii string")
	f.StringVar(&c.face, "face", "", "typeface of the ascii")
}

func (c *asciiArtCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if c.list {
		font.List()
		return subcommands.ExitSuccess
	}
	if c.viewAll {
		font.EchoAll(c.name)
		return subcommands.ExitSuccess
	}
	logo, err := content.ReadFile(path.Join("logo", c.name))
	if err != nil {
		logo = []byte(font.AsciiArt(c.name, c.face))
	}
	fmt.Printf("%s", logo)
	return subcommands.ExitSuccess
}
