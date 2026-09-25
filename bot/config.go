package resolvebot

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ruyisdk-test/ruyi-index-resolve-bot/bot/nvchecker"
	testbot "github.com/ruyisdk-test/ruyi-index-test-bot/bot"
	"go.yaml.in/yaml/v3"
)

type Model struct {
	Provider  string `yaml:"provider"`
	BaseUrl   string `yaml:"base_url"`
	ApiKey    string `yaml:"api_key" json:"-"`
	ModelName string `yaml:"model"`
}

type Github struct {
	Pat string `yaml:"token" json:"-"`
}

type Config struct {
	configPath string `yaml:"-"`

	Server struct {
		ListenAddr        string `yaml:"listen_addr" json:"-"`
		ControlListenAddr string `yaml:"control_addr" json:"-"`
	} `yaml:"server" json:"-"`

	TestBot struct {
		Url    string         `yaml:"control_url" json:"-"`
		Config testbot.Config `yaml:"-" json:"config"`
	} `yaml:"test_bot"`

	Model Model `yaml:"model" json:"-"`

	Github Github `yaml:"github" json:"-"`
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
	config.configPath = filepath.Join(currentPath, configName)
	if _, err := os.Stat(config.configPath); err != nil {
		data, err := yaml.Marshal(config)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(config.configPath, data, 0644)
		return nil, errors.New("no config file, created")
	}

	data, err := os.ReadFile(config.configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	if config.Server.ListenAddr == "" {
		config.Server.ListenAddr = "127.0.0.1:9875"
	}
	if config.Server.ControlListenAddr == "" {
		config.Server.ControlListenAddr = "127.0.0.1:9874"
	}
	slog.Info("service listening on address:", "addr", config.Server.ListenAddr)
	slog.Info("control listening on address:", "addr", config.Server.ControlListenAddr)

	if config.TestBot.Url == "" {
		config.TestBot.Url = "http://127.0.0.1:9876/"
	}
	err = pingTestBot(&config)
	if err != nil {
		return nil, err
	}

	err = pingGithubApi(&config)
	if err != nil {
		return nil, err
	}

	slog.Info("use testbot:", "url", config.TestBot.Url)

	return &config, nil
}

func pingTestBot(config *Config) error {
	u, err := url.JoinPath(config.TestBot.Url, "/config")
	if err != nil {
		return err
	}
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status code from testbot: " + resp.Status)
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

func pingGithubApi(config *Config) error {
	if config.Github.Pat == "" {
		return errors.New("no github pat configured")
	}

	err := nvchecker.InitGithubClient(config.Github.Pat)
	if err != nil {
		return err
	}

	return nvchecker.ListFoxOrgs()
}
