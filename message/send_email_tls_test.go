/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2022-05-13 10:11:14
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-08 14:16:26
 * @Description:
 */
package message_test

import (
	"testing"

	"github.com/cnjeffliu/core/message"
)

func TestSendEmailTLS(t *testing.T) {
	to := []string{"zhifeng172@163.com" /* , "zhifeng172@163.com" */}
	title := "GPU reset"
	body := `1211905-lmkd -> 10
	43877-lmkd -> 10
	43256-lmkd -> 10
	106424-lmkd -> 10
	3402127-lmkd -> 10`

	message.SendEmailTLS(to, title, body)
}
