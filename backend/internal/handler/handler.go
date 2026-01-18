// Package handler provides HTTP handlers for the PieMDM API endpoints.
// It contains the main Handler struct and common utilities for request processing,
// user authentication, and pagination handling.
package handler

import (
	"strconv"

	"piemdm/pkg/jwt"
	"piemdm/pkg/log"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// GetUserIDFromCtx extracts the user ID from the Gin context.
// It returns the user ID if it exists, otherwise returns an empty string.
// It uses the jwt.CustomClaims struct to extract the user ID.
func GetUserIDFromCtx(ctx *gin.Context) string {
	v, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return v.(*jwt.CustomClaims).ID
}

// GetPage extracts page and pageSize parameters from the Gin context query parameters.
// It returns page number and page size with default values if not provided.
func GetPage(ctx *gin.Context) (page, pageSize int) {
	// TODO use viper config
	page, _ = strconv.Atoi(ctx.Query("page"))
	pageSize, _ = strconv.Atoi(ctx.Query("pageSize"))
	if pageSize == 0 {
		pageSize = 15
	}
	if page == 0 {
		page = 1
	}
	return
}
