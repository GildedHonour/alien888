import * as dotenv from 'dotenv';
import minimist from 'minimist';
import * as ethers from 'ethers';
import * as fs from 'fs';

dotenv.config();

const ALIEN888_ITEMS_ABI_JSON_FILE_NAME = "../artifacts/contracts/Alien888Item.sol/Alien888Item.json";
const ALIEN888_ITEMS_ABI = JSON.parse(fs.readFileSync(ALIEN888_ITEMS_ABI_JSON_FILE_NAME)).abi;
const ALIEN888_ITEMS_CONTRACT_ADDRESS = process.env.ALIEN888_ITEMS_CONTRACT_ADDRESS;

async function main() {
  const provider = new ethers.providers.JsonRpcProvider(process.env.NETWORK_PROVIDER_API_URL);
  const wallet = new ethers.Wallet(process.env.PRIVATE_KEY, provider);
  const contract = new ethers.Contract(ALIEN888_ITEMS_CONTRACT_ADDRESS, ALIEN888_ITEMS_ABI, wallet);

  //1
  const tokenId = 3;
  const amount = 2;
  const tx = await contract.mint(tokenId, amount);
  const receipt = await tx.wait();
  console.log(`ok; tx.hash: ${tx.hash}`);
}

main()
  .catch(err => {
    console.log(err);
  })
  .finally(() => {
    process.exit();
});