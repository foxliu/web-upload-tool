package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type Claims struct {
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 24 * 365  // 1 year
var MySecrt = GenSecrt()
var tokenString, _ = GenToken()
// GenSecrt Generate secrt key
func GenSecrt() []byte  {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("Generate random secrt key error: %s\n", err.Error())
		os.Exit(1)
	}
	return b
}
// GenToken Generate JWT
func GenToken() (string, error)  {
	c := Claims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenS, err :=  token.SignedString(MySecrt)
	if err != nil {
		fmt.Printf("Generate token error: %v\nExit and byebye\n", err.Error())
		os.Exit(1)
	}
	return tokenS, nil
}
// ParseToken Parse JWT
func ParseToken(tokenString string) (*Claims, error)  {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecrt, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTMiddleware() func(context *gin.Context) {
	return func(context *gin.Context) {
		authHeader := context.Request.Header.Get("Authorization")
		if authHeader == "" {
			context.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg": "Request header Authorization is nil",
			})
			context.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			context.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg": "Request header Authorizations error",
			})
			context.Abort()
			return
		}
		_, err := ParseToken(parts[1])
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg": "Invalid Token",
			})
			context.Abort()
			return
		}
		context.Next()
	}
}

func main() {
	args := os.Args
	var port string
	if len(args) <= 1 {
		port = "8030"
	} else {
		port = args[1]
	}
	gin.SetMode("release")
	router := gin.Default()
	router.Use(JWTMiddleware())
	router.MaxMultipartMemory = 1024 << 20
	router.POST("/", func(context *gin.Context) {
		file, err := context.FormFile("file")
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"msg": err.Error()})
		} else {
			dst := path.Join("./", file.Filename)
			_ = context.SaveUploadedFile(file, dst)
			fmt.Printf("%s | ", file.Filename)
			context.JSON(http.StatusOK, gin.H{"msg": "upload success"})
		}
	})
	fmt.Printf("Use curl like this to upload file:\ncurl --location --request POST 'localhost:%s/' " +
		"--header 'Authorization: Bearer %s' --form 'file=@\"/your-file-path\"'\n", port, tokenString)
	err := router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		fmt.Printf("Start failed: %v", err.Error())
	}
}