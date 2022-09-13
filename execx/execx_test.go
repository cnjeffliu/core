package execx

import "testing"

func TestExecCMD(t *testing.T) {

	cmd := "ls"
	out, code, err := ExecCMD(cmd, WithWorkDir("/"))
	t.Log("code:", code)
	if err != nil {
		t.Error("msg:", err.Error())
	}
	t.Log("data:", out)
}
