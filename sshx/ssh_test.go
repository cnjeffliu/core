/*
 * @Author: Jeffrey.Liu
 * @Date: 2022-01-06 10:29:57
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2022-01-06 11:31:44
 * @Description:
 */
package sshx

import (
	"fmt"
	"testing"
)

func TestSSH(t *testing.T) {
	cli := NewSSHCli("10.0.2.113:22", "motech", "123456")
	defer cli.Close()

	cli.Run("pwd")

	// fmt.Println(string(cli.Output("pwd")))
	fmt.Println(string(cli.Output("ps -ef |head")))
}
