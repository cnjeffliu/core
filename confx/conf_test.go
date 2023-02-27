/*
 * @Author: cnzf1
 * @Date: 2021-12-15 16:21:51
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-02-27 17:49:52
 * @Description:
 */

package confx_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cnzf1/gocore/confx"
)

type Env struct {
	Key   string
	Value int
}

type EnvGroup struct {
	Dev Env
	Pro Env
}

type Config struct {
	Root_str string
	Groups   map[string]EnvGroup
	Group    Env
}

const toml = `
root_str="root_string"

[group]
# 测试注释1
key = "driver11"
value = 1000 # 测试注释2

[groups.grp1]
[groups.grp1.dev]
key = "grp1111"
value = 1111

[groups.grp1.pro]
key = "grp2222"
value = 2222

[groups.grp2]
[groups.grp2.dev]
key = 'grp3333'
value=3333
`

func TestTomlFile(t *testing.T) {
	fullpath := "./tmp_file.toml"
	file, _ := os.Create(fullpath)
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	file.Write([]byte(toml))

	var cfg Config
	confx.Parse(&cfg, fullpath)

	fmt.Printf("%v", cfg)

	/*
		output:
		{
			root_string map[
				grp1:{{grp1111 1111} {grp2222 2222}}
				grp2:{{grp3333 3333}
				{ 0}}]
			{ 0}}
	*/
}

func TestTomlBytes(t *testing.T) {
	var cfg Config
	confx.ParseStr([]byte(toml), confx.TYPE_TOML, &cfg)

	fmt.Printf("%#v", cfg)
	/*
		output:
		confx.Config{Root:"root_string", Groups:map[string]confx.EnvGroup{
			"grp1":confx.EnvGroup{
				Dev:confx.Env{Key:"grp1111", Value:1111},
				Pro:confx.Env{Key:"grp2222", Value:2222}},
			"grp2":confx.EnvGroup{
				Dev:confx.Env{Key:"grp3333", Value:3333},
				Pro:confx.Env{Key:"", Value:0}}
		},
		Env:confx.Env{Key:"", Value:0}}
	*/

}
