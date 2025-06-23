package main

import (
	"github.com/gin-gonic/gin"
	"kluisz-object-storage/config"
	"kluisz-object-storage/handlers"
	"kluisz-object-storage/middleware"
	_ "kluisz-object-storage/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// @title           Object Storage API
// @version         1.0
// @description     This is a sample server for object storage -- bucket and object manipulations
// @host            localhost:8080
// @BasePath        /
// @contact.name    Ritu Priyadarshini
// @contact.email   ritu.priyadarshini@kluisz.ai
func main() {
	config.LoadConfig()

	r := gin.Default()

	//logger middleware-with log rotation
	zapLoggerR := middleware.NewZapLogger()
	r.Use(middleware.ZapLogger(zapLoggerR, true))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/bucket", handlers.CreateBucket)
	r.DELETE("/bucket/:name", handlers.DeleteBucket)
	r.GET("/buckets", handlers.ListBuckets)
	r.POST("/upload/:bucket", handlers.UploadFile)
    r.GET("/download/:bucket/:file", handlers.DownloadFile)
	r.GET("/objects/:bucket", handlers.ListObjects)
	r.DELETE("/objects/:bucket/:file", handlers.DeleteObject)


	r.Run(":8080")
}
