package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func NewERC1155(address common.Address, client *ethclient.Client) (*ERC1155, error) {
	abiBytes, err := ioutil.ReadFile("./erc1155_partial_abi.json")
	if err != nil {
		log.Fatal(err)
	}

	abiString := string(abiBytes)
	contractABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		log.Fatal(err)
	}

	contractInstance, err := NewERC1155(address, client)
	if err != nil {
		log.Fatal(err)
	}

	return &ERC1155{
		Address:  address,
		Client:   client,
		ABI:      contractABI,
		Instance: contractInstance,
	}, nil
}

type ERC1155 struct {
	Address  common.Address
	Client   *ethclient.Client
	ABI      abi.ABI
	Instance *ERC1155
}

func (erc1155 *ERC1155) balanceOf(owner common.Address, tokenID *big.Int) (*big.Int, error) {
	callData, err := erc1155.ABI.Pack("balanceOf", owner, tokenID)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &erc1155.Address,
		Data: callData,
	}

	result, err := erc1155.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var balance *big.Int
	err = erc1155.ABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (erc1155 *ERC1155) balanceOfBatch(owners []common.Address, tokenIDs []*big.Int) ([]*big.Int, error) {
	callData, err := erc1155.ABI.Pack("balanceOfBatch", owners, tokenIDs)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &erc1155.Address,
		Data: callData,
	}

	result, err := erc1155.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var balances []*big.Int
	err = erc1155.ABI.UnpackIntoInterface(&balances, "balanceOfBatch", result)
	if err != nil {
		return nil, err
	}

	return balances, nil
}
