package main

// only need mysql OR sqlite
// both are included here for reference
import (
	"context"

	"github.com/razorpay/devstack/example/apps/webapp/logger"

	"github.com/razorpay/devstack/example/apps/webapp/config"

	"github.com/razorpay/devstack/example/apps/webapp/controllers"
	"github.com/razorpay/devstack/example/apps/webapp/database"
	"github.com/razorpay/devstack/example/apps/webapp/middlewares"
	"github.com/razorpay/devstack/example/apps/webapp/models"
	"github.com/razorpay/devstack/example/apps/webapp/pkg/tracing"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func setupMiddlewares(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.ErrorLogger())
	r.Use(gin.Recovery())
	r.Use(middlewares.TracingMiddleware())
	r.Use(middlewares.SimpleMiddleware())
}

func main() {
	appConfig := getAppConfig()

	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	tracing.InitConfig(appConfig.Tracing)
	appServiceName := appConfig.Tracing.ServiceName

	tracing.InitAppTracer(appServiceName)
	defer tracing.CloseTracers()

	database.Init()
	db := database.GetDB(context.Background())
	db.AutoMigrate(&models.Person{})
	defer db.Close()

	setupMiddlewares(r)
	registerRoutes(r)

	r.Run(":9090")
}

func registerRoutes(r *gin.Engine) {
	r.GET("/status", controllers.Status)
	r.GET("/people/", controllers.GetPeople)
	r.GET("/people/:id", controllers.GetPerson)
	r.POST("/people", controllers.CreatePerson)
	r.PUT("/people/:id", controllers.UpdatePerson)
	r.DELETE("/people/:id", controllers.DeletePerson)
}

func getAppConfig() config.AppConfig {
	config.LoadConfig(".", "default")
	appConfig := config.GetConfig()
	appConfig.Tracing.Logger = &logger.JaegerLogger{}
	return appConfig
}
