package helpers

import (
    "os"
    "log"
    "strings"
    "io/ioutil"
    "github.com/xhit/go-simple-mail/v2"
)

func SendEmail(to string, token string) {
    htmlContent, err := ioutil.ReadFile("markup/verification_email.html")
    if err != nil {
        log.Println("Error reading the HTML file:", err)
        return
    }

    htmlContentString := string(htmlContent)
    htmlContentParsed := strings.Replace(htmlContentString, "{% token %}", token, -1)

    from := os.Getenv("EMAIL")
    password := os.Getenv("EMAIL_HOST_PASSWORD")
    smtpHost := os.Getenv("EMAIL_HOST")
    smtpPort := GetEnvInt("EMAIL_PORT")

    server := mail.NewSMTPClient()
    server.Host = smtpHost
    server.Port = smtpPort
    server.Username = from
    server.Password = password
	server.Encryption = mail.EncryptionTLS

    smtpClient, err := server.Connect()
    if err != nil {
        log.Println("Failed to connect to SMTP server:", err)
        return
    }

    email := mail.NewMSG()
    email.SetFrom("From Me <oltlatifi2003@gmail.com>")
    email.AddTo(to)
    email.SetSubject("New Go Email")

    // this is a permanent solution, the token won't be sent on the email
    email.SetBody(mail.TextHTML, htmlContentParsed)

    err = email.Send(smtpClient)
    if err != nil {
        log.Println("Failed to send email:", err)
    }
}
