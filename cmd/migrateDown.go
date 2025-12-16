package cmd

import (
	"fmt"
	"log"

	"github.com/saleh-ghazimoradi/Cartopher/config"
	"github.com/saleh-ghazimoradi/Cartopher/infra/migrations"
	"github.com/saleh-ghazimoradi/Cartopher/infra/postgresql"
	"github.com/saleh-ghazimoradi/Cartopher/internal/logger"

	"github.com/spf13/cobra"
)

// migrateDownCmd represents the migrateDown command
var migrateDownCmd = &cobra.Command{
	Use:   "migrateDown",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrateDown called")

		cfg, err := config.GetInstance()
		if err != nil {
			log.Fatalf("config.NewInstance err: %v", err)
		}

		log := logger.NewLogger(cfg)

		postDB := postgresql.NewPostgresql(
			postgresql.WithHost(cfg.Postgresql.Host),
			postgresql.WithPort(cfg.Postgresql.Port),
			postgresql.WithUser(cfg.Postgresql.User),
			postgresql.WithPassword(cfg.Postgresql.Password),
			postgresql.WithName(cfg.Postgresql.Name),
			postgresql.WithMaxOpenConn(cfg.Postgresql.MaxOpenConn),
			postgresql.WithMaxIdleConn(cfg.Postgresql.MaxIdleConn),
			postgresql.WithMaxIdleTime(cfg.Postgresql.MaxIdleTime),
			postgresql.WithSSLMode(cfg.Postgresql.SSLMode),
			postgresql.WithTimeout(cfg.Postgresql.Timeout),
			postgresql.WithLogger(&log),
		)

		_, sqlDB, err := postDB.Connect()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to database")
		}

		migrate, err := migrations.NewMigrate(sqlDB, postDB.Name)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create migration instance")
		}

		if err := migrate.Rollback(); err != nil {
			log.Fatal().Err(err).Msg("failed to run migration")
		}

		defer func() {
			if err := migrate.Close(); err != nil {
				log.Fatal().Err(err).Msg("failed to close database connection")
			}
		}()
	},
}

func init() {
	rootCmd.AddCommand(migrateDownCmd)
}
