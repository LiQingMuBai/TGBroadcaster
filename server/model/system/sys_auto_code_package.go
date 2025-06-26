package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type SysAutoCodePackage struct {
	global.GVA_MODEL
	Desc        string `json:"desc" gorm:"comment:描述"`
	Label       string `json:"label" gorm:"comment:展示名"`
	Template    string `json:"template"  gorm:"comment:模版"`
	PackageName string `json:"packageName" gorm:"comment:包名"`
	Module      string `json:"-" example:"模块"`
}

func (s *SysAutoCodePackage) TableName() string {
	return "sys_auto_code_packages"
}

type SysTgUser struct {
	ID          uint   `gorm:"primarykey" json:"ID"` // 主键ID
	UserID      string `json:"user_id" gorm:"comment:用户id"`
	Times       int64  `json:"times" gorm:"comment:次数"`
	Username    string `json:"username" gorm:"comment:用户名"`
	Amount      string `json:"amount" gorm:"comment:金额"`
	Associates  string `json:"associates" gorm:"comment:关联人"`
	TronAmount  string `json:"tron_amount" gorm:"comment:波场余额"`
	TronAddress string `json:"tron_address" gorm:"comment:波场地址"`
	EthAddress  string `json:"eth_address" gorm:"comment:以太坊地址"`
	EthAmount   string `json:"eth_amount" gorm:"comment:以太坊余额"`
}

func (s *SysTgUser) TableName() string {
	return "tg_users"
}
