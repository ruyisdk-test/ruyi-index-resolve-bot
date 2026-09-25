package repo

import reporepo "github.com/ruyisdk-test/ruyi-index-test-bot/bot/repo"

var repoMirrors = make(map[string][]string)

func MirrorLoad(repoPath string) error {
	repoConfig, err := reporepo.IndexConfigLoad(repoPath)
	if err != nil {
		return err
	}

	for _, mirror := range repoConfig.Mirrors {
		repoMirrors[mirror.Id] = mirror.Urls
	}

	return nil
}

func MirrorApply(urls []string) ([]string, error) {
	return reporepo.ApplyIndexConfigUrl("", urls, true, repoMirrors)
}
