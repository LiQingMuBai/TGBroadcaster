package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/pkg"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	currentKeyIndex uint32
	userService     = service.ServiceGroupApp.SystemServiceGroup.UserService
	sysOrderService = service.ServiceGroupApp.SystemServiceGroup.SysOrderService
)

type App struct {
	ticker *time.Ticker
	done   chan bool
	logger *log.Logger
}

func main() {
	global.GVA_VP = core.Viper() // 初始化Viper

	buddha := `============================================
                       代理一上线
   `
	fmt.Println(buddha)
	global.GVA_DB = initialize.Gorm() // gorm连接数据库
	initialize.DBList()

	global.GVA_LOG = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.GVA_LOG)

	if global.GVA_DB != nil {
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}
	global.TRONGRID_KEYS = strings.Split(global.GVA_CONFIG.System.TRONGRID_KEYS, ",")

	log.Println(global.TRONGRID_KEYS)
	// 初始化应用
	app := &App{
		done:   make(chan bool),
		logger: log.New(os.Stdout, "ORDER-EXCHANGE-ENERGY-SCHEDULER: ", log.LstdFlags),
	}

	// 每隔1min启动定时任务
	app.startScheduler(1 * time.Minute)

	// 等待关闭信号
	app.waitForShutdown()

	app.logger.Println("应用程序已关闭")

}

// 定时任务调度器
func (a *App) startScheduler(interval time.Duration) {
	a.ticker = time.NewTicker(interval)

	// 立即执行第一次任务
	go a.executeTask()

	go func() {
		for {
			select {
			case <-a.ticker.C:
				a.executeTask()
			case <-a.done:
				a.logger.Println("定时任务调度器已停止")
				return
			}
		}
	}()
}
func (a *App) executeTask() {
	a.logger.Println("开始能量兑换系统-执行定时任务...")
	startTime := time.Now()
	tronClient := pkg.NewTronClient(global.GVA_CONFIG.System.TRON_FULL_NODE)
	sendAmount := utils.ConvertFloatToBigInt(global.GVA_CONFIG.System.DEPOSIT_TRX_AMOUNT, 6)

	apiSecret := global.GVA_CONFIG.System.TRXFEE_APISECRET
	apiKey := global.GVA_CONFIG.System.TRXFEE_APIKEY
	baseUrl := global.GVA_CONFIG.System.TRXFEE_BASE_URL
	trxfeeClient := pkg.NewTrxfeeClient(baseUrl, apiKey, apiSecret)

	accountResp, err := trxfeeClient.Account()
	trxFeeAccountAddress := accountResp.Data.RechargeAddr
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("获取trxfee账户失败: %v\n", err))
		return
	}

	trxFeeAccountBalance := accountResp.Data.Balance

	if trxFeeAccountBalance <= global.GVA_CONFIG.System.MAX_TRX_AMOUNT {
		//
		global.GVA_LOG.Error(fmt.Sprintf("需要充值trxfee余额不够，余额：%f，最低余额%f\n", trxFeeAccountBalance, global.GVA_CONFIG.System.MAX_TRX_AMOUNT))

		global.GVA_LOG.Info(fmt.Sprintf("需要充值trxfee地址：%s\n", trxFeeAccountAddress))

		//telegram通知
		go notifyInsufficientGas(global.GVA_CONFIG.System.ChatID, global.GVA_CONFIG.System.BotToken, accountResp.Data.RechargeAddr, trxFeeAccountBalance)
		global.GVA_LOG.Info(fmt.Sprintf("telegram通知：%s\n", trxFeeAccountAddress))

		//调用接口去充值

		log.Println("=======================================")
		log.Println("sendAmount:", sendAmount)
		log.Println("tronClient:", global.GVA_CONFIG.System.TRON_FULL_NODE)
		log.Println("pk:", global.GVA_CONFIG.System.MasterPK)
		log.Println("address:", trxFeeAccountAddress)
		log.Println("=======================================")
		go func() {
			_, err := tronClient.TransferNative(context.Background(), global.GVA_CONFIG.System.MasterPK, trxFeeAccountAddress, sendAmount)
			if err != nil {

			}
		}()
		global.GVA_LOG.Info(fmt.Sprintf("主地址進行地址充值：%s\n", trxFeeAccountAddress))

		return
	}

	users, total, err := userService.GetUserInfoListAndAddressNotNull()
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("GetUserInfoListAndAddressNotNull失败: %v\n", err))
	}
	if total > 0 {

		for _, user := range users {
			transactions, err := getTRXTransactionsByAddress(user.Address, global.GVA_CONFIG.System.TRON_FULL_NODE, "200")
			if err != nil {
				global.GVA_LOG.Error(fmt.Sprintf("Error fetching bussiness's transactions: %v\n", err))
				continue
			}

			for _, transaction := range transactions {
				if transaction.Amount <= global.GVA_CONFIG.System.LIMIT_TRANSFER_AMOUNT {
					global.GVA_LOG.Info(fmt.Sprintf("订单金额太小，交易: %s，金额: %f\n", transaction.TxID, transaction.Amount))
					continue
				}

				order, err := sysOrderService.GetSysOrderByTxID(transaction.TxID)
				if err != nil {
					global.GVA_LOG.Error(fmt.Sprintf("获取数据订单失败: %v\n", err))
					continue
				}

				if order.ID > 0 {
					global.GVA_LOG.Info(fmt.Sprintf("订单已经发送无需重复: %s\n", order.TxID))
					continue
				} else {
					var sysOrder system.SysOrder
					orderNo, _ := pkg.GenerateOrderID(transaction.From, 4)
					fmt.Printf("  OrderNo: %s\n", orderNo)
					sysOrder.OrderNo = orderNo
					sysOrder.TxID = transaction.TxID
					sysOrder.FromAddress = transaction.From
					sysOrder.ToAddress = transaction.To
					sysOrder.Amount = transaction.Amount

					//添加一条记录
					err := sysOrderService.CreateSysOrder(&sysOrder)

					if err != nil {
						global.GVA_LOG.Error(fmt.Sprintf("添加一条记录订单失败: %v\n", err))
						continue
					}

					count := int(transaction.Amount / global.GVA_CONFIG.System.LIMIT_TRANSFER_AMOUNT)

					if count*int(global.GVA_CONFIG.System.LIMIT_TRANSFER_AMOUNT) > int(trxFeeAccountBalance) {
						global.GVA_LOG.Error(fmt.Sprintf("需要(%d)笔数，金额不够需要充值\n", count))
						go notifyInsufficientGas(global.GVA_CONFIG.System.ChatID, global.GVA_CONFIG.System.BotToken, accountResp.Data.RechargeAddr, trxFeeAccountBalance)

						go func() {
							_, err := tronClient.TransferNative(context.Background(), global.GVA_CONFIG.System.MasterPK, trxFeeAccountAddress, sendAmount)
							if err != nil {

							}
						}()

						continue
					}
					global.GVA_LOG.Info(fmt.Sprintf("发送（%d）笔能量给（%s），订单号 %s\n", count, sysOrder.FromAddress, orderNo))
					trxfeeClient.Order(sysOrder.OrderNo, sysOrder.FromAddress, 65_000*count)
				}
			}
		}
	}

	a.logger.Printf("任务完成， 耗时: %v", time.Since(startTime))

}

// 等待关闭信号并关闭
func (a *App) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待关闭信号
	sig := <-quit
	a.logger.Printf("接收到信号: %v, 开始关闭...", sig)

	// 停止定时任务
	a.ticker.Stop()
	a.done <- true

	// 等待所有任务完成
	time.Sleep(1 * time.Second) // 可根据需要调整
}

type TransactionTRXResp struct {
	TxID   string  `json:"txID"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

func notifyInsufficientGas(_chatID string, _botToken string, _address string, _amount float64) {
	//var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	//	tgbotapi.NewInlineKeyboardRow(
	//		tgbotapi.NewInlineKeyboardButtonURL("交易详情", "https://tronscan.org/#/address/"+_address),
	//	),
	//)

	message := map[string]interface{}{
		"chat_id": _chatID, // 或直接用 chat_id 如 "123456789"=
		"text": "⚠【主地址Trx余额不足警告提醒】\n\n" +
			"📢地址：" + _address + "\n\n" +
			"📢平台余额：      " + fmt.Sprintf("%f", _amount) + "\n\n",
	}
	// 转换为 JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		fmt.Println("JSON  parse error...:", err)
		return
	}

	// 发送 POST 请求到 Telegram Bot API
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", _botToken)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("发送消息失败:", err)
		return
	}
	defer resp.Body.Close()

	// 打印响应结果
	//fmt.Println("消息发送状态:", resp.Status)
}

// 获取指定地址的交易列表
func getTRXTransactionsByAddress(address string, apiURL string, pageSize string) ([]TransactionTRXResp, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("%s/v1/accounts/%s/transactions?only_to=true&limit="+pageSize, apiURL, address)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	//// 简单轮询
	keyIndex := atomic.AddUint32(&currentKeyIndex, 1) % uint32(len(global.TRONGRID_KEYS))
	currentKey := global.TRONGRID_KEYS[keyIndex]

	req.Header.Set("TRON-PRO-API-KEY", currentKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("API endpoint not found (404)")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 解析JSON响应
	var txResponse TransactionTRXResponse
	err = json.Unmarshal(body, &txResponse)
	if err != nil {
		//fmt.Printf("解析JSON失败: %v\n", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// 打印TRX交易列表
	fmt.Printf("地址 %s 的TRX交易列表:\n", address)
	trxCount := 0
	var result []TransactionTRXResp // nil slice
	for _, tx := range txResponse.Data {
		// 检查是否为TRX交易
		if len(tx.RawData.Contract) > 0 {
			contract := tx.RawData.Contract[0]

			// 只处理TransferContract类型的交易(TRX转账)
			if contract.Type == "TransferContract" {
				// 确保不是TRC20交易(TRC20交易会有contract_address字段)
				if contract.Parameter.Value.ContractAddress == "" {
					trxCount++
					fmt.Printf("\n交易 #%d\n", trxCount)
					fmt.Printf("交易ID: %s\n", tx.TxID)
					//fmt.Printf("发送方: %s\n", contract.Parameter.Value.OwnerAddress)
					//fmt.Printf("接收方: %s\n", contract.Parameter.Value.ToAddress)
					amount := float64(contract.Parameter.Value.Amount) / 1000000
					fmt.Printf("金额: %.6f TRX\n", amount)

					ownerAddress, _ := TronHexToBase58(contract.Parameter.Value.OwnerAddress)
					toAddress, _ := TronHexToBase58(contract.Parameter.Value.ToAddress)

					fmt.Printf("发送方: %s\n", ownerAddress)
					fmt.Printf("接收方: %s\n", toAddress)
					var resource TransactionTRXResp
					resource.TxID = tx.TxID
					resource.From = ownerAddress
					resource.To = toAddress
					resource.Amount = amount
					result = append(result, resource)

				}
			}
		}
	}

	if trxCount == 0 {
		fmt.Println("未找到TRX交易记录")
	}
	return result, nil
}

type TransactionTRXResponse struct {
	Data []TransactionTRX `json:"data"`
}
type TransactionTRX struct {
	TxID    string `json:"txID"`
	RawData struct {
		Contract []struct {
			Type      string `json:"type"`
			Parameter struct {
				Value struct {
					Amount       int64  `json:"amount"`
					OwnerAddress string `json:"owner_address"`
					ToAddress    string `json:"to_address"`
					// 对于TRC20交易会有不同的字段
					ContractAddress string `json:"contract_address"` // TRC20交易会有此字段
					Data            string `json:"data"`             // TRC20交易会有此字段
				} `json:"value"`
			} `json:"parameter"`
		} `json:"contract"`
	} `json:"raw_data"`
}

func TronHexToBase58(hexAddr string) (string, error) {
	addrBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex address: %w", err)
	}
	first := sha256.Sum256(addrBytes)
	second := sha256.Sum256(first[:])
	checksum := second[:4]
	fullBytes := append(addrBytes, checksum...)
	return base58.Encode(fullBytes), nil
}
