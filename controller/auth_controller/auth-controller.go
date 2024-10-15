package authcontroller

import (
	"net/http"
	"os"
	"time"

	mongo "my-project/mongo"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var hmacSampleSecret []byte

func Register(c *gin.Context, db *mongo.UserDB) {
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
		c.JSON(http.StatusOK, gin.H{"message": "Register successfully"})
	}
}

func Login(c *gin.Context, db *mongo.UserDB) {
	// cash body to login struct
	var login LoginBody
	if err := c.ShouldBindBodyWithJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// find user in database
	userExist, err := db.GetUserByUserName(login.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userExist != nil {
		err = bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(login.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		} else {
			// hmacSampleSecret := []byte("my_secret_key")
			hmacSampleSecret := []byte(os.Getenv("JWT_SECRET_KEY"))

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"userID": userExist.ID,
				"exp":    time.Now().Add(time.Minute * 5).Unix(),
				"iat":    time.Now().Unix(),
			})

			tokenString, err := token.SignedString(hmacSampleSecret)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to generate token",
					"details": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Login successfully", "token": tokenString})
			return
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
}

func Logout(c *gin.Context, db *mongo.UserDB) {

}
