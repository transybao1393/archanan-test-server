package main

import (
	"fmt"
	"time"

	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-contrib/cors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AmplifyLogin struct{
	IDTOKEN string `json:"idToken,omitempty"`
	ERR string `json:"errorCode,omitempty"`
	ERR_MESSAGE string `json:"errorMessage,omitempty"`
}

type JWTToken struct {
	Message string `json:"message"`
	jwt.StandardClaims
}

func main() {
	r := gin.Default()
	r.POST("/internal", handleInternalRequest)
	r.GET("/healthcheck", handleHealthCheck)
	r.Use(cors.Default())
	r.Run(":6001")
}

func handleInternalRequest(c *gin.Context) {
	//- bind to struct
	var login AmplifyLogin
	c.BindJSON(&login)

	//- headers
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	//- check if error is null
	if login.ERR == "" {
		//- response jwt token
		jwtToken := jwtGenerate("User authenticated")
		c.JSON(200, gin.H{"error":false, "jwtToken": jwtToken})
		fmt.Println("=> success")
	} else{
		c.JSON(401, gin.H{"error": true ,"data": nil, "reason": login.ERR_MESSAGE})
		fmt.Println("=> has error")
	}
}

func jwtGenerate(message string) (string) {
	var jwtKey = []byte("qwertyuiopasdfghjkl;'zxcvbnm,!@#$%^&*()")
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &JWTToken{
		Message: message,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	fmt.Println("token string", tokenString)
	return tokenString
}

func handleHealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "health",
	})
}
