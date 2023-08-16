package config

type Pack struct {
	Path        string `mapstructure:"path" json:"path" yaml:"path"`                         // 项目路径
	PublishPath string `mapstructure:"publish-path" json:"publish-path" yaml:"publish-path"` // 发布项目路径
}
