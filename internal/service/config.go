package service

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/straydragon/bookxnote-local-ocr/internal/common/settings"
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

	// 添加配置文件搜索路径
	for _, dir := range settings.GetUserConfigDirs() {
		v.AddConfigPath(dir)
	}

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
