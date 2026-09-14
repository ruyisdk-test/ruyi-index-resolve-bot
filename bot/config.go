package resolvebot

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	testbot "github.com/ruyisdk-test/ruyi-index-test-bot/bot"
	"go.yaml.in/yaml/v3"
)

type Model struct {
	Provider  string `yaml:"provider"`
	BaseUrl   string `yaml:"base_url"`
	ApiKey    string `yaml:"api_key"`
	ModelName string `yaml:"model"`
}

type Config struct {
	TestBot struct {
		Addr   string         `yaml:"control_addr" json:"-"`
		Config testbot.Config `yaml:"-" json:"config"`
	} `yaml:"test_bot"`

	Model Model `yaml:"model"`
}

const configName = "config.yaml"

func ConfigLoad() (*Config, error) {
	// currentPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	currentPath, err := filepath.Abs(".")
	if err != nil {
		return nil, err
	}

	config := Config{
		Model: Model{
			Provider:  "<provider>",
			BaseUrl:   "https://llama.iscas.ac.cn/v1/",
			ApiKey:    "<api-key>",
			ModelName: "gpt-4o",
		},
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

	if config.TestBot.Addr == "" {
		return nil, errors.New("no testbot addr provided")
	}
	err = pingTestBot(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func pingTestBot(config *Config) error {
	u, err := url.JoinPath(config.TestBot.Addr, "/config")
	if err != nil {
		return err
	}
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("bad statuc code from testbot: " + resp.Status)
	}

	buf := make([]byte, resp.ContentLength)
	_, err = resp.Body.Read(buf)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && err != io.EOF {
		return err
	}

	err = json.Unmarshal(buf, &(config.TestBot))
	if err != nil {
		return err
	}

	return nil
}
