package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"ecommerce-microservices/pkg/middleware"
	"ecommerce-microservices/services/product/config"
	"ecommerce-microservices/services/product/controller"
	"ecommerce-microservices/services/product/repository"
	"ecommerce-microservices/services/product/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()

	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		log.Fatalf("unable to parse database config: %v", err)
	}
	poolConfig.MaxConns = 25
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = 5 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("product database connection error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("product database ping failed: %v", err)
	}

	repo := repository.NewPostgresProductRepository(db)
	svc := service.NewProductService(repo)
	ctrl := controller.NewProductController(svc)

	r := gin.New()
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RecoveryMiddleware("product-service"))
	r.Use(middleware.StructuredLogger("product-service"))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "product-service", "status": "UP"})
	})
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "product-service", "status": "ALIVE"})
	})
	r.GET("/health/ready", func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "NOT_READY", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "READY"})
	})

	r.POST("/products", ctrl.Create)
	r.GET("/products/:id", ctrl.GetByID)
	r.GET("/products", ctrl.GetAll)

	log.Println("Product Service listening on :" + cfg.Port)
	_ = r.Run(":" + cfg.Port)
}
