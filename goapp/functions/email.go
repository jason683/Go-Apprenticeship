package functions

import (
	"bytes"
	"fmt"
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

	random := "This is to inform you that you have an outstanding action to be completed on the contract management system. The system link is provided below."
	b := new(bytes.Buffer)
	Tpl.ExecuteTemplate(b, "email.html", random) //should insert 3rd argument as myUser for email list

	//HTMLBody := string(byteFile)
	HTMLBody := b.String()
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
