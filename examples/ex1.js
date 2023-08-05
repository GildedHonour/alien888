import * as dotenv from 'dotenv';
import minimist from 'minimist';
import * as ethers from 'ethers';
import * as fs from 'fs';

dotenv.config();

const ALIEN888_ITEMS_ABI_JSON_FILE_NAME = "../artifacts/contracts/Alien888Item.sol/Alien888Item.json";
const ALIEN888_ITEMS_ABI = JSON.parse(fs.readFileSync(ALIEN888_ITEMS_ABI_JSON_FILE_NAME)).abi;
const ALIEN888_ITEMS_CONTRACT_ADDRESS = process.env.ALIEN888_ITEMS_CONTRACT_ADDRESS;

async function main() {
  const provider = ethers.getDefaultProvider(process.env.ETHEREUM_NETWORK, {
    etherscan: process.env.ETHERSCAN_API_KEY,
  });

  const contract = new ethers.Contract(ALIEN888_ITEMS_CONTRACT_ADDRESS, ALIEN888_ITEMS_ABI, provider);
  const signer = new ethers.Wallet(process.env.PRIVATE_KEY, provider);


  //1
  const tokenId1 = 1;
  let a1 = await contract.getPrice(tokenId1);
  console.log(`token #1 > price: ${a1}`);

  //2
  let newPrice = 4.2;
  console.log(`token #1 > set price: ${newPrice}`);
  await contract.setPrice(tokenId1, ethers.utils.parseEther(newPrice));
  let a2 = await contract.getPrice(tokenId1);
  console.log(`token #1 > price: ${ethers.utils.formatEther(a2)}`);

}

main()
  .catch(err => {
    console.log(err);
  })
  .finally(() => {
    process.exit();
});