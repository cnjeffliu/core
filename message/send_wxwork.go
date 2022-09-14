/*
 * @Author: Jeffrey Liu
 * @Date: 2020-10-21 15:36:56
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-09-08 14:55:29
 * @Description: 发送企业微信消息通知
 */

package message

import (
	"fmt"

	"github.com/cnjeffliu/gocore/httpx"
)

type WxWorkClient struct {
	corpid     string // 企业ID
	corpsecret string // 企业凭证密钥
	agentID    string // 企业应用ID
}

func NewWxWorkClient(corpid, corpsecret, agentID string) *WxWorkClient {
	return &WxWorkClient{
		corpid,
		corpsecret,
		agentID,
	}
}

type WxWorkOption func(map[string]interface{})

// 指定接收消息的成员，成员ID列表（多个接收者用‘|’分隔，最多支持1000个）。
// 特殊情况：指定为"@all"，则向该企业应用的全部成员发送
func WithUser(users string) WxWorkOption {
	return func(m map[string]interface{}) {
		m["touser"] = users
	}
}

// 指定接收消息的部门，部门ID列表，多个接收者用‘|’分隔，最多支持100个。
// 当touser为"@all"时忽略本参数
func WithParty(users string) WxWorkOption {
	return func(m map[string]interface{}) {
		m["toparty"] = users
	}
}

// 指定接收消息的标签，标签ID列表，多个接收者用‘|’分隔，最多支持100个。
// 当touser为"@all"时忽略本参数
func WithTag(users string) WxWorkOption {
	return func(m map[string]interface{}) {
		m["totag"] = users
	}
}

// 文本卡片消息
func (c *WxWorkClient) SendCardToUser(title, content, toUrl string, opts ...WxWorkOption) error {
	// btntxt := "详情" //可自定义卡片底下的文字

	// 卡片消息，可以点击详情跳转到URL
	req := map[string]interface{}{
		// "touser": "liuzhifeng",
		// "toparty" : "PartyID1|PartyID2",
		// "totag" : "TagID1 | TagID2",
		"msgtype": "textcard",
		"agentid": c.agentID,
		"textcard": map[string]interface{}{
			"title":       title,
			"description": content,
			"url":         toUrl,
			// "btntext":     btntxt,
		},
	}

	for _, opt := range opts {
		opt(req)
	}

	return sendMsg(c.corpid, c.corpsecret, req)
}

// 文本消息
func (c *WxWorkClient) SendTextToUser(content string, opts ...WxWorkOption) error {
	req := map[string]interface{}{
		// "touser": "liuzhifeng",
		// 	"toparty" : "PartyID1|PartyID2",
		//  "totag" : "TagID1 | TagID2",
		"msgtype": "text",
		"agentid": c.agentID,
		"text": map[string]interface{}{
			"content": content,
		},
	}
	for _, opt := range opts {
		opt(req)
	}

	return sendMsg(c.corpid, c.corpsecret, req)
}

// markdown消息格式
func (c *WxWorkClient) SendMarkdownToUser(content string, opts ...WxWorkOption) error {
	req := map[string]interface{}{
		// "touser": "liuzhifeng",
		// 	"toparty" : "PartyID1|PartyID2",
		//  "totag" : "TagID1 | TagID2",
		"msgtype": "markdown",
		"agentid": c.agentID,
		"markdown": map[string]interface{}{
			"content": content,
		},
	}

	for _, opt := range opts {
		opt(req)
	}

	return sendMsg(c.corpid, c.corpsecret, req)
}

func getToken(corpid, secret string) (string, error) {
	qyurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpid, secret)
	data, err := httpx.Client().GetToMap(qyurl)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	errcode := data["errcode"].(float64)
	if errcode != 0 {
		fmt.Println(errcode)
		return "", nil
	}
	access_token := data["access_token"]
	return access_token.(string), nil
}

func sendMsg(corpid, secret string, req map[string]interface{}) error {
	access_token, err := getToken(corpid, secret)
	if err != nil {
		return err
	}

	sendurl := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", access_token)
	_, err = httpx.Client().PostToMap(sendurl, req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
