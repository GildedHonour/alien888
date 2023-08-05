package main

type Config struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	SessionKey string `yaml:"session_key"`

	ChainID                        int    `yaml:"chain_id"`
	MainNftContractAddress         string `yaml:"main_nft_contract_address"`
	SecondaryNftContractAddress    string `yaml:"secondary_nft_contract_address"`
	MainNftToSecondaryNftMaxAmount int    `yaml:"main_nft_to_secondary_nft_max_amount"`

	InfuraApiUrls struct {
		EthMainnet string `yaml:"eth_mainnet"`
		EthGoerli  string `yaml:"eth_goerli"`
		EthSepolia string `yaml:"eth_sepolia"`
	} `infura_api_urls`

	AlchemyApiUrls struct {
		EthMainnet string `yaml:"eth_mainnet"`
		EthGoerli  string `yaml:"eth_goerli"`
		EthSepolia string `yaml:"eth_sepolia"`

		ArbMainnetV2 string `yaml:"arb_mainnet_v2"`
		ArbGoerliV2  string `yaml:"arb_goerli_v2"`
	} `alchemy_api_urls`

	ArbiscanApi struct {
		ArbMainnetApiKey string `yaml:"arb_mainnet_api_key"`
		ArbGoerliApiKey  string `yaml:"arb_goerli_api_key"`
	} `arbiscan_api`

	BlueberryGBCContractAddress      string `yaml:"blueberry_gbc_contract_address"`
	BlueberryLabItemsContractAddress string `yaml:"blueberry_lab_items_contract_address"`

	TestWalletAddress1 string `yaml:"test_wallet_address_1"`
	TestWalletAddress2 string `yaml:"test_wallet_address_2"`
}
