package pkg

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
	"time"
)

// GenerateTronOrderID 生成波场订单号（年月日时分 + 波场地址后4位）
func GenerateOrderID(tronAddress string, suffix int) (string, error) {
	// 1. 校验波场地址格式
	tronAddress = strings.TrimSpace(tronAddress)
	//if len(tronAddress) != 34 || !strings.HasPrefix(tronAddress, "T") {
	//	return "", fmt.Errorf("无效的波场地址（必须34位且以T开头）")
	//}

	// 2. 获取当前时间的 "年月日时分"（格式：200601021504）
	timestamp := time.Now().Format("200601021504")

	// 3. 截取波场地址后4位
	addressSuffix := tronAddress[len(tronAddress)-suffix:]

	// 4. 拼接时间 + 地址后4位
	orderID := timestamp + addressSuffix
	return orderID, nil
}

// HexToTronBase58 将 TRON Hex 地址（41...）转换为 Base58（T...）
func HexToTronBase58(tronHex string) (string, error) {

	hexStr := "41" + tronHex[24:64]
	// 检查长度和前缀
	if len(hexStr) != 42 || !strings.HasPrefix(hexStr, "41") {
		return "", fmt.Errorf("无效的 TRON 地址格式")
	}

	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}

	// TRON 使用 Base58Check 编码（第一个字节是版本号 0x41）
	return base58.CheckEncode(data[1:], data[0]), nil
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

func Base58ToTronHex(base58Addr string) (string, error) {
	decoded := base58.Decode(base58Addr)
	if len(decoded) != 25 {
		return "", fmt.Errorf("invalid address length")
	}
	checksum := decoded[len(decoded)-4:]
	payload := decoded[:len(decoded)-4]

	hash0 := sha256.Sum256(payload)
	hash1 := sha256.Sum256(hash0[:])

	if string(checksum) != string(hash1[:4]) {
		return "", fmt.Errorf("invalid checksum")
	}

	return fmt.Sprintf("%x", payload), nil
}

func ConstructTronTokenTxData(recipientHex string, amount *big.Int) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}]`))
	if err != nil {
		return "", fmt.Errorf("error while txn build")
	}
	data, err := parsedABI.Pack("transfer", common.HexToAddress(recipientHex), amount)
	if err != nil {
		return "", fmt.Errorf("error while txn build")
	}
	return hex.EncodeToString(data[4:]), nil
}

func GetTronAddressFromPrivateKey(privateKey string) (string, error) {
	senderPrivateKey := strings.TrimPrefix(privateKey, "0x")
	privKey, err := crypto.HexToECDSA(senderPrivateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}
	pubKey := privKey.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("failed to cast public key to ECDSA")
	}
	pubBytes := crypto.FromECDSAPub(publicKeyECDSA)
	hash := crypto.Keccak256(pubBytes[1:])
	addrBytes := hash[len(hash)-20:]
	tronAddrBytes := append([]byte{0x41}, addrBytes...)
	hexAddr := hex.EncodeToString(tronAddrBytes)
	base58Addr, err := TronHexToBase58(hexAddr)
	if err != nil {
		return "", fmt.Errorf("failed to convert hex address to base58: %w", err)
	}
	return base58Addr, nil
}
