package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v3"
)

const (
	customAvatarsPath = "./assets/img/custom_avatars"
)

func loadConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func verifySignedData(data []byte, signature string, address string) bool {
	// Decode the signature and data
	sigBytes := common.FromHex(signature)

	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	// dataHash := crypto.Keccak256Hash(data)
	dataHash := accounts.TextHash(data)

	// Recover the public key from the signature
	// pubKey, err := crypto.SigToPub(dataHash.Bytes(), sigBytes)
	pubKey, err := crypto.SigToPub(dataHash, sigBytes)
	if err != nil {
		log.Println("Failed to recover public key from signature:", err)
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()
	return strings.ToLower(recoveredAddress) == strings.ToLower(address)
}

// verifies the signed data with a nonce
func verifySignedDataWithNonce(data map[string]interface{}, signature string, address string, expectedNonce int64) bool {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println("Failed to marshal data to JSON:", err)
		return false
	}

	nonceBytes := int64IntoBytes(expectedNonce)
	dataHash := crypto.Keccak256Hash(dataBytes, nonceBytes)
	sigBytes := common.FromHex(signature)
	pubKey, err := crypto.SigToPub(dataHash.Bytes(), sigBytes)
	if err != nil {
		log.Println("Failed to recover public key from signature:", err)
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey).Hex()
	return strings.ToLower(recoveredAddress) == strings.ToLower(address) && expectedNonce == data["nonce"].(int64)
}

func int64IntoBytes(num int64) []byte {
	byteSlice := make([]byte, 8) // int64 occupies 8 bytes
	binary.LittleEndian.PutUint64(byteSlice, uint64(num))
	return byteSlice
}

func walletAddressIntoShortVersion(s1 string) string {
	const maxCharsLen = 4
	currentLen := len(s1)

	if currentLen <= maxCharsLen {
		return s1
	}

	if (currentLen > maxCharsLen) && (currentLen < (maxCharsLen * 2)) {
		return s1[:maxCharsLen] + "...."
	}

	return s1[:maxCharsLen] + "...." + s1[len(s1)-maxCharsLen:]
}

func getLocalFileUrlOfAvatar(fileName string) string {
	if fileName == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s", customAvatarsPath, fileName)
}

func writeBase64IntoFile(fileName string, data []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func writeBase64IntoFile2(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, 0644)
}

func extractBase64Data(dataURI string) string {
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return dataURI
}

func getBaseLayoutTemplatePath() string {
	return path.Join("templates", "layouts", "base.html")
}

func getPageTemplatePath(tpl string) string {
	return path.Join("templates", "pages", tpl)
}
