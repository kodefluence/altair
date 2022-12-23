package command_test

import (
	"os"
	"testing"

	"github.com/kodefluence/altair/module/app"
	"github.com/kodefluence/altair/module/controller"
	"github.com/kodefluence/altair/module/projectgenerator"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCommandNew(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}
	appController := controller.Provide(nil, nil, cmd)
	appModule := app.Provide(appController)
	projectgenerator.Load(appModule)

	t.Run("Given embed file system, when command is executed then it would create a new folder contain all of altair config", func(t *testing.T) {
		os.Args = []string{"test", "new", "kuma"}
		err := cmd.Execute()
		assert.Nil(t, err)
		// testhelper.RemoveTempTestFiles("kuma")
	})
}
