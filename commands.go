package main

import (
	"github.com/devstik0407/shenron/cleaner"
	"github.com/devstik0407/shenron/config"
	"github.com/devstik0407/shenron/server"
	"github.com/devstik0407/shenron/store"
	"github.com/spf13/cobra"
)

func startServer() *cobra.Command {
	return &cobra.Command{
		Use:     "server",
		Short:   "starts http server",
		Aliases: []string{"start", "serve"},
		Run: func(_ *cobra.Command, _ []string) {
			cfg := &server.Config{
				Address: config.App.Address,
				Port:    config.App.Port,
			}
			s := store.New(config.App.Store)
			deps := server.Dependencies{Store: s}

			go func() {
				cleaner.Clean(config.App.Cleaner, s)
			}()

			server.Start(cfg, deps)
		},
	}
}
