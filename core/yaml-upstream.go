package core


type YamlUpstream struct {
	Name string
	Server string
	RepoPath string `yaml:"repo-path"`
	Protocol string
}