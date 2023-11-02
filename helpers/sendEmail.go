package helpers

import (
    "os"
    "log"
    "github.com/xhit/go-simple-mail/v2"
)

var htmlBody = `
<html>
<head>
   <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
   <title>Hello, World</title>
</head>
<body>
   <p>This is an email using Go</p>
</body>
`

func SendEmail(to string, token string) {
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
    email.SetBody(mail.TextHTML, htmlBody + token)

    err = email.Send(smtpClient)
    if err != nil {
        log.Println("Failed to send email:", err)
    }
}
