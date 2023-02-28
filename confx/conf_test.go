/*
 * @Author: cnzf1
 * @Date: 2021-12-15 16:21:51
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-02-28 15:42:57
 * @Description:
 */

package confx_test

import (
	"fmt"
	"io/ioutil"
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

func TestTomlFile(t *testing.T) {
	fullpath := "./config.toml"

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
	fullpath := "./config.toml"
	data, _ := ioutil.ReadFile(fullpath)

	var cfg Config
	confx.ParseStr([]byte(data), confx.TYPE_TOML, &cfg)

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
