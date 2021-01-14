package main

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/codefluence-x/altair/cmd"
	"github.com/codefluence-x/journal"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		journal.Error("Error running altair:", err).SetTags("altair", "main").Log()
	}
}
