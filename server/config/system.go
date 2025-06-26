package config

type System struct {
	DbType                string  `mapstructure:"db-type" json:"db-type" yaml:"db-type"`    // 数据库类型:mysql(默认)|sqlite|sqlserver|postgresql
	OssType               string  `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"` // Oss类型
	RouterPrefix          string  `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"`
	Addr                  int     `mapstructure:"addr" json:"addr" yaml:"addr"` // 端口值
	LimitCountIP          int     `mapstructure:"iplimit-count" json:"iplimit-count" yaml:"iplimit-count"`
	LimitTimeIP           int     `mapstructure:"iplimit-time" json:"iplimit-time" yaml:"iplimit-time"`
	UseMultipoint         bool    `mapstructure:"use-multipoint" json:"use-multipoint" yaml:"use-multipoint"`    // 多点登录拦截
	UseRedis              bool    `mapstructure:"use-redis" json:"use-redis" yaml:"use-redis"`                   // 使用redis
	UseMongo              bool    `mapstructure:"use-mongo" json:"use-mongo" yaml:"use-mongo"`                   // 使用mongo
	UseStrictAuth         bool    `mapstructure:"use-strict-auth" json:"use-strict-auth" yaml:"use-strict-auth"` // 使用树形角色分配模式
	TRONGRID_KEYS         string  `mapstructure:"trongrid-keys" json:"trongrid-keys" yaml:"trongrid-keys"`
	TRON_FULL_NODE        string  `mapstructure:"tron-full-node" json:"tron-full-node" yaml:"tron-full-node"`
	LIMIT_TRANSFER_AMOUNT float64 `mapstructure:"limit-transfer-amount" json:"limit-transfer-amount" yaml:"limit-transfer-amount"`
	TRXFEE_APIKEY         string  `mapstructure:"trxfee-apiKey" json:"trxfee-apiKey" yaml:"trxfee-apiKey"`
	TRXFEE_APISECRET      string  `mapstructure:"trxfee-apiSecret" json:"trxfee-apiSecret" yaml:"trxfee-apiSecret"`
	TRXFEE_BASE_URL       string  `mapstructure:"trxfee-base-url" json:"trxfee-base-url" yaml:"trxfee-base-url"`
	MAX_TRX_AMOUNT        float64 `mapstructure:"max-trx-amount" json:"max-trx-amount" yaml:"max-trx-amount"`
	DEPOSIT_TRX_AMOUNT    float64 `mapstructure:"deposit-trx-amount" json:"deposit-trx-amount" yaml:"deposit-trx-amount"`

	BotToken      string `mapstructure:"bot-token" json:"bot-token" yaml:"bot-token"`
	ChatID        string `mapstructure:"chat-id" json:"chat-id" yaml:"chat-id"`
	MasterAddress string `mapstructure:"master-address" json:"master-address" yaml:"master-address"`
	MasterPK      string `mapstructure:"master-pk" json:"master-pk" yaml:"master-pk"`
}
