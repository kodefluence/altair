package usecase

import (
	"github.com/spf13/cobra"

	"github.com/kodefluence/altair/module"
)

// TODO: Return error when command already registered
func (c *Controller) InjectCommand(commands ...module.CommandController) {
	for _, command := range commands {
		cmd := &cobra.Command{
			Use:     command.Use(),
			Short:   command.Short(),
			Example: command.Example(),
			Run:     command.Run,
		}
		command.ModifyFlags(cmd.Flags())
		c.rootCommand.AddCommand(cmd)
	}
}
