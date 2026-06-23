package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"error": msg})
}

// parseIDParam парсит числовой path-параметр, при ошибке сразу отвечает 400.
func parseIDParam(c *gin.Context, param, errMsg string) (int, bool) {
	id, err := strconv.Atoi(c.Param(param))
	if err != nil {
		respondError(c, 400, errMsg)
		return 0, false
	}
	return id, true
}

// parsePagination парсит page/limit из query, при ошибке сразу отвечает 400.
func parsePagination(c *gin.Context, defaultLimit, maxLimit int) (page, limit int, ok bool) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		respondError(c, 400, "invalid page")
		return 0, 0, false
	}
	limit, err = strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(defaultLimit)))
	if err != nil || limit < 1 || limit > maxLimit {
		respondError(c, 400, "invalid limit")
		return 0, 0, false
	}
	return page, limit, true
}

// requireOwnResource проверяет, что resourceID принадлежит текущему пользователю.
func requireOwnResource(c *gin.Context, resourceID int, forbiddenMsg string) bool {
	if resourceID != c.GetInt("user_id") {
		respondError(c, 403, forbiddenMsg)
		return false
	}
	return true
}
