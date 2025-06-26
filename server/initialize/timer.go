package initialize

import (
	"boost/data/server/global"
	"fmt"
	"github.com/robfig/cron/v3"
)

const (
	// TRON主网节点
	tronGRPCEndpoint = "grpc.trongrid.io:50051"
	// USDT TRC20合约地址
	usdtContractAddress = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	// 主地址(归集目标地址)
	mainAddress = "TARGET_ADDRESS_HERE"
	// 主地址私钥(实际生产环境应从安全存储获取)
	mainPrivateKey = "PRIVATE_KEY_HERE"
	// 最小归集金额(单位: 最小精度)
	minCollectAmount = 1000000 // 1 USDT (USDT有6位小数)
)

func Timer() {
	go func() {
		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		//_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
		//	err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
		//	if err != nil {
		//		fmt.Println("timer error:", err)
		//	}
		//}, "定时清理数据库【日志，黑名单】内容", option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}

		// 其他定时任务定在这里 参考上方使用方法

		_, err := global.GVA_Timer.AddTaskByFunc("CollectUSDT", "@every 00h15m00s", func() {
			//err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			//if err != nil {
			//	fmt.Println("timer error:", err)
			//}

			global.GVA_LOG.Info("开始归集USDT")
		}, "每隔15分钟检查下用户余额，进行归集", option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}
	}()
}
