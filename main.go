package main

import (
	"github.com/spf13/cobra"
)

func main() {
	cli := cobra.Command{
		Use:   "shenron",
		Short: "A powerful key-value store",
	}
	cli.AddCommand(startServer())

	if err := cli.Execute(); err != nil {
		panic(err)
	}
}
