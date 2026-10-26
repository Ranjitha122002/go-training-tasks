package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"ecommerce-microservices/pkg/middleware"
	"ecommerce-microservices/services/order/config"
	"ecommerce-microservices/services/order/controller"
	"ecommerce-microservices/services/order/repository"
	"ecommerce-microservices/services/order/service"

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
		log.Fatalf("order database connection error: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("order database ping failed: %v", err)
	}

	repo := repository.NewPostgresOrderRepository(db)
	svc := service.NewOrderService(repo, cfg.ProductServiceURL)
	ctrl := controller.NewOrderController(svc)

	r := gin.New()
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.RecoveryMiddleware("order-service"))
	r.Use(middleware.StructuredLogger("order-service"))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "order-service", "status": "UP"})
	})
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "order-service", "status": "ALIVE"})
	})
	r.GET("/health/ready", func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "NOT_READY", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "READY"})
	})

	r.POST("/orders", ctrl.Create)
	r.GET("/orders", ctrl.ListAll)
	r.GET("/orders/:id", ctrl.GetByID)
	r.PUT("/orders/:id/status", ctrl.UpdateStatus)
	r.GET("/orders/:id/product", ctrl.GetOrderWithProduct)

	log.Println("Order Service listening on :" + cfg.Port)
	_ = r.Run(":" + cfg.Port)
}
