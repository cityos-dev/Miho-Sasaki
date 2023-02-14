package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"videoservice/handler"
	"videoservice/infra"
	"videoservice/service"
)

func main() {
	g := gin.Default()

	engine, err := infra.Init()
	if err != nil {
		log.Fatalf("error occurs when setting up mysql: %v", err)
	}

	factory := service.NewVideoService(infra.NewVideDatabase(engine,
		infra.NewFileServer("/video")))

	g.Use(service.SetFactoryMiddleware(factory))

	log.Println("start server...")
	r := g.Group("/v1")
	RegisterHandlersWithOptions(r, handler.NewServerHandler())

	log.Fatal(g.Run(":8080"))
}

// RegisterHandlersWithOptions creates a handler
func RegisterHandlersWithOptions(rg *gin.RouterGroup, si handler.ServerHandler) {
	rg.GET("/files", si.GetFiles)

	rg.POST("/files", si.PostFiles)

	rg.DELETE("/files/:id", si.DeleteFilesFileId)

	rg.GET("/files/:id", si.GetFilesFileId)

	rg.GET("/health", si.GetHealth)
}
