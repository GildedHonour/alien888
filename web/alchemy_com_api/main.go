package alchemy_com_api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	helper "alien888/helper"
	model "alien888/model"
)

// returns all NFTs currently owned by a given address
func GetNFTs(baseApiUrl string, contractAddress string, ownerAddress string) []byte {
	url := fmt.Sprintf(
		"%s/getNFTs?&contractAddresses[]=%s&owner=%s&withMetadata=true&pageSize=100",
		baseApiUrl,
		contractAddress,
		ownerAddress,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", helper.JsonHeader)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

// it returns all NFTs for a given NFT contract
func GetNFTsForCollection(baseApiUrl string, contractAddress string) []byte {
	url := fmt.Sprintf(
		"%s/getNFTsForCollection?contractAddress=%s&withMetadata=true",
		baseApiUrl,
		contractAddress,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", helper.JsonHeader)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

// returns the metadata associated with a given NFT.
func GetNFTMetadata(baseApiUrl string, contractAddress string, tokenID int64) []byte {
	url := fmt.Sprintf(
		"%s/getNFTMetadata?contractAddress=%s&tokenId=%d",
		baseApiUrl,
		contractAddress,
		tokenID,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", helper.JsonHeader)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func ParseNftsInfo(nftByteData []byte, rootKey string) []model.NftTokenInfo {
	var nftData map[string]interface{}
	err := json.Unmarshal(nftByteData, &nftData)
	checkErr(err)

	var nftInfos []model.NftTokenInfo
	for _, item := range nftData[rootKey].([]interface{}) {
		nftInfoItem := parseSingleNftImplV2(item)
		nftInfos = append(nftInfos, nftInfoItem)
	}

	return nftInfos
}

func ParseSingleNftInfo(nftByteData []byte) model.NftTokenInfo {
	var nftData map[string]interface{}
	err := json.Unmarshal(nftByteData, &nftData)
	checkErr(err)
	nftInfo := parseSingleNftImplV2(nftData)
	return nftInfo
}

// API v2
func parseSingleNftImplV2(item interface{}) model.NftTokenInfo {
	var nftInfoItem model.NftTokenInfo

	var itemDict = item.(map[string]interface{})
	var tokenID = itemDict["id"].(map[string]interface{})["tokenId"]
	var thumbnail = itemDict["media"].([]interface{})[0].(map[string]interface{})["thumbnail"]

	tid, err := strconv.ParseInt(tokenID.(string), 0, 64)
	checkErr(err)
	nftInfoItem.TokenID = tid

	nftInfoItem.Title = itemDict["title"].(string)
	nftInfoItem.Description = itemDict["description"].(string)

	if itemDict["balance"] != nil {
		bal, err := strconv.ParseInt(itemDict["balance"].(string), 0, 64)
		checkErr(err)
		nftInfoItem.Balance = bal
	} else {
		nftInfoItem.Balance = 1
	}

	if thumbnail != nil {
		nftInfoItem.Thumbnail = thumbnail.(string)
	} else {
		nftInfoItem.Thumbnail = ""
	}

	if itemDict["metadata"] != nil {
		md := itemDict["metadata"].(map[string]any)
		if md["parent_relation_trait"] != nil {
			s := md["parent_relation_trait"].(string)
			parRelTr := model.ParentRelationTraitType(s)
			nftInfoItem.ParentRelationTrait = parRelTr
		}
	}

	return nftInfoItem
}
