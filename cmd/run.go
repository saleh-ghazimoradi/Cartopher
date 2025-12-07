package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/config"
	"github.com/saleh-ghazimoradi/Cartopher/infra/postgresql"
	"github.com/saleh-ghazimoradi/Cartopher/internal/logger"
	"log"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		cfg, err := config.GetInstance()
		if err != nil {
			log.Fatalf("Error getting config: %v", err)
		}

		log := logger.NewLogger(cfg)

		db, err := postgresql.NewPostgresql(cfg)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to connect to database")
		}

		mainDB, err := db.DB()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get database connection")
		}

		defer func() {
			if err := mainDB.Close(); err != nil {
				log.Fatal().Err(err).Msg("failed to close database connection")
			}
		}()

		gin.SetMode(cfg.Server.GinMode)

		log.Info().Msg("starting server")

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
