package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"

	alchemy_com_api "alien888/alchemy_com_api"
	helper "alien888/helper"
	model "alien888/model"
)

const (
	mainNftCollectionKey      = "main_nft_collection"
	secondaryNftCollectionKey = "secondary_nft_collection"
	relationsKey              = "relations"
	defaultSessionName        = "default_session"
	currentWalletAddressKey   = "wallet_address"
	tokenIDParam              = "token_id"
	walletAddressDashes       = "----"
)

type UserWalletViewData struct {
	IsAuthenticated bool
	Address         string `json:"address"`
	ShortAddress    string `json:"short_address"`
}

type ViewData struct {
	UserWallet          UserWalletViewData `json:"user_wallet"`
	PageTitle           string
	ActiveTopNavBarMenu string
	UserAvatarInBase64  string
	UserAvatarFileName  string
	Data                map[string]any
}

type WardrobeUpdateViewData struct {
	MainNftID       int64  `json:"main_nft_id"`
	SecondaryNft1ID int64  `json:"secondary_nft1_id"`
	SecondaryNft2ID int64  `json:"secondary_nft2_id"`
	SecondaryNft3ID int64  `json:"secondary_nft3_id"`
	WalletAddress   string `json:"wallet_address"`
	Signature       string `json:"signature"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, defaultSessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vd := ViewData{}
	currentWalletAdress := walletAddressDashes
	isAuth := false
	if session.Values[currentWalletAddressKey] != nil {
		rawAddr := session.Values[currentWalletAddressKey].(string)
		if isValidEthereumWalletAddress(rawAddr) {
			isAuth = true
			currentWalletAdress = rawAddr

			avatarFileName := getCustomAvatarAsFileName(currentWalletAdress)
			if avatarFileName != "" {
				vd.UserAvatarFileName = getLocalFileUrlOfAvatar(avatarFileName)
			}
		}
	}

	vd.UserWallet = UserWalletViewData{
		isAuth,
		currentWalletAdress,
		walletAddressIntoShortVersion(currentWalletAdress),
	}

	html_templates := []string{
		getBaseLayoutTemplatePath(),
		getPageTemplatePath("index.html"),
	}

	tmpl, err := template.ParseFiles(html_templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", vd)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func walletHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, defaultSessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vd := ViewData{}
	currentWalletAdress := walletAddressDashes
	isAuth := false
	if session.Values[currentWalletAddressKey] != nil {
		rawAddr := session.Values[currentWalletAddressKey].(string)
		if isValidEthereumWalletAddress(rawAddr) {
			isAuth = true
			currentWalletAdress = rawAddr

			avatarFileName := getCustomAvatarAsFileName(currentWalletAdress)
			if avatarFileName != "" {
				vd.UserAvatarFileName = getLocalFileUrlOfAvatar(avatarFileName)
			}

			//TODO make them async
			//
			mainNftsRaw := alchemy_com_api.GetNFTs(
				cfg.AlchemyApiUrls.ArbGoerliV2,
				cfg.MainNftContractAddress,
				currentWalletAdress,
			)
			mainNfts := alchemy_com_api.ParseNftsInfo(mainNftsRaw, "ownedNfts")

			secondaryNftsRaw := alchemy_com_api.GetNFTs(
				cfg.AlchemyApiUrls.ArbGoerliV2,
				cfg.SecondaryNftContractAddress,
				currentWalletAdress,
			)
			secondaryNfts := alchemy_com_api.ParseNftsInfo(secondaryNftsRaw, "ownedNfts")

			for i := range secondaryNfts {
				secondaryNfts[i].BalanceOfCurrentWallet = secondaryNfts[i].Balance
			}

			vd.Data = map[string]any{
				"MainNftCollection":      mainNfts,
				"SecondaryNftCollection": secondaryNfts,
			}
		}
	}

	vd.UserWallet = UserWalletViewData{
		isAuth,
		currentWalletAdress,
		walletAddressIntoShortVersion(currentWalletAdress),
	}

	html_templates := []string{
		getBaseLayoutTemplatePath(),
		getPageTemplatePath("wallet.html"),
	}

	tmpl, err := template.ParseFiles(html_templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", vd)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	session, err := cookieStore.Get(r, defaultSessionName)
	helper.CheckErr(err)

	html_templates := []string{
		getBaseLayoutTemplatePath(),
		getPageTemplatePath("gallery.html"),
	}

	tmpl, err := template.ParseFiles(html_templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	erc1155NftsRaw := alchemy_com_api.GetNFTsForCollection(
		cfg.AlchemyApiUrls.ArbGoerliV2,
		cfg.SecondaryNftContractAddress,
	)
	erc1155Nfts := alchemy_com_api.ParseNftsInfo(erc1155NftsRaw, "nfts")

	vd := ViewData{}
	isAuth := false
	currentWalletAdress := walletAddressDashes
	if session.Values[currentWalletAddressKey] != nil {
		rawAddr := session.Values[currentWalletAddressKey].(string)
		if isValidEthereumWalletAddress(rawAddr) {
			isAuth = true
			currentWalletAdress = rawAddr

			avatarFileName := getCustomAvatarAsFileName(currentWalletAdress)
			if avatarFileName != "" {
				vd.UserAvatarFileName = getLocalFileUrlOfAvatar(avatarFileName)
			}
		}
	}

	vd.UserWallet = UserWalletViewData{
		isAuth,
		currentWalletAdress,
		walletAddressIntoShortVersion(currentWalletAdress),
	}

	vd.Data = map[string]any{
		"Nfts": erc1155Nfts,
	}

	err = tmpl.ExecuteTemplate(w, "base", vd)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func wardrobeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		session, err := cookieStore.Get(r, defaultSessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		htmlTplName := "empty_wardrobe.html"

		vd := ViewData{}
		isAuth := false
		currentWalletAdress := walletAddressDashes
		if session.Values[currentWalletAddressKey] != nil {
			rawAddr := session.Values[currentWalletAddressKey].(string)
			if isValidEthereumWalletAddress(rawAddr) {
				currentWalletAdress = rawAddr
				isAuth = true

				avatarFileName := getCustomAvatarAsFileName(currentWalletAdress)
				if avatarFileName != "" {
					vd.UserAvatarFileName = getLocalFileUrlOfAvatar(avatarFileName)
				}

				//TODO
				// var wg sync.WaitGroup
				// wg.Add(2)
				// resultChs := make([]chan string, 2)
				// resultChs[0] = make(chan string)
				// go getMainNftCollection(&wg, resultChs[0], currentWalletAdress)
				// resultChs[1] = make(chan string)
				// go getSecondaryNftCollection(&wg, resultChs[1])
				// wg.Wait()

				// fmt.Printf("test: %v", "#2")

				// templateData0 := []NftTokenInfo{}
				// templateData1 := []NftTokenInfo{}

				// for i, resultCh := range resultChs {
				// 	data := <-resultCh
				// 	nti := alchemy_com_api.ParseNftsInfo(data, "??")
				// 	switch i {
				// 	case 0:
				// 		templateData0 = nti
				// 	case 1:
				// 		templateData1 = nti
				// 	}

				// 	close(resultCh)
				// }
				//

				mainNfts := getMainNftCollectionSync(currentWalletAdress)
				if len(mainNfts) > 0 {
					htmlTplName = "wardrobe.html"

					secondaryNfts := getSecondaryNftCollectionSync()
					//get secondary NFTs owned by the current wallet
					ownedSecondaryNfts := getOwnedSecondaryNftCollectionSync(currentWalletAdress)

					rawTokenIDParam := r.URL.Query().Get(tokenIDParam)
					selectedMainNft := getSelectedOrDefaultMainNft(mainNfts, rawTokenIDParam)

					mainNftIDs := []int64{}
					for _, item := range mainNfts {
						mainNftIDs = append(mainNftIDs, item.TokenID)
					}

					//map[int64][]int64
					relBwNfts := getRelationsBetweenNfts(mainNftIDs)
					secondaryNft1, secondaryNft2, secondaryNft3 := getSelectedSecondaryNfts(selectedMainNft.TokenID, secondaryNfts, relBwNfts)

					//mark the secondary NFTs that are owned by the current wallet
					//set the balance of the current wallet
					for idx2, item2 := range secondaryNfts {
						idx := slices.IndexFunc(ownedSecondaryNfts, func(it model.NftTokenInfo) bool { return it.TokenID == item2.TokenID })
						if idx != -1 {
							secondaryNfts[idx2].IsOWnedByCurrentWallet = true
							secondaryNfts[idx2].BalanceOfCurrentWallet = ownedSecondaryNfts[idx].Balance
						}
					}


					var selMainNftTransThmb string
					if len(mainNfts) > 0 {
						selMainNftTransThmb = fmt.Sprintf("./assets/img/transparent_heros/hero__%d.png", selectedMainNft.TokenID)
					}

					vd.Data = map[string]any{
						"MainNftCollection":            mainNfts,
						"MainNftCollectionHasElements": len(mainNfts) > 0,
						"SelectedMainNft":              selectedMainNft,

						"SecondaryNftCollection":            secondaryNfts,
						"SecondaryNftCollectionHasElements": len(mainNfts) > 0,

						"ParentRelationBackroundKey": model.ParentNftRelationBackround,
						"ParentRelationItem1Key":     model.ParentNftRelationItem1,
						"ParentRelationItem2Key":     model.ParentNftRelationItem2,

						"SelectedSecondaryNftForBackround": secondaryNft1,
						"SelectedSecondaryNftForItem1":     secondaryNft2,
						"SelectedSecondaryNftForItem2":     secondaryNft3,

      					"SelectedMainNftTransparentThumbnail": selMainNftTransThmb,
					}
				}
			}
		}

		vd.UserWallet = UserWalletViewData{
			isAuth,
			currentWalletAdress,
			walletAddressIntoShortVersion(currentWalletAdress),
		}

		html_templates := []string{
			getBaseLayoutTemplatePath(),
			getPageTemplatePath(htmlTplName),
		}

		tmpl, err := template.ParseFiles(html_templates...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, "base", vd)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}

	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var rawParams map[string]interface{}
		err = json.Unmarshal([]byte(body), &rawParams)
		if err != nil {
			http.Error(w, "JSON is in invalid format", http.StatusInternalServerError)
			return
		}

		w.Header().Set(helper.ContentTypeHeader, helper.JsonHeader)

		//TODO
		data, ok1 := rawParams["data"].(map[string]interface{})
		walletAddress, ok2 := rawParams["wallet_address"].(string)
		signature, ok3 := rawParams["signature"].(string)
		if !ok1 || !ok2 || !ok3 {
			jsonErrResp, err := json.Marshal(map[string]string{"message": "one of the mandatory params is absent"})
			helper.CheckErr(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonErrResp)
			return
		}

		message, err := json.Marshal(data)
		helper.CheckErr(err)
		if verifySignedData(message, signature, walletAddress) {
			mainNftTokenID, err := strconv.ParseInt((data["main_nft_id"].(string)), 0, 64)
			helper.CheckErr(err)

			secondaryNftTokenIDs := []int64{}
			for _, x := range data["secondary_nft_ids"].([]any) {
				val, err := strconv.ParseInt((x.(string)), 0, 64)
				if err == nil {
					secondaryNftTokenIDs = append(secondaryNftTokenIDs, val)
				}
			}

			//TODO
			// (1) request data from blockchain:
			// 		verify that these NFT-attributes belong to this wallet as well as the main NFT

			//TODO
			// (2)
			//		verify their traits by retrieving data from the blockchain

			// (3) then
			//		update data in the Db
			updateRelationsBetweenNfts(mainNftTokenID, secondaryNftTokenIDs)

			//returns ok
			jsonResp, err := json.Marshal(map[string]string{"message": "ok"})
			if err == nil {
				w.WriteHeader(http.StatusOK)
				w.Write(jsonResp)
			} else {
				//returns error
				jsonErrResp, err := json.Marshal(map[string]string{"message": "generic error"})
				log.Print(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonErrResp)
			}
		} else {
			jsonErrResp, err := json.Marshal(map[string]string{"message": "signature is invalid"})
			helper.CheckErr(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonErrResp)
			return
		}
	}
}

func currentWalletHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		contentType := r.Header.Get(helper.ContentTypeHeader)
		if strings.Contains(contentType, helper.JsonHeader) {
			decoder := json.NewDecoder(r.Body)
			var params map[string]string
			err := decoder.Decode(&params)
			if err != nil {
				http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
				return
			}

			session, err := cookieStore.Get(r, defaultSessionName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			newWallettAddr := params[currentWalletAddressKey]
			oldWalletAddr := session.Values[currentWalletAddressKey]
			session.Values[currentWalletAddressKey] = newWallettAddr
			err = session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			isPageReloadNeeded := false
			if (oldWalletAddr != nil && oldWalletAddr.(string) != newWallettAddr) || (oldWalletAddr == nil) {
				isPageReloadNeeded = true
			}

			msg := fmt.Sprintf("current wallet address: %s", newWallettAddr)
			jsonResp, err := json.Marshal(map[string]any{
				"message":               msg,
				"is_page_reload_needed": isPageReloadNeeded,
			})

			if err == nil {
				w.Header().Set(helper.ContentTypeHeader, helper.JsonHeader)
				w.WriteHeader(http.StatusOK)
				w.Write(jsonResp)
			}
		}
	}
}

func whiteListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		session, err := cookieStore.Get(r, defaultSessionName)
		helper.CheckErr(err)

		tokenIDRaw := r.URL.Query().Get(tokenIDParam)
		if tokenIDRaw != "" {
			tokenID, err := strconv.ParseInt(tokenIDRaw, 0, 64)
			helper.CheckErr(err)

			//TODO make async
			res := alchemy_com_api.GetNFTMetadata(
				cfg.AlchemyApiUrls.ArbGoerliV2,
				cfg.SecondaryNftContractAddress,
				tokenID,
			)

			secondaryNft := alchemy_com_api.ParseSingleNftInfo(res)
			html_templates := []string{
				getBaseLayoutTemplatePath(),
				getPageTemplatePath("whitelist.html"),
			}

			tmpl, err := template.ParseFiles(html_templates...)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			vd := ViewData{}
			currentWalletAdress := walletAddressDashes
			isAuth := false
			isInWhiteListVal := false
			if session.Values[currentWalletAddressKey] != nil {
				rawAddr := session.Values[currentWalletAddressKey].(string)
				if isValidEthereumWalletAddress(rawAddr) {
					currentWalletAdress = rawAddr
					isAuth = true

					//TODO make async
					isInWhiteListVal, err = IsWalletInWhiteList(tokenID, currentWalletAdress)
					helper.CheckErr(err)

					avatarFileName := getCustomAvatarAsFileName(currentWalletAdress)
					if avatarFileName != "" {
						vd.UserAvatarFileName = getLocalFileUrlOfAvatar(avatarFileName)
					}
				}
			}

			vd.UserWallet = UserWalletViewData{
				isAuth,
				currentWalletAdress,
				walletAddressIntoShortVersion(currentWalletAdress),
			}

			//TODO make async
			amountMinted, amountAllowed, err := GetCountersOfMint(tokenID)
			helper.CheckErr(err)

			vd.Data = map[string]any{
				"SecondaryNft":               secondaryNft,
				"ContractAddress":            cfg.SecondaryNftContractAddress,
				"CurrenTokenID":              tokenID,
				"AmountMinted":               amountMinted,
				"AmountAllowed":              amountAllowed,
				"IsCurrentWalletInWhiteList": isInWhiteListVal,
				"ChainID":                    cfg.ChainID,
			}

			err = tmpl.ExecuteTemplate(w, "base", vd)
			if err != nil {
				log.Print(err.Error())
				http.Error(w, "Internal Server Error", 500)
			}
		} else {
			errMsg := "The param 'token_id' must not be empty"
			log.Print(errMsg)
			http.Error(w, errMsg, 403)
		}
	case http.MethodPost:
		break
	}
}

func tokensERC1155MetadataHandler(w http.ResponseWriter, r *http.Request) {
	const baseFilesPath = "./tokens_metadata/erc1155"

	vars := mux.Vars(r)
	tokenID := vars["id"]
	fileName := fmt.Sprintf("%s/%s.json", baseFilesPath, tokenID)
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set(helper.ContentTypeHeader, helper.JsonHeader)
	w.Write(content)
}

func tokensERC721MetadataHandler(w http.ResponseWriter, r *http.Request) {
	const baseFilesPath = "./tokens_metadata/erc721"

	vars := mux.Vars(r)
	tokenID := vars["id"]
	fileName := fmt.Sprintf("%s/%s.json", baseFilesPath, tokenID)
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set(helper.ContentTypeHeader, helper.JsonHeader)
	w.Write(content)
}

func getMainNftCollection(wg *sync.WaitGroup, results chan<- []byte, ownerWalletAddress string) {
	defer wg.Done()
	res := alchemy_com_api.GetNFTs(
		cfg.AlchemyApiUrls.ArbGoerliV2,
		cfg.MainNftContractAddress,
		ownerWalletAddress,
	)

	results <- res
}

func getSecondaryNftCollection(wg *sync.WaitGroup, results chan<- []byte) {
	defer wg.Done()

	res := alchemy_com_api.GetNFTsForCollection(
		cfg.AlchemyApiUrls.ArbGoerliV2,
		cfg.SecondaryNftContractAddress,
	)

	results <- res
}

func getMainNftCollectionSync(ownerWalletAddress string) []model.NftTokenInfo {
	resBytes := alchemy_com_api.GetNFTs(
		cfg.AlchemyApiUrls.ArbGoerliV2,
		cfg.MainNftContractAddress,
		ownerWalletAddress,
	)

	res := alchemy_com_api.ParseNftsInfo(resBytes, "ownedNfts")
	return res
}

func getSecondaryNftCollectionSync() []model.NftTokenInfo {
	resBytes := alchemy_com_api.GetNFTsForCollection(
		cfg.AlchemyApiUrls.ArbGoerliV2,
		cfg.SecondaryNftContractAddress,
	)

	res := alchemy_com_api.ParseNftsInfo(resBytes, "nfts")
	return res
}

func getSelectedSecondaryNfts(
	mainNftTokenID int64,
	secondaryNfts []model.NftTokenInfo,
	relationBetweenNfts map[int64][]int64,
) (*model.NftTokenInfo, *model.NftTokenInfo, *model.NftTokenInfo) {
	var (
		secondaryNft1 *model.NftTokenInfo
		secondaryNft2 *model.NftTokenInfo
		secondaryNft3 *model.NftTokenInfo
	)

	secondaryNftIDs := relationBetweenNfts[mainNftTokenID]
	for idx2, item2 := range secondaryNfts {
		if slices.Contains(secondaryNftIDs, item2.TokenID) {
			secondaryNfts[idx2].ParentTokenID = mainNftTokenID

			switch secondaryNfts[idx2].ParentRelationTrait {

			case model.ParentNftRelationBackround:
				secondaryNft1 = &secondaryNfts[idx2]
				break

			case model.ParentNftRelationItem1:
				secondaryNft2 = &secondaryNfts[idx2]
				break

			case model.ParentNftRelationItem2:
				secondaryNft3 = &secondaryNfts[idx2]
				break

			default:
				break
			}
		}
	}

	return secondaryNft1, secondaryNft2, secondaryNft3
}

func getSelectedOrDefaultMainNft(mainNfts []model.NftTokenInfo, rawTokenIDParam string) model.NftTokenInfo {
	selectedMainNft := model.NftTokenInfo{}
	if len(mainNfts) > 0 {
		if rawTokenIDParam != "" {
			tokenID, err := strconv.ParseInt(rawTokenIDParam, 0, 64)
			helper.CheckErr(err)
			for _, item := range mainNfts {
				if item.TokenID == tokenID {
					selectedMainNft = item
				}
			}
		} else {
			selectedMainNft = mainNfts[0]
		}
	}

	return selectedMainNft
}

func getOwnedSecondaryNftCollectionSync(currentWalletAdress string) []model.NftTokenInfo {
	ownedSecondaryNftBytes := alchemy_com_api.GetNFTs(
		cfg.AlchemyApiUrls.ArbGoerliV2,
		cfg.SecondaryNftContractAddress,
		currentWalletAdress,
	)

	ownedSecondaryNfts := alchemy_com_api.ParseNftsInfo(ownedSecondaryNftBytes, "ownedNfts")
	return ownedSecondaryNfts
}

func currentCustomAvatarHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		contentType := r.Header.Get(helper.ContentTypeHeader)
		if strings.Contains(contentType, helper.JsonHeader) {
			decoder := json.NewDecoder(r.Body)
			var params map[string]any
			err := decoder.Decode(&params)
			if err != nil {
				http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
				return
			}

			insertQuery := `
			        INSERT INTO custom_avatars(
			            chain_id,
			            wallet_address,
			            image_in_base64,
			            image_file_name,

			            main_nft_token_id,
			            secondary_nft1_token_id,
			            secondary_nft2_token_id,
			            secondary_nft3_token_id
			        )
			        VALUES (
			            ?, ?, ?, ?,
			            ?, ?, ?, ?
			        )
    		`

			//save base64 as jpeg file
			base64Data := params["image_in_base64"].(string)
			decodedData, err := base64.StdEncoding.DecodeString(extractBase64Data(base64Data))
			if err != nil {
				fmt.Println("Error decoding base64:", err)
				return
			}

			newUUID := uuid.New()
			imgFileName := fmt.Sprintf("%s.jpeg", newUUID.String())
			exePath, err := os.Executable()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			currDir := filepath.Dir(exePath)
			fullImgFileName := filepath.Join(currDir, customAvatarsPath, imgFileName)
			err = writeBase64IntoFile(fullImgFileName, decodedData)
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}

			_, err = dbConn.Exec(
				insertQuery,

				cfg.ChainID,
				params["wallet_address"],
				params["image_in_base64"],
				imgFileName,

				params["main_nft_token_id"],
				params["secondary_nft1_token_id"],
				params["secondary_nft2_token_id"],
				params["secondary_nft3_token_id"],
			)

			if err != nil {
				panic(err)
			}

			jsonRespData := map[string]any{
				"message": "avatar updated",
			}

			w.Header().Set(helper.ContentTypeHeader, helper.JsonHeader)
			w.WriteHeader(http.StatusOK)

			jsonResp, err := json.Marshal(jsonRespData)
			helper.CheckErr(err)
			w.Write(jsonResp)
		}
	}
}
