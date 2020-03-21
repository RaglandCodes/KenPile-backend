package ken

import "github.com/gin-gonic/gin"

//SendResponse handles sendign response to client
func SendResponse(c *gin.Context, status string, message interface{}) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	c.JSON(200, gin.H{
		"status":  status,
		"message": message,
	})

	c.Abort()
	return
}
