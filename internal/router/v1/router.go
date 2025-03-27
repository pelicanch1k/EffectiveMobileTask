package router_v1

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/pelicanch1k/EffectiveMobileTestTask/docs"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/handler"
)


func NewRouter(h *handler.Handler) *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	songs := router.Group("/api/v1")
	{
		songs.GET("/songs", h.GetSongs)
		songs.GET("/song/:id/lyrics", h.GetSongLyrics)

		songs.DELETE("/song/:id", h.DeleteSong)

		songs.PUT("/song", h.UpdateSong)
		songs.POST("/song", h.AddSong)
	}

	return router
}