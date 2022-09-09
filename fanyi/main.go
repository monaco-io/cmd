package main

import (
	"os"
	"strings"

	"github.com/monaco-io/cmd/fanyi/src"
)

func main() {
	args := strings.Join(os.Args[1:], " ")
	src.Translate(args)
}
