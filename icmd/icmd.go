package icmd

import (
	"context"
	"flag"

	"github.com/google/subcommands"
)

type Interface struct{}

var _ subcommands.Command = (*Interface)(nil)

func (*Interface) Name() string               { panic("Interface-Name") }
func (*Interface) Synopsis() string           { return "Interface-Synopsis" }
func (*Interface) Usage() string              { return "Interface-Usage" }
func (p *Interface) SetFlags(f *flag.FlagSet) {}
func (p *Interface) Execute(_ context.Context, f *flag.FlagSet, _ ...any) subcommands.ExitStatus {
	panic("Interface-Execute\n")
}
