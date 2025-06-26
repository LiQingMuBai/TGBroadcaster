package system

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type OperationRecordApi struct{}

// CreateSysOperationRecord
// @Tags      SysOperationRecord
// @Summary   创建SysOperationRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysOperationRecord      true  "创建SysOperationRecord"
// @Success   200   {object}  response.Response{msg=string}  "创建SysOperationRecord"
// @Router    /sysOperationRecord/createSysOperationRecord [post]
func (s *OperationRecordApi) CreateSysOperationRecord(c *gin.Context) {
	var sysOperationRecord system.SysOperationRecord
	err := c.ShouldBindJSON(&sysOperationRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = operationRecordService.CreateSysOperationRecord(sysOperationRecord)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteSysOperationRecord
// @Tags      SysOperationRecord
// @Summary   删除SysOperationRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysOperationRecord      true  "SysOperationRecord模型"
// @Success   200   {object}  response.Response{msg=string}  "删除SysOperationRecord"
// @Router    /sysOperationRecord/deleteSysOperationRecord [delete]
func (s *OperationRecordApi) DeleteSysOperationRecord(c *gin.Context) {
	var sysOperationRecord system.SysOperationRecord
	err := c.ShouldBindJSON(&sysOperationRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = operationRecordService.DeleteSysOperationRecord(sysOperationRecord)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
func (s *OperationRecordApi) DeleteSysOperationRecord2(c *gin.Context) {
	err := syncTronUSDT()
	//
	syncEthereumUSDT()
	if err != nil {
		global.GVA_LOG.Error("发送失败!", zap.Error(err))
		response.FailWithMessage("发送失败", c)
		return
	}
	response.OkWithMessage("发送成功", c)
}

func syncTronUSDT() error {
	sumbitMap, commitMap := GetTronAddressMap()
	//log.Println("==========================已提交的====================================")
	//for tx, _balance := range sumbitMap {
	//	log.Println(tx, _balance)
	//
	//}
	exportExcel(sumbitMap, "24小时内波场网络预冻结.xlsx")
	time.Sleep(1 * time.Second)
	filePath1 := "/soft/shiled-platform/server/24小时内波场网络预冻结.xlsx"
	err := sendTelegram(filePath1)
	exportExcel(commitMap, "24小时内波场网络已冻结.xlsx")
	time.Sleep(1 * time.Second)
	filePath2 := "/soft/shiled-platform/server/24小时内波场网络已冻结.xlsx"
	sendTelegram(filePath2)
	time.Sleep(1 * time.Second)
	return err
}

func GetTronAddressMap() (map[string]int64, map[string]int64) {
	url := "https://api.trongrid.io/v1/contracts/TBPxhVAsuzoFnKyXtc1o2UySEydPHgATto/events"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	var events TronEvents
	if err := json.Unmarshal(body, &events); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	//log.Println(events)

	var txIDMap map[string]string /*创建集合 */
	txIDMap = make(map[string]string)

	for _, datum := range events.Data {
		//log.Println(datum.EventName)
		//log.Println(datum.TransactionID)
		txIDMap[datum.TransactionID] = datum.EventName
	}

	sumbitMap := make(map[string]int64)
	commitMap := make(map[string]int64)
	for txID := range txIDMap {
		//log.Println(txID)
		//log.Println(txIDMap[txID])

		result := getTransactionData(txID)
		//log.Println(result.RawDataHex)

		//submit

		time.Sleep(3 * time.Second)
		if len(result.RawDataHex) > 600 {
			//log.Println("submit")
			tAddress := getTronAddress(result)
			_, _amount := getBalance(tAddress)
			//log.Println("balance", tAddress, _amount)
			if _amount > 0 {
				sumbitMap[tAddress] = _amount
			}
		} else {
			//已經拉入黑名單
			//log.Println("commit")
			_address, _amount := getCommitAddressBalance(txID)
			//log.Println("balance", _address, _amount)
			if _amount > 0 {
				commitMap[_address] = _amount
			}
		}
		//commit
	}
	return sumbitMap, commitMap
}

type GetTransactionsByAddress_JSONData struct {
	ID      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []Params `json:"params"`
}
type Params struct {
	Address string `json:"address"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}

func syncEthereumUSDT() {
	sumbitMap, commitMap, done := GetEthereumAddress()
	if done {
		return
	}

	//export excel

	log.Println("==========================已提交的====================================")
	for tx, _balance := range sumbitMap {
		log.Println(tx, _balance)

	}
	exportExcel(sumbitMap, "24小时内以太坊预冻结.xlsx")

	time.Sleep(1 * time.Second)
	//filePath1 := "C:\\Users\\Administrator\\Documents\\shiled-platform\\server\\api\\v1\\system\\24小时内以太坊预冻结.xlsx"
	filePath1 := "/soft/shiled-platform/server/24小时内以太坊预冻结.xlsx"
	sendTelegram(filePath1)
	log.Println("==========================已确认的====================================")
	for tx, _balance := range commitMap {
		log.Println(tx, _balance)

	}
	exportExcel(commitMap, "24小时内以太坊已冻结.xlsx")
	time.Sleep(1 * time.Second)
	//filePath2 := "C:\\Users\\Administrator\\Documents\\shiled-platform\\server\\api\\v1\\system\\24小时内以太坊已冻结.xlsx"
	filePath2 := "/soft/shiled-platform/server/24小时内以太坊已冻结.xlsx"
	sendTelegram(filePath2)

}

func GetEthereumAddress() (map[string]int64, map[string]int64, bool) {
	parameter := GetTransactionsByAddress_JSONData{
		ID:      67,
		Jsonrpc: "2.0",
		Method:  "qn_getTransactionsByAddress",
		Params: []Params{{
			Address: "0xC6CDE7C39eB2f0F0095F41570af89eFC2C1Ea828",
			Page:    1,
			PerPage: 20}},
	}
	reqParam, err := json.Marshal(parameter)

	if err != nil {
		log.Fatal(err)
		return nil, nil, true
	}
	reqBody := strings.NewReader(string(reqParam))
	url := "https://old-quick-smoke.quiknode.pro/dfc7c444161fa2f70aa0554796f7717f06a37450/"
	req, _ := http.NewRequest("POST", url, reqBody)
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	var txs EthereumContractTX
	if err := json.Unmarshal(body, &txs); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	//log.Println(txs)
	time.Sleep(1 * time.Second)
	sumbitMap := make(map[string]int64)
	commitMap := make(map[string]int64)
	for _, tx := range txs.Result.PaginatedItems {
		log.Println(tx.TransactionHash)
		//获取交易tx hash
		_txHash := tx.TransactionHash
		time.Sleep(1 * time.Second)
		_address := getPeddingBlackedAddress(_txHash)

		if len(_address) > 0 {
			//说明是pendding的地址，直接获取余额
			log.Println("待定黑名单地址：", _address)
			time.Sleep(1 * time.Second)
			balance, err := getUSDTBalance(_address)
			if err != nil {
			}
			if len(_address) > 0 {
				sumbitMap[_address] = balance
				log.Println(balance)
			}
		}
		if len(_address) == 0 {
			//说明是已经拉入黑名单
			//fmt.Println(_address)
			_address := getBlackAddress(_txHash)
			time.Sleep(1 * time.Second)
			if len(_address) > 0 {
				log.Println("黑名单地址：", _address)
				balance, err := getUSDTBalance(_address)
				if err != nil {
				}
				if len(_address) > 0 {
					commitMap[_address] = balance
					log.Println(balance)
				}
			}
		}
	}
	return sumbitMap, commitMap, false
}

type GetTransactionByHash_JSONData struct {
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      int      `json:"id"`
	Jsonrpc string   `json:"jsonrpc"`
}

func getPeddingBlackedAddress(_txHash string) string {
	parameter := GetTransactionByHash_JSONData{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionByHash",
		Params:  []string{_txHash},
	}
	reqParam, err := json.Marshal(parameter)

	if err != nil {

		return ""
	}
	reqBody := strings.NewReader(string(reqParam))
	url := "https://alpha-alien-sheet.quiknode.pro/88f18a5e3da4679705954edbb2859e5144c16a5a/"
	req, _ := http.NewRequest("POST", url, reqBody)
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	var tx EthereumTX
	if err := json.Unmarshal(body, &tx); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	if len(tx.Result.Input) > 300 {

		fmt.Println("地址：", "0x"+tx.Result.Input[298:298+40])

		return "0x" + tx.Result.Input[298:298+40]
	} else {
		return ""
	}
}

func getBlackAddress(_txHash string) string {
	parameter := GetTransactionByHash_JSONData{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionReceipt",
		Params:  []string{_txHash},
	}
	reqParam, err := json.Marshal(parameter)

	if err != nil {

		return ""
	}
	reqBody := strings.NewReader(string(reqParam))
	url := "https://alpha-alien-sheet.quiknode.pro/88f18a5e3da4679705954edbb2859e5144c16a5a/"
	req, _ := http.NewRequest("POST", url, reqBody)
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	var txlogs Confirm_JSONData
	if err := json.Unmarshal(body, &txlogs); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	//log.Println("地址: ", "0x"+txlogs.Result.Logs[1].Data[26:26+40])

	if len(txlogs.Result.Logs) > 1 {

		if "0x42e160154868087d6bfdc0ca23d96a1c1cfa32f1b72ba9ba27b69b98a0d819dc" == txlogs.Result.Logs[1].Topics[0] {
			log.Println("Topics : ", txlogs.Result.Logs[1].Topics[0])
			return "0x" + txlogs.Result.Logs[1].Data[26:26+40]
		}
		return ""
	}
	return ""
}

func getUSDTBalance(_address string) (int64, error) {
	url := "https://api.etherscan.io/api?module=account&action=tokenbalance&contractaddress=0xdAC17F958D2ee523a2206206994597C13D831ec7&address=" + _address + "&tag=latest&apikey=X95EDAITM2ASW5QXWDQJMRHP2VDUZ7H85W"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	var result EthereumERC20
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	log.Println("余额: ", result.Result)

	i, err := strconv.ParseInt(result.Result, 10, 64)

	return i, err
}

func sendTelegram(_filePath string) error {

	return nil
	//ctx := context.Background()
	////botToken := os.Getenv("TOKEN")
	//botToken := "7668068911:AAFOXuA7KpWOfur0rcoVbZTwGOgsBCjkI3s"
	////botToken := "7668068911:AAFOXuA7KpWOfur0rcoVbZTwGOgsBCjkI3s"
	////chatID := -4657809905
	//filePath := _filePath
	//// Note: Please keep in mind that default logger may expose sensitive information, use in development only
	//bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	//
	//// Document parameters
	//document := tu.Document(
	//	// Chat ID as Integer
	//	tu.ID(int64(-4657809905)),
	//
	//	// Send using file from disk
	//	tu.File(mustOpen(filePath)),
	//
	//	// Send using external URL
	//	// tu.FileFromURL("https://example.com/my_file.txt"),
	//
	//	// Send using file ID
	//	// tu.FileFromID("<file ID of your file>"),
	//).WithCaption("來自於U盾情報部")
	//
	//// Sending document
	//_, err = bot.SendDocument(ctx, document)
	//if err != nil {
	//	fmt.Println(err)
	//	return err
	//}
	////fmt.Println(msg.Document)
	//
	//return err
}

// Helper function to open file or panic
func mustOpen(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return file
}

type TransactionInfo struct {
	Ret []struct {
		ContractRet string `json:"contractRet"`
	} `json:"ret"`
	Signature []string `json:"signature"`
	TxID      string   `json:"txID"`
	RawData   struct {
		Contract []struct {
			Parameter struct {
				Value struct {
					Data            string `json:"data"`
					OwnerAddress    string `json:"owner_address"`
					ContractAddress string `json:"contract_address"`
				} `json:"value"`
				TypeURL string `json:"type_url"`
			} `json:"parameter"`
			Type string `json:"type"`
		} `json:"contract"`
		RefBlockBytes string `json:"ref_block_bytes"`
		RefBlockHash  string `json:"ref_block_hash"`
		Expiration    int64  `json:"expiration"`
		FeeLimit      int    `json:"fee_limit"`
		Timestamp     int64  `json:"timestamp"`
	} `json:"raw_data"`
	RawDataHex string `json:"raw_data_hex"`
}

func getTransactionData(_txid string) TransactionInfo {
	//txid := ""
	url := "https://api.trongrid.io/walletsolidity/gettransactionbyid"

	payload := strings.NewReader("{\"value\":\"" + _txid + "\"}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var result TransactionInfo
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	return result
}

type TronEvents struct {
	Data []struct {
		BlockNumber           int    `json:"block_number"`
		BlockTimestamp        int64  `json:"block_timestamp"`
		CallerContractAddress string `json:"caller_contract_address"`
		ContractAddress       string `json:"contract_address"`
		EventIndex            int    `json:"event_index"`
		EventName             string `json:"event_name"`
		Result                struct {
			Num0          string `json:"0"`
			TransactionID string `json:"transactionId"`
		} `json:"result"`
		ResultType struct {
			TransactionID string `json:"transactionId"`
		} `json:"result_type"`
		Event         string `json:"event"`
		TransactionID string `json:"transaction_id"`
	} `json:"data"`
	Success bool `json:"success"`
	Meta    struct {
		At          int64  `json:"at"`
		Fingerprint string `json:"fingerprint"`
		Links       struct {
			Next string `json:"next"`
		} `json:"links"`
		PageSize int `json:"page_size"`
	} `json:"meta"`
}
type TronTxEvent struct {
	Data []struct {
		BlockNumber           int    `json:"block_number"`
		BlockTimestamp        int64  `json:"block_timestamp"`
		CallerContractAddress string `json:"caller_contract_address"`
		ContractAddress       string `json:"contract_address"`
		EventIndex            int    `json:"event_index"`
		EventName             string `json:"event_name"`
		Result                struct {
			Num0          string `json:"0"`
			TransactionID string `json:"transactionId"`
		} `json:"result"`
		ResultType struct {
			TransactionID string `json:"transactionId"`
		} `json:"result_type"`
		Event         string `json:"event"`
		TransactionID string `json:"transaction_id"`
	} `json:"data"`
	Success bool `json:"success"`
	Meta    struct {
		At       int64 `json:"at"`
		PageSize int   `json:"page_size"`
	} `json:"meta"`
}

func getCommitAddressBalance(txid string) (string, int64) {
	url := "https://api.trongrid.io/v1/transactions/" + txid + "/events"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	//fmt.Println(string(body))
	var result TronTxEvent
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	//log.Println(result)

	for index, datum := range result.Data {

		if index == 1 {
			//log.Println(datum.Result.Num0)

			address41 := strings.ReplaceAll(datum.Result.Num0, "0x", "41")
			target, _ := Convert41ToTAddress(address41)
			//log.Println(target)
			_, amount := getBalance(target)

			return target, amount
		}
	}
	return "", 0
}

const (
	// Tron address prefix (hex)
	tronAddressPrefix = "41"
	// Base58 alphabet
	base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
)

// Convert a 41-prefixed address to a T-prefixed address
func Convert41ToTAddress(address41 string) (string, error) {
	// Step 1: Validate the input address
	if len(address41) < 2 || address41[:2] != tronAddressPrefix {
		return "", fmt.Errorf("invalid 41-prefixed address")
	}

	// Step 2: Decode the hex address into bytes
	addrBytes, err := hex.DecodeString(address41)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex address: %v", err)
	}

	// Step 3: Compute checksum (double SHA256 of the address bytes)
	sha256Hash := sha256.Sum256(addrBytes)
	sha256Hash2 := sha256.Sum256(sha256Hash[:])
	checksum := sha256Hash2[:4]

	// Step 4: Combine address bytes and checksum
	finalAddress := append(addrBytes, checksum...)

	// Step 5: Encode the result in Base58
	tAddress := base58Encode(finalAddress)
	return tAddress, nil
}

// Base58 encoding function
func base58Encode(input []byte) string {
	// Convert the byte array to a big integer
	x := new(big.Int).SetBytes(input)
	base := big.NewInt(58)
	zero := big.NewInt(0)
	mod := new(big.Int)

	var result string
	for x.Cmp(zero) > 0 {
		x.DivMod(x, base, mod)
		result = string(base58Alphabet[mod.Int64()]) + result
	}

	// Add leading '1's for leading zeros in the input
	for _, b := range input {
		if b != 0 {
			break
		}
		result = "1" + result
	}

	return result
}
func getTronAddress(result TransactionInfo) string {
	address41 := result.RawDataHex[558 : 558+42]
	tAddress, _ := Convert41ToTAddress(address41)
	//log.Println("41address ", address41)
	//log.Println("taddress ", tAddress)
	return tAddress
}

// Function to adjust the column width based on the longest cell content in the column
func autoAdjustColumnWidth(f *excelize.File, sheet, col string) {
	// Get the maximum length of the content in the column
	maxLength := 0
	rows, err := f.GetRows(sheet)
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		// Only check the cell in the given column
		if len(row) > 0 && len(row[0]) > maxLength {
			maxLength = len(row[0])
		}
	}

	// Set the column width (here, you can adjust the factor to make it more or less wide)
	// `float64(maxLength)` sets the column width based on the longest string's length
	if maxLength > 0 {
		f.SetColWidth(sheet, col, col, float64(maxLength)+2) // "+2" adds a bit of padding
	}
}
func exportExcel(source map[string]int64, fileName string) {
	data := map[string]float64{}
	for k, v := range source {
		data[k] = float64(v) / 1000000
	}

	log.Println(data)

	// 创建一个新的 Excel 文件
	f := excelize.NewFile()

	// 设置表头
	headers := []string{"Address", "Value"}
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1) // 列从1开始，行从1开始
		f.SetCellValue("Sheet1", cell, header)
	}

	// 填充数据
	row := 1 // 从第二行开始填充数据
	for address, value := range data {
		// 写入地址
		cell, _ := excelize.CoordinatesToCellName(1, row)
		f.SetColWidth("Sheet1", "1", "2", 50)
		f.SetCellValue("Sheet1", cell, address)

		// 写入值
		cell, _ = excelize.CoordinatesToCellName(2, row)
		f.SetCellValue("Sheet1", cell, value)

		row++
	}
	_cell1, _ := excelize.CoordinatesToCellName(1, row)
	f.SetCellValue("Sheet1", _cell1, "总计")
	_cell2, _ := excelize.CoordinatesToCellName(2, row)

	_total := 0.0
	for _, _value := range data {
		_total = _total + _value
	}

	f.SetCellValue("Sheet1", _cell2, _total)
	autoAdjustColumnWidth(f, "Sheet1", "A")
	// 创建折线图
	if err := createLineChart(f, len(data)); err != nil {
		log.Fatalf("创建折线图失败: %v", err)
	}

	// 保存文件
	if err := f.SaveAs(fileName); err != nil {
		log.Fatalf("保存文件失败: %v", err)
	}

	fmt.Println("Excel 文件已生成: output_with_chart.xlsx")
}

// 创建折线图
func createLineChart(f *excelize.File, dataLength int) error {
	// 定义图表数据范围
	categoriesRange := fmt.Sprintf("Sheet1!$A$2:$A$%d", dataLength+1) // Address 数据范围
	valuesRange := fmt.Sprintf("Sheet1!$B$2:$B$%d", dataLength+1)     // Value 数据范围

	// 定义图表系列
	series := []excelize.ChartSeries{
		{
			Name:       "Sheet1!$B$1",   // 系列名称
			Categories: categoriesRange, // X 轴数据范围
			Values:     valuesRange,     // Y 轴数据范围
		},
	}

	// 定义图表
	chart := excelize.Chart{
		Type:   excelize.Line, // 折线图
		Series: series,        // 系列数据
		Title: []excelize.RichTextRun{
			{
				Text: "统计今日冻结金额",
			},
		}, // 图表标题
		Legend: excelize.ChartLegend{Position: "bottom"}, // 图例位置
		XAxis: excelize.ChartAxis{Title: []excelize.RichTextRun{
			{
				Text: "地址",
			},
		}}, // X 轴标题
		YAxis: excelize.ChartAxis{Title: []excelize.RichTextRun{
			{
				Text: "冻结金额",
			},
		}}, // Y 轴标题
	}

	// 添加图表到工作表
	if err := f.AddChart("Sheet1", "D1", &chart); err != nil {
		return err
	}

	return nil
}

func getBalance(tAddress string) (error, int64) {
	trc20Contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT
	address := tAddress

	conn := client.NewGrpcClient("grpc.trongrid.io:50051")
	err := conn.Start(grpc.WithInsecure())

	balance, err := conn.TRC20ContractBalance(address, trc20Contract)

	if err != nil {
		return err, 0
	}
	log.Println(err)
	log.Println("余额：", balance)
	return err, balance.Int64()
}

// DeleteSysOperationRecordByIds
// @Tags      SysOperationRecord
// @Summary   批量删除SysOperationRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "批量删除SysOperationRecord"
// @Success   200   {object}  response.Response{msg=string}  "批量删除SysOperationRecord"
// @Router    /sysOperationRecord/deleteSysOperationRecordByIds [delete]
func (s *OperationRecordApi) DeleteSysOperationRecordByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = operationRecordService.DeleteSysOperationRecordByIds(IDS)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// FindSysOperationRecord
// @Tags      SysOperationRecord
// @Summary   用id查询SysOperationRecord
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     system.SysOperationRecord                                  true  "Id"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "用id查询SysOperationRecord"
// @Router    /sysOperationRecord/findSysOperationRecord [get]
func (s *OperationRecordApi) FindSysOperationRecord(c *gin.Context) {
	var sysOperationRecord system.SysOperationRecord
	err := c.ShouldBindQuery(&sysOperationRecord)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(sysOperationRecord, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	reSysOperationRecord, err := operationRecordService.GetSysOperationRecord(sysOperationRecord.ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"reSysOperationRecord": reSysOperationRecord}, "查询成功", c)
}

// GetSysOperationRecordList
// @Tags      SysOperationRecord
// @Summary   分页获取SysOperationRecord列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     request.SysOperationRecordSearch                        true  "页码, 每页大小, 搜索条件"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取SysOperationRecord列表,返回包括列表,总数,页码,每页数量"
// @Router    /sysOperationRecord/getSysOperationRecordList [get]
func (s *OperationRecordApi) GetSysOperationRecordList(c *gin.Context) {
	var pageInfo systemReq.SysOperationRecordSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//list, total, err := operationRecordService.GetSysOperationRecordInfoList(pageInfo)
	//if err != nil {
	//	global.GVA_LOG.Error("获取失败!", zap.Error(err))
	//	response.FailWithMessage("获取失败", c)
	//	return
	//}

	list2 := make([]system.SysOperationRecord, 0, 4)
	total := 1

	var record1 system.SysOperationRecord
	record1.Method = "24小时内网络预冻结以及已冻结"
	record1.ID = 1
	record1.Path = "24小时内网络预冻结以及已冻结"

	var record2 system.SysOperationRecord
	record2.Method = "24小时内波场网络已冻结"
	record2.ID = 2
	record2.Path = "24小时内波场网络已冻结.xlsx"

	var record3 system.SysOperationRecord
	record3.Method = "24小时内以太坊网络预冻结"
	record3.ID = 3
	record3.Path = "24小时内波场网络预冻结.xlsx"

	var record4 system.SysOperationRecord
	record4.Method = "24小时内以太坊网络已冻结"
	record4.ID = 4
	record4.Path = "24小时内以太坊网络已冻结.xlsx"
	//list2 = append(list2, record1, record2, record3, record4)
	list2 = append(list2, record1)
	response.OkWithDetailed(response.PageResult{
		List:     list2,
		Total:    int64(total),
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

type EthereumTX struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		BlockHash            string        `json:"blockHash"`
		BlockNumber          string        `json:"blockNumber"`
		From                 string        `json:"from"`
		Gas                  string        `json:"gas"`
		GasPrice             string        `json:"gasPrice"`
		MaxFeePerGas         string        `json:"maxFeePerGas"`
		MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
		Hash                 string        `json:"hash"`
		Input                string        `json:"input"`
		Nonce                string        `json:"nonce"`
		To                   string        `json:"to"`
		TransactionIndex     string        `json:"transactionIndex"`
		Value                string        `json:"value"`
		Type                 string        `json:"type"`
		AccessList           []interface{} `json:"accessList"`
		ChainID              string        `json:"chainId"`
		V                    string        `json:"v"`
		R                    string        `json:"r"`
		S                    string        `json:"s"`
		YParity              string        `json:"yParity"`
	} `json:"result"`
}
type EthereumContractTX struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		Address        string      `json:"address"`
		EnsName        interface{} `json:"ensName"`
		PaginatedItems []struct {
			BlockTimestamp   time.Time   `json:"blockTimestamp"`
			TransactionHash  string      `json:"transactionHash"`
			BlockNumber      string      `json:"blockNumber"`
			TransactionIndex int         `json:"transactionIndex"`
			FromAddress      string      `json:"fromAddress"`
			ToAddress        string      `json:"toAddress"`
			ContractAddress  interface{} `json:"contractAddress"`
			Value            string      `json:"value"`
			Status           string      `json:"status"`
		} `json:"paginatedItems"`
		TotalPages int `json:"totalPages"`
		TotalItems int `json:"totalItems"`
		PageNumber int `json:"pageNumber"`
	} `json:"result"`
}
type Confirm_JSONData struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  struct {
		BlockHash         string      `json:"blockHash"`
		BlockNumber       string      `json:"blockNumber"`
		ContractAddress   interface{} `json:"contractAddress"`
		CumulativeGasUsed string      `json:"cumulativeGasUsed"`
		EffectiveGasPrice string      `json:"effectiveGasPrice"`
		From              string      `json:"from"`
		GasUsed           string      `json:"gasUsed"`
		Logs              []struct {
			Address          string   `json:"address"`
			Topics           []string `json:"topics"`
			Data             string   `json:"data"`
			BlockNumber      string   `json:"blockNumber"`
			TransactionHash  string   `json:"transactionHash"`
			TransactionIndex string   `json:"transactionIndex"`
			BlockHash        string   `json:"blockHash"`
			LogIndex         string   `json:"logIndex"`
			Removed          bool     `json:"removed"`
		} `json:"logs"`
		LogsBloom        string `json:"logsBloom"`
		Status           string `json:"status"`
		To               string `json:"to"`
		TransactionHash  string `json:"transactionHash"`
		TransactionIndex string `json:"transactionIndex"`
		Type             string `json:"type"`
	} `json:"result"`
}
type EthereumERC20 struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}
