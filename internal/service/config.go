package service

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/straydragon/bookxnote-local-ocr/internal/common/settings"
)

type ConfigLLM struct {
	Models  []*ConfigLLMModel  `mapstructure:"models"`
	Prompts []*ConfigLLMPrompt `mapstructure:"prompts"`
}

type ConfigLLMPrompt struct {
	Ident  string `mapstructure:"ident"`
	Prompt string `mapstructure:"prompt"`
}

type ConfigLLMModel struct {
	Provider   string `mapstructure:"provider"`
	ApiKey     string `mapstructure:"api_key"`
	ApiBaseUrl string `mapstructure:"api_base_url"`
	Name       string `mapstructure:"name"`
	Ident      string `mapstructure:"ident"`
}

type ConfigTranslate struct {
	Enabled    bool   `mapstructure:"enabled"`
	TargetLang string `mapstructure:"target_language"`
	By         string `mapstructure:"by"`
	UseIdent   string `mapstructure:"use_ident"`
}

type ConfigGenerateByLLM struct {
	Enabled     bool   `mapstructure:"enabled"`
	PromptIdent string `mapstructure:"prompt_ident"`
}

type ConfigAfterOCR struct {
	AutoFixContent *ConfigAutoFixContent `mapstructure:"auto_fix_content"`
	Translate      *ConfigTranslate      `mapstructure:"translate"`
	GenerateByLLM  *ConfigGenerateByLLM  `mapstructure:"generate_by_llm"`
}

type ConfigAutoFixContent struct {
	Enabled bool                        `mapstructure:"enabled"`
	Rules   []*ConfigAutoFixContentRule `mapstructure:"rules"`
}

type ConfigAutoFixContentRule struct {
	Type  string `mapstructure:"type"`
	Ident string `mapstructure:"ident"`
}

type ConfigOCRUmiOCR struct {
	APIURL string `mapstructure:"api_url"`
}

type ConfigOCRCustom struct {
	APIURL string `mapstructure:"api_url"`
	APIKey string `mapstructure:"api_key"`
}

type ConfigOCR struct {
	Selected string           `mapstructure:"selected"`
	UmiOCR   *ConfigOCRUmiOCR `mapstructure:"umiocr"`
	Custom   *ConfigOCRCustom `mapstructure:"custom"`
}

type Config struct {
	OCR      *ConfigOCR      `mapstructure:"ocr"`
	LLM      *ConfigLLM      `mapstructure:"llm"`
	AfterOCR *ConfigAfterOCR `mapstructure:"after_ocr"`
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
		return nil, errors.Wrap(err, "获取用户配置失败")
	}
	return configMgr.GetConfig()
}
