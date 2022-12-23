package usecase_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/suite"
)

type CommandSuiteTest struct {
	*ControllerSuiteTest
}

func TestCommand(t *testing.T) {
	suite.Run(t, &CommandSuiteTest{
		&ControllerSuiteTest{},
	})
}

type fakeCommand struct{}

func (*fakeCommand) Use() string                           { return "fake" }
func (*fakeCommand) Short() string                         { return "fake it" }
func (*fakeCommand) Example() string                       { return "fake it" }
func (*fakeCommand) Run(cmd *cobra.Command, args []string) {}
func (*fakeCommand) ModifyFlags(flags *pflag.FlagSet)      {}

func (suite *HttpSuiteTest) TestInjectCommand() {
	suite.controller.InjectCommand(&fakeCommand{}, &fakeCommand{}, &fakeCommand{}, &fakeCommand{})
}
