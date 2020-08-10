package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-contrib/cors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type LOGIN struct{
    EMAIL string `json:"email" binding:"required"`
    PASSWORD string `json:"password" binding:"required"`
}

type Claims struct {
	Email string `json:"email"`
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
	var login LOGIN
	c.BindJSON(&login)

	var jwtKey = []byte("qwertyuiopasdfghjkl;'zxcvbnm,!@#$%^&*()")
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: login.EMAIL,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// userPoolId := "ap-southeast-1_ihWPvETFg"
	// var clientId, poolRegion, email, password string = "5mhq4tjs8nuscdjr1cee65ev9i", "ap-southeast-1", "transybao28@gmail.com", "transybao93"
	var clientId, poolRegion, email, password string = "5mhq4tjs8nuscdjr1cee65ev9i", "ap-southeast-1", login.EMAIL, login.PASSWORD
	//- authentication here
	svc := cognitoidentityprovider.New(session.New(), &aws.Config{Region: aws.String(poolRegion)})
	fmt.Println("Created new session of aws services...")
	params := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(email),
			"EMAIL":    aws.String(email),
			"PASSWORD": aws.String(password),
		},
		ClientId: aws.String(clientId),
		// UserPoolId: aws.String(userPoolId),
	}

	resp, err := svc.InitiateAuth(params)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(200, gin.H{"challengeName": resp.ChallengeName, "tokenString": tokenString, "expirationTime": expirationTime})
	}
}

func handleJWT() string {
	return ""
}

func handleHealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "health",
	})
}
