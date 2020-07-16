package common

import (
	"bytes"
	"server/model"
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/spf13/viper"
	"html/template"
	"log"
)
var(
	host string
	port int
	email string
	password string
)

func InitSmtp(){
	host=viper.GetString("smtp.host")
	port=viper.GetInt("smtp.port")
	email=viper.GetString("smtp.email")
	password=viper.GetString("smtp.password")
}

func SendCaptcha(toEmail string)(string,error){

	m:=gomail.NewMessage()
	m.SetAddressHeader("From", email /*"发件人地址"*/, "clf") // 发件人
	m.SetHeader("To", m.FormatAddress(toEmail, "收件人")) // 收件人
	m.SetHeader("Subject", "系统验证码")
	var buf bytes.Buffer
	t, err := template.ParseFiles("public/email/captcha.html")
	if err != nil {
		fmt.Println("parse file err:", err)
		return "",err
	}
	c:=Captcha(6)
	if err := t.Execute(&buf, c); err != nil {
		fmt.Println("There was an error:", err.Error())
	}
	m.SetBody("text/html",buf.String()) // 正文
	d := gomail.NewPlainDialer(host, port, email, password) // 发送邮件服务器、端口、发件人账号、发件人密码
	if err := d.DialAndSend(m); err != nil {
		log.Println("发送失败", err)
		return "",err
	}
	return  c,nil
	//header   :=  make(map[string]string)
	//header["From"] = "test"+"<" +email+">"
	//header["To"] = toEmail
	//header["Subject"] = "dgut_iot_app"
	//header["Content-Type"] = "text/html;charset=UTF-8"
	//text := "验证码"
	//Captcha:=Captcha(6)
	////邮件内容
	//body := `
    //<html>
    //<body>
    //<h3>
    //"dgut_iot"` + text + `
    //</h3>
	//<h2>
	//`+Captcha+`
	//</h2>
    //</body>
    //</html>
    //`
	//message := ""
	//for k,v :=range header{
	//	message  += fmt.Sprintf("%s:%s\r\n",k,v)
	//}
	//message +="\r\n"+body
	//auth :=smtp.PlainAuth(
	//	"",
	//	email,
	//	password,
	//	host,
	//)
	//err := SendMailUsingTLS(
	//	fmt.Sprintf("%s:%d", host, port),
	//	auth,
	//	email,
	//	toEmail,
	//	[]byte(message),
	//)
	//if err  !=nil {
	//	//fmt.Println("发送邮件失败!")
	//	panic(err)
	//}
	//return Captcha,err
}


func SendConfirmAdmin(admin model.Admin) error{
	m:=gomail.NewMessage()
	m.SetAddressHeader("From", email /*"发件人地址"*/, "clf") // 发件人
	m.SetHeader("To", m.FormatAddress(admin.Email, "收件人")) // 收件人
	m.SetHeader("Subject", "系统验证码")


	var buf bytes.Buffer
	t, err := template.ParseFiles("public/email/confirmAdmin.html")
	if err != nil {
		fmt.Println("parse file err:", err)
		return err
	}
	type p struct {
		Url string
		NickName string
		Email string
	}

	a:=p{
		Url:"http://127.0.0.1:8888/api/admin/confirm",
		NickName: admin.NickName,
		Email: admin.Email,
	}
	if err := t.Execute(&buf, a); err != nil {
		fmt.Println("There was an error:", err.Error())
	}
	m.SetBody("text/html",buf.String()) // 正文
	d := gomail.NewPlainDialer(host, port, email, password) // 发送邮件服务器、端口、发件人账号、发件人密码
	if err := d.DialAndSend(m); err != nil {
		log.Println("发送失败", err)
		return err
	}
	return  nil

}