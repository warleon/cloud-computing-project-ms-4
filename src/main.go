    package main

    import (
    	"net/http"

    	"github.com/gin-gonic/gin"
    )

    func main() {
    	// Initialize a new Gin router
    	router := gin.Default()

    	// Define a GET route
    	router.GET("/", func(c *gin.Context) {
    		c.JSON(http.StatusOK, gin.H{
    			"message": "Hello, World!",
    		})
    	})

    	// Run the server on port 8080
    	router.Run(":3000")
    }