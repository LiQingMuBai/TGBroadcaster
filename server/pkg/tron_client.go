package pkg

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/shopspring/decimal"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type TronClient struct {
	rpcURL     string
	httpClient *http.Client
}

func NewTronClient(rpcURL string) *TronClient {
	return &TronClient{
		rpcURL:     strings.TrimRight(rpcURL, "/"),
		httpClient: &http.Client{},
	}
}

var (
	currentKeyIndex uint32
)

func (c *TronClient) doRequest(ctx context.Context, method, path string, payload interface{}) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.rpcURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("received non-OK status: %d, url, %s, body: %s", resp.StatusCode, path, string(respBody))
	}
	return respBody, nil
}
func (c *TronClient) GetAccountResources(address string) (*AccountResources, error) {
	url := fmt.Sprintf("%s/wallet/getaccountresource?address=%s", c.rpcURL, address)

	fmt.Println(url)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		FreeNetLimit int64 `json:"freeNetLimit"`
		AssetNetUsed []struct {
			Key   string `json:"key"`
			Value int64  `json:"value"`
		} `json:"assetNetUsed"`
		AssetNetLimit []struct {
			Key   string `json:"key"`
			Value int64  `json:"value"`
		} `json:"assetNetLimit"`
		TotalNetLimit     int64 `json:"TotalNetLimit"`
		TotalNetWeight    int64 `json:"TotalNetWeight"`
		TotalEnergyLimit  int64 `json:"TotalEnergyLimit"`
		TotalEnergyWeight int64 `json:"TotalEnergyWeight"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &AccountResources{
		Energy:         result.TotalEnergyLimit - result.TotalNetWeight,
		EnergyLimit:    result.TotalEnergyLimit,
		Bandwidth:      result.TotalNetLimit - result.TotalNetWeight,
		BandwidthLimit: result.FreeNetLimit,
	}, nil
}

type TronResource struct {
	EnergyLimit    int64 // 总能量
	FreeEnergy     int64 // 可用能量
	EnergyUsed     int64 // 已用能量
	BandwidthLimit int64 // 总带宽
	BandwidthUsed  int64 // 已用带宽
	FreeBandwidth  int64 // 剩余免费带宽
}

func (c *TronClient) GetAccountResourcesViaHTTP(address string) (*TronResource, error) {
	url := fmt.Sprintf("https://api.trongrid.io/wallet/getaccountresource?address=%s", address)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		EnergyLimit  int64 `json:"EnergyLimit"`
		EnergyUsed   int64 `json:"EnergyUsed"`
		NetLimit     int64 `json:"NetLimit"`
		NetUsed      int64 `json:"NetUsed"`
		FreeNetLimit int64 `json:"FreeNetLimit"`
		FreeNetUsed  int64 `json:"FreeNetUsed"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &TronResource{
		EnergyLimit:    result.EnergyLimit,
		EnergyUsed:     result.EnergyUsed,
		FreeEnergy:     result.EnergyLimit - result.EnergyUsed,
		BandwidthLimit: result.NetLimit,
		BandwidthUsed:  result.NetUsed,
		FreeBandwidth:  result.FreeNetLimit - result.FreeNetUsed,
	}, nil
}

func (c *TronClient) FetchAccountData(ctx context.Context, address string) (*AccountResponse, error) {
	respBody, err := c.doRequest(ctx, "GET", fmt.Sprintf("/v1/accounts/%s", address), nil)
	if err != nil {
		return nil, err
	}
	var accountResp AccountResponse
	if err := json.Unmarshal(respBody, &accountResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}
	if !accountResp.Success {
		return nil, errors.New("no account data returned or success false")
	}
	return &accountResp, nil
}

func (c *TronClient) CreateTransaction(ctx context.Context, fromAddress, toAddress string, amount *big.Int) (*Transaction, error) {
	reqPayload := TransactionCreateRequest{
		OwnerAddress: fromAddress,
		ToAddress:    toAddress,
		Amount:       int(amount.Int64()),
		Visible:      true,
	}
	respBody, err := c.doRequest(ctx, "POST", "/wallet/createtransaction", reqPayload)
	if err != nil {
		return nil, err
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
		return nil, errors.New(errResp.Error)
	}
	var tx Transaction
	if err := json.Unmarshal(respBody, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction response: %w", err)
	}
	if tx.TxID == "" {
		return nil, errors.New("transaction ID is empty in response")
	}
	return &tx, nil
}

func (c *TronClient) SignTransaction(tx *Transaction, privateKey string) error {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}
	rawDataBytes, err := hexutil.Decode("0x" + tx.RawDataHex)
	if err != nil {
		return fmt.Errorf("unable to decode rawDataHex: %w", err)
	}
	hash := sha256.New()
	hash.Write(rawDataBytes)
	txHash := hash.Sum(nil)
	signature, err := crypto.Sign(txHash, privKey)
	if err != nil {
		return fmt.Errorf("unable to sign tx: %w", err)
	}
	tx.Signature = append(tx.Signature, hexutil.Encode(signature)[2:])
	return nil
}

func (c *TronClient) BroadcastTransaction(ctx context.Context, tx *Transaction) (common.Hash, error) {
	respBody, err := c.doRequest(ctx, "POST", "/wallet/broadcasttransaction", tx)
	if err != nil {
		return common.Hash{}, err
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
		return common.Hash{}, errors.New(errResp.Error)
	}
	var broadcastResp BroadcastTransactionResponse
	if err := json.Unmarshal(respBody, &broadcastResp); err != nil {
		return common.Hash{}, errors.New("failed to unmarshal broadcast response")
	}
	return common.HexToHash(broadcastResp.TxID), nil
}

// Returns wallet's TRX balance
func (c *TronClient) GetNativeBalance(ctx context.Context, address string) (*big.Int, error) {
	accountResp, err := c.FetchAccountData(ctx, address)
	if err != nil {
		return nil, err
	}
	if len(accountResp.Data) == 0 {
		return big.NewInt(0), nil
	}
	return big.NewInt(accountResp.Data[0].Balance), nil
}

// Returns wallet balance for given TRC20 contract.
func (c *TronClient) GetTokenBalance(ctx context.Context, address, tokenAddress string) (*big.Int, error) {
	accountResp, err := c.FetchAccountData(ctx, address)
	if err != nil {
		return nil, err
	}
	if !accountResp.Success || len(accountResp.Data) == 0 {
		return nil, errors.New("no account data returned or success false")
	}
	for _, token := range accountResp.Data[0].TRC20 {
		if val, ok := token[tokenAddress]; ok {
			balance := new(big.Int)
			_, ok := balance.SetString(val, 10)
			if !ok {
				return nil, fmt.Errorf("failed to parse token balance: %s", val)
			}
			return balance, nil
		}
	}
	return nil, fmt.Errorf("token %s not found in account %s", tokenAddress, address)
}

func (c *TronClient) TransferNative(
	ctx context.Context,
	senderPrivateKey,
	toAddress string,
	amount *big.Int,
) (common.Hash, error) {
	senderAddress, err := GetTronAddressFromPrivateKey(senderPrivateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get sender address: %w", err)
	}
	tx, err := c.CreateTransaction(ctx, senderAddress, toAddress, amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction: %w", err)
	}
	if err = c.SignTransaction(tx, senderPrivateKey); err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %w", err)
	}
	return c.BroadcastTransaction(ctx, tx)
}

func (c *TronClient) TransferToken(
	ctx context.Context,
	senderPrivateKey,
	tokenContractAddress,
	toAddress string,
	amount *big.Int,
) (common.Hash, error) {
	senderAddress, err := GetTronAddressFromPrivateKey(senderPrivateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while getting address from private key")
	}

	toAddressHex, err := Base58ToTronHex(toAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while address conversion")
	}
	parameterHex, err := ConstructTronTokenTxData(toAddressHex, amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while constructing txn data: %w", err)
	}

	triggerPayload := TriggerSmartContractRequest{
		OwnerAddress:     senderAddress,
		ContractAddress:  tokenContractAddress,
		FunctionSelector: "transfer(address,uint256)",
		Parameter:        parameterHex,
		FeeLimit:         10_000_000, // @TODO
		Visible:          true,
	}

	respBody, err := c.doRequest(ctx, "POST", "/wallet/triggersmartcontract", triggerPayload)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to trigger smart contract: %w", err)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
		return common.Hash{}, fmt.Errorf("trigger smart contract error: %s", errResp.Error)
	}

	var resp TriggerSmartContractResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return common.Hash{}, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	if !resp.Result.Result {
		return common.Hash{}, errors.New("error while smart contract trigger")
	}
	tx := resp.Transaction
	if tx.TxID == "" {
		return common.Hash{}, errors.New("transaction id is empty in response")
	}

	if err = c.SignTransaction(&tx, senderPrivateKey); err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return c.BroadcastTransaction(ctx, &tx)
}

func (c *TronClient) BatchGetAddressBalances(addresses []string) (map[string]decimal.Decimal, error) {
	data := make(map[string]decimal.Decimal)
	for _, address := range addresses {
		//time.Sleep(1 * time.Second)
		senderWalletBalance, err := c.GetTokenBalance(context.Background(), address, "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t")
		if err != nil {
			data[address] = decimal.NewFromFloat(0)
		}

		//log.Println("================================================")
		//log.Println("address", address)
		//log.Println("senderWalletBalance", senderWalletBalance)
		//log.Println("================================================")
		if senderWalletBalance == nil {
			data[address] = decimal.NewFromFloat(0)
		} else {
			senderWalletBalanceFloat := utils.ConvertBigIntToFloat(
				senderWalletBalance, 6,
			)

			data[address] = decimal.NewFromFloat(senderWalletBalanceFloat)
		}
	}

	return data, nil
}

// Transaction 表示Tron交易数据结构
type TxTransaction struct {
	TxID        string          `json:"txID"`
	BlockNumber int64           `json:"blockNumber"`
	FromAddress string          `json:"from_address"`
	ToAddress   string          `json:"to_address"`
	Amount      decimal.Decimal `json:"amount,omitempty"` // USDT金额
	Timestamp   time.Time       `json:"timestamp"`
	Contract    string          `json:"contract_address,omitempty"` // 合约地址(对于TRC20交易)
}

// GetIncomingTransactions 获取地址的转出交易
func (c *TronClient) GetOutgoingTransactions(address string, limit int, since time.Time) ([]TxTransaction, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s/transactions/trc20?only_from=true&limit=%d&min_timestamp=%d&contract_address=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		c.rpcURL, address, limit, since.Unix()*1000)

	//log.Println(url)
	//
	//req, err := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	return nil, fmt.Errorf("创建请求失败: %v", err)
	//}
	//
	//req.Header.Add("Accept", "application/json")
	////req.Header.Add("TRON-PRO-API-KEY", apiKey)
	//
	//client := &http.Client{Timeout: 30 * time.Second}
	//resp, err := client.Do(req)

	// 简单轮询
	keyIndex := atomic.AddUint32(&currentKeyIndex, 1) % uint32(len(global.TRONGRID_KEYS))
	currentKey := global.TRONGRID_KEYS[keyIndex]

	//log.Println(currentKey)

	time.Sleep(100 * time.Millisecond)

	req, err := http.NewRequest("GET", url, nil)
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TRON-PRO-API-KEY", currentKey) // 添加API密钥
	client := &http.Client{Timeout: 30 * time.Second}
	//resp, err := client.Get(url)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误: %s, 响应: %s", resp.Status, string(body))
	}

	var result struct {
		Data []struct {
			TransactionID string `json:"transaction_id"`
			BlockNumber   int64  `json:"block_number"`
			From          string `json:"from"`
			To            string `json:"to"`
			Value         string `json:"value"`
			Timestamp     int64  `json:"block_timestamp"`
			TokenInfo     struct {
				Address string `json:"address"`
			} `json:"token_info"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	var transactions []TxTransaction
	for _, tx := range result.Data {
		// 将字符串金额转换为float64 (USDT通常是6位小数)
		value, err := strconv.ParseFloat(tx.Value, 64)
		if err != nil {
			continue // 跳过金额解析失败的交易
		}
		amount := value / 1e6 // 假设是USDT(6位小数)

		transactions = append(transactions, TxTransaction{
			TxID:        tx.TransactionID,
			BlockNumber: tx.BlockNumber,
			FromAddress: tx.From,
			ToAddress:   tx.To,
			Amount:      decimal.NewFromFloat(amount),
			Timestamp:   time.Unix(tx.Timestamp/1000, 0),
			Contract:    tx.TokenInfo.Address,
		})
	}

	return transactions, nil
}

// GetIncomingTransactions 获取地址的转入交易
func (c *TronClient) GetIncomingTransactions(address string, limit int, since time.Time) ([]TxTransaction, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s/transactions/trc20?only_to=true&limit=%d&min_timestamp=%d&contract_address=TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		c.rpcURL, address, limit, since.Unix()*1000)

	//log.Println(url)

	//req, err := http.NewRequest("GET", url, nil)
	//if err != nil {
	//	return nil, fmt.Errorf("创建请求失败: %v", err)
	//}
	//
	//req.Header.Add("Accept", "application/json")
	////req.Header.Add("TRON-PRO-API-KEY", apiKey)
	//
	//client := &http.Client{Timeout: 30 * time.Second}
	//resp, err := client.Do(req)
	//

	keyIndex := atomic.AddUint32(&currentKeyIndex, 1) % uint32(len(global.TRONGRID_KEYS))
	currentKey := global.TRONGRID_KEYS[keyIndex]

	//log.Println(currentKey)

	time.Sleep(100 * time.Millisecond)

	req, err := http.NewRequest("GET", url, nil)
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TRON-PRO-API-KEY", currentKey) // 添加API密钥
	client := &http.Client{Timeout: 30 * time.Second}
	//resp, err := client.Get(url)
	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误: %s, 响应: %s", resp.Status, string(body))
	}

	var result struct {
		Data []struct {
			TransactionID string `json:"transaction_id"`
			BlockNumber   int64  `json:"block_number"`
			From          string `json:"from"`
			To            string `json:"to"`
			Value         string `json:"value"`
			Timestamp     int64  `json:"block_timestamp"`
			TokenInfo     struct {
				Address string `json:"address"`
			} `json:"token_info"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	var transactions []TxTransaction
	for _, tx := range result.Data {
		// 将字符串金额转换为float64 (USDT通常是6位小数)
		value, err := strconv.ParseFloat(tx.Value, 64)
		if err != nil {
			continue // 跳过金额解析失败的交易
		}
		amount := value / 1e6 // 假设是USDT(6位小数)

		transactions = append(transactions, TxTransaction{
			TxID:        tx.TransactionID,
			BlockNumber: tx.BlockNumber,
			FromAddress: tx.From,
			ToAddress:   tx.To,
			Amount:      decimal.NewFromFloat(amount),
			Timestamp:   time.Unix(tx.Timestamp/1000, 0),
			Contract:    tx.TokenInfo.Address,
		})
	}

	return transactions, nil
}

type Trc20Response struct {
	Data []struct {
		TokenInfo struct {
			Address string `json:"address"`
			Name    string `json:"name"`
			Symbol  string `json:"symbol"`
		} `json:"token_info"`
	} `json:"data"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

func (c *TronClient) GetUSDTTransferCount(address, token string) (int, error) {
	params := url.Values{}
	params.Add("sort", "-timestamp")
	params.Add("count", "true")
	params.Add("limit", "1") // 只需要总数，不需要具体数据
	params.Add("start", "0")
	params.Add("address", address)
	params.Add("token", token)
	params.Add("only_to", "true")

	resp, err := http.Get("https://api.tronscan.org/api/transaction?only_from=true" + "&" + params.Encode())
	if err != nil {
		return 0, fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	var result Trc20Response

	log.Println(&result)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, fmt.Errorf("解析JSON失败: %v", err)
	}

	return result.Meta.Total, nil
}

type BlockResponse struct {
	BlockHeader struct {
		RawData struct {
			Number int64 `json:"number"`
		} `json:"raw_data"`
	} `json:"block_header"`
}

func (c *TronClient) GetLatestBlock() int64 {
	resp, err := http.Get(c.rpcURL + "/wallet/getnowblock")
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
		return 0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应失败: %v\n", err)
		return 0
	}

	var block BlockResponse
	err = json.Unmarshal(body, &block)
	if err != nil {
		fmt.Printf("解析JSON失败: %v\n", err)
		return 0
	}

	fmt.Printf("最新区块高度: %d\n", block.BlockHeader.RawData.Number)
	return block.BlockHeader.RawData.Number
}
