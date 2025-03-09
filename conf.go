package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

//go:embed template_config.json
var TemplateConfig string

type Platform struct {
	URL    string `json:"url"`
	APIKey string `json:"api_key"`
	Models map[string]struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"models"`
}

type Config struct {
	DefaultModel      string              `json:"default_model"`
	DefaultTranslator string              `json:"default_translator"`
	SystemMsg         string              `json:"system_msg"`
	Platforms         map[string]Platform `json:"platforms"`
}

// contains full information to call a model
type Model struct {
	Name   string
	Type   string
	Url    string
	APIKey string
}

type ModelList map[string]Model

// 获取配置文件目录
func GetConfDir() string {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("APPDATA") + "\\minall\\"
	case "linux", "freebsd", "openbsd", "netbsd":
		return os.Getenv("HOME") + "/.config/minall/"
	case "darwin":
		return os.Getenv("HOME") + "/Library/Application Support/minall/"
	default: // 对于其他系统，保存至可执行文件所在目录
		exePath, err := os.Executable()
		if err != nil {
			Fatal(err.Error())
		}
		return filepath.Dir(exePath) + "/minall/"
	}
}

// 初始化配置文件
func Init(confdir string) {
	// make sure the directory exists
	if _, err := os.Stat(confdir); os.IsNotExist(err) {
		err := os.MkdirAll(confdir, 0755) // create directory
		if err != nil {
			Fatal(err.Error())
		}
	} else if err != nil { // other errors
		Fatal(err.Error())
	}
	file, err := os.Create(confdir + "config.json") // create config file
	if err != nil {
		Fatal(err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(TemplateConfig)
	if err != nil {
		Fatal(err.Error())
	}
	fmt.Println("Config template saved to", confdir+"config.json")
	fmt.Println("Please make necessary modifications.")
}

// 读取配置文件
func Load(confdir string) *Config {
	var c Config
	if _, err := os.Stat(confdir + "config.json"); os.IsNotExist(err) {
		fmt.Printf("Config file not found. Creating a new one at %s\n", confdir)
		Init(confdir)
		os.Exit(0)
	} else if err != nil {
		Fatal(err.Error())
	}

	file, err := os.Open(confdir + "config.json")
	if err != nil {
		Fatal(fmt.Sprintf("打开%s失败: %s\n", confdir+"config.json", err))
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		Fatal(fmt.Sprintf("解析%s失败: %s\n", confdir+"config.json", err))
	}
	return &c
}

// 生成模型列表
func GetModel(c *Config) ModelList {
	models := make(map[string]Model)
	for _, platform := range c.Platforms {
		for abbr, model := range platform.Models {
			models[abbr] = Model{
				Name:   model.Name,
				Type:   model.Type,
				Url:    platform.URL,
				APIKey: platform.APIKey,
			}
		}
	}
	return models
}

func (m *ModelList) IsValidModel(model string, types []string) bool {
	_, ok := (*m)[model]
	if !ok {
		return false
	}
	return slices.Contains(types, (*m)[model].Type)
}
