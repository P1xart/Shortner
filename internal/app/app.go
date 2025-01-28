package app

import (
	"os"

	"github.com/p1xart/shortner-service/internal/controller"
	"github.com/p1xart/shortner-service/internal/repo"
	"github.com/p1xart/shortner-service/internal/service"
	"github.com/p1xart/shortner-service/pkg/postgres"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("Logger init error")
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	sugar.Info("router init")
	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	router.SetTrustedProxies(nil)

	sugar.Info("postgresql init")
	postgresql, err := postgres.New(sugar)
	if err != nil {
		sugar.Fatal(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	sugar.Info("repository init")
	repository := repo.NewRepo(sugar, postgresql)

	sugar.Info("service init")
	services := service.NewService(sugar, repository)

	sugar.Info("router init")
	controller.NewRouter(sugar, router, services)

	sugar.Info("starting app...")
	err = router.Run(":8080")
	if err != nil {
		sugar.Fatal(os.Stderr, "error running router", err)
		os.Exit(1)
	}
}
