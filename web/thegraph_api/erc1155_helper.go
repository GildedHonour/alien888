package thegraph_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// TODO
type ERC1155Token struct {
	ID string `json:"id"`
	// TokenID    int `json:"tokenId"`
	TokenID string `json:"tokenId"`
	//TODO rename to ContractAddress
	Contract string `json:"contract"`
	Owner    string `json:"owner"`
	// Balance string `json:"balance"`
}

type ERC1155TokensResult struct {
	Tokens []ERC1155Token `json:"tokens"`
}

func GetAllERC1155Tokens() []ERC1155Token {
	// const subgraphSlug = "blueberry_erc1155"
	// endpoint := fmt.Sprintf(theGraphBaseEndpointUrl, subgraphSlug)

	endpoint := "https://api.studio.thegraph.com/query/49439/blueberry_erc1155/v0.0.3"
	allTokensQuery := `
        query {
            tokens {
                id
                tokenId
                contract
                owner
            }
        }
    `

	allTokensResponse := executeGraphQuery(endpoint, allTokensQuery)
	allTokens := ERC1155TokensResult{}
	if err := json.Unmarshal(allTokensResponse, &allTokens); err != nil {
		log.Fatal("Error unmarshaling response:", err)
	}

	return allTokens.Tokens
}

func getAllERC1155TokensByOwner(walletAddress string) {
	const subgraphSlug = ""
	endpoint := fmt.Sprintf(theGraphBaseEndpointUrl, subgraphSlug)

	tokensByOwnerQuery := fmt.Sprintf(`
        query {
            tokens(where: { owner: "%s" }) {
                id
                tokenId
                contract
                owner
            }
        }
    `, walletAddress)
	tokensByOwnerResponse := executeGraphQuery(endpoint, tokensByOwnerQuery)
	tokensByOwner := ERC1155TokensResult{}
	if err := json.Unmarshal(tokensByOwnerResponse, &tokensByOwner); err != nil {
		log.Fatal("Error unmarshaling response:", err)
	}

	fmt.Println("Tokens held by wallet:")
	for _, token := range tokensByOwner.Tokens {
		fmt.Printf("Token ID: %s, Contract: %s\n", token.TokenID, token.Contract)
	}
}

func executeGraphQuery(endpoint, query string) []byte {
	requestBody := fmt.Sprintf(`{ "query": "%s" }`, strings.ReplaceAll(query, "\n", "\\n"))
	resp, err := http.Post(endpoint, jsonHeader, strings.NewReader(requestBody))
	checkErr(err)
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	checkErr(err)
	return responseBody
}
