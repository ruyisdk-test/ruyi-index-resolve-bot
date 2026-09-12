package bot

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	promptd "github.com/ruyisdk-test/ruyi-index-resolve-bot/bot/prompt"
)

func TestMain(m *testing.M) {
	config, err := ConfigLoad()
	if err != nil {
		slog.Error("config error will skip test:", "error", err)
		return
	}
	err = ModelHello(config)
	if err != nil {
		slog.Error("model hello error will skip test:", "error", err)
		return
	}
	os.Exit(m.Run())
}

type testUpstreamVersion struct {
	Webs   []string              `json:"web"`
	Oldv   string                `json:"oldv"`
	Result UpstreamVersionResult `json:"result"`
}

func TestAskUpstreamVersion(t *testing.T) {
	const dataDir = "./testdata/upstream_version"
	if _, err := os.Stat(dataDir); err != nil {
		t.Fatal(err)
	}

	ans, err := filepath.Glob(filepath.Join(dataDir, "ans*.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range ans {
		data, err := os.ReadFile(a)
		if err != nil {
			t.Fatal(err)
		}
		ad := testUpstreamVersion{}
		err = json.Unmarshal(data, &ad)
		if err != nil {
			t.Fatal(err)
		}

		prompt := fmt.Sprintf("%s\n", promptd.UpstreamRawPrompt)
		for _, w := range ad.Webs {
			p, err := os.ReadFile(filepath.Join(dataDir, w))
			if err != nil {
				t.Fatal(err)
			}

			prompt += fmt.Sprintf("```\n%s\n```\n", string(p))
		}

		prompt += promptd.UpstreamOldVersionPrompt + "\n"
		p, err := os.ReadFile(filepath.Join(dataDir, ad.Oldv))
		if err != nil {
			t.Fatal(err)
		}

		prompt += string(p)

		nad, _, err := ModelAskUpstreamVersion(prompt)
		if err != nil {
			t.Fatal(err)
		}

		if len(nad.Version) != len(ad.Result.Version) {
			t.Fatalf("%s want %d results, got %d", a, len(ad.Result.Version), len(nad.Version))
		}
		for i, v := range nad.Version {
			if v != ad.Result.Version[i] {
				t.Fatalf("%s want %s, got %s", a, ad.Result.Version, nad.Version)
			}
		}
	}
}
