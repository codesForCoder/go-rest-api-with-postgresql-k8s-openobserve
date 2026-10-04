package controllers

import (
	"fmt"
	"go-rest-api/models"
	"go-rest-api/utils"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func FindProducts(c *gin.Context) {
	ctx := c.Request.Context()
	slog.InfoContext(ctx, "Finding products")
	var products []models.Product
	models.DB.WithContext(ctx).Find(&products)
	slog.InfoContext(ctx, "Products found", slog.Int("count", len(products)))
	utils.RespondJSON(c, http.StatusOK, products)
}

func CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	slog.InfoContext(ctx, "Creating product")
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.ErrorContext(ctx, "Failed to bind JSON", slog.Any("error", err))
		utils.RespondJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product := models.Product{Name: input.Name, Price: input.Price}
	slog.InfoContext(ctx, "Product creating with payload", slog.Any("product", product))
	models.DB.WithContext(ctx).Create(&product)
	slog.InfoContext(ctx, "Product created", product)
	utils.RespondJSON(c, http.StatusCreated, product)
}

func FindProduct(c *gin.Context) {
	fmt.Println(
		"RECEIVED traceparent =====================",
		c.GetHeader("traceparent"),
	)
	ctx := c.Request.Context()
	span := trace.SpanFromContext(c.Request.Context())
	sc := span.SpanContext()

	fmt.Println(
		"INCOMING TRACE ID ========================",
		sc.TraceID().String(),
	)
	fmt.Println(
		"INCOMING SPAN ID =========================",
		sc.SpanID().String(),
	)
	slog.InfoContext(ctx, "Finding product with id", slog.String("id", c.Param("id")))
	var product models.Product
	if err := models.DB.WithContext(ctx).Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		slog.ErrorContext(ctx, "Failed to find product", slog.Any("err", err))
		utils.RespondJSON(c, http.StatusNotFound, gin.H{"error": "Product not found!"})
		return
	}
	slog.InfoContext(ctx, "Product found with id", slog.String("id", c.Param("id")), slog.Any("product", product))
	utils.RespondJSON(c, http.StatusOK, product)
}

func UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	slog.InfoContext(ctx, "Updating product with id", slog.String("id", c.Param("id")))
	var product models.Product
	if err := models.DB.WithContext(ctx).Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		slog.ErrorContext(ctx, "Failed to find product", slog.Any("err", err))
		utils.RespondJSON(c, http.StatusNotFound, gin.H{"error": "Product not found!"})
		return
	}
	slog.InfoContext(ctx, "Updating product ", slog.Any("product", product))
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.ErrorContext(ctx, "Failed to bind JSON", slog.Any("err", err))
		utils.RespondJSON(c, http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.WithContext(ctx).Model(&product).Updates(input)
	slog.InfoContext(ctx, "Product updated", slog.Any("product", product))
	utils.RespondJSON(c, http.StatusOK, product)
}

func DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var product models.Product
	if err := models.DB.WithContext(ctx).Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		utils.RespondJSON(c, http.StatusNotFound, gin.H{"error": "Product not found!"})
		return
	}
	models.DB.WithContext(ctx).Delete(&product)
	utils.RespondJSON(c, http.StatusNoContent, nil)
}
