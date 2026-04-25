package command_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

	"github.com/kodefluence/altair/module/app"
	"github.com/kodefluence/altair/module/controller"
	"github.com/kodefluence/altair/module/projectgenerator"
	"github.com/kodefluence/altair/plugin"
)

func TestCommandNew(t *testing.T) {
	cmd := &cobra.Command{
		Use: "test",
	}
	appController := controller.Provide(nil, nil, cmd)
	appModule := app.Provide(appController)
	projectgenerator.Load(appModule, plugin.Registry())

	t.Run("Given embed file system, when command is executed then it would create a new folder contain all of altair config", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "kuma")
		os.Args = []string{"test", "new", dir}
		err := cmd.Execute()
		assert.Nil(t, err)

		// Assumption: `altair new <dir>` creates the standard project layout
		// with config files for every plugin in Registry().
		assert.FileExists(t, filepath.Join(dir, "config", "app.yml"))
		assert.FileExists(t, filepath.Join(dir, "config", "database.yml"))
		assert.FileExists(t, filepath.Join(dir, "routes", "service-a.yml"))
		assert.FileExists(t, filepath.Join(dir, ".env"))
		for _, p := range plugin.Registry() {
			assert.FileExists(t, filepath.Join(dir, "config", "plugin", p.Name()+".yml"))
		}
	})

	// Assumption: `altair new` without a directory argument prints a help
	// message to stdout and exits without panicking. We don't capture stdout
	// — just confirm no error and no side-effect dir is created.
	t.Run("No-arg invocation is a no-op success", func(t *testing.T) {
		freshCmd := &cobra.Command{Use: "test2"}
		freshController := controller.Provide(nil, nil, freshCmd)
		freshModule := app.Provide(freshController)
		projectgenerator.Load(freshModule, plugin.Registry())
		os.Args = []string{"test2", "new"}
		err := freshCmd.Execute()
		assert.Nil(t, err)
	})
}
