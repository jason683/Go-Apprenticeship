package functions

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	mail "github.com/xhit/go-simple-mail/v2"
)

//SendEmail works on sending email
func SendEmail(recipient string) {

	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587

	serverUserEmail := os.Getenv("USEREMAIL")
	serverPassword := os.Getenv("USERPW")

	server.Username = serverUserEmail
	server.Password = serverPassword
	server.Encryption = mail.EncryptionTLS

	smtpClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}
	email := mail.NewMSG()
	sender := fmt.Sprintf("From system admin <%s>", serverUserEmail)
	email.SetFrom(sender)
	email.AddTo(recipient)
	email.SetSubject("Notification")

	byteFile, err := ioutil.ReadFile("templates/email.html")
	if err != nil {
		log.Fatal(err)
	}

	HTMLBody := string(byteFile)

	email.SetBody(mail.TextHTML, HTMLBody)

	err = email.Send(smtpClient)
	if err != nil {
		log.Fatal(err)
	}
}
