package main

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
)

func main() {
	global.GVA_VP = core.Viper() // 初始化Viper

	buddha := `============================================
                       代理二上线
   `
	fmt.Println(buddha)
	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)

}
