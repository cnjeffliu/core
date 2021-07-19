package httpx

import "testing"

func TestGet(t *testing.T) {
	result, err := NewHttpxCli().Get("http://www.baidu.com", WithJSONContent())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(result))
}
