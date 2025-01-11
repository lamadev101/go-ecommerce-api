package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lamadev101/ecommerce-api/auth"
	"golang.org/x/time/rate"
)

// ========== Types ==============
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*gin.Context)
}

type routes struct {
	router *gin.Engine
}
type Routes []Route

// =========== Route Grouping ===================
func (r *routes) SeverStatusCheck(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/api")
	orderRouteGrouping.Use(CorsMiddleware())

	for _, route := range statusCheckRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(404, gin.H{"message": "Method not allowed"})
			})
		}
	}
}

func (r routes) User(rg *gin.RouterGroup) {
	setupRoutes(rg, "/api/user", userRoutes, CorsMiddleware())
}

func (r routes) EcommerceProduct(rg *gin.RouterGroup) {
	setupRoutes(rg, "/api/ecommerce", ecommerceRoutes, CorsMiddleware())
}

func (r routes) EcommerceGlobalProductRoutes(rg *gin.RouterGroup) {
	setupRoutes(rg, "/api/ecommerce", ecommerceGlobalRoutes, CorsMiddleware())
}

func ClientRoutes() {
	r := routes{
		router: gin.Default(),
	}
	r.router.Use(CorsMiddleware())
	r.router.Use(rateLimiterMiddleware)

	ver := r.router.Group(os.Getenv("API_VERSION"))

	r.SeverStatusCheck(ver)
	r.User(ver)
	r.EcommerceGlobalProductRoutes(ver)

	// Protected Routes
	ver.Use(auth.Auth())
	r.EcommerceProduct(ver)

	if err := r.router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Printf("Error starting the server: %v", err)
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Use specific origin in production
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		// Process the request
		c.Next()
	}
}

var limiter = rate.NewLimiter(1, 5) // 1 request per second with a burst of 5.

func rateLimiterMiddleware(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		c.Abort()
		return
	}
	c.Next()
}

// ====================== Helper functions ====================
func setupRoutes(rg *gin.RouterGroup, prefix string, routes []Route, middleware ...gin.HandlerFunc) {
	routeGroup := rg.Group(prefix)
	routeGroup.Use(middleware...)

	for _, route := range routes {
		switch route.Method {
		case "GET":
			routeGroup.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			routeGroup.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			routeGroup.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			routeGroup.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			routeGroup.DELETE(route.Pattern, route.HandlerFunc)
		default:
			routeGroup.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(200, gin.H{
					"result": "Specify a valid http method with this route.",
				})
			})
		}
	}
}
