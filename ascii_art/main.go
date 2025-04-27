package main

import (
	"embed"
	"flag"
	"fmt"
	"path"

	"github.com/monaco-io/cmd/ascii_art/font"
)

var (
	//go:embed logo/*
	content embed.FS

	name     string
	face     string
	list     bool
	viewwAll bool
)

func init() {
	flag.BoolVar(&list, "list", false, "list all fonts")
	flag.BoolVar(&viewwAll, "view-all", false, "view all fonts")
	flag.StringVar(&name, "name", "bilibili.com", "content of the ascii string")
	flag.StringVar(&face, "face", "", "typeface of the ascii")
	flag.Parse()
}

func main() {
	if list {
		font.List()
		return
	}
	if viewwAll {
		font.EchoAll(name)
		return
	}
	logo, err := content.ReadFile(path.Join("logo", name))
	if err != nil {
		logo = []byte(font.AsciiArt(name, face))
	}
	fmt.Printf("%s", logo)
}
