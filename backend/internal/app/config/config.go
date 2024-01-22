package config

type Server struct {
	Zap       Zap       `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis     Redis     `mapstructure:"redis" json:"redis" yaml:"redis"`
	Casbin    Casbin    `mapstructure:"casbin" json:"casbin" yaml:"casbin"`
	System    System    `mapstructure:"system" json:"system" yaml:"system"`
	Captcha   Captcha   `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Pgsql     Pgsql     `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	JWT       JWT       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	IPFS      []IPFS    `mapstructure:"ipfs" json:"ipfs" yaml:"ipfs"`
	Pack      Pack      `mapstructure:"pack" json:"pack" yaml:"pack"`
	Airdrop   Airdrop   `mapstructure:"airdrop" json:"airdrop" yaml:"airdrop"`
	Local     Local     `mapstructure:"local" json:"local" yaml:"local"`
	Quest     *Quest    `mapstructure:"quest" json:"quest" yaml:"quest"`
	Translate Translate `mapstructure:"translate" json:"translate" yaml:"translate"`
	ZCloak    *ZCloak   `mapstructure:"zcloak" json:"zcloak" yaml:"zcloak"`
	NFT       NFT       `mapstructure:"nft" json:"nft" yaml:"nft"`
}
