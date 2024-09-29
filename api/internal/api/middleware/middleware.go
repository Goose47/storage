package middleware

import (
	"Goose47/storage/internal/api/errs"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handle404(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		var notFoundErr *errs.NotFoundError
		if errors.As(err.Err, &notFoundErr) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": notFoundErr.Error(),
			})
			c.Abort()
			return
		}
	}
}
