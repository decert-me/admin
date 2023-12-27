package config

type Translate struct {
	GithubRepo   string `mapstructure:"github-repo" json:"github-repo" yaml:"github-repo"`       // github 仓库地址
	GithubBranch string `mapstructure:"github-branch" json:"github-branch" yaml:"github-branch"` // github 仓库分支
}
