package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CreateUser struct {
	Email    string `form:"title" json:"title" binding:"required"`
	Password string `form:"body" json:"body" binding:"required"`
}

var json CreatePost

func main() {

	r := gin.Default()
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"PUT", "PATCH", "POST"},
	// 	AllowHeaders:     []string{"Origin"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == "https://github.com"
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))

	// r.GET("/internal", handleInternalRequest)
	r.POST("/internal", handleInternalRequest)
	r.GET("/healthcheck", handleHealthCheck)
	r.Use(cors.Default())
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func handleInternalRequest(c *gin.Context) {

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	fmt.Println("query params...", c.Request.URL.Query())
	fmt.Println("params...", c.Query("firstName"))
	fmt.Println("raw data...", c.Request.Body)
	// userPoolId := "ap-southeast-1_ihWPvETFg"
	var clientId, poolRegion, email, password string = "5mhq4tjs8nuscdjr1cee65ev9i", "ap-southeast-1", "transybao93@gmail.com", "LvFasCK5"
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
		fmt.Println(err.Error())
		return
	}
	fmt.Println(resp)

	c.JSON(401, gin.H{
		"message": "pong",
	})
}

func handleHealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "health",
	})
}
