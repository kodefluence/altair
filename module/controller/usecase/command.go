package usecase

import (
	"github.com/kodefluence/altair/module"
	"github.com/spf13/cobra"
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
