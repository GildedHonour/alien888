package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

//
//todo remove this file
//

func getERC1155NftsByWallet(contractAddressStr string, walletAddressStr string) {
	// Replace this with the contract address of your ERC1155 contract
	contractAddress := common.HexToAddress(contractAddressStr)

	// Replace this with the Ethereum addresses you want to check balances for
	// addresses := []common.Address{
	// 	common.HexToAddress("ADDRESS_1"),
	// 	common.HexToAddress("ADDRESS_2"),
	// }

	// var walletAddress common.Address = common.HexToAddress(walletAddressStr)
	addresses := []common.Address{
		common.HexToAddress(walletAddressStr),
	}

	// Initialize the Ethereum client
	client, err := ethclient.Dial(cfg.InfuraApiUrls.Goerli)
	checkErr(err)

	// Replace this with the ID(s) of the tokens you want to check balances for
	// Assuming you want to check balances for token IDs 1, 2, and 3
	tokenIDs := []*big.Int{
		big.NewInt(1),
		big.NewInt(2),
		big.NewInt(3),
	}

	// Call the `balanceOfBatch` function on the ERC1155 contract
	callData, err := encodeBalanceOfBatchCall(contractAddress, addresses, tokenIDs)
	checkErr(err)

	// Make the call to the Ethereum network
	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}, nil)

	checkErr(err)

	// Decode the result
	balances, err := decodeBalances(result)
	checkErr(err)

	// Display the balances
	for i, address := range addresses {
		fmt.Printf("Balances for address %s:\n", address.Hex())
		for j, tokenID := range tokenIDs {
			fmt.Printf("  Token ID %d: %d\n", tokenID.Int64(), balances[i][j].Int64())
		}
	}
}

// Helper function to encode the `balanceOfBatch` function call
func encodeBalanceOfBatchCall(contractAddress common.Address, addresses []common.Address, tokenIDs []*big.Int) ([]byte, error) {
	// Load the ERC1155 contract ABI
	// contractAbi, err := abi.JSON(strings.NewReader(erc1155ABI))
	// if err != nil {
	//     return nil, err
	// }

	abiBytes, err := ioutil.ReadFile("./erc1155_partial_abi.json")
	checkErr(err)

	abiString := string(abiBytes)
	contractAbi, err := abi.JSON(strings.NewReader(abiString))
	checkErr(err)

	// Generate the data for the `balanceOfBatch` function call
	data, err := contractAbi.Pack("balanceOfBatch", addresses, tokenIDs)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Helper function to decode the result of the `balanceOfBatch` function call
func decodeBalances(result []byte) ([][]*big.Int, error) {
	// Load the ERC1155 contract ABI
	// contractAbi, err := abi.JSON(strings.NewReader(erc1155ABI))
	// if err != nil {
	//     return nil, err
	// }

	abiBytes, err := ioutil.ReadFile("./erc1155_partial_abi.json")
	if err != nil {
		log.Fatal(err)
	}

	abiString := string(abiBytes)
	contractAbi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		log.Fatal(err)
	}

	// Prepare the outputs for the `balanceOfBatch` function
	balances := make([][]*big.Int, 0)
	outputs := make([]interface{}, 1)
	outputs[0] = &balances

	// Unpack the result into the outputs
	err = contractAbi.UnpackIntoInterface(outputs, "balanceOfBatch", result)
	if err != nil {
		return nil, err
	}

	return balances, nil
}
