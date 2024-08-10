package main

import (
	"context"
	"github.com/devstik0407/shenron/cleaner"
	"github.com/devstik0407/shenron/config"
	"github.com/devstik0407/shenron/persistence"
	"github.com/devstik0407/shenron/server"
	"github.com/devstik0407/shenron/service"
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

			ctx := context.Background()

			storageClient := store.GetInstance()
			persistenceClient := persistence.New(config.App.Persistence)
			cleanupWorker := cleaner.New(config.App.CleanupWorker, storageClient)
			svc := service.New(config.App.Service, storageClient, persistenceClient, cleanupWorker)
			svc.Start(ctx)

			deps := server.Dependencies{Service: svc}
			server.Start(cfg, deps)
		},
	}
}
