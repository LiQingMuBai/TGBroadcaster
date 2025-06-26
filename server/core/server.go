package core

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
		initialize.RedisList()
	}

	if global.GVA_CONFIG.System.UseMongo {
		err := initialize.Mongo.Initialization()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}
	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)

	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	buddha := `
                  _ooOoo_
                o8888888o
                  88" . "88
              	    (| -_- |)
                  O\  =  /O
               ____/'---'\____
             .'  \\|     |//  '.
            /  \\|||  :  |||//  \
           /  _||||| -:- |||||_  \
           |   | \\\  -  /'| |   |
           | \_|  '\'---'//  |_/ |
           \  .-\__ '-. -' __/-.  /
         ___'. .'  /--.--\  '. .'___
      ."" '<  '.___\_<|>_/___.' _> \"".
     | | :  '- \'. ;'. _/; .'/ /  .' ; |
     \  \ '-.   \_\_'. _.'_/_/  -' _.' /
   ====='-.____'.___ \_____/___.-'____.-'=====
                   '=---='
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
            佛祖保佑        永无BUG
   `
	global.GVA_LOG.Info(buddha)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
