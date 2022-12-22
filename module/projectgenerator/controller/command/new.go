package command

import (
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type New struct {
	fs embed.FS
}

func NewNew(fs embed.FS) *New {
	return &New{fs: fs}
}

func (n *New) Use() string {
	return "new"
}

func (n *New) Short() string {
	return "Initiate altair API gateway project"
}

func (n *New) Example() string {
	return "altair new [project-directory]"
}

func (n *New) Run(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Invalid number of arguments, expected 1. Example `altair new [project-directory]`.")
		return
	}

	path := args[0]

	appYml, err := n.fs.ReadFile("template/app.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	databaseYml, err := n.fs.ReadFile("template/database.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	serviceYml, err := n.fs.ReadFile("template/routes/service-a.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	metricYml, err := n.fs.ReadFile("template/plugin/metric.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	oauthYml, err := n.fs.ReadFile("template/plugin/oauth.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, dir := range []string{path, fmt.Sprintf("%s/routes", path), fmt.Sprintf("%s/config", path), fmt.Sprintf("%s/config/plugin", path)} {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Directory does not exist, creating directory...")
			err = os.Mkdir(dir, 0755)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	err = os.WriteFile(fmt.Sprintf("%s/config/app.yml", path), appYml, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("%s/config/database.yml", path), databaseYml, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("%s/routes/service-a.yml", path), serviceYml, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("%s/config/plugin/metric.yml", path), metricYml, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("%s/config/plugin/oauth.yml", path), oauthYml, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (n *New) ModifyFlags(flags *pflag.FlagSet) {

}
