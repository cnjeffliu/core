package execx

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"syscall"
)

type CMDOption func(*exec.Cmd)

func WithDir(d string) CMDOption {
	return func(c *exec.Cmd) {
		c.Dir = d
	}
}

func ExecCMD(commond string, opts ...CMDOption) (out string, code int, err error) {
	var cmd *exec.Cmd
	goos := runtime.GOOS

	switch goos {
	case "linux", "darwin":
		cmd = exec.Command("sh", "-c", commond)
	case "windows":
		cmd = exec.Command("cmd.exe", "/c", commond)
	default:
		return "", -1, errors.New("unsupported os")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return out, -1, err
	}

	for _, opt := range opts {
		opt(cmd)
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return out, -1, err
	}

	ret, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println(err)
		return out, -1, err
	}
	out = string(ret)

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
		if ex, ok := err.(*exec.ExitError); ok {
			//获取命令执行返回状态，相当于shell: echo $?
			code = ex.Sys().(syscall.WaitStatus).ExitStatus()
			return out, code, err
		} else {
			return out, -1, err
		}
	}

	return out, 0, nil
}
