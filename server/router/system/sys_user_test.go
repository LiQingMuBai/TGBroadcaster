package system

//
//import (
//	"context"
//	"crypto/sha256"
//	"encoding/hex"
//	"encoding/json"
//	"fmt"
//	"github.com/fbsobreira/gotron-sdk/pkg/abi"
//	"os"
//
//	//"github.com/fbsobreira/gotron-sdk/pkg/abi"
//	"github.com/fbsobreira/gotron-sdk/pkg/address"
//	"github.com/fbsobreira/gotron-sdk/pkg/client"
//	"github.com/fbsobreira/gotron-sdk/pkg/common"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//	"math/big"
//	"time"
//
//	"github.com/xuri/excelize/v2"
//	"google.golang.org/grpc"
//	"io"
//	"log"
//	"net/http"
//	"strings"
//	"testing"
//
//	"github.com/mymmrac/telego"
//	tu "github.com/mymmrac/telego/telegoutil"
//)
//
//func TestTRC20_Balance(t *testing.T) {
//	trc20Contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT
//	address := "TGonTLjFnjKhA7bVrPuWfnDKSfSvEBqxUy"
//
//	conn := client.NewGrpcClient("grpc.trongrid.io:50051")
//	err := conn.Start(grpc.WithInsecure())
//	require.Nil(t, err)
//
//	balance, err := conn.TRC20ContractBalance(address, trc20Contract)
//
//	log.Println(balance)
//	assert.Nil(t, err)
//	assert.Greater(t, balance.Int64(), int64(0))
//}
//func TestSystemApi_GetSystemConfig(t *testing.T) {
//	rpcClient := client.NewGrpcClient("grpc.trongrid.io:50051")
//	err := rpcClient.Start(grpc.WithInsecure())
//
//	addr, err := address.Base58ToAddress("TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
//	log.Println(addr.String())
//
//	// 获得方法名的SHA-3
//	nameSign := common.BytesToHexString(abi.Signature("name()"))
//	result, err := rpcClient.TRC20Call("", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", nameSign, true, 0)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	data, _ := rpcClient.ParseTRC20StringProperty(common.BytesToHexString(result.GetConstantResult()[0]))
//	log.Println(data)
//	//
//	//
//	//chokhoo2024@gmail.com leo779868@gmail.com wx87110776@gmail.com
//}
//func TestUserRouter_InitUserRouter(t *testing.T) {
//	url := "https://api.trongrid.io/v1/contracts/TBPxhVAsuzoFnKyXtc1o2UySEydPHgATto/events"
//	req, _ := http.NewRequest("GET", url, nil)
//	req.Header.Add("accept", "application/json")
//	res, _ := http.DefaultClient.Do(req)
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//	fmt.Println(string(body))
//}
//
//func TestUserRouter_InitUserRouter_Login(t *testing.T) {
//	url := "https://api.trongrid.io/v1/transactions/c41e0e5cdbd6289000a36c293d97a3b868f9b8d8dbe98a7acec93d9743e465b3/events"
//	req, _ := http.NewRequest("GET", url, nil)
//	req.Header.Add("accept", "application/json")
//	res, _ := http.DefaultClient.Do(req)
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//	fmt.Println(string(body))
//	var result TronTxEvent
//	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
//		fmt.Println("Can not unmarshal JSON")
//	}
//	log.Println(result)
//
//	for index, datum := range result.Data {
//
//		if index == 1 {
//			log.Println(datum.Result.Num0)
//
//			address41 := strings.ReplaceAll(datum.Result.Num0, "0x", "41")
//			target, _ := Convert41ToTAddress(address41)
//			log.Println(target)
//
//		}
//
//	}
//}
//
//type TronEvents struct {
//	Data []struct {
//		BlockNumber           int    `json:"block_number"`
//		BlockTimestamp        int64  `json:"block_timestamp"`
//		CallerContractAddress string `json:"caller_contract_address"`
//		ContractAddress       string `json:"contract_address"`
//		EventIndex            int    `json:"event_index"`
//		EventName             string `json:"event_name"`
//		Result                struct {
//			Num0          string `json:"0"`
//			TransactionID string `json:"transactionId"`
//		} `json:"result"`
//		ResultType struct {
//			TransactionID string `json:"transactionId"`
//		} `json:"result_type"`
//		Event         string `json:"event"`
//		TransactionID string `json:"transaction_id"`
//	} `json:"data"`
//	Success bool `json:"success"`
//	Meta    struct {
//		At          int64  `json:"at"`
//		Fingerprint string `json:"fingerprint"`
//		Links       struct {
//			Next string `json:"next"`
//		} `json:"links"`
//		PageSize int `json:"page_size"`
//	} `json:"meta"`
//}
//type TronTxEvent struct {
//	Data []struct {
//		BlockNumber           int    `json:"block_number"`
//		BlockTimestamp        int64  `json:"block_timestamp"`
//		CallerContractAddress string `json:"caller_contract_address"`
//		ContractAddress       string `json:"contract_address"`
//		EventIndex            int    `json:"event_index"`
//		EventName             string `json:"event_name"`
//		Result                struct {
//			Num0          string `json:"0"`
//			TransactionID string `json:"transactionId"`
//		} `json:"result"`
//		ResultType struct {
//			TransactionID string `json:"transactionId"`
//		} `json:"result_type"`
//		Event         string `json:"event"`
//		TransactionID string `json:"transaction_id"`
//	} `json:"data"`
//	Success bool `json:"success"`
//	Meta    struct {
//		At       int64 `json:"at"`
//		PageSize int   `json:"page_size"`
//	} `json:"meta"`
//}
//
//func TestUserRouter_InitUserRouter_Logout(t *testing.T) {
//	url := "https://api.trongrid.io/v1/contracts/TBPxhVAsuzoFnKyXtc1o2UySEydPHgATto/events"
//	req, _ := http.NewRequest("GET", url, nil)
//	req.Header.Add("accept", "application/json")
//	res, _ := http.DefaultClient.Do(req)
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//	fmt.Println(string(body))
//	var result TronEvents
//	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
//		fmt.Println("Can not unmarshal JSON")
//	}
//	log.Println(result)
//
//	var txIDMap map[string]string /*创建集合 */
//	txIDMap = make(map[string]string)
//
//	for _, datum := range result.Data {
//		log.Println(datum.EventName)
//		log.Println(datum.TransactionID)
//		txIDMap[datum.TransactionID] = datum.EventName
//
//	}
//	log.Println("===============================================================================")
//
//	for txID := range txIDMap {
//		log.Println(txID)
//		log.Println(txIDMap[txID])
//	}
//	//https://www.trongrid.io/v1/api/trc20/getAddressBalance?address=TGDsEr2cSRC98Zo9WnwNDik2Y5rdboPRvd
//	//https://api.trongrid.io/v1/accounts/TYPrKF2sevXuE86Xo3Y2mhFnjseiUcybny/transactions/trc20?limit=100&contract_address=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
//}
//
//type TransactionInfo struct {
//	Ret []struct {
//		ContractRet string `json:"contractRet"`
//	} `json:"ret"`
//	Signature []string `json:"signature"`
//	TxID      string   `json:"txID"`
//	RawData   struct {
//		Contract []struct {
//			Parameter struct {
//				Value struct {
//					Data            string `json:"data"`
//					OwnerAddress    string `json:"owner_address"`
//					ContractAddress string `json:"contract_address"`
//				} `json:"value"`
//				TypeURL string `json:"type_url"`
//			} `json:"parameter"`
//			Type string `json:"type"`
//		} `json:"contract"`
//		RefBlockBytes string `json:"ref_block_bytes"`
//		RefBlockHash  string `json:"ref_block_hash"`
//		Expiration    int64  `json:"expiration"`
//		FeeLimit      int    `json:"fee_limit"`
//		Timestamp     int64  `json:"timestamp"`
//	} `json:"raw_data"`
//	RawDataHex string `json:"raw_data_hex"`
//}
//
//func TestGetTransactionById(t *testing.T) {
//
//	url := "https://api.trongrid.io/v1/contracts/TBPxhVAsuzoFnKyXtc1o2UySEydPHgATto/events"
//	req, _ := http.NewRequest("GET", url, nil)
//	req.Header.Add("accept", "application/json")
//	res, _ := http.DefaultClient.Do(req)
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//	//fmt.Println(string(body))
//	var events TronEvents
//	if err := json.Unmarshal(body, &events); err != nil { // Parse []byte to go struct pointer
//		fmt.Println("Can not unmarshal JSON")
//	}
//	//log.Println(events)
//
//	var txIDMap map[string]string /*创建集合 */
//	txIDMap = make(map[string]string)
//
//	for _, datum := range events.Data {
//		//log.Println(datum.EventName)
//		//log.Println(datum.TransactionID)
//		txIDMap[datum.TransactionID] = datum.EventName
//	}
//
//	sumbitMap := make(map[string]int64)
//	commitMap := make(map[string]int64)
//	for txID := range txIDMap {
//		//log.Println(txID)
//		//log.Println(txIDMap[txID])
//
//		result := getTransactionData(txID)
//		//log.Println(result.RawDataHex)
//
//		//submit
//
//		time.Sleep(3 * time.Second)
//		if len(result.RawDataHex) > 600 {
//			//log.Println("submit")
//			tAddress := getTronAddress(result)
//			_, _amount := getBalance(tAddress)
//			//log.Println("balance", tAddress, _amount)
//
//			sumbitMap[tAddress] = _amount
//		} else {
//			//已經拉入黑名單
//			//log.Println("commit")
//			_address, _amount := getCommitAddressBalance(txID)
//			//log.Println("balance", _address, _amount)
//			commitMap[_address] = _amount
//		}
//		//commit
//	}
//	log.Println("==========================已提交的====================================")
//	for tx, _balance := range sumbitMap {
//		log.Println(tx, _balance)
//
//	}
//	exportExcel(sumbitMap, "今日预冻结.xlsx")
//	log.Println("==========================接下来是确认的===============================")
//	for _committx, _commitbalance := range commitMap {
//		log.Println(_committx, _commitbalance)
//
//	}
//
//	time.Sleep(3 * time.Second)
//	//filePath1 := "C:\\Users\\Administrator\\Documents\\shiled-platform\\server\\router\\system\\今日预冻结.xlsx"
//	//sendTelegram(filePath1)
//
//	exportExcel(commitMap, "今日已冻结.xlsx")
//	time.Sleep(3 * time.Second)
//	//filePath2 := "C:\\Users\\Administrator\\Documents\\shiled-platform\\server\\router\\system\\今日已冻结.xlsx"
//	//sendTelegram(filePath2)
//}
//
//func sendTelegram(_filePath string) {
//	ctx := context.Background()
//	//botToken := os.Getenv("TOKEN")
//	botToken := "7668068911:AAFOXuA7KpWOfur0rcoVbZTwGOgsBCjkI3s"
//	//botToken := "7668068911:AAFOXuA7KpWOfur0rcoVbZTwGOgsBCjkI3s"
//	//chatID := -4657809905
//	filePath := _filePath
//	// Note: Please keep in mind that default logger may expose sensitive information, use in development only
//	bot, err := telego.NewBot(botToken, telego.WithDefaultDebugLogger())
//	if err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//
//	// Document parameters
//	document := tu.Document(
//		// Chat ID as Integer
//		tu.ID(int64(-4657809905)),
//
//		// Send using file from disk
//		tu.File(mustOpen(filePath)),
//
//		// Send using external URL
//		// tu.FileFromURL("https://example.com/my_file.txt"),
//
//		// Send using file ID
//		// tu.FileFromID("<file ID of your file>"),
//	).WithCaption("來自於U盾情報部")
//
//	// Sending document
//	msg, err := bot.SendDocument(ctx, document)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Println(msg.Document)
//
//}
//
//// Helper function to open file or panic
//func mustOpen(filename string) *os.File {
//	file, err := os.Open(filename)
//	if err != nil {
//		panic(err)
//	}
//	return file
//}
//
//func getTransactionData(_txid string) TransactionInfo {
//	//txid := ""
//	url := "https://api.trongrid.io/walletsolidity/gettransactionbyid"
//
//	payload := strings.NewReader("{\"value\":\"" + _txid + "\"}")
//
//	req, _ := http.NewRequest("POST", url, payload)
//
//	req.Header.Add("accept", "application/json")
//	req.Header.Add("content-type", "application/json")
//
//	res, _ := http.DefaultClient.Do(req)
//
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//
//	var result TransactionInfo
//	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
//		fmt.Println("Can not unmarshal JSON")
//	}
//	return result
//}
//
//func getCommitAddressBalance(txid string) (string, int64) {
//	url := "https://api.trongrid.io/v1/transactions/" + txid + "/events"
//	req, _ := http.NewRequest("GET", url, nil)
//	req.Header.Add("accept", "application/json")
//	res, _ := http.DefaultClient.Do(req)
//	defer res.Body.Close()
//	body, _ := io.ReadAll(res.Body)
//	//fmt.Println(string(body))
//	var result TronTxEvent
//	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
//		fmt.Println("Can not unmarshal JSON")
//	}
//	//log.Println(result)
//
//	for index, datum := range result.Data {
//
//		if index == 1 {
//			//log.Println(datum.Result.Num0)
//
//			address41 := strings.ReplaceAll(datum.Result.Num0, "0x", "41")
//			target, _ := Convert41ToTAddress(address41)
//			//log.Println(target)
//			_, amount := getBalance(target)
//
//			return target, amount
//		}
//	}
//	return "", 0
//}
//
//func getTronAddress(result TransactionInfo) string {
//	address41 := result.RawDataHex[558 : 558+42]
//	tAddress, _ := Convert41ToTAddress(address41)
//	//log.Println("41address ", address41)
//	//log.Println("taddress ", tAddress)
//	return tAddress
//}
//
//func getBalance(tAddress string) (error, int64) {
//	trc20Contract := "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" // USDT
//	address := tAddress
//
//	conn := client.NewGrpcClient("grpc.trongrid.io:50051")
//	err := conn.Start(grpc.WithInsecure())
//
//	balance, err := conn.TRC20ContractBalance(address, trc20Contract)
//
//	log.Println("余额：", balance)
//	return err, balance.Int64()
//}
//func TestSysExportTemplateRouter_InitSysExportTemplateRouter(t *testing.T) {
//	f := excelize.NewFile()
//	defer func() {
//		if err := f.Close(); err != nil {
//			fmt.Println(err)
//		}
//	}()
//	for idx, row := range [][]interface{}{
//		{nil, "Apple", "Orange", "Pear"}, {"Small", 2, 3, 3},
//		{"Normal", 5, 2, 4}, {"Large", 6, 7, 8},
//	} {
//		cell, err := excelize.CoordinatesToCellName(1, idx+1)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		f.SetSheetRow("Sheet1", cell, &row)
//	}
//	if err := f.AddChart("Sheet1", "E1", &excelize.Chart{
//		Type: excelize.Col3DClustered,
//		Series: []excelize.ChartSeries{
//			{
//				Name:       "Sheet1!$A$2",
//				Categories: "Sheet1!$B$1:$D$1",
//				Values:     "Sheet1!$B$2:$D$2",
//			},
//			{
//				Name:       "Sheet1!$A$3",
//				Categories: "Sheet1!$B$1:$D$1",
//				Values:     "Sheet1!$B$3:$D$3",
//			},
//			{
//				Name:       "Sheet1!$A$4",
//				Categories: "Sheet1!$B$1:$D$1",
//				Values:     "Sheet1!$B$4:$D$4",
//			}},
//		Title: []excelize.RichTextRun{
//			{
//				Text: "Fruit 3D Clustered Column Chart",
//			},
//		},
//	}); err != nil {
//		fmt.Println(err)
//		return
//	}
//	// Save spreadsheet by the given path.
//	if err := f.SaveAs("Book1.xlsx"); err != nil {
//		fmt.Println(err)
//	}
//}
//
//func TestConvertTronAddress(t *testing.T) {
//	//address41 := "41a614f803b6fd780986a42c78ec9c7f77e6ded13c" // Replace with your 41-prefixed address
//	//address41 := "41e196f2f5e8eae1139dab2fb9d683c7cac6736b57"
//	address41 := strings.ReplaceAll("0x51ec29238e6dbbddcda8c1ba4f46f698323e5cf8", "0x", "41")
//
//	tAddress, err := Convert41ToTAddress(address41)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	fmt.Println("41 Address:", address41)
//	fmt.Println("T Address:", tAddress)
//}
//
//const (
//	// Tron address prefix (hex)
//	tronAddressPrefix = "41"
//	// Base58 alphabet
//	base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
//)
//
//// Convert a 41-prefixed address to a T-prefixed address
//func Convert41ToTAddress(address41 string) (string, error) {
//	// Step 1: Validate the input address
//	if len(address41) < 2 || address41[:2] != tronAddressPrefix {
//		return "", fmt.Errorf("invalid 41-prefixed address")
//	}
//
//	// Step 2: Decode the hex address into bytes
//	addrBytes, err := hex.DecodeString(address41)
//	if err != nil {
//		return "", fmt.Errorf("failed to decode hex address: %v", err)
//	}
//
//	// Step 3: Compute checksum (double SHA256 of the address bytes)
//	sha256Hash := sha256.Sum256(addrBytes)
//	sha256Hash2 := sha256.Sum256(sha256Hash[:])
//	checksum := sha256Hash2[:4]
//
//	// Step 4: Combine address bytes and checksum
//	finalAddress := append(addrBytes, checksum...)
//
//	// Step 5: Encode the result in Base58
//	tAddress := base58Encode(finalAddress)
//	return tAddress, nil
//}
//
//// Base58 encoding function
//func base58Encode(input []byte) string {
//	// Convert the byte array to a big integer
//	x := new(big.Int).SetBytes(input)
//	base := big.NewInt(58)
//	zero := big.NewInt(0)
//	mod := new(big.Int)
//
//	var result string
//	for x.Cmp(zero) > 0 {
//		x.DivMod(x, base, mod)
//		result = string(base58Alphabet[mod.Int64()]) + result
//	}
//
//	// Add leading '1's for leading zeros in the input
//	for _, b := range input {
//		if b != 0 {
//			break
//		}
//		result = "1" + result
//	}
//
//	return result
//}
//
//func TestExportExcel(t *testing.T) {
//
//	// 示例数据：map[string]interface{}
//	sumbitMap := map[string]interface{}{
//		"address1": 33333,
//		"address2": 22222,
//		"address3": 11111,
//	}
//
//	exportExcel(sumbitMap, "24小时内以太坊预冻结.xlsx")
//
//}
//
////
////func TestExport2(t *testing.T) {
////	// 示例数据：map[string]int64
////
////	source := map[string]int64{
////		"address1": 333333333,
////		"address2": 22222222,
////		"address3": 1111111111,
////		"address4": 88888888,
////		"address5": 77777777,
////		"address6": 5555555,
////	}
////
////	exportExcel(source, "output.xlsx")
////}
//
////
////func exportExcel(source map[string]int64, fileName string) {
////	data := map[string]float64{}
////	for k, v := range source {
////		data[k] = float64(v) / 1000000
////	}
////
////	log.Println(data)
////
////	// 创建一个新的 Excel 文件
////	f := excelize.NewFile()
////
////	// 设置表头
////	headers := []string{"Address", "Value"}
////	for col, header := range headers {
////		cell, _ := excelize.CoordinatesToCellName(col+1, 1) // 列从1开始，行从1开始
////		f.SetCellValue("Sheet1", cell, header)
////	}
////
////	// 填充数据
////	row := 2 // 从第二行开始填充数据
////	for address, value := range data {
////		// 写入地址
////		cell, _ := excelize.CoordinatesToCellName(1, row)
////		f.SetCellValue("Sheet1", cell, address)
////
////		// 写入值
////		cell, _ = excelize.CoordinatesToCellName(2, row)
////		f.SetCellValue("Sheet1", cell, value)
////
////		row++
////	}
////	_cell1, _ := excelize.CoordinatesToCellName(1, row)
////	f.SetCellValue("Sheet1", _cell1, "总计")
////	_cell2, _ := excelize.CoordinatesToCellName(2, row)
////
////	_total := 0.0
////	for _, _value := range data {
////		_total = _total + _value
////	}
////
////	f.SetCellValue("Sheet1", _cell2, _total)
////
////	// 创建折线图
////	if err := createLineChart(f, len(data)); err != nil {
////		log.Fatalf("创建折线图失败: %v", err)
////	}
////
////	// 保存文件
////	if err := f.SaveAs(fileName); err != nil {
////		log.Fatalf("保存文件失败: %v", err)
////	}
////
////	fmt.Println("Excel 文件已生成: output_with_chart.xlsx")
////}
//
//// 创建折线图
//func createLineChart(f *excelize.File, dataLength int) error {
//	// 定义图表数据范围
//	categoriesRange := fmt.Sprintf("Sheet1!$A$2:$A$%d", dataLength+1) // Address 数据范围
//	valuesRange := fmt.Sprintf("Sheet1!$B$2:$B$%d", dataLength+1)     // Value 数据范围
//
//	// 定义图表系列
//	series := []excelize.ChartSeries{
//		{
//			Name:       "Sheet1!$B$1",   // 系列名称
//			Categories: categoriesRange, // X 轴数据范围
//			Values:     valuesRange,     // Y 轴数据范围
//		},
//	}
//
//	// 定义图表
//	chart := excelize.Chart{
//		Type:   excelize.Line, // 折线图
//		Series: series,        // 系列数据
//		Title: []excelize.RichTextRun{
//			{
//				Text: "统计今日冻结金额",
//			},
//		}, // 图表标题
//		Legend: excelize.ChartLegend{Position: "bottom"}, // 图例位置
//		XAxis: excelize.ChartAxis{Title: []excelize.RichTextRun{
//			{
//				Text: "地址",
//			},
//		}}, // X 轴标题
//		YAxis: excelize.ChartAxis{Title: []excelize.RichTextRun{
//			{
//				Text: "冻结金额",
//			},
//		}}, // Y 轴标题
//	}
//
//	// 添加图表到工作表
//	if err := f.AddChart("Sheet1", "D1", &chart); err != nil {
//		return err
//	}
//
//	return nil
//}
