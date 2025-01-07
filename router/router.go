package router

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

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
	orderRouteGrouping := rg.Group("/api/user")
	orderRouteGrouping.Use(CorsMiddleware())
	for _, route := range userRoutes {
		switch route.Method {
		case "GET":
			orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
		case "OPTIONS":
			orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
		case "PUT":
			orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
		default:
			orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
				c.JSON(200, gin.H{
					"result": "Specify a valid http method with this route.",
				})
			})
		}
	}
}

func ClientRoutes() {
	r := routes{
		router: gin.Default(),
	}
	r.router.Use(CorsMiddleware())
	ver := r.router.Group(os.Getenv("API_VERSION"))
	r.SeverStatusCheck(ver)
	r.User(ver)

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

// func CorsMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	}
// }
