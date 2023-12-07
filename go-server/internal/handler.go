package internal

import (
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrBadRequest           = errors.New("error with api request body")
	ErrCannotPerformRequest = errors.New("error cannot perform request currently")
)

type Handler struct {
	service *Service
	logger  Logger
}

func NewHandler(service *Service, logger Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) EmbedTexts(c *gin.Context) {
	var input EmbeddingRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.LogError(err.Error())
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return
	}

	err := h.service.GenerateAndStoreTextEmbeddings(c, input)
	if err != nil {
		h.logger.LogError(err.Error())
		c.JSON(http.StatusInternalServerError, ErrCannotPerformRequest)
		return
	}

	c.JSON(http.StatusOK, "text embedding successful")
}

func (h *Handler) PerformSemanticSearch(c *gin.Context) {
	var input EmbeddingRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.LogError(err.Error())
		c.JSON(http.StatusBadRequest, ErrBadRequest)
		return
	}

	response, err := h.service.RetrieveFiveSimilarEmbeddingService(c, input.Text)
	if err != nil {
		h.logger.LogError(err.Error())
		c.JSON(http.StatusInternalServerError, ErrCannotPerformRequest)
		return
	}

	c.JSON(http.StatusOK, response)
}
