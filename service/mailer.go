package service

import (
	"log"

	"github.com/flosch/pongo2"
	"github.com/smartwalle/pongo2render"
	"gopkg.in/gomail.v2"
)

type mailer struct {
	*gomail.Dialer
}

func InitMailer() *mailer {
	return &mailer{
		Dialer: gomail.NewDialer("smtp.exmail.qq.com", 465, "noreplay@yitum.com", "Yt@888888"),
	}
}

type email struct {
	tmpl    string
	subject string
}

var RegisterEmail = email{
	tmpl:    "register.tmpl",
	subject: "注册成功",
}

func (m *mailer) Send(e email, to string, data map[string]interface{}) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "noreplay@yitum.com")
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", e.subject)
	ctx := pongo2.Context(data)
	html, err := parseTpl(e.tmpl, ctx)
	if err != nil {
		log.Println("parse tmpl fail,", err)
		return err
	}
	msg.SetBody("text/html", html)
	// Send the email to Bob, Cora and Dan.
	if err := m.Dialer.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}

func parseTpl(filename string, ctx pongo2.Context) (html string, err error) {
	var render = pongo2render.NewRender("./conf")
	html, err = render.Template(filename).Execute(ctx)
	if err != nil {
		return
	}
	return
}
