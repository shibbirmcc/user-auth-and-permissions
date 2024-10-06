package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func SendEmail(to, name, password string) error {
	from := os.Getenv("SMTP_USER")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, os.Getenv("SMTP_PASSWORD"), smtpHost)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = "Welcome to ABS Trafikskola"
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += `<!DOCTYPE html>
<html>
<head>
    <title>Welcome to Abs Trafikskola</title>
    <style>
        img {
            display: block;
            max-width: 50px; /* Or any other size you prefer */
            margin: 0 auto;
        }
        header {
            text-align: center;
        }
    </style>
</head>
<body>
    <header>
        <nav>
            <div><img src="https://abstrafikskola.se/img/logo.png" alt="Abs Trafikskola"></div>
        </nav>
    </header>
    <main>
        <div>
            <p><b>Hello ` + name + `</b>,</p>
            <p>&nbsp;We are pleased to have you onboard. As a new user, you might find the following credentials handy to get started.</p>
            <p>&nbsp;Your new password: ` + password + `</p>
            <p>&nbsp;We recommend changing your password after your first login to ensure your account's security.</p>
			<p>&nbsp;Follow this link to login: https://trafikskola.vercel.app/login</p>
        </div>
        <div>
            <p>See you on the road!</p>
            <p><b>Best Regards,<br>ABS Team</b></p>
        </div>
    </main>
</body>
</html>`

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
}

func SendPasswordResetEmail(to, name, password string) error {
	from := os.Getenv("SMTP_USER")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, os.Getenv("SMTP_PASSWORD"), smtpHost)

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = "ABS Password Reset"
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	message += `<!DOCTYPE html>
<html>
<head>
    <title>Password reset complete</title>
    <style>
        img {
            display: block;
            max-width: 50px; /* Or any other size you prefer */
            margin: 0 auto;
        }
        header {
            text-align: center;
        }
    </style>
</head>
<body>
    <header>
        <nav>
            <div><img src="https://abstrafikskola.se/img/logo.png" alt="Abs Trafikskola"></div>
        </nav>
    </header>
    <main>
        <div>
            <p><b>Hello ` + name + `</b>,</p>
            <p>&nbsp;We have created a new password for you.</p>
            <p>&nbsp;Your new password: ` + password + `</p>
            <p>&nbsp;We recommend changing your password after you login to ensure your account's security.</p>
			<p>&nbsp;Follow this link to login: https://trafikskola.vercel.app/login</p>
        </div>
        <div>
            <p>See you on the road!</p>
            <p><b>Best Regards,<br>ABS Team</b></p>
        </div>
    </main>
</body>
</html>`

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
}
