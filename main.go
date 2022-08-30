package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	gomail "gopkg.in/mail.v2"
)

func sendEmail(sender string, recipients []string, subject string, message string) error {
	var mailerToGoHost = os.Getenv("MAILERTOGO_SMTP_HOST")
	var mailerToGoPort, _ = strconv.ParseInt(os.Getenv("MAILERTOGO_SMTP_PORT"), 10, 64)
	var mailerToGoUser = os.Getenv("MAILERTOGO_SMTP_USER")
	var mailerToGoPassword = os.Getenv("MAILERTOGO_SMTP_PASSWORD")
	var mailerToGoDomain = os.Getenv("MAILERTOGO_DOMAIN")
	smtp_message := gomail.NewMessage()

	// Set headers.
	smtp_message.SetHeader("From", sender+"@"+mailerToGoDomain)
	smtp_message.SetHeader("To", strings.Join(recipients[:], ","))
	smtp_message.SetHeader("Subject", subject)
	// Set html body.
	smtp_message.SetBody("text/html", message)
	// Connect to SMTP server.
	smtp_dialer := gomail.NewDialer(mailerToGoHost, int(mailerToGoPort), mailerToGoUser, mailerToGoPassword)
	// Send EMail.
	err := smtp_dialer.DialAndSend(smtp_message)
	return err
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	var mailerToGoDomain = os.Getenv("MAILERTOGO_DOMAIN")
	var message = "The Server is Starting... "
	err := sendEmail("paulsimonk2", []string{"admins@" + mailerToGoDomain, "paulsimonk2@gmail.com"}, "Server", message)
	if err != nil {
		log.Fatal(err)
	}
	router.Run(":" + port)
}
