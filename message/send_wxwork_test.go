/*
 * @Author: Jeffrey Liu <zhifeng172@163.com>
 * @Date: 2021-10-21 15:36:56
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-13 10:41:04
 * @Description:
 */
package message_test

import (
	"testing"

	"github.com/cnjeffliu/core/message"
)

func TestSendWxWrk(t *testing.T) {
	corpid := "ww412ca4ab5f99xxxx"                          //企业ID
	secret := "Fw4-WhdQj9KF_S9nfD1k44qhmGOgJOKb8HRGJeJxxxx" //企业凭证密钥
	agentId := "1000000"                                    //应用ID

	client := message.NewWxWorkClient(corpid, secret, agentId)

	// content := fmt.Sprintf("<a>title</a>\ncontent")
	// client.SendTextToUser(content, message.WithUser("liuzhifeng"))
	// client.SendTextToUser("系统通知测试=》成员", message.WithUser("liuzhifeng"))
	// client.SendTextToUser("系统通知测试=》标签", WithTag("2"))
	// client.SendTextToUser("系统通知测试=》部门", WithParty("3))

	// client.SendCardToUser("系统通知测试", "test", "www.baidu.co")

	// client.SendMarkdownToUser("test", "2021-036.158.255.16036.158.255.16036.158.255.160")
	content := `**this is title**
this is content
`
	client.SendMarkdownToUser(content, message.WithUser("liuzhifeng"))
}
