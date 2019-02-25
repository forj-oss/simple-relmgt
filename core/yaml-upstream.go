package core


type yamlUpstream struct {
	Name string
	Server string
	RepoPath string `yaml:"repo-path"`
	Protocol string
}