package cmd

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	sshlogin "github.com/mattkasun/ssh-login"
	"golang.org/x/crypto/ssh"
)

var users map[string]string

func run(p int) {
	users = make(map[string]string)
	router := setupRouter()
	router.Run(fmt.Sprintf("127.0.0.1:%d", p))
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(200, randomString(14))
	})
	r.POST("/login", func(c *gin.Context) {
		var login sshlogin.Login
		if err := c.ShouldBindJSON(&login); err != nil {
			log.Println("login ", err)
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(users)
		pub, ok := users[login.User]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(pub))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := pubKey.Verify([]byte(login.Message), &login.Sig); err != nil {
			log.Println("login verify ", err)
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Hello World"})
	})
	r.POST("/register", func(c *gin.Context) {
		var reg sshlogin.Registation
		if err := c.ShouldBindJSON(&reg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		users[reg.User] = reg.Key
		c.String(http.StatusOK, "registration successfull")
	})
	return r
}

func randomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("randomString", err)
	}
	return base32.StdEncoding.EncodeToString(b)[:n]
}
