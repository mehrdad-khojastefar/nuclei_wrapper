package utils

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"hamravesh.ir/mehrdad-khojastefar/database"
	"hamravesh.ir/mehrdad-khojastefar/models"
)

func GetUserFromToken(token string) (*models.User, error) {
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
	bsonUser, err := database.Db.GetUser(username)
	if err != nil {
		return nil, err
	}
	userByteSlice, _ := bson.Marshal(bsonUser)
	var user models.User
	err = bson.Unmarshal(userByteSlice, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
