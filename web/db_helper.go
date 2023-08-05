package main

import (
	"fmt"
	"log"
	"strings"

	helper "alien888/helper"
)

//
// ID of main NFT token => IDs of secondary ones
// {
//   1 => [3, 8, 33],
//   49 => [5, 167, 12, 99, 633],
// }
//

func getRelationsBetweenNfts(mainTokenIDs []int64) map[int64][]int64 {
	placeHolders := generatePlaceholders(len(mainTokenIDs))
	q := fmt.Sprintf(`
        SELECT
        	main_nft_token_id,
            secondary_nft_token_id,
            secondary_nft_token_trait_id
        FROM nft_collections
            WHERE
            	main_nft_contract_address = ?
            	AND secondary_nft_contract_address = ?
            	AND chain_id = ?
            	AND main_nft_token_id IN (%s)
    `, placeHolders,
	)

	stmtSel, err := dbConn.Prepare(q)
	helper.CheckErr(err)
	defer stmtSel.Close()

	whereCond := make([]any, len(mainTokenIDs)+3)
	whereCond[0] = cfg.MainNftContractAddress
	whereCond[1] = cfg.SecondaryNftContractAddress
	whereCond[2] = cfg.ChainID
	for i, v := range mainTokenIDs {
		whereCond[3+i] = v
	}

	rows, err := stmtSel.Query(whereCond...)
	helper.CheckErr(err)
	defer rows.Close()

	res := make(map[int64][]int64)
	for rows.Next() {
		var a1, a2, a3 int64
		err := rows.Scan(&a1, &a2, &a3)

		//a1 -> main_nft_token_id
		//a2 -> secondary_nft_token_id
		//a3 -> trait_id

		helper.CheckErr(err)
		if _, ok := res[a1]; !ok {
			res[a1] = []int64{}
		}
		res[a1] = append(res[a1], a2)
	}

	err = rows.Err()
	helper.CheckErr(err)

	return res
}

func updateRelationsBetweenNfts(
	mainNftTokenID int64,
	secondaryNftTokenIDs []int64,
) {
	whereCond := []any{
		mainNftTokenID,
		cfg.MainNftContractAddress,
		cfg.SecondaryNftContractAddress,
	}

	// begin a transaction
	tx, err := dbConn.Begin()
	helper.CheckErr(err)

	//1
	//copy to 'history' table
	stmtInsSel, err := dbConn.Prepare(`
        INSERT INTO nft_collections_history(
            chain_id,
            main_nft_contract_address,
            main_nft_token_id,

            secondary_nft_contract_address,
            secondary_nft_token_id,
            secondary_nft_token_trait_id,

            inserted_at
        )
        SELECT
            chain_id,
            main_nft_contract_address,
            main_nft_token_id,

            secondary_nft_contract_address,
            secondary_nft_token_id,
            secondary_nft_token_trait_id,

            DATETIME('now')
        FROM nft_collections
            WHERE main_nft_token_id = ?
            AND main_nft_contract_address = ?
            AND secondary_nft_contract_address = ?
    `)

	helper.CheckErr(err)
	defer stmtInsSel.Close()

	_, err = stmtInsSel.Exec(whereCond...)
	helper.CheckErr(err)

	//2
	//detele from the main table
	stmtDel, err := dbConn.Prepare(`
        DELETE FROM nft_collections
            WHERE main_nft_token_id = ?
            AND main_nft_contract_address = ?
            AND secondary_nft_contract_address = ?
    `)

	helper.CheckErr(err)
	defer stmtDel.Close()

	_, err = stmtDel.Exec(whereCond...)
	helper.CheckErr(err)

	//3
	//insert into the main table
	insValues := []string{}
	for _, tokIDItem := range secondaryNftTokenIDs {
		sqlStmt := fmt.Sprintf(
			"(%d, '%s', %d,    '%s', %d, %d)",
			cfg.ChainID, cfg.MainNftContractAddress, mainNftTokenID,
			cfg.SecondaryNftContractAddress, tokIDItem, 0,
		)

		insValues = append(insValues, sqlStmt)
	}

	stmtIns := fmt.Sprintf(`
        INSERT INTO nft_collections(
            chain_id,
            main_nft_contract_address,
            main_nft_token_id,

            secondary_nft_contract_address,
            secondary_nft_token_id,
            secondary_nft_token_trait_id
        ) VALUES %s`,
		strings.Join(insValues, ", "))

	_, err = dbConn.Exec(stmtIns)
	helper.CheckErr(err)

	// commit a transaction
	err = tx.Commit()
	helper.CheckErr(err)
}

func getCustomAvatarAsBase64(walletAddress string) string {
	stmtSel, err := dbConn.Prepare(`
        SELECT
            image_in_base64
        FROM custom_avatars
            WHERE
                wallet_address = ?
                AND chain_id = ?
        LIMIT 1
    `)

	helper.CheckErr(err)
	defer stmtSel.Close()

	whereCond := []any{
		walletAddress,
		cfg.ChainID,
	}

	rows, err := stmtSel.Query(whereCond...)
	helper.CheckErr(err)
	defer rows.Close()

	var res string
	rowCount := 0
	for rows.Next() {
		err := rows.Scan(&res)
		helper.CheckErr(err)
		rowCount++
	}

	err = rows.Err()
	helper.CheckErr(err)

	if rowCount == 0 {
		return ""
	} else if rowCount == 1 {
		return res
	} else {
		log.Fatalf("DB > custom_avatars > select > more than one row returned; %v\n", stmtSel)
		return res
	}
}

func getCustomAvatarAsFileName(walletAddress string) string {
	stmtSel, err := dbConn.Prepare(`
        SELECT
            image_file_name
        FROM custom_avatars
            WHERE
                wallet_address = ?
                AND chain_id = ?
        ORDER BY id DESC
        LIMIT 1
    `)

	helper.CheckErr(err)
	defer stmtSel.Close()

	whereCond := []any{
		walletAddress,
		cfg.ChainID,
	}

	rows, err := stmtSel.Query(whereCond...)
	helper.CheckErr(err)
	defer rows.Close()

	var res string
	rowCount := 0
	for rows.Next() {
		err := rows.Scan(&res)
		helper.CheckErr(err)
		rowCount++
	}

	err = rows.Err()
	helper.CheckErr(err)

	if rowCount == 0 {
		return ""
	} else if rowCount == 1 {
		return res
	} else {
		log.Fatalf("DB > custom_avatars > select > more than one row returned; %v\n", stmtSel)
		return res
	}
}

func generatePlaceholders(count int) string {
	if count == 0 {
		return ""
	}

	return strings.Repeat("?, ", count-1) + "?"
}
