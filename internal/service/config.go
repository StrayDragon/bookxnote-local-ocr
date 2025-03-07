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

// GetUserConfigCtrl 获取用户配置控制器
func GetUserConfigCtrl() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	for _, dir := range settings.GetUserConfigDirs() {
		v.AddConfigPath(dir)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	return v, nil

}

// GetUserConfig 获取用户配置
func GetUserConfig() (*Config, error) {
	configMgr, err := NewConfigManager()
	if err != nil {
		return nil, fmt.Errorf("Failed to load config: %v, using default settings", err)
	}
	return configMgr.GetConfig()
}
