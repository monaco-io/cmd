package fanyi

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/google/subcommands"
	"github.com/monaco-io/cmd/fanyi/src"
)

type fanyiCmd struct {
	targetLang  string
	filePath    string
	outputPath  string
	interactive bool
	noCache     bool
	showVersion bool
}

// New returns a new fanyi command.
func New() subcommands.Command {
	return &fanyiCmd{}
}

func (*fanyiCmd) Name() string     { return "fanyi" }
func (*fanyiCmd) Synopsis() string { return "AI-Powered CLI Translation Tool" }
func (*fanyiCmd) Usage() string {
	return `fanyi [OPTIONS] [TEXT]

OPTIONS:
  -t, --target-lang <LANG>    Target language (zh, en, ja, ko, es, etc.)
  -l, --language <LANG>       Alias for --target-lang
  -f, --file <FILE>           Input file path
  -o, --output <FILE>         Output file path
  -i, --interactive           Interactive mode
  --no-cache                  Skip cache (force fresh translation)
  --version                   Show version information

EXAMPLES:
  fanyi hello world
  fanyi -t zh "hello world"
  fanyi -f input.txt
  echo "hello" | fanyi -t ko
  fanyi -i

`
}

func (c *fanyiCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.targetLang, "t", "", "Target language (zh, en, ja, etc.)")
	f.StringVar(&c.targetLang, "target-lang", "", "Target language")
	f.StringVar(&c.targetLang, "l", "", "Language alias for -t")
	f.StringVar(&c.filePath, "f", "", "Input file path")
	f.StringVar(&c.filePath, "file", "", "Input file")
	f.StringVar(&c.outputPath, "o", "", "Output file path")
	f.StringVar(&c.outputPath, "output", "", "Output file")
	f.BoolVar(&c.interactive, "i", false, "Interactive mode")
	f.BoolVar(&c.interactive, "interactive", false, "Interactive mode")
	f.BoolVar(&c.noCache, "no-cache", false, "Skip cache")
	f.BoolVar(&c.showVersion, "version", false, "Show version")
}

func (c *fanyiCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Show version
	if c.showVersion {
		fmt.Println("Fanyi v1.0.0 - AI-Powered CLI Translation Tool")
		return subcommands.ExitSuccess
	}

	// Load configuration
	cfg, err := src.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nQuick fix: Set your API key with:\n")
		fmt.Fprintf(os.Stderr, "  export FANYI_API_KEY=\"sk-your-api-key-here\"\n\n")
		return subcommands.ExitFailure
	}

	// Disable cache if requested
	if c.noCache {
		cfg.Cache.Enabled = false
	}

	// Create translator
	trans, err := src.NewTranslator(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize translator: %v\n", err)
		return subcommands.ExitFailure
	}
	defer trans.Close()

	// Interactive mode
	if c.interactive {
		c.runInteractive(trans)
		return subcommands.ExitSuccess
	}

	// Get text from arguments, file, or stdin
	var text string

	if c.filePath != "" {
		// Read from file
		data, err := os.ReadFile(c.filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			return subcommands.ExitFailure
		}
		text = string(data)
	} else if f.NArg() > 0 {
		// Use command line arguments
		text = strings.Join(f.Args(), " ")
	} else {
		// Check if stdin has data
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Reading from pipe
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			text = strings.Join(lines, "\n")
		}
	}

	// Validate input
	if text == "" {
		fmt.Fprint(os.Stderr, c.Usage())
		return subcommands.ExitUsageError
	}

	// Trim whitespace
	text = strings.TrimSpace(text)

	// Translate
	result, err := trans.Translate(text, c.targetLang)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Translation error: %v\n", err)
		return subcommands.ExitFailure
	}

	// Output
	if c.outputPath != "" {
		err := os.WriteFile(c.outputPath, []byte(result), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			return subcommands.ExitFailure
		}
		fmt.Printf("✓ Translation saved to: %s\n", c.outputPath)
	} else {
		fmt.Println(result)
	}

	return subcommands.ExitSuccess
}

func (c *fanyiCmd) runInteractive(trans *src.Translator) {
	fmt.Println("Fanyi - Interactive Mode")
	if c.targetLang != "" {
		fmt.Printf("Translating to: %s\n", c.targetLang)
	} else {
		fmt.Println("Translating to default languages")
	}
	fmt.Println("Enter text to translate (type 'exit' or 'quit' to quit):")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		// Check for exit commands
		lowerText := strings.ToLower(text)
		if lowerText == "exit" || lowerText == "quit" || lowerText == "q" {
			fmt.Println("Goodbye!")
			break
		}

		// Translate
		result, err := trans.Translate(text, c.targetLang)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}

		fmt.Println(result)
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
