package config

type Server struct {
	Zap     Zap     `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   Redis   `mapstructure:"redis" json:"redis" yaml:"redis"`
	Casbin  Casbin  `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	System  System  `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	// gorm
	Pgsql Pgsql `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	// JWT
	JWT JWT `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	// OSS
	Local Local `mapstructure:"local" json:"local" yaml:"local"`
	// IPFS
	IPFS []IPFS `mapstructure:"ipfs" json:"ipfs" yaml:"ipfs"`
}
