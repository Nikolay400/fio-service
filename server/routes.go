package server

import (
	"fio-service/controller"
	"fio-service/graph"
)

func (s *Server) SetRoutes(c *controller.PersonController) {
	s.Gin.GET("/people", c.GetPeopleByFilters)
	s.Gin.GET("/people/:id", c.GetPersonById)
	s.Gin.POST("/people", c.AddPerson)
	s.Gin.PUT("/people/:id", c.UpdatePerson)
	s.Gin.DELETE("/people/:id", c.DeletePerson)
}

func (s *Server) SetGqlRoutes(h *graph.GqlHandler) {
	s.Gin.POST("/query", h.Query)
	s.Gin.GET("/", h.Playground)
}
