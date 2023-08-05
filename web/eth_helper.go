package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	//ERC721
	//TODO - read from file
	contractABI = `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]`

	//TODO take from config
	apiUrl = "https://arbitrum-goerli.publicnode.com"
)

func GetERC721NftsByWallet(walletAddress string) {
	client, err := ethclient.Dial(cfg.InfuraApiUrls.EthGoerli)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	address := common.HexToAddress(cfg.MainNftContractAddress)
	contract, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	balance, err := GetERC721NftBalance(client, address, walletAddress, contract)
	if err != nil {
		log.Fatalf("Failed to retrieve NFT balance: %v", err)
	}

	for i := big.NewInt(0); i.Cmp(balance) < 0; i.Add(i, big.NewInt(1)) {
		tokenID, err := getTokenIDByIndex(client, address, walletAddress, contract, i)
		if err != nil {
			log.Printf("Failed to retrieve token ID: %v", err)
			continue
		}

		fmt.Printf("Token ID: %s\n", tokenID.String())
		fmt.Printf("Owner: %s\n", walletAddress)
	}
}

func GetERC721NftBalance(client *ethclient.Client, contractAddress common.Address, ownerAddress string, contract abi.ABI) (*big.Int, error) {
	callData, err := contract.Pack("balanceOf", common.HexToAddress(ownerAddress))
	if err != nil {
		return nil, fmt.Errorf("failed to pack balanceOf function call: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	var balance *big.Int
	err = contract.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack balanceOf result: %v", err)
	}

	return balance, nil
}

func getTokenIDByIndex(client *ethclient.Client, contractAddress common.Address, ownerAddress string, contract abi.ABI, index *big.Int) (*big.Int, error) {
	callData, err := contract.Pack("tokenOfOwnerByIndex", common.HexToAddress(ownerAddress), index)
	if err != nil {
		return nil, fmt.Errorf("failed to pack tokenOfOwnerByIndex function call: %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %v", err)
	}

	var tokenID *big.Int
	err = contract.UnpackIntoInterface(&tokenID, "tokenOfOwnerByIndex", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack tokenOfOwnerByIndex result: %v", err)
	}

	return tokenID, nil
}

// return 2 integers as:
// amout_minted / max_amount_allowed
func GetCountersOfMint(tokenId int64) (int64, int64, error) {
	const abiFile = "./erc_abi/erc1155_partial2.json"

	client, err := ethclient.Dial(apiUrl)
	if err != nil {
		log.Fatal(err)
		return 0, 0, err
	}

	contractAddress := common.HexToAddress(cfg.SecondaryNftContractAddress)
	abiBytes, err := ioutil.ReadFile(abiFile)
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
		return 0, 0, err
	}

	// Parse the contract's ABI
	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
		return 0, 0, err
	}

	// Prepare the context for the call
	ctx := context.Background()

	// Call the first function (getTotalCopiesMinted)
	var totalCopiesMinted *big.Int
	callData, err := contractABI.Pack("getTotalCopiesMinted", big.NewInt(tokenId))
	if err != nil {
		log.Fatalf("Failed to pack function call 'getTotalCopiesMinted' data: %v", err)
		return 0, 0, err
	}

	responseData, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		log.Fatalf("Failed to call function getTotalCopiesMinted: %v", err)
		return 0, 0, err
	}

	err = contractABI.UnpackIntoInterface(&totalCopiesMinted, "getTotalCopiesMinted", responseData)
	if err != nil {
		log.Fatalf("Failed to unpack function response data: %v", err)
		return 0, 0, err
	}

	// Call the second function (getMaxCopiesPerToken)
	var maxCopiesPerToken *big.Int
	callData, err = contractABI.Pack("getMaxCopiesPerToken")
	if err != nil {
		log.Fatalf("Failed to pack function call 'getMaxCopiesPerToken' data: %v", err)
		return 0, 0, err
	}

	responseData, err = client.CallContract(ctx, ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		log.Fatalf("Failed to call function getMaxCopiesPerToken: %v", err)
		return 0, 0, err
	}

	err = contractABI.UnpackIntoInterface(&maxCopiesPerToken, "getMaxCopiesPerToken", responseData)
	if err != nil {
		log.Fatalf("Failed to unpack function response data: %v", err)
		return 0, 0, err
	}

	return totalCopiesMinted.Int64(), maxCopiesPerToken.Int64(), nil
}

// FIXME
func IsWalletInWhiteList(tokenId int64, addr string) (bool, error) {
	const abiFile = "./erc_abi/erc1155_partial2.json"

	client, err := ethclient.Dial(apiUrl)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	contractAddress := common.HexToAddress(cfg.SecondaryNftContractAddress)
	abiBytes, err := ioutil.ReadFile(abiFile)
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
		return false, err
	}

	// Parse the contract's ABI
	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
		return false, err
	}

	// Prepare the context for the call
	ctx := context.Background()

	// Call the third function (isInWhiteList)
	var isInWhitelistVal bool
	tokenIdAsBigInt := &big.Int{}
	callData, err := contractABI.Pack("isInWhiteList", tokenIdAsBigInt.SetInt64(tokenId), common.HexToAddress(addr))

	if err != nil {
		log.Fatalf("Failed to pack function call 'isInWhiteList' data: %v", err)
		return false, err
	}

	responseData, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		log.Fatalf("Failed to call function isInWhiteList: %v", err)
		return false, err
	}

	err = contractABI.UnpackIntoInterface(&isInWhitelistVal, "isInWhiteList", responseData)
	if err != nil {
		log.Fatalf("Failed to unpack function response data: %v", err)
		return false, err
	}

	return isInWhitelistVal, nil
}

// Function to load the contract's ABI from a JSON file
func loadABI(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	abi, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return abi, nil
}

func isValidEthereumWalletAddress(s1 string) bool {
	s1 = strings.Trim(s1, " ")
	a1 := s1 != ""
	a2 := strings.HasPrefix(s1, "0x")
	// a3 := len(s1) == 40 || len(s1) == 42
	a3 := len(s1) == 42

	if a1 && a2 && a3 {
		return true
	}

	return false
}
