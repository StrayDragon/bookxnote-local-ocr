package service

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"github.com/straydragon/bookxnote-local-ocr/internal/common/settings"
	"gopkg.in/yaml.v3"
)

type ConfigManager struct {
	viper      *viper.Viper
	configPath string
	yamlNode   *yaml.Node
}

func NewConfigManager() (*ConfigManager, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	for _, dir := range settings.GetUserConfigDirs() {
		v.AddConfigPath(dir)
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	configPath := v.ConfigFileUsed()

	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var node yaml.Node
	if err := yaml.Unmarshal(content, &node); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	if node.Kind != yaml.DocumentNode || len(node.Content) == 0 {
		return nil, fmt.Errorf("invalid YAML structure")
	}

	return &ConfigManager{
		viper:      v,
		configPath: configPath,
		yamlNode:   &node,
	}, nil
}

func updateYAMLNode(node *yaml.Node, path []string, value interface{}) error {
	if len(path) == 0 {
		var valueNode yaml.Node
		if err := valueNode.Encode(value); err != nil {
			return fmt.Errorf("failed to encode value: %w", err)
		}
		if node.Style > 0 {
			valueNode.Style = node.Style
		}
		*node = valueNode
		return nil
	}

	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("expected mapping node")
	}

	key := path[0]
	restPath := path[1:]

	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return updateYAMLNode(node.Content[i+1], restPath, value)
		}
	}

	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!str",
		Value: key,
	}
	valueNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Tag:  "!!map",
	}
	node.Content = append(node.Content, keyNode, valueNode)
	return updateYAMLNode(valueNode, restPath, value)
}

func (cm *ConfigManager) Get(key string) interface{} {
	return cm.viper.Get(key)
}

func (cm *ConfigManager) Set(key string, value interface{}) error {
	cm.viper.Set(key, value)

	rootNode := cm.yamlNode.Content[0]

	path := strings.Split(key, ".")
	if err := updateYAMLNode(rootNode, path, value); err != nil {
		return fmt.Errorf("failed to update YAML: %w", err)
	}

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(cm.yamlNode); err != nil {
		return fmt.Errorf("failed to encode YAML: %w", err)
	}

	if err := os.WriteFile(cm.configPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (cm *ConfigManager) GetConfig() (*Config, error) {
	var config Config
	if err := cm.viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &config, nil
}
