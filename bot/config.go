package bot

import (
	"errors"
	"os"
	"path/filepath"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Provider  string `yaml:"provider"`
	BaseUrl   string `yaml:"base_url"`
	ApiKey    string `yaml:"api_key"`
	ModelName string `yaml:"model"`
}

const configName = "config.yaml"

func ConfigLoad() (*Config, error) {
	// currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	currentPath, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	config := Config{
		Provider:  "<provider>",
		BaseUrl:   "https://llama.iscas.ac.cn/v1/",
		ApiKey:    "<api-key>",
		ModelName: "gpt-4o",
	}
	configPath := filepath.Join(currentPath, configName)
	if _, err := os.Stat(configPath); err != nil {
		data, err := yaml.Marshal(config)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(configPath, data, 0644)
		return nil, errors.New("no config file, created")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
