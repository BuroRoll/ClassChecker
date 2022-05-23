package handlers

import (
	"ClassChecker/service"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() {
	router := gin.Default()
	router.Use(corsMiddleware())
	sseRouter := sse.NewServer(nil)
	defer sseRouter.Shutdown()
	service.InitSseServe(sseRouter)

	notifications := router.Group("/notifications")
	{
		notifications.GET("/:userId", func(c *gin.Context) {
			sseRouter.ServeHTTP(c.Writer, c.Request)
		})
	}
	router.Run(":8001")
}
