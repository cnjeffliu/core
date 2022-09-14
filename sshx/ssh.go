/*
 * @Author: Jeffrey.Liu
 * @Date: 2022-01-06 10:21:36
 * @LastEditors: Jeffrey Liu
 * @LastEditTime: 2022-08-04 15:42:03
 * @Description: ssh远程执行命令工具
 */
package sshx

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHCli struct {
	cli    *ssh.Client
	user   string
	passwd string
	in     chan<- string //mux_shell
	out    chan<- string // mux_shell
}

/**
 * @description:  需要调用Close关闭client
 * @param {string} addr
 * @param {string} user
 * @param {string} passwd
 * @return {*}
 */
func NewSSHCli(addr string, user string, passwd string) *SSHCli {
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		fmt.Printf("SSH dial error: %s", err.Error())
	}

	return &SSHCli{
		cli:    client,
		user:   user,
		passwd: passwd,
	}
}

/**
 * @description:  需要调用Close关闭client
 * @param {string} addr
 * @param {string} user
 * @param {string} passwd
 * @return {*}
 */
func NewSSHCliWithError(addr string, user string, passwd string) (*SSHCli, error) {
	client, err := ssh.Dial("tcp", addr, &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(passwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return nil, err
	}

	return &SSHCli{
		cli:    client,
		user:   user,
		passwd: passwd,
	}, nil
}

func (s *SSHCli) Run(cmd string) {
	sess, err := s.cli.NewSession()
	if err != nil {
		fmt.Printf("new session error: %s", err.Error())
	}
	defer sess.Close()

	if err := sess.Run(cmd); err != nil {
		fmt.Print("Failed to run: " + err.Error())
	}
}

func (s *SSHCli) Output(cmd string) (output []byte) {
	sess, err := s.cli.NewSession()
	if err != nil {
		fmt.Printf("new session error: %s", err.Error())
	}
	defer sess.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = sess.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return nil
	}

	in, err := sess.StdinPipe()
	if err != nil {
		fmt.Print(err)
	}

	out, err := sess.StdoutPipe()
	if err != nil {
		fmt.Print(err)
	}

	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)
			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)
			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(s.passwd + "\n"))
				if err != nil {
					break
				}
			}
		}
	}(in, out, &output)

	_, err = sess.Output(cmd)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Failed to run command, Err:%s", err.Error())
		return nil
	}

	return output
}

func (s *SSHCli) Pipe() {
	sess, err := s.cli.NewSession()
	if err != nil {
		fmt.Printf("new session error: %s", err.Error())
	}
	defer sess.Close()

	sess.Stdout = os.Stdout // 会话输出关联到系统标准输出设备
	sess.Stderr = os.Stderr // 会话错误输出关联到系统标准错误输出设备
	sess.Stdin = os.Stdin   // 会话输入关联到系统标准输入设备

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
	}

	if err := sess.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Printf("request pty error: %s", err.Error())
	}

	if err := sess.Shell(); err != nil {
		fmt.Printf("start shell error: %s", err.Error())
	}

	if err := sess.Wait(); err != nil {
		fmt.Printf("return error: %s", err.Error())
	}
}

func (s *SSHCli) Close() {
	s.cli.Close()
}

/**
 * @description: // todo 还需要改造
 * @param {*}
 * @return {*}
 */
func (s *SSHCli) NewMultiSess() {
	sess, err := s.cli.NewSession()
	if err != nil {
		fmt.Printf("unable to create sess: %s", err)
	}
	defer sess.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
	}

	if err := sess.RequestPty("xterm", 80, 40, modes); err != nil {
		fmt.Print(err)
	}

	w, err := sess.StdinPipe()
	if err != nil {
		fmt.Print(err)
	}

	r, err := sess.StdoutPipe()
	if err != nil {
		fmt.Print(err)
	}

	in, out := muxshell(w, r)
	if err := sess.Start("/bin/sh"); err != nil {
		fmt.Print(err)
	}

	<-out //ignore the shell output

	in <- "ls -lhav"
	fmt.Printf("ls output: %s", <-out)

	in <- "whoami"
	fmt.Printf("whoami: %s", <-out)

	in <- "exit"
	sess.Wait()
}

func muxshell(w io.Writer, r io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 1)
	out := make(chan string, 1)

	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + ""))
			wg.Wait()
		}
	}()

	go func() {
		var (
			buf [65 * 1024]byte
			t   int
		)

		for {
			n, err := r.Read(buf[t:])
			if err != nil {
				close(in)
				close(out)
				return
			}

			t += n
			if buf[t-2] == '$' { //assuming the $PS1 == 'sh-4.3$ '
				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()

	return in, out
}
