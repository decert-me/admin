package config

type Pack struct {
	Server      string `mapstructure:"server" json:"server" yaml:"server"`                   // 打包服务
	PublishPath string `mapstructure:"publish-path" json:"publish-path" yaml:"publish-path"` // 发布项目路径
}
