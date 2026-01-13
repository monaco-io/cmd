package timestamp

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/subcommands"
)

type TimestampCmd struct {
	input string
}

func New() subcommands.Command {
	return &TimestampCmd{}
}

func (*TimestampCmd) Name() string     { return "timestamp" }
func (*TimestampCmd) Synopsis() string { return "Convert unix timestamp or show now." }
func (*TimestampCmd) Usage() string {
	return `timestamp [-input string]:
	Convert unix timestamp or show now.`
}

func (c *TimestampCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.input, "input", "", "timestamp or 'now'")
}

func (c *TimestampCmd) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	input := c.input
	if input == "" && len(f.Args()) > 0 {
		input = f.Args()[0]
	}
	if input == "" {
		fmt.Println("unix timestamp now is:", time.Now().Unix())
		return subcommands.ExitSuccess
	}
	log.Println(input)
	isTs := len(input) == 10
	if isTs {
		ts, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			log.Fatalln("parse input failed:", err)
			return subcommands.ExitFailure
		}
		tstr := time.Unix(ts, 0).Format(time.RFC3339)
		fmt.Println("tims is:", tstr)
		return subcommands.ExitSuccess
	}
	switch input {
	case "now":
		fmt.Println("unix timestamp now is:", time.Now().Unix())
	default:
		fmt.Println("unknown input")
	}
	return subcommands.ExitSuccess
}
