package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"hamravesh.ir/mehrdad-khojastefar/utils"
)

func GetUser(c *gin.Context) {
	user, err := utils.GetUserFromToken(c.GetHeader("Authorization"))
	if err != nil {
		utils.ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetJobs(c *gin.Context) {
	user, err := utils.GetUserFromToken(c.GetHeader("Authorization"))
	if err != nil {
		utils.ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	if user.Jobs == nil {
		c.JSON(http.StatusOK, gin.H{
			"jobs": nil,
		})
		return
	}
	c.JSON(http.StatusOK, user.Jobs)
}

func AddJob(c *gin.Context) {
	user, err := utils.GetUserFromToken(c.GetHeader("Authorization"))
	if err != nil {
		utils.ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		utils.ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	var jsonMap map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonMap)
	if err != nil {
		utils.ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
		return
	}
	if v, ok := jsonMap["domain"].(string); ok {
		id, err := user.AddJob(v)
		if err != nil {
			utils.ThrowJsonError(http.StatusInternalServerError, err.Error(), c)
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"job": id,
		})
		return
	}
	utils.ThrowJsonError(http.StatusBadRequest, "invalid domain", c)
}
