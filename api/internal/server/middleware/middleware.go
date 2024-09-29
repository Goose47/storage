package middleware

import (
	"Goose47/storage/internal/utils/jwt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strings"
)

type PermsProvider interface {
	IsAdmin(userID int64) (bool, error)
}

func NewAuthMiddleware(
	l *slog.Logger,
	secret string,
	permsProvider PermsProvider,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		const op = "server.middleware.AuthMiddleware"
		token := extractBearerToken(c)

		log := l.With(slog.String("op", op), slog.String("token", token))

		log.Info("trying to authenticate request")

		if token == "" {
			log.Warn("empty Bearer token")

			c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
			c.Abort()
			return
		}

		claims, err := jwt.Parse(token, secret)

		if err != nil {
			log.Warn("failed to parse token", slog.Any("error", err))

			c.JSON(http.StatusUnauthorized, gin.H{"message": "bad token"})
			c.Abort()
			return
		}

		log.Info("token parsed successfully")

		log.Info("authorizing request")

		ok, err := permsProvider.IsAdmin(int64(claims["uid"].(float64)))
		if err != nil {
			log.Warn("failed to authorize request", slog.Any("error", err))

			c.JSON(http.StatusInternalServerError, gin.H{"message": "server error"})
			c.Abort()
			return
		}

		if !ok {
			log.Warn("unauthorized")

			c.JSON(http.StatusForbidden, gin.H{"message": "forbidden"})
			c.Abort()
			return
		}

		log.Info("request authorized successfully")

		c.Set("jwt", claims)
	}
}

func extractBearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	splitToken := strings.Split(header, "Bearer ")

	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}
