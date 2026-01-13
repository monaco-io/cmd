package fanyi

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/subcommands"
	"github.com/monaco-io/cmd/fanyi/src"
)

type fanyiCmd struct {
	init    bool
	noCache bool
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
	--init                      Initialize config at ~/.config/fanyi/config.yaml
	--no-cache                  Skip cache (force fresh translation)

EXAMPLES:
  fanyi hello world
	echo "hello" | fanyi

`
}

func (c *fanyiCmd) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.init, "init", false, "Initialize config file (~/.config/fanyi/config.yaml)")
	f.BoolVar(&c.noCache, "no-cache", false, "Skip cache")
}

func (c *fanyiCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// Init config and exit
	if c.init {
		if err := initConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Init failed: %v\n", err)
			return subcommands.ExitFailure
		}
		fmt.Println("✓ Config initialized at ~/.config/fanyi/config.yaml")
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

	// Get text from arguments or stdin
	var text string

	if f.NArg() > 0 {
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
	result, err := trans.Translate(text, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Translation error: %v\n", err)
		return subcommands.ExitFailure
	}

	fmt.Println(result)

	return subcommands.ExitSuccess
}

// initConfig creates ~/.config/fanyi/config.yaml from the bundled example if it does not exist.
func initConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot find home dir: %w", err)
	}
	dstDir := filepath.Join(home, ".config", "fanyi")
	dstFile := filepath.Join(dstDir, "config.yaml")

	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("cannot create config dir: %w", err)
	}

	if _, err := os.Stat(dstFile); err == nil {
		return fmt.Errorf("config already exists at %s", dstFile)
	}

	srcCandidates := []string{
		filepath.Join("fanyi", "config.example.yaml"),
		filepath.Join("config.example.yaml"),
	}

	var content []byte
	for _, p := range srcCandidates {
		if data, err := os.ReadFile(p); err == nil {
			content = data
			break
		}
	}
	if len(content) == 0 {
		return fmt.Errorf("cannot find config.example.yaml")
	}

	if err := os.WriteFile(dstFile, content, 0644); err != nil {
		return fmt.Errorf("cannot write config: %w", err)
	}
	return nil
}
