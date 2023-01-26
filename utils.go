package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

func ThrowJsonError(code int, message string, c *gin.Context) {
	c.JSON(code, gin.H{
		"error": message,
	})
	c.Abort()
}

func GetUserFromToken(token string) (*User, error) {
	tokenArr := strings.Split(token, " ")
	if len(tokenArr) == 1 {
		return nil, fmt.Errorf("no token were provided")
	}
	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	_, _, err := parser.ParseUnverified(tokenArr[1], claims)
	if err != nil {
		return nil, err
	}
	username := claims["username"].(string)
	bsonUser, err := Db.GetUser(username)
	if err != nil {
		return nil, err
	}
	userByteSlice, _ := bson.Marshal(bsonUser)
	var user User
	err = bson.Unmarshal(userByteSlice, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
