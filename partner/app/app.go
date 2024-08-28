package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"partner/config"
	"partner/internal/consts"
	"partner/internal/controllers"

	"partner/internal/entities"
	"partner/internal/repo"
	"partner/internal/repo/driver"
	"partner/internal/usecases"
	"syscall"
	"time"

	cacheConf "gitlab.com/tuneverse/toolkit/core/cache"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"gitlab.com/tuneverse/toolkit/core/activitylog"
	"gitlab.com/tuneverse/toolkit/core/logger"
	"gitlab.com/tuneverse/toolkit/middleware"
)

// method run
// env configuration
// logrus, zap
// use case intia
// repo initalization
// controller init

// Run function
func Run() {
	// init the env config
	cfg, err := config.LoadConfig(consts.AppName)
	if err != nil {
		panic(err)
	}

	file := &logger.FileMode{
		LogfileName:  "partner.log",
		LogPath:      "logs",
		LogMaxAge:    consts.LogMaxAge,
		LogMaxSize:   consts.LogMaxSize,
		LogMaxBackup: consts.LogMaxBackup,
	}

	clientOpt := &logger.ClientOptions{
		Service:             consts.AppName,
		LogLevel:            "info",
		IncludeRequestDump:  false,
		IncludeResponseDump: false,
	}

	if cfg.Debug {
		logger.InitLogger(clientOpt, file)
	} else {
		db := &logger.CloudMode{

			URL:    consts.LoggerServiceURL,
			Secret: consts.LoggerSecret,
		}
		logger.InitLogger(clientOpt, db, file)
	}
	// Check if the application is in debug mode.
	if cfg.Debug {
		logger.InitLogger(clientOpt, file)
	} else {
		db := &logger.CloudMode{
			URL:    consts.LoggerServiceURL,
			Secret: consts.LoggerSecret,
		}

		// Initialize the logger with the specified configurations for database, file, and console logging.
		logger.InitLogger(clientOpt, db, file)
	}
	// Check if the application is in debug mode.
	if cfg.Debug {
		logger.InitLogger(clientOpt, file)
	} else {
		// Release Mode: Logs will print to a database, file, and console.

		// Create a database logger configuration with the specified URL and secret.
		db := &logger.CloudMode{
			// Database API endpoint (for best practice, load this from an environment variable).
			URL: consts.LoggerServiceURL,
			// Secret for authentication.
			Secret: consts.LoggerSecret,
		}

		// Initialize the logger with the specified configurations for database, file, and console logging.
		logger.InitLogger(clientOpt, db, file)
	}
	activitylog, err := activitylog.Init(consts.ActivityLogServiceURL)
	if err != nil {
		log.Fatalf("unable to connect the activity log service : %v", err)
		return
	}

	// database connection
	pgsqlDB, err := driver.ConnectDB(cfg.Db)
	if err != nil {
		log.Fatalf("unable to connect the database")
		return
	}
	redisClient, err := cacheConf.New(&cacheConf.RedisCacheOptions{
		Host:     cfg.Redis.Host,
		UserName: cfg.Redis.UserName,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	if err != nil {
		log.Fatalf("unable to connect  redis server: err: %s", err)
		return
	}

	// here initalizing the router
	router := initRouter()
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

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
			LocalisationServiceURL: fmt.Sprintf("%s/localization/error", consts.ErrorLocalizationURL),
		},
	))
	api.Use(middleware.EndpointExtraction(
		middleware.EndPointOptions{
			Cache:            cache.New(5*time.Minute, 10*time.Minute),
			CacheExpiration:  time.Duration(time.Hour * 24),
			CacheKeyLabel:    consts.CacheEndpointsKey,
			ContextEndPoints: consts.ContextEndPoints,
			EndPointsURL:     fmt.Sprintf("%s/localization/endpointname", consts.ErrorLocalizationURL),
		},
	))

	// complete user related initialization
	{

		// repo initialization
		partnerRepo := repo.NewPartnerRepo(pgsqlDB, redisClient)

		// initilizing usecases
		partnerUseCases := usecases.NewPartnerUseCases(partnerRepo, redisClient)

		// initalizing controllers
		partnerControllers := controllers.NewPartnerController(api, partnerUseCases, cfg, activitylog)

		// init the routes
		partnerControllers.InitRoutes()
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
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
			log.Fatalf("listen: %s\n", err.Error())
		}
	}()
	_, err := fmt.Println("Server listening in...", cfg.Port)
	if err != nil {
		return
	}
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown:%s", err.Error())
	}
	// catching ctx.Done(). timeout of 5 seconds.

	<-ctx.Done()
	log.Printf("timeout of 5 seconds.")

	log.Printf("Server exiting")
}
