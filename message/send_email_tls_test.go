/*
 * @Author: Jeffrey Liu
 * @Date: 2022-05-13 10:11:14
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-14 18:09:15
 * @Description:
 */
package message_test

import (
	"testing"

	"github.com/cnjeffliu/gocore/message"
)

func TestSendEmailTLS(t *testing.T) {
	host := "smtp.163.com"
	port := 465
	from := "user_send_from@163.com"
	pwd := "this_is_from_user_pwd"

	client := message.NewEmailClient(host, port, from, pwd)

	to := []string{"send_to_1@163.com" /* , "send_to_2@163.com" */}
	title := "this is title"
	body := `1211905-lmkd -> 10
	43877-lmkd -> 10
	43256-lmkd -> 10
	106424-lmkd -> 10
	3402127-lmkd -> 10`

	client.SendEmailTLS(to, title, body)
}
