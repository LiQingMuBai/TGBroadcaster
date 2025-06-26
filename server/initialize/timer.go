package initialize

import (
	model "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"log"
)

func TimeSync() {

	address := model.SysAddress{
		Address: "TWePRsCpramXCyDrZ61j3oqyFTr2SCDoDk",
		Status:  "pending",
		Network: "tron",
	}

	//err := global.GVA_DB.Create(address).Error

	//if err != nil {
	//	log.Println(err)
	//}

	log.Println(address)
}

func Timer() {
	//go func() {
	//	var option []cron.Option
	//	option = append(option, cron.WithSeconds())
	//	//// 清理DB定时任务
	//	//_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
	//	//	err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
	//	//	if err != nil {
	//	//		fmt.Println("timer error:", err)
	//	//	}
	//	//}, "定时清理数据库【日志，黑名单】内容", option...)
	//	//if err != nil {
	//	//	fmt.Println("add timer error:", err)
	//	//}
	//
	//	// 其他定时任务定在这里 参考上方使用方法
	//
	//	//_, err := global.GVA_Timer.AddTaskByFunc("定时任务标识", "corn表达式", func() {
	//	//	具体执行内容...
	//	//  ......
	//	//}, option...)
	//	//if err != nil {
	//	//	fmt.Println("add timer error:", err)
	//	//}
	//	// 清理DB定时任务
	//	_, err := global.GVA_Timer.AddTaskByFunc("Ping", "@every 00h04m59s", func() {
	//		//err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
	//		//if err != nil {
	//		//	fmt.Println("timer error:", err)
	//		//}
	//
	//		sumbitMap, commitMap := system.GetTronAddressMap()
	//
	//		for tronAddress, _ := range sumbitMap {
	//
	//			log.Println("预冻结", tronAddress)
	//
	//			if strUtil.IsNotEmpty(tronAddress) {
	//
	//				var info model.SysAddress
	//				err := global.GVA_DB.Where("address = ? and status = ?", tronAddress, "pending").First(&info).Error
	//
	//				if info.Address == "" {
	//					address := model.SysAddress{
	//						Address: tronAddress,
	//						Status:  "pending",
	//						Network: "tron",
	//					}
	//
	//					err = global.GVA_DB.Create(&address).Error
	//
	//					if err != nil {
	//						log.Println(err)
	//					}
	//				}
	//
	//			}
	//		}
	//
	//		for tronAddress, _ := range commitMap {
	//			log.Println("已冻结", tronAddress)
	//			if strUtil.IsNotEmpty(tronAddress) {
	//				var info model.SysAddress
	//				err := global.GVA_DB.Where("address = ? and status = ?", tronAddress, "completed").First(&info).Error
	//
	//				log.Println("地址     ", info)
	//				if info.Address == "" {
	//					address := model.SysAddress{
	//						Address: tronAddress,
	//						Status:  "completed",
	//						Network: "tron",
	//					}
	//
	//					err = global.GVA_DB.Create(&address).Error
	//
	//					if err != nil {
	//						log.Println(err)
	//					}
	//				}
	//			}
	//
	//		}
	//
	//		sumbitEthereumMap, commitEthereumMap, _ := system.GetEthereumAddress()
	//
	//		for ethereumAddress, _ := range sumbitEthereumMap {
	//
	//			log.Println("预冻结", ethereumAddress)
	//			if strUtil.IsNotEmpty(ethereumAddress) {
	//
	//				var info model.SysAddress
	//				err := global.GVA_DB.Where("address = ? and status = ?", ethereumAddress, "pending").First(&info).Error
	//
	//				log.Println("地址     ", info)
	//				if info.Address == "" {
	//					address := model.SysAddress{
	//						Address: ethereumAddress,
	//						Status:  "pending",
	//						Network: "ethereum",
	//					}
	//
	//					err = global.GVA_DB.Create(&address).Error
	//
	//					if err != nil {
	//						log.Println(err)
	//					}
	//				}
	//			}
	//		}
	//
	//		for ethereumAddress, _ := range commitEthereumMap {
	//			log.Println("已冻结", ethereumAddress)
	//			if strUtil.IsNotEmpty(ethereumAddress) {
	//				var info model.SysAddress
	//				err := global.GVA_DB.Where("address = ? and status = ?", ethereumAddress, "completed").First(&info).Error
	//
	//				if info.Address == "" {
	//					address := model.SysAddress{
	//						Address: ethereumAddress,
	//						Status:  "completed",
	//						Network: "ethereum",
	//					}
	//
	//					err = global.GVA_DB.Create(&address).Error
	//
	//					if err != nil {
	//						log.Println(err)
	//					}
	//
	//				}
	//			}
	//		}
	//
	//	}, "定时清理数据库【日志，黑名单】内容", option...)
	//	if err != nil {
	//		fmt.Println("add timer error:", err)
	//	}
	//
	//}()
}
