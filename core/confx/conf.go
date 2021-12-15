/*
 * @Author: Jeffrey.Liu <zhifeng172@163.com>
 * @Date: 2021-12-15 16:01:51
 * @LastEditors: Jeffrey.Liu
 * @LastEditTime: 2021-12-15 18:00:00
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

type params struct {
	name string
	typo string
	path string
}

type ConfOption func(cfg *params)

func WithName(name string) ConfOption {
	return func(cfg *params) {
		cfg.name = name
	}
}

func WithType(typo string) ConfOption {
	return func(cfg *params) {
		cfg.typo = typo
	}
}

func WithPath(path string) ConfOption {
	return func(cfg *params) {
		cfg.path = path
	}
}

// Init support type:JSON, TOML, YAML, INI
func Init(object interface{}, opts ...ConfOption) {
	cfg := params{
		name: "config",
		typo: "toml",
		path: ".",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	confDir, _ := filepath.Abs(path.Dir(cfg.path))
	viper.AddConfigPath(confDir)
	viper.SetConfigType(cfg.typo)
	viper.SetConfigName(cfg.name)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&object); err != nil {
		panic(err)
	}
}

func InitStr(content []byte, object interface{}, opts ...ConfOption) {
	cfg := params{
		typo: "toml",
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	viper.SetConfigType(cfg.typo)

	in := bytes.NewReader(content)
	if err := viper.ReadConfig(in); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&object); err != nil {
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
