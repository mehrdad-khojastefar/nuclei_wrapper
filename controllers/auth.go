package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"hamravesh.ir/mehrdad-khojastefar/database"
	"hamravesh.ir/mehrdad-khojastefar/models"
)

// log a user in
func Login(c *gin.Context) {
	// create a default username password
	var user models.User
	c.BindJSON(&user)

	var data map[string]interface{}
	err := database.Db.UsersCollection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if val, ok := data["password"].(string); ok && user.Password == val {
		// generating a new token
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user.Username
		tokenString, err := token.SignedString([]byte(uuid.New().String()))
		// for development only ( no need to sign the jwt )
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "something went wrong",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "user not found or password incorrect",
	})
}

// register a user
func Register(c *gin.Context) {
	var user models.User
	c.BindJSON(&user)
	// check if a user is already registered with this username
	usernameAvailable, err := database.Db.CheckUsernameAvailability(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if !usernameAvailable {
		c.JSON(http.StatusConflict, gin.H{
			"error": "username already taken",
		})
		return
	}
	user.Id = uuid.New().String()
	res, err := database.Db.UsersCollection.InsertOne(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id": res.InsertedID,
	})
}
