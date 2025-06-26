package system

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"gorm.io/gorm"
)

type SysOrderService struct{}

// CreateSysOrder 创建参数记录
// Author [Mr.奇淼](https://github.com/pixelmaxQm)
func (sysOrderService *SysOrderService) CreateSysOrder(sysOrder *system.SysOrder) (err error) {
	err = global.GVA_DB.Create(sysOrder).Error
	return err
}

// DeleteSysOrder 删除参数记录
// Author [Mr.奇淼](https://github.com/pixelmaxQm)
func (sysOrderService *SysOrderService) DeleteSysOrder(ID string) (err error) {
	err = global.GVA_DB.Delete(&system.SysOrder{}, "id = ?", ID).Error
	return err
}

// DeleteSysOrderByIds 批量删除参数记录
// Author [Mr.奇淼](https://github.com/pixelmaxQm)
func (sysOrderService *SysOrderService) DeleteSysOrderByIds(IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]system.SysOrder{}, "id in ?", IDs).Error
	return err
}

// UpdateSysOrder 更新参数记录
// Author [Mr.奇淼](https://github.com/pixelmaxQm)
func (sysOrderService *SysOrderService) UpdateSysOrder(sysOrder system.SysOrder) (err error) {
	err = global.GVA_DB.Model(&system.SysOrder{}).Where("id = ?", sysOrder.ID).Updates(&sysOrder).Error
	return err
}

// GetSysOrder 根据ID获取参数记录
// Author [Mr.奇淼](https://github.com/pixelmaxQm)
func (sysOrderService *SysOrderService) GetSysOrder(ID string) (sysOrder system.SysOrder, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&sysOrder).Error
	return
}

// GetSysOrder 根据ID获取参数记录
// Author [Mr.奇淼](https://github.com/pixelmaxQm)
func (sysOrderService *SysOrderService) GetSysOrderByTxID(txID string) (sysOrder system.SysOrder, err error) {
	err = global.GVA_DB.Where("tx_id = ?", txID).First(&sysOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 处理没有找到记录的情况
		return sysOrder, nil
	} else if err != nil {
		return sysOrder, err
	}
	return
}
