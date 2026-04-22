package command

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kodefluence/altair/module"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type New struct {
	fs      embed.FS
	plugins []module.Plugin
}

// NewNew builds the `altair new` command. plugins is the compiled-in plugin
// registry; the command writes one config/plugin/<name>.yml per plugin via
// its SampleConfig() method.
func NewNew(fs embed.FS, plugins []module.Plugin) *New {
	return &New{fs: fs, plugins: plugins}
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

	dotEnv, err := n.fs.ReadFile("template/env.sample")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, dir := range []string{path, filepath.Join(path, "routes"), filepath.Join(path, "config"), filepath.Join(path, "config", "plugin")} {
		if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
			fmt.Println("Directory does not exist, creating directory...")
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Println(err)
				return
			}
		}
	}

	if err := os.WriteFile(filepath.Join(path, "config", "app.yml"), appYml, 0644); err != nil {
		fmt.Println(err)
		return
	}

	if err := os.WriteFile(filepath.Join(path, "config", "database.yml"), databaseYml, 0644); err != nil {
		fmt.Println(err)
		return
	}

	if err := os.WriteFile(filepath.Join(path, "routes", "service-a.yml"), serviceYml, 0644); err != nil {
		fmt.Println(err)
		return
	}

	for _, p := range n.plugins {
		sample := p.SampleConfig()
		if len(sample) == 0 {
			continue
		}
		target := filepath.Join(path, "config", "plugin", p.Name()+".yml")
		if err := os.WriteFile(target, sample, 0644); err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := os.WriteFile(filepath.Join(path, ".env"), dotEnv, 0644); err != nil {
		fmt.Println(err)
		return
	}
}

func (n *New) ModifyFlags(flags *pflag.FlagSet) {}
