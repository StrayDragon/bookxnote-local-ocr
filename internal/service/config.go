package service

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type ConfigOCRUmiOCR struct {
	APIURL string `mapstructure:"api_url"`
}

type ConfigOCR struct {
	UmiOCR *ConfigOCRUmiOCR `mapstructure:"umiocr"`
}

type Config struct {
	OCR *ConfigOCR `mapstructure:"ocr"`
}

var DefaultConfig = Config{
	OCR: &ConfigOCR{
		UmiOCR: &ConfigOCRUmiOCR{
			APIURL: "http://127.0.0.1:1224",
		},
	},
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	// 配置默认值
	v.SetDefault("ocr.umiocr.api_url", DefaultConfig.OCR.UmiOCR.APIURL)

	// 平台优先级获取配置信息
	home, err := os.UserHomeDir()
	if err == nil {
		switch runtime.GOOS {
		case "linux":
			v.AddConfigPath(filepath.Join(home, ".config/bookxnote-local-ocr"))
			v.AddConfigPath(filepath.Join(home, ".local/share/bookxnote-local-ocr"))
		case "darwin":
			v.AddConfigPath(filepath.Join(home, "Library/Application Support/bookxnote-local-ocr"))
		case "windows":
			v.AddConfigPath(filepath.Join(os.Getenv("APPDATA"), "bookxnote-local-ocr"))
		}
	}
	v.AddConfigPath("config")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &config, nil
}
