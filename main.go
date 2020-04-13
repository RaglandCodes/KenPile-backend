package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/ken"
	_ "github.com/heroku/x/hmetrics/onload"
	_ "github.com/lib/pq"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:2020")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		c.Next()
	}
}
func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.Default()
	// TODO see the difference bettween gin.New() and gin.Default()
	//router.Use(gin.Logger())
	router.Use(CORSMiddleware())
	router.LoadHTMLGlob("templates/*.html")
	//router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/new", func(c *gin.Context) {
		c.String(http.StatusOK, "namCC@s")
	})

	router.POST("/createNewNote", ken.RouteCreateNewNote)
	router.POST("/verifyIdToken", ken.RouteVerifyIDToken)
	router.POST("/logBackIn", ken.RouteLogBackIn)

	router.Run(":" + port)
}

// clear && go build -o bin/go-getting-started -v . && heroku local web
