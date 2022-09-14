/*
 * @Author: Jeffrey Liu
 * @Date: 2021-10-21 15:36:56
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-04-22 12:19:14
 * @Description:
 */
package message

import (
	"net/smtp"
	"strings"
)

func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + ">\r\nSubject: " + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

// func main() {
// 	user := "1057xxxx@qq.com"
// 	password := "sxlvdbtswvmpbcej"
// 	host := "smtp.qq.com:25"
// 	to := "xxxx@qq.com"
// 	subject := "邮件标题tls"

// 	body := `
// 		<html>
// 		<body>
// 		<h3>
// 		"Test send to email"
// 		</h3>
// 		</body>
// 		</html>
// 		`
// 	err := SendToMail(user, password, host, to, subject, body, "html")
// 	if err != nil {
// 		fmt.Println("Send mail error!")
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("Send mail success!")
// 	}
// }
