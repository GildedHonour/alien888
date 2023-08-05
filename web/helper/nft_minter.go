package helper

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	// "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func mintNft(
	contractAddress string,
	amount int64,
	privateKey string,
	receiverAddress string,
	gatewayApiUrl string,
) {
	// client, err := ethclient.Dial("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID")
	client, err := ethclient.Dial(gatewayApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Replace with the ABI of your ERC1155 contract
	contractABI, err := abi.JSON(strings.NewReader(`...`))
	if err != nil {
		log.Fatal(err)
	}

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// publicKey := privateKey.PublicKey
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	//     log.Fatal("error casting public key to ECDSA")
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fromAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKeyECDSA)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	// Replace the 'to' address, 'id', 'amount', and 'data' parameters with your desired values.
	tokenId := big.NewInt(1)
	// amount := big.NewInt(10)

	// Generate the data to call the mint function of the ERC1155 contract
	data, err := contractABI.Pack(
		"mint",
		common.HexToAddress(receiverAddress),
		tokenId,
		amount,
		[]byte("YOUR_DATA"),
	)

	if err != nil {
		log.Fatal(err)
	}

	// tx := ethereum.CallMsg{
	//     To:   &contractAddress,
	//     Data: data,
	// }

	// signedTx, err := auth.Signer(auth.From, &tx)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// err = client.SendTransaction(context.Background(), signedTx)
	// if err != nil {
	//     log.Fatal(err)
	// }

	tx := types.NewTransaction(
		auth.Nonce.Uint64(),
		common.HexToAddress(contractAddress),
		big.NewInt(0),
		auth.GasLimit,
		auth.GasPrice,
		data,
	)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyECDSA)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		// Transaction failed, handle the error
		// You might want to decode the logs to check for error messages or events
		for _, l := range receipt.Logs {
			// Assuming there's an event with the name "Error" that gets emitted when an exception is thrown
			if l.Topics[0] == crypto.Keccak256Hash([]byte("Error(uint256)")) {
				// Do something with the error event data if needed
				fmt.Println("An error occurred:", l.Data)
				break
			}
		}

		// Additional error handling if required
		log.Fatalf("Transaction failed with status: %d", receipt.Status)
	}

	fmt.Printf("Transaction hash: 0x%x\n", receipt.TxHash)
}

func mintNft2(
	contractAddress string,
	amount int64,
	privateKey string,
	receiverAddress string,
	gatewayApiUrl string,
) {
	client, err := ethclient.Dial(gatewayApiUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Replace with the ABI of your ERC1155 contract
	contractABI, err := abi.JSON(strings.NewReader(`...`))
	if err != nil {
		log.Fatal(err)
	}

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKeyECDSA)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(5000000) // in units (adjust accordingly)
	auth.GasPrice = gasPrice

	// Replace the 'to' address, 'token_id', and the number of copies to be minted.
	tokenId := big.NewInt(1)
	numCopies := 10 // Specify the number of copies to mint

	// Generate the data to call the mint function of the ERC1155 contract
	data, err := contractABI.Pack(
		"mint",
		common.HexToAddress(receiverAddress),
		tokenId,
		big.NewInt(int64(numCopies)),
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(auth.Nonce.Uint64(), common.HexToAddress(contractAddress), big.NewInt(0), auth.GasLimit, auth.GasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyECDSA)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	if receipt.Status == types.ReceiptStatusFailed {
		// Transaction failed, handle the error
		// You might want to decode the logs to check for error messages or events
		for _, l := range receipt.Logs {
			// Assuming there's an event with the name "Error" that gets emitted when an exception is thrown
			if l.Topics[0] == crypto.Keccak256Hash([]byte("Error(uint256)")) {
				// Do something with the error event data if needed
				fmt.Println("An error occurred:", l.Data)
				break
			}
		}

		// Additional error handling if required
		log.Fatalf("Transaction failed with status: %d", receipt.Status)
	}

	fmt.Printf("Transaction hash: 0x%x\n", receipt.TxHash)
}
