package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// TODO
func getNFTsForCollection(nftCollectionAddress string) []byte {
	alchemyApiUrl := fmt.Sprintf(
		"https://arb-mainnet.g.alchemy.com/nft/v2/%s/getNFTsForCollection?contractAddress=%s&withMetadata=true",
		cfg.AlchemyComApiKeyForMainnet,
		nftCollectionAddress,
	)

	req, _ := http.NewRequest("GET", alchemyApiUrl, nil)
	req.Header.Add("accept", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	//
	//
	return body
}
