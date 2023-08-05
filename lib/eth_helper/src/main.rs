use ethers::contract::abigen;
use ethers::providers::{Provider, Http};
use ethers::types::Address;
use std::str::FromStr;
use std::convert::TryFrom;

abigen!(
    MyContract,
    "./../../artifacts/contracts/Alien888Item.sol/Alien888Item.json",
    methods {
        balance_of_batch(
            address: &Address,
            ids: Vec<u64>
        ) -> Vec<u256>;
    }
);

#[tokio::main]
async fn main() {
    let wallet_address = Address::from_str("0x...").unwrap(); // The wallet address you want to check
    let erc1155_contract_address = Address::from_str("0x...").unwrap(); // The ERC1155 contract address
    let token_ids = vec![123, 456, 789]; // An array of token IDs you want to check

    // Connect to the Ethereum network
    let provider = Http::new("https://mainnet.infura.io/v3/YOUR_INFURA_PROJECT_ID").unwrap();

    // Create an instance of the ERC1155 contract
    let contract = MyContract::new(erc1155_contract_address, provider);

    // Check the balance of the wallet address for multiple token IDs in a single request
    let balances = contract
        .balance_of_batch(&[wallet_address], token_ids.clone())
        .await
        .unwrap();

    for i in 0..token_ids.len() {
        if balances[i] > 0.into() {
            println!(
                "{} is the owner of token ID {}",
                wallet_address, token_ids[i]
            );
        } else {
            println!(
                "{} does not own token ID {}",
                wallet_address, token_ids[i]
            );
        }
    }
}
