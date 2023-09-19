package server

import (
	"fio-service/iface"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Gin    *gin.Engine
	logger iface.Ilogger
}

func NewServer(logger iface.Ilogger) *Server {
	gin.SetMode(gin.ReleaseMode)
	return &Server{gin.Default(), logger}
}

func (server *Server) Start() error {
	server.logger.Info("Server started")
	return server.Gin.Run(":8080")
}
