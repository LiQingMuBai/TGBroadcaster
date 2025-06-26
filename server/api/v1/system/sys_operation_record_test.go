package system

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"
)

func TestGetEthereumAddress(t *testing.T) {
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
		return
	}
	reqBody := strings.NewReader(string(reqParam))
	url := "https://alpha-alien-sheet.quiknode.pro/88f18a5e3da4679705954edbb2859e5144c16a5a/"
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
	// time.Sleep(1 * time.Second)
	sumbitMap := make(map[string]int64)
	commitMap := make(map[string]int64)
	for _, tx := range txs.Result.PaginatedItems {
		log.Println(tx.TransactionHash)
		//获取交易tx hash
		_txHash := tx.TransactionHash
		// time.Sleep(1 * time.Second)
		_address := getPeddingBlackedAddress(_txHash)

		if len(_address) > 0 {
			//说明是pendding的地址，直接获取余额
			log.Println("待定黑名单地址：", _address)
			// time.Sleep(1 * time.Second)
			balance, err := getUSDTBalance(_address)
			if err != nil {
			}
			if balance > 0 {
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
				if balance > 0 {
					commitMap[_address] = balance
					log.Println(balance)
				}
			}
		}
	}

	//export excel

	log.Println("==========================已提交的====================================")
	for tx, _balance := range sumbitMap {
		log.Println(tx, _balance)

	}
	exportExcel(sumbitMap, "24小时内以太坊预冻结.xlsx")

	// time.Sleep(1 * time.Second)
	//filePath1 := "C:\\Users\\Administrator\\Documents\\shiled-platform\\server\\api\\v1\\system\\24小时内以太坊预冻结.xlsx"
	//filePath1 := "/soft/shiled-platform/server/今日预冻结.xlsx"
	//sendTelegram(filePath1)
	log.Println("==========================已确认的====================================")
	for tx, _balance := range commitMap {
		log.Println(tx, _balance)

	}
	exportExcel(commitMap, "24小时内以太坊已冻结.xlsx")
	// time.Sleep(1 * time.Second)
	//filePath2 := "C:\\Users\\Administrator\\Documents\\shiled-platform\\server\\api\\v1\\system\\24小时内以太坊已冻结.xlsx"
	//filePath2 := "/soft/shiled-platform/server/今日已冻结.xlsx"
	//sendTelegram(filePath2)

}
func TestSysExportTemplateApi_FindSysExportTemplate(t *testing.T) {
	commitMap := make(map[string]int64)
	commitMap["0xd3c7223ccEb888620DD90141e543FB62E0Dc7c7b"] = 1238811111
	commitMap["0xd3c7223ccEb888620DD92141e543FB62E0Dc7c7b"] = 32388122
	commitMap["0xd3c7223ccEb8886201D90141e543FB62E0Dc7c7b"] = 16388111112
	ExportExcel2(commitMap, "online 24hours.xlsx")
}

func TestExport2(t *testing.T) {
	// 示例数据：map[string]int64

	source := map[string]int64{
		"address1": 333333333,
		"address2": 22222222,
		"address3": 1111111111,
		"address4": 88888888,
		"address5": 77777777,
		"address6": 5555555,
	}

	exportExcel(source, "output.xlsx")
}
func ExportExcel2(source map[string]int64, fileName string) {
	data := map[string]float64{}
	for k, v := range source {
		data[k] = float64(v) / 1000000
	}

	log.Println(data)

	// 创建一个新的 Excel 文件
	f := excelize.NewFile()

	// 设置表头
	headers := []string{"Address", "Value"}

	//f.SetColWidth("Sheet1", "A1", "B2", 200)
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1) // 列从1开始，行从1开始
		f.SetCellValue("Sheet1", cell, header)

		//f.SetColWidth("Sheet1", "1", "1", 550)
	}

	// 填充数据
	row := 2 // 从第二行开始填充数据
	for address, value := range data {
		// 写入地址
		cell, _ := excelize.CoordinatesToCellName(1, row)

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

func TestGetTransactionsByAddress(t *testing.T) {

	_txHash := "0x6beba993b3c70e7df7724506ad8a1882b10bbf4bbcf42213bc4341e6af1549e2"
	peddingAddress := getPeddingBlackedAddress(_txHash)

	log.Println(peddingAddress)

}

func TestConfirm(t *testing.T) {

	_txHash := "0x435f6aa68d847daec0db68620f3fefffe77192d13b017a5493d175c736ca81f7"
	address := getBlackAddress(_txHash)

	log.Println(address)
}

func TestGetBalance(t *testing.T) {

	_address := "0x3cadC4144CCDC15891C1f6F77A09E5d341fcC01A"

	balance, err := getUSDTBalance(_address)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(balance)
}
