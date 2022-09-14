/*
 * @Author: Jeffrey Liu
 * @Date: 2022-05-13 10:11:14
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-14 10:35:38
 * @Description:
 */
package message_test

import (
	"testing"

	"github.com/cnjeffliu/gocore/message"
)

func TestSendEmailTLS(t *testing.T) {
	to := []string{"user1@163.com" /* , "user2@163.com" */}
	title := "this is title"
	body := `1211905-lmkd -> 10
	43877-lmkd -> 10
	43256-lmkd -> 10
	106424-lmkd -> 10
	3402127-lmkd -> 10`

	message.SendEmailTLS(to, title, body)
}
