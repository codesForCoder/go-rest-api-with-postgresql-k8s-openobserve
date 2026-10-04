package utils

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RespondJSON(c *gin.Context, status int, payload interface{}) {
	ctx := c.Request.Context()
	slog.InfoContext(ctx, "Responding with JSON", slog.Int("status", status), slog.Any("payload", payload))
	c.JSON(status, payload)
}
