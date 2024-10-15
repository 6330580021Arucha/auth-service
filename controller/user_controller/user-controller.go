package controller

import (
	"fmt"
	"net/http"
	_ "strings"

	mongo "my-project/mongo"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(c *gin.Context, db *mongo.UserDB) {
	users, err := db.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(users) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users in database"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_count": len(users), "users": users})
}

func GetUserByID(c *gin.Context, db *mongo.UserDB) {
	id := c.Param("id")
	user, err := db.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func GetProfile(c *gin.Context, db *mongo.UserDB) {
	id := c.MustGet("userID").(string)
	user, err := db.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNoContent, gin.H{"message": "No users in database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func InsertUser(c *gin.Context, db *mongo.UserDB) {
	// cash body to struct
	var user mongo.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// encrypt password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	user.Password = string(encryptedPassword)

	// check user exist
	isExist, err := db.UserExist(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isExist {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	} else {
		err = db.InsertUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Create user successfully"})
		return
	}
}

func UpdateUser(c *gin.Context, db *mongo.UserDB) {
	// cash body to struct
	var user mongo.User
	if err := c.ShouldBindBodyWithJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "can not bind body with JSON"})
		return
	}

	// encrypt password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	user.Password = string(encryptedPassword)

	id := c.Param("id")

	err = db.UpdateUser(id, user)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user found with ID: %s", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user was updated"})
}

func DeleteUser(c *gin.Context, db *mongo.UserDB) {
	id := c.Param("id")
	err := db.DeleteUser(id)
	if err != nil {
		if err.Error() == fmt.Sprintf("no user found with ID: %s", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user was deleted"})
}
