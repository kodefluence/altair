package usecase

import (
	"github.com/kodefluence/altair/module"
	"github.com/spf13/cobra"
)

func (c *Controller) InjectCommand(commands ...module.CommandController) {
	for _, command := range commands {
		c.rootCommand.AddCommand(
			&cobra.Command{
				Use:     command.Use(),
				Short:   command.Short(),
				Example: command.Example(),
				Run:     command.Run,
			},
		)
	}
}
