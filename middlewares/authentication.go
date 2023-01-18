package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"hamravesh.ir/mehrdad-khojastefar/database"
	"hamravesh.ir/mehrdad-khojastefar/utils"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		// no token were provided
		if len(token) == 1 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "no auth token were provided",
			})
			c.Abort()
			return
		}
		user, err := utils.GetUserFromToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}
		// check if the user is present in the database or not
		if available, err := database.Db.CheckUsernameAvailability(user.Username); available || err != nil {
			fmt.Println(available, err)
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid auth token",
			})
			c.Abort()
			return
		}
	}
}
