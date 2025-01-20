package controller

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Shortner interface {
	ReduceLink(ctx context.Context, link string) (string, error)
}

type shortnerRoutes struct {
	log *zap.SugaredLogger

	service Shortner
}

func NewRouter(log *zap.SugaredLogger, router *gin.Engine, service Shortner) {
	r := &shortnerRoutes{
		log: log,

		service: service,
	}

	router.POST("/:link", r.reduceLink)
}

func (r *shortnerRoutes) reduceLink(c *gin.Context) {
	link := c.Param("link")

	reduceLink, err := r.service.ReduceLink(c, link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"reduceLink": reduceLink,
	})
}
