package system

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type SysAddress struct {
	global.GVA_MODEL
	Address string `json:"address" gorm:"comment:地址"` // api路径
	Network string `json:"network" gorm:"comment:网络"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
	Status  string `json:"status" gorm:"comment:状态"`  // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE
}

func (SysAddress) TableName() string {
	return "sys_address"
}

type SysOrder struct {
	global.GVA_MODEL
	OrderNo     string ` json:"order_no" form:"order_no" gorm:"column:order_no;comment:order_no"`
	TxID        string ` json:"tx_id" form:"tx_id" gorm:"column:tx_id;comment:tx_id"`
	FromAddress string `json:"from_address" form:"from_address" gorm:"column:from_address;comment:from_address"`
	ToAddress   string `json:"to_address" form:"to_address" gorm:"column:to_address;comment:to_address"`
	Remark      string `json:"remark" gorm:"comment:备注"` // 方法:创建POST(默认)|查看GET|更新PUT|删除DELETE

	Amount float64 `json:"amount" form:"amount" gorm:"column:amount;comment:数量;size:64;"` //数量
}

func (SysOrder) TableName() string {
	return "sys_order"
}
