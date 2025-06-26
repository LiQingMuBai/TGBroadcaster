package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/magiconair/properties/assert"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

//var currentKeyIndex uint32

func TestGetTRXBalance(t *testing.T) {

	key := "2e81d257-d742-4c96-be75-baf131c1e12f,87de4692-0464-4ba4-910c-b9d160d7e2c2,bebe7987-c941-4dff-95d1-b396b1ce81a5,a700e0d9-5693-481d-bc5a-5a4104e63096,ba950e04-c2b6-429c-86f7-1924ab5bb8a0,b815c662-6d1b-40b3-b9d4-fae8031de8dc,812f0bf0-8508-4f30-967b-7794e8e5ebd7,0a197976-6a7b-49a4-91ca-138467d4dc82,4d410e56-02c5-4180-be9f-99f853efe54e,e6b8b248-2ce5-4bd6-ad31-51960b97d575,3b766c89-9192-4b8e-8e13-756da7a43991,11ab7729-05ca-45b0-8c26-9632ce3ed3f4,d6e84e55-69d3-4035-bacf-7f84fa487c14,a78647e3-bb80-4655-8c70-a89e8f23ca6b,b23f5b61-3368-4002-94d3-70602571ffb6,40c7c7bf-6816-436d-acc9-58ceef9f4bfb,da724852-53de-435b-b77a-20770f5491d8,15a26b20-8fcf-45a2-a9bc-6ed3ecd00822,db9c3af7-9548-42c4-ad2a-83b056b895c1,6aee69ed-02bc-4839-854c-3e2e1bc73680,87de4692-0464-4ba4-910c-b9d160d7e2c2"

	apiKeys := strings.Split(key, ",")

	for i := 0; i < 100; i++ {
		// 简单轮询
		keyIndex := atomic.AddUint32(&currentKeyIndex, 1) % uint32(len(apiKeys))
		currentKey := apiKeys[keyIndex]

		log.Println(currentKey)

		time.Sleep(100 * time.Millisecond)

		req, err := http.NewRequest("GET", "https://api.trongrid.io/wallet/gettransactioninfobyid?value=559d90ff21c85fe5764bffb1a2af98b84f6af3e5aadd208a141608084c576871", nil)
		// 设置请求头
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("TRON-PRO-API-KEY", currentKey) // 添加API密钥
		client := &http.Client{}
		//resp, err := client.Get(url)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Errorf("API请求失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
		}
		//time.Sleep(3 * time.Second)
		//json.NewDecoder(resp.Body).Decode(result)

		fmt.Println(resp.StatusCode)
		// 读取响应
		//body, _ := ioutil.ReadAll(resp.Body)
		var response struct {
			BlockNumber int64 `json:"blockNumber"`
			Logs        []struct {
				Address string   `json:"address"`
				Topics  []string `json:"topics"`
				Data    string   `json:"data"`
			} `json:"log"`
		}
		var result = response
		json.NewDecoder(resp.Body).Decode(&result)

		fmt.Println(result)
	}

	//// 创建HTTP请求

}

func TestTronHexToBase58(t *testing.T) {

	usdt, _ := TronHexToBase58("a614f803b6fd780986a42c78ec9c7f77e6ded13c")
	//contract, _ := tron.TronHexToBase58(log.Address)
	assert.Equal(t, usdt, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t", "success")
}

func TestTronBase58ToHex(t *testing.T) {

	tronHex := "000000000000000000000000302ac98a7e1db6d55cb77188c4528ac0cd42379b"[24:64]
	fmt.Println(tronHex)
	base58Addr, _ := HexToTronBase58("41" + tronHex)
	fmt.Println(base58Addr)

}

func TestTronClient_GetAccountResources(t *testing.T) {
	//TRON_FULL_NODE := "https://api.trongrid.io"
	TRON_FULL_NODE := "https://young-clean-orb.tron-mainnet.quiknode.pro/9283c9ddb51102d9d22cf1ac5a6e6fc898eeaf77"
	tronClient := NewTronClient(TRON_FULL_NODE)

	_address, _ := Base58ToTronHex("TRLi4KskHxfACnxKzLrf2wNHauWUQpwziN")
	resources, _ := tronClient.GetAccountResources(_address)
	//resources2, _ := tronClient.GetAccountResources("TPiTNfSJdk2Ao12eCDESahv3ZN4npdGBWQ")

	fmt.Println(resources)
	//fmt.Println(resources2)

}

func TestTransferNative(t *testing.T) {
	TRON_FULL_NODE := "https://nile.trongrid.io"
	sendAmount := utils.ConvertFloatToBigInt(100.01, 6)

	tronClient := NewTronClient(TRON_FULL_NODE)

	tronClient.TransferNative(context.Background(), "43b8e682fd65cfc5fd0a67d0caf6c5451e271aacb2f055d6c5f2c429470e0e23", "TQtp1QHk5H1ccnKy1eorR3UDbFNj5UDjpp", sendAmount)
}
func TestTronClient_BatchGetAddressBalances(t *testing.T) {
	//TRON_FULL_NODE := "https://api.trongrid.io"
	TRON_FULL_NODE := "https://sleek-thrumming-orb.tron-mainnet.quiknode.pro/d96d2a61c8c096e66c906225f972c44c779199f8/"
	//TRON_FULL_NODE := "https://young-clean-orb.tron-mainnet.quiknode.pro/9283c9ddb51102d9d22cf1ac5a6e6fc898eeaf77/"
	tronClient := NewTronClient(TRON_FULL_NODE)

	addresses := []string{"TRLi4KskHxfACnxKzLrf2wNHauWUQpwziN", "TXtNKWqibAqaFE3HRTKBovUyEiaD8S3Kb9", "TVaHLSfHvdqxCAaiYY5eY5tm1ouHL9zDU1"}
	addressesValues, err := tronClient.BatchGetAddressBalances(addresses)
	if err != nil {
	}
	for _address, _amount := range addressesValues {

		log.Println(_address)
		log.Println(_amount.String())

	}
}
func TestTronClient_GetNativeBalance(t *testing.T) {
	//TRON_FULL_NODE := "https://api.trongrid.io"
	TRON_FULL_NODE := "https://young-clean-orb.tron-mainnet.quiknode.pro/9283c9ddb51102d9d22cf1ac5a6e6fc898eeaf77/"
	tronClient := NewTronClient(TRON_FULL_NODE)

	//okex hot wallet
	_address, _ := Base58ToTronHex("TRLi4KskHxfACnxKzLrf2wNHauWUQpwziN")
	fmt.Println(_address)
	resources, _ := tronClient.GetNativeBalance(context.Background(), _address)

	fmt.Println(resources.Int64())

	balance := utils.DivideWithPrecision(resources, 4)
	fmt.Println(balance)

	trxBalance, err := tronClient.GetNativeBalance(context.Background(), "4196c784e3985d35a7133b2a33670871b09e8f86ea")
	if err != nil {
		fmt.Errorf("query master's address failed : %w", err)
	}

	fmt.Println(trxBalance)
	//fmt.Println(resources.Int64() / 1_000_000)
}

func TestGetTronAddressFromPrivateKey(t *testing.T) {
	names := "TG4nheHr9ZsBms6tn3seAsKwGXyG9XWJsV"

	// Split into a slice
	addresses := strings.Split(names, ",")

	for _, address := range addresses {

		log.Println(address)
	}
}

func TestTronClient_GetUSDTTransferCount(t *testing.T) {
	TRON_FULL_NODE := "https://api.trongrid.io"
	//TRON_FULL_NODE := "https://young-clean-orb.tron-mainnet.quiknode.pro/9283c9ddb51102d9d22cf1ac5a6e6fc898eeaf77/"
	tronClient := NewTronClient(TRON_FULL_NODE)

	_count, err := tronClient.GetUSDTTransferCount("TG4nheHr9ZsBms6tn3seAsKwGXyG9XWJsV", "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")

	if err != nil {

	}

	fmt.Println(_count)

	_total, err := getUSDTTransferCount("TDoXUNZ6PajKuiUkcYg3EDSV9bnqGqsbcf")

	if err != nil {
	}

	fmt.Println(_total)
}

func TestTronClient_GetLatestBlock(t *testing.T) {
	TRON_FULL_NODE := "https://api.trongrid.io"

	tronClient := NewTronClient(TRON_FULL_NODE)

	nowBlock := tronClient.GetLatestBlock()

	fmt.Println(nowBlock)
}

// Base TronGrid API endpoint
const apiTemplate = "https://api.trongrid.io/v1/accounts/%s/transactions/trc20"

// USDT TRC20 contract address on TRON
const usdtContract = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"

// getUSDTTransferCount returns the number of USDT transfers *to* the given address.
func getUSDTTransferCount(address string) (int, error) {
	apiURL := fmt.Sprintf(apiTemplate, address)
	params := url.Values{}
	params.Set("only_from", "true") // Only incoming transfers
	params.Set("limit", "1")        // We only care about the total
	params.Set("contract_address", usdtContract)

	fullURL := apiURL + "?" + params.Encode()

	log.Println(fullURL)
	resp, err := http.Get(fullURL)
	if err != nil {
		return 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	var response Trc20TxResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return 0, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return response.Meta.Total, nil
}

// Structure for the API response
type Trc20TxResponse struct {
	Data []struct {
		TokenInfo struct {
			Address string `json:"address"`
			Symbol  string `json:"symbol"`
		} `json:"token_info"`
	} `json:"data"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}
