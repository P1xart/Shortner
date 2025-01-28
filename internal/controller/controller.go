package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/p1xart/shortner-service/internal/entity"
	"github.com/p1xart/shortner-service/internal/config"
	"github.com/p1xart/shortner-service/internal/controller/request"
	"github.com/p1xart/shortner-service/internal/controller/response"
	"github.com/p1xart/shortner-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Shortner interface {
	ReduceLink(ctx context.Context, srcLink string) (string, error)
	GetSourceByShort(ctx context.Context, shortLink string) (entity.LinkDTO, error)
	IncrementVisitsByShort(ctx context.Context, shortLink string) error
}

type shortnerRoutes struct {
	log *zap.SugaredLogger

	service Shortner
}

func NewRouter(log *zap.SugaredLogger, router *gin.Engine, service Shortner) {
	log = log.With("component", "controller")

	r := &shortnerRoutes{
		log: log,

		service: service,
	}

	router.POST("/", r.reduceLink)
	router.GET("/:shortLink", r.getShortLink)
}

func (r *shortnerRoutes) reduceLink(c *gin.Context) {
	request := request.CreateLink{}
	if err := c.BindJSON(&request); err != nil {
		r.log.Errorln("failed to bind json", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	matched, err := regexp.MatchString(`^(https?:\/\/)?([\w-]{1,32}\.[\w-]{1,32})[^\s@]*$`, request.SrcLink)
	if err != nil {
		r.log.Errorln("failed to match", zap.String("string", request.SrcLink), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !matched {
		r.log.Debugln("link not matched", zap.String("link", request.SrcLink))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error of validate link",
		})
		return
	}

	reduceLink, err := r.service.ReduceLink(c, request.SrcLink)
	if err != nil {
		r.log.Errorln("failed to reduce link", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := response.GetLink{
		SrcLink:   request.SrcLink,
		ShortLink: fmt.Sprintf("%s/%s", config.DOMAIN, reduceLink),
	}
	c.JSON(http.StatusCreated, response)
}

func (r *shortnerRoutes) getShortLink(c *gin.Context) {
	shortLink := c.Param("shortLink")

	sourceLink, err := r.service.GetSourceByShort(c, shortLink)
	if err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			r.log.Debugln("link not exists", zap.String("short link", shortLink))
			c.JSON(http.StatusNotFound, gin.H{
				"error": service.ErrLinkNotFound.Error(),
			})
			return
		}
		r.log.Errorln("failed to get short link", zap.String("short link", shortLink))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err = r.service.IncrementVisitsByShort(c, shortLink); err != nil {
		if errors.Is(err, service.ErrLinkNotFound) {
			r.log.Debugln("link not exists", zap.String("short link", shortLink))
			c.JSON(http.StatusNotFound, gin.H{
				"error": service.ErrLinkNotFound.Error(),
			})
			return
		}
		r.log.Errorln("failed to increment link", zap.String("short link", shortLink))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := response.GetLink{
		SrcLink:   sourceLink.SourceLink,
		ShortLink: fmt.Sprintf("%s/%s", config.DOMAIN, shortLink),
		Visits: sourceLink.Visits,
	}
	c.JSON(http.StatusCreated, response)
}
