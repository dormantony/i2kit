package main

import (
	"log"

	"bitbucket.org/riberaproject/cli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	commands := &cobra.Command{
		Use:   "i2kit COMMAND [ARG...]",
		Short: "Manage i2kit applications",
	}
	commands.AddCommand(
		cmd.Deploy(),
		cmd.Destroy(),
	)
	if err := commands.Execute(); err != nil {
		log.Fatal(err)
	}
}