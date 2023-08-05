package thegraph_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

func _subgraph_test1() {
	client := &http.Client{}

	query := `{
	  collections(where: { id: "COLLECTION_ID" }) {
	    tokens {
	      id
	      wallet {
	        id
	      }
	    }
	  }

	  wallets(where: { id: "WALLET_ID" }) {
	    tokens {
	      id
	      collection {
	        id
	      }
	    }
	  }
	}`

	requestBody, err := json.Marshal(GraphQLRequest{Query: query})
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", "https://api.thegraph.com/subgraphs/name/YOUR_GITHUB_USERNAME/YOUR_SUBGRAPH_NAME", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var graphQLResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphQLResp); err != nil {
		log.Fatal(err)
	}

	if len(graphQLResp.Errors) > 0 {
		log.Fatal(graphQLResp.Errors[0].Message)
	}

	fmt.Printf("%+v\n", graphQLResp.Data)
}
