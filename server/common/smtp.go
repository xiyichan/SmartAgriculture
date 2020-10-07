package common

import (
	"bytes"
	"github.com/go-gomail/gomail"
	"github.com/spf13/viper"
	"html/template"
)

var (
	host     string
	port     int
	email    string
	password string
)

func InitSmtp() {
	host = viper.GetString("smtp.host")
	port = viper.GetInt("smtp.port")
	email = viper.GetString("smtp.email")
	password = viper.GetString("smtp.password")
}

func SendCaptcha(toEmail string) (string, error) {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", email /*"发件人地址"*/, "不简单公司") // 发件人
	m.SetHeader("To", m.FormatAddress(toEmail, "收件人"))     // 收件人
	m.SetHeader("Subject", "验证码")
	var buf bytes.Buffer
	t, err := template.ParseFiles("public/email/captcha.html")
	if err != nil {
		return "", err
	}
	c := Captcha(6)
	if err := t.Execute(&buf, c); err != nil {
		return "", err
	}
	m.SetBody("text/html", buf.String())               // 正文
	d := gomail.NewDialer(host, port, email, password) // 发送邮件服务器、端口、发件人账号、发件人密码
	d.SSL = true
	if err := d.DialAndSend(m); err != nil {
		return "", err
	}
	return c, nil
}
