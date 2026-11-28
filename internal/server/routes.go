package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	api := r.Group("/api/v1")
	api.POST("/tracks", gin.WrapF(s.music.CreateTrack))
	api.GET("/tracks", gin.WrapF(s.music.GetTracks))
	api.GET("/tracks/:id", func(c *gin.Context) {
		c.Request.SetPathValue("id", c.Param("id"))
		s.music.GetTrack(c.Writer, c.Request)
	})
	api.GET("/tracks/:id/stream", func(c *gin.Context) {
		c.Request.SetPathValue("id", c.Param("id"))
		s.music.StreamTrack(c.Writer, c.Request)
	})

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
