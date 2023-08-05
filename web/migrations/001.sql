CREATE TABLE nft_collections(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    chain_id INTEGER,

    -- main NFT collection - heros
    main_nft_contract_address VARCHAR(70),
    main_nft_token_id INTEGER,

    -- secondary NFT collection - attributes
    secondary_nft_contract_address VARCHAR(70),
    secondary_nft_token_id INTEGER,
    secondary_nft_token_trait_id INTEGER,
    -- secondary_nft_token_trait_slug VARCHAR(20),

    -- meta
    inserted_at INTEGER NOT NULL default (cast(strftime('%s','now') as int)),
    updated_at INTEGER
);
