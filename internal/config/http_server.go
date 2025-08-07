package config

import "github.com/gin-gonic/gin"

type HttpServer struct {
	gin_object *gin.Engine
	port       string
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		gin_object: gin.Default(),
	}
}
func (h *HttpServer) SetPort(port string) {
	h.port = port
}
func (h *HttpServer) Start() {
	h.gin_object.Run(h.port)
}
