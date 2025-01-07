package utils

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/lamadev101/ecommerce-api/constant"
	"github.com/lamadev101/ecommerce-api/types"
	"gopkg.in/gomail.v2"
)

// "github.com/sendgrid/sendgrid-go"
// "github.com/sendgrid/sendgrid-go/helpers/mail"

func SendEmail(req types.UserVerification) (types.UserVerification, error) {
	// Create a new message
	m := gomail.NewMessage()
	otp := generateOTP()
	req.Otp = int64(otp)
	m.SetHeader("From", constant.EMAIL_SENDER)                                  // Sender's email
	m.SetHeader("To", req.Email)                                                // Recipient's email
	m.SetHeader("Subject", "OTP for email verification")                        // Email subject
	m.SetBody("text/html", "<strong>Your OTP is: "+fmt.Sprint(otp)+"</strong>") // Email body (HTML format)

	// SMTP server configuration
	d := gomail.NewDialer("smtp.gmail.com", 587, constant.EMAIL_SENDER, os.Getenv("GOOGLE_PASSWORD"))

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Println("Failed to send email:", err)
	} else {
		log.Println("Email sent successfully!")
	}
	return req, nil
}

// func SendEmail(req types.UserVerification) (types.UserVerification, error) {
// 	// Send email
// 	apiKey := os.Getenv("SENDGRID_API_KEY")
// 	if apiKey == "" {
// 		return req, errors.New("SENDGRID_API_KEY not found")
// 	}

// 	// Create a Send Grid client
// 	client := sendgrid.NewSendClient(apiKey)

// 	// Set up the email message
// 	from := mail.NewEmail("Sender Name", constant.EMAIL_SENDER)
// 	to := mail.NewEmail("Recipient Name", req.Email)
// 	subject := "OTP for email verification"
// 	otp := generateOTP()
// 	req.Otp = int64(otp)
// 	htmlContent := "<strong>Your OTP is: " + string(otp) + "</strong>"
// 	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

// 	// Send the email message
// 	_, err := client.Send(message)

// 	fmt.Println("Send email: ", err)
// 	if err != nil {
// 		return req, err
// 	}
// 	return req, nil
// }

func generateOTP() int {
	// Initialize the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Generates a random number between 100000 and 999999
	otp := r.Intn(900000) + 100000
	return otp
}
