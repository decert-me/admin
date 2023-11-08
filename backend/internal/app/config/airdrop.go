package config

type Airdrop struct {
	VerifyKey string `mapstructure:"verify-key" json:"verify-key" yaml:"verify-key"`
	Api       string `mapstructure:"api" json:"api" yaml:"api"`
}
