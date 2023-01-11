/*
 * @Author: cnzf1
 * @Date: 2022-08-10 17:17:32
 * @LastEditors: cnzf1
 * @LastEditTime: 2022-09-08 10:43:45
 * @Description:
 */
package httpx_test

import (
	"testing"

	"github.com/cnzf1/gocore/httpx"
)

func TestGet(t *testing.T) {
	result, err := httpx.Client().Get("http://www.baidu.com", httpx.WithJSONContent())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(result))
}
