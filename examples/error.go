/**
 * @Author: Jeffrey.Liu
 * @Date: 2021-11-17 15:37:27
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2021-11-17 15:37:27
 * @Description: 错误码统一处理函数
 */
package demo

import (
	"fmt"
	"gitee.com/cnjeffliu/core/errorx"
)

type ErrMsg struct {
	code  int64
	msg   string
	msgCN string
	err   error // 往往存新加入的error
}

type ErrLang uint8

const (
	ErrLangCN ErrLang = iota
	ErrLangEN
)

var StatusOK = NewErrMsg(200, "success", "成功")

func NewErrMsg(code int64, msg, msgCN string) *ErrMsg {
	return &ErrMsg{code, msg, msgCN, nil}
}

func (p *ErrMsg) ErrCode() int64 {
	return p.code
}

func (p *ErrMsg) Error() string {
	return String(p, ErrLangEN)
}

func String(err error, lang ErrLang) string {
	var e *ErrMsg
	var ok bool
	if e, ok = err.(*ErrMsg); !ok {
		// 非ErrMsg错误
		fmt.Println("jeff ", err.Error())
		return err.Error()
	}

	errmsg := ""
	if lang == ErrLangCN {
		errmsg = e.msgCN
	} else {
		errmsg = e.msg
	}

	if e.err != nil {
		if len(errmsg) > 0 {
			errmsg += ":"
		}

		errmsg += e.err.Error()
	}
	return errmsg
}

// 输出最表层的错误信息,只有最外层的会定义中英文
func Cause(err error, lang ErrLang) string {
	errmsg := ""
	if e, ok := err.(*ErrMsg); ok {
		if lang == ErrLangCN {
			errmsg = e.msgCN
		} else {
			errmsg = e.msg
		}

		if len(errmsg) > 0 {
			return errmsg
		}

		if e.err != nil {
			return e.err.Error()
		}
	}

	return err.Error()
}

/**
 * @description: 检查错误码是否代表处理成功
 * @param {error} err
 * @return {*}
 */
func CheckErrOK(err error) bool {
	if e, ok := err.(*ErrMsg); ok {
		return e.ErrCode() == StatusOK.ErrCode()
	}
	return err == nil
}

/**
 * @description: 保留入参错误信息+原错误信息
 * @param {error} err
 * @return {*}
 */
func (p *ErrMsg) WithErr(err error) *ErrMsg {
	return &ErrMsg{
		code:  p.code,
		msg:   p.msg,
		msgCN: p.msgCN,
		err:   errorx.WithMessage(err, ""),
	}
}

/**
 * @description: 保留入参错误信息+原错误信息
 * @param {error} err
 * @return {*}
 */
func (p *ErrMsg) WithMsg(err error, msg string) *ErrMsg {
	return &ErrMsg{
		code:  p.code,
		msg:   p.msg,
		msgCN: p.msgCN,
		err:   errorx.WithMessage(err, msg),
	}
}
