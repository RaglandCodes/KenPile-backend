package ken

import "github.com/gin-gonic/gin"

//SendResponse handles sendign response to client
func SendResponse(c *gin.Context, status string, message interface{}) {

	c.JSON(200, gin.H{
		"status":  status,
		"message": message,
	})

	c.Abort()
	return
}
