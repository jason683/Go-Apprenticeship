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

	email := mail.NewMSG()
	sender := fmt.Sprintf("From system admin <%s>", serverUserEmail)
	email.SetFrom(sender)
	email.AddTo(recipient)
	email.SetSubject("Notification")

	byteFile, err := ioutil.ReadFile("templates/email.html")
	if err != nil {
		log.Fatal(err)
	}
	// random := "testing 123"
	// b := new(bytes.Buffer)
	// Tpl.ExecuteTemplate(b, "templates/email.html", random) //should insert 3rd argument as myUser for email list

	HTMLBody := string(byteFile)
	//HTMLBody := b.String()
	email.SetBody(mail.TextHTML, HTMLBody)

	SMTPClient, err := server.Connect()
	if err != nil {
		log.Fatal(err)
	}
	err = email.Send(SMTPClient)
	if err != nil {
		log.Fatal(err)
	}
}
