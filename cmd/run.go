package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Cartopher/infra/postgresql"
	"github.com/saleh-ghazimoradi/Cartopher/internal/repository"
	"github.com/saleh-ghazimoradi/Cartopher/internal/service"
	"github.com/saleh-ghazimoradi/Cartopher/pkg/uploadProvider"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/saleh-ghazimoradi/Cartopher/config"
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

		gormDB, _, err := postDB.Connect()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to database")
		}

		var uploadProviders uploadProvider.UploadProvider
		if cfg.Upload.UploadProviders == "s3" {
			uploadProviders = uploadProvider.NewS3Provider(cfg)
		} else {
			uploadProviders = uploadProvider.NewLocalUploadProvider(cfg.Upload.Path)
		}

		middleware := middlewares.NewMiddlewares()
		authenticationMiddleware := middlewares.NewAuthentication(cfg)

		healthHandler := handlers.NewHealthHandler()
		healthRoutes := routes.NewHealthRoutes(healthHandler)

		userRepository := repository.NewUserRepository(gormDB, gormDB)
		cartRepository := repository.NewCartRepository(gormDB, gormDB)
		productRepository := repository.NewProductRepository(gormDB, gormDB)

		authService := service.NewAuthService(cfg, userRepository, cartRepository)
		userService := service.NewUserService(userRepository)
		productService := service.NewProductService(productRepository)
		uploadService := service.NewUploadService(uploadProviders)

		authHandler := handlers.NewAuthHandler(authService)
		userHandler := handlers.NewUserHandler(userService)
		productHandler := handlers.NewProductHandler(productService, uploadService)

		authRoutes := routes.NewAuthRoutes(authHandler)
		userRoutes := routes.NewUserRoutes(userHandler, authenticationMiddleware)
		productRoutes := routes.NewProductRoutes(productHandler, authenticationMiddleware)
		registerRoutes := routes.NewRegister(
			routes.WithHealthRoute(healthRoutes),
			routes.WithAuthRoute(authRoutes),
			routes.WithUserRoute(userRoutes),
			routes.WithProductRoute(productRoutes),
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
