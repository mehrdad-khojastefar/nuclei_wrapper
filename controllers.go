package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

// ###### AUTH

// log a user in
func Login(c *gin.Context) {
	// create a default username password
	var user User
	c.BindJSON(&user)

	var data map[string]interface{}
	err := Db.UsersCollection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&data)
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
	var user User
	c.BindJSON(&user)
	// check if a user is already registered with this username
	usernameAvailable, err := Db.CheckUsernameAvailability(user.Username)
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
	res, err := Db.UsersCollection.InsertOne(context.Background(), user)
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

// ###### USER

func GetUser(c *gin.Context) {
	user, err := GetUserFromToken(c.GetHeader("Authorization"))
	if err != nil {
		ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetJob(c *gin.Context) {
	jobId := strings.Split(c.Param("id"), "/")[1]
	user, err := GetUserFromToken(c.GetHeader("Authorization"))
	if err != nil {
		ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	hasJob, job := user.HasJob(jobId)
	if !hasJob && jobId != "" {
		ThrowJsonError(http.StatusNotFound, "no jobs were found with this id.", c)
		return
	}

	if jobId == "" || jobId == " " {
		byteUser, err := user.MarshalJSON()
		if err != nil {
			ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
			return
		}
		var u User
		err = json.Unmarshal(byteUser, &u)
		if err != nil {
			ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
			return
		}
		c.JSON(http.StatusOK, u.Jobs)
		return
	}
	c.JSON(http.StatusOK, job)
}

func AddJob(c *gin.Context) {
	user, err := GetUserFromToken(c.GetHeader("Authorization"))
	if err != nil {
		ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	if v, ok := jsonMap["domain"].(string); ok {
		id, err := user.AddJob(v)
		if err != nil {
			ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"job": id,
		})
		return
	}
	ThrowJsonError(http.StatusBadRequest, "invalid domain", c)
}

// 82489203-29f8-41fc-ab24-00e233439e98
func ManageJob(c *gin.Context) {
	user, err := GetUserFromToken(c.GetHeader("Authorization"))
	if err != nil {
		ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	jobId := c.Param("id")
	hasJob, job := user.HasJob(jobId)
	if !hasJob && jobId != "" {
		ThrowJsonError(http.StatusNotFound, "no jobs were found with this id.", c)
		return
	}
	action := strings.Split(c.Param("action"), "/")[1]
	switch action {
	case "":
		c.JSON(http.StatusOK, job)
		return
	case "start":
		err := job.StartNewJob(user)
		if err != nil {
			job.Status = JOBSTATUS_ERROR
			err := Db.UpdateJob(job.Id, user)
			if err != nil {
				ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
				return
			}
			ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
			return
		}

		c.JSON(http.StatusAccepted, job)
		return

	}
	// job, err := Db.GetJob(user.Username)
}
