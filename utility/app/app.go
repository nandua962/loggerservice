package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"utility/config"
	"utility/internal/consts"
	"utility/internal/controllers"
	"utility/internal/entities"
	"utility/internal/middlewares"
	"utility/internal/repo"
	"utility/internal/repo/driver"
	"utility/internal/usecases"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/middleware"
)

// Run initializes environment configuration, logging, database connection, and API routing.
// It sets up the necessary components and routes for the application and launches it.
func Run() {
	// init the env config
	cfg, err := config.LoadConfig(consts.AppName)
	if err != nil {
		panic(err)
	}

	//creating a new logger.
	var log *logger.Logger

	// Configuring client options for the logger.
	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  true,
		IncludeResponseDump: true,
		JSONFormater:        true,
	}
	if cfg.Debug {
		log = logger.InitLogger(clientOpt)
	} else {
		// Create a database logger configuration with the specified URL and secret.
		db := &logger.CloudMode{
			URL:    consts.LoggerServiceURL,
			Secret: consts.LoggerSecret,
		}
		// Initialize the logger with the specified configurations for database, file, and console logging.
		log = logger.InitLogger(clientOpt, db)
	}

	// database connection
	pgsqlDB, err := driver.ConnectDB(cfg.Db)
	if err != nil {
		log.Fatalf("unable to connect the database")
		return
	}

	// here initalizing the router
	router := initRouter()
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// middleware initialization
	m := middlewares.NewMiddlewares(cfg)
	api := router.Group("/api")
	api.Use(middleware.LogMiddleware(map[string]interface{}{}))
	api.Use(middleware.APIVersionGuard(middleware.VersionOptions{
		AcceptedVersions: cfg.AcceptedVersions,
	}))
	api.Use(middleware.Localize())
	api.Use(middleware.ErrorLocalization(
		middleware.ErrorLocaleOptions{
			Cache:                  cache.New(5*time.Minute, 10*time.Minute),
			CacheExpiration:        time.Duration(time.Hour * 24),
			CacheKeyLabel:          consts.CacheErrorKey,
			LocalisationServiceURL: fmt.Sprintf("%s/localization/error", consts.LocalisationServiceURL),
		},
	))
	api.Use(middleware.EndpointExtraction(
		middleware.EndPointOptions{
			Cache:            cache.New(5*time.Minute, 10*time.Minute),
			CacheExpiration:  time.Duration(time.Hour * 24),
			CacheKeyLabel:    consts.CacheEndpointsKey,
			ContextEndPoints: consts.ContextEndPoints,
			EndPointsURL:     fmt.Sprintf("%s/localization/endpointname", consts.LocalisationServiceURL),
		},
	))

	api.Use(m.QueryParams(
		middlewares.QueryOptions{
			Key:             consts.PaginationKey,
			DefaultLimit:    consts.DefaultLimit,
			DefaultPage:     consts.DefaultPage,
			MaxAllowedLimit: consts.MaxLimit,
		},
	))

	// complete user related initialization
	{

		// repo initialization
		genreRepo := repo.NewGenreRepo(pgsqlDB)
		languageRepo := repo.NewLanguageRepo(pgsqlDB)
		currencyRepo := repo.NewCurrencyRepo(pgsqlDB)
		CountryRepo := repo.NewCountryRepo(pgsqlDB)
		roleRepo := repo.NewRoleRepo(pgsqlDB)
		themeRepo := repo.NewThemeRepo(pgsqlDB)
		lookupRepo := repo.NewLookupRepo(pgsqlDB)
		gatewayRepo := repo.NewPaymentGatewayRepo(pgsqlDB)

		// initilizing usecases
		genreUseCase := usecases.NewGenreUseCases(genreRepo)
		languageUseCase := usecases.NewLanguageUseCases(languageRepo)
		currencyUseCase := usecases.NewCurrencyUseCases(currencyRepo)
		CountryUseCase := usecases.NewCountryUseCases(CountryRepo)
		roleUseCase := usecases.NewRoleUseCases(roleRepo)
		themeUseCase := usecases.NewThemeUseCases(themeRepo)
		lookupUseCase := usecases.NewLookupUseCases(lookupRepo)
		gatewayUseCase := usecases.NewPaymentGatewayUseCases(gatewayRepo)

		// initalizing controllers
		genreController := controllers.NewGenreController(api, genreUseCase, cfg)
		languageController := controllers.NewLanguageController(api, languageUseCase, cfg)
		currencyController := controllers.NewCurrencyController(api, currencyUseCase, cfg)
		CountryController := controllers.NewCountryController(api, CountryUseCase, cfg)
		roleController := controllers.NewRoleController(api, roleUseCase, cfg)
		themeController := controllers.NewThemeController(api, themeUseCase, cfg)
		lookupController := controllers.NewLookupController(api, lookupUseCase, cfg)
		gatewayController := controllers.NewPaymentGatewayController(api, gatewayUseCase, cfg)

		// init the routes
		genreController.InitRoutes()
		languageController.InitRoutes()
		currencyController.InitRoutes()
		CountryController.InitRoutes()
		roleController.InitRoutes()
		themeController.InitRoutes()
		lookupController.InitRoutes()
		gatewayController.InitRoutes()

	}

	// run the app
	launch(cfg, router)
}

func initRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	// CORS
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// },
		MaxAge: 12 * time.Hour,
	}))

	// common middlewares should be added here

	return router
}

// launch
func launch(cfg *entities.EnvConfig, router *gin.Engine) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", cfg.Port),
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("Server listening in...", cfg.Port)
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")

	log.Println("Server exiting")
}
