CREATE TABLE custom_avatars(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wallet_address VARCHAR(70),
    chain_id INTEGER,
    image_in_base64 TEXT,
    image_file_name TEXT,

    main_nft_token_id INTEGER,
    secondary_nft1_token_id INTEGER,
    secondary_nft2_token_id INTEGER,
    secondary_nft3_token_id INTEGER,

    -- meta
    inserted_at INTEGER NOT NULL default (cast(strftime('%s','now') as int)),
    updated_at INTEGER
);
