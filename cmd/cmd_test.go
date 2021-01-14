package cmd_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/codefluence-x/altair/cmd"
	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv("DATABASE_NAME", "altair_development")
	os.Setenv("DATABASE_USERNAME", "root")
	os.Setenv("DATABASE_PASSWORD", "root")
	os.Setenv("DATABASE_HOST", "localhost")
	os.Setenv("DATABASE_PORT", "33060")
	os.Setenv("BASIC_AUTH_PASSWORD", "1234")
}

func TestRootCmd(t *testing.T) {
	root := cmd.RootCmd

	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)

	expectedResult := "Light Weight and Robust API Gateway.\n\nUsage:\n  altair [command]\n\nAvailable Commands:\n  help             Help about any command\n  migrate          Do a migration from current version into latest versions.\n  migrate:down     Down the migration from current version into earliest versions.\n  migrate:rollback Do a migration rollback from current versions into previous versions.\n  server           Run API gateway services.\n\nFlags:\n  -h, --help   help for altair\n\nUse \"altair [command] --help\" for more information about a command.\n"

	_, err := root.ExecuteC()
	require.NoError(t, err)
	require.Equal(t, expectedResult, buf.String())
}

func TestServerCmd(t *testing.T) {
	root := cmd.RootCmd
	root.SetArgs([]string{"server"})

	_, err := root.ExecuteC()
	require.NoError(t, err)
}

func TestMigrateCmd(t *testing.T) {
	root := cmd.RootCmd
	root.SetArgs([]string{"migrate", "altair_development"})

	_, err := root.ExecuteC()
	require.NoError(t, err)
}

func TestMigrateDownCmd(t *testing.T) {
	root := cmd.RootCmd
	root.SetArgs([]string{"migrate:down", "altair_development"})

	_, err := root.ExecuteC()
	require.NoError(t, err)
}

func TestMigrateRollbackCmd(t *testing.T) {
	root := cmd.RootCmd
	root.SetArgs([]string{"migrate:rollback", "altair_development"})

	_, err := root.ExecuteC()
	require.NoError(t, err)
}
