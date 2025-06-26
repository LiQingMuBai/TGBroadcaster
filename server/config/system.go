package config

type System struct {
	RouterPrefix string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Addr         int    `mapstructure:"addr" json:"addr" yaml:"addr"` // 端口值
}
