package cmd

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/config"
	"github.com/saleh-ghazimoradi/Cartopher/infra/postgresql"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Cartopher/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/Cartopher/internal/logger"
	"github.com/saleh-ghazimoradi/Cartopher/internal/server"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")

		errorLog := slog.New(slog.NewJSONHandler(os.Stdout, nil))

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

		middleware := middlewares.NewMiddlewares()

		healthHandler := handlers.NewHealthHandler()
		healthRoutes := routes.NewHealthRoutes(healthHandler)
		registerRoutes := routes.NewRegister(
			routes.WithHealthRoute(healthRoutes),
			routes.WithMiddlewares(middleware),
		)

		gin.SetMode(cfg.Server.GinMode)

		wg := &sync.WaitGroup{}
		httpServer := server.NewServer(
			server.WithHost(cfg.Server.Host),
			server.WithPort(cfg.Server.Port),
			server.WithHandler(registerRoutes.RegisterRoutes()),
			server.WithIdleTimeout(cfg.Server.IdleTimeout),
			server.WithReadTimeout(cfg.Server.ReadTimeout),
			server.WithWriteTimeout(cfg.Server.WriteTimeout),
			server.WithLogger(&log),
			server.WithWG(wg),
			server.WithErrorLog(slog.NewLogLogger(errorLog.Handler(), slog.LevelError)),
		)

		log.Info().Str("port", cfg.Server.Port).Msg("starting http server")
		if err := httpServer.Connect(); err != nil {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
