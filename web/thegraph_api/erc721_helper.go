package thegraph_api

import (
	"encoding/json"
	"fmt"
)

type ERC721Token struct {
	ID      string `json:"id"`
	TokenID string `json:"tokenID"`
	//TODO rename to ContractAddress
	Contract string `json:"contract"`
	Owner    string `json:"owner"`
	// Metadata map[string]string `json:"metadata"`
}

type ERC721TokensResult struct {
	Tokens []ERC721Token `json:"tokens"`
}

func getAllERC721TokensByOwner(walletAddress string) {
	const subgraphSlug = ""
	endpoint := fmt.Sprintf(theGraphBaseEndpointUrl, subgraphSlug)
	tokensByOwnerQuery := fmt.Sprintf(`{
		tokensByOwner(owner: "%s") {
			id
			contract
			tokenID
			owner
		}
	}`, walletAddress)

	tokensByOwnerResponse := executeGraphQuery(endpoint, tokensByOwnerQuery)
	var result ERC721TokensResult
	err := json.Unmarshal(tokensByOwnerResponse, &result)
	if err != nil {
		fmt.Println("Failed to parse response JSON:", err)
		return
	}

	// Print the retrieved tokens
	for _, token := range result.Tokens {
		fmt.Printf("Token ID: %s\n", token.ID)
		fmt.Printf("Contract: %s\n", token.Contract)
		fmt.Printf("Token ID: %s\n", token.TokenID)
		fmt.Printf("Owner: %s\n", token.Owner)
		// fmt.Printf("Balance: %d\n", token.Balance)
		// fmt.Printf("Metadata: %+v\n", token.Metadata)
		fmt.Println("--------------------")
	}
}
