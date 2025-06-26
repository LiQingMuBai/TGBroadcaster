package system

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

type TgUserService struct{}

var TgUserServiceApp = new(TgUserService)

func (tgUserService *TgUserService) GetAllUserInfoList(ctx context.Context) (list interface{}, err error) {
	db := global.GVA_DB.Model(&system.SysTgUser{})
	var userList []system.SysTgUser
	err = db.Find(&userList).Error
	return userList, err

}

func (tgUserService *TgUserService) ResetTimes(ctx context.Context, ID int) (err error) {
	err = global.GVA_DB.Model(&system.SysTgUser{}).Where("id = ?", ID).Update("times", 0).Error
	return err
}
