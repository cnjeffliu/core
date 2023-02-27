/*
 * @Author: cnzf1
 * @Date: 2021-12-15 16:01:51
 * @LastEditors: cnzf1
 * @LastEditTime: 2023-02-27 17:47:19
 * @Description: viper解析配置文件
 */
package confx

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type TYPE string

const (
	TYPE_JSON = "json"
	TYPE_YAML = "yaml"
	TYPE_TOML = "toml"
	TYPE_INI  = "ini"
)

func check(t TYPE) bool {
	if t == TYPE_JSON {
		return true
	}

	if t == TYPE_YAML {
		return true
	}

	if t == TYPE_TOML {
		return true
	}

	if t == TYPE_INI {
		return true
	}
	return false
}

// Parse support file type:JSON, TOML, YAML, INI
func Parse(target interface{}, fullpath string) {
	dir, file := filepath.Split(fullpath)
	name := filepath.Base(file)
	typo := filepath.Ext(file)

	if !check(TYPE(typo)) {
		return
	}

	confDir, _ := filepath.Abs(path.Dir(dir))
	viper.AddConfigPath(confDir)
	viper.SetConfigName(name)
	viper.SetConfigType(typo)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&target); err != nil {
		panic(err)
	}
}

// ParseStr support content type:JSON, TOML, YAML, INI
func ParseStr(content []byte, typo TYPE, target interface{}) {
	if !check(typo) {
		return
	}

	viper.SetConfigType(string(typo))

	in := bytes.NewReader(content)
	if err := viper.ReadConfig(in); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&target); err != nil {
		panic(err)
	}
}

func OnConfigChanged(f func()) {
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if f != nil {
			f()
		}
	})
	viper.WatchConfig()
}
