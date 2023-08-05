CREATE TABLE nft_collections_history(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chain_id INTEGER,

    -- main NFT collection - heros
    main_nft_contract_address VARCHAR(70),
    main_nft_token_id INTEGER,

    -- secondary NFT collection - attributes
    secondary_nft_contract_address VARCHAR(70),
    secondary_nft_token_id INTEGER,
    secondary_nft_token_trait_id INTEGER,

    -- meta
    inserted_at INTEGER NOT NULL,
    updated_at INTEGER
);
