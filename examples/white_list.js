import * as dotenv from 'dotenv';
import * as ethers from 'ethers';
import * as fs from 'fs';
import yargs from "yargs";
import { hideBin } from 'yargs/helpers';
import { getArgs } from './cmd_args.js';

dotenv.config();

const ALIEN888_ITEMS_ABI_JSON_FILE_NAME = "../artifacts/contracts/Alien888Item.sol/Alien888Item.json";
const ALIEN888_ITEMS_ABI = JSON.parse(fs.readFileSync(ALIEN888_ITEMS_ABI_JSON_FILE_NAME)).abi;
const ALIEN888_ITEMS_CONTRACT_ADDRESS = process.env.ALIEN888_ITEMS_CONTRACT_ADDRESS;

function isValidEthereumAddress(address) {
  if (!address || typeof address !== 'string') {
    return false;
  }

  // Check for the correct length and starting with '0x'
  if (!/^0x[a-fA-F0-9]{40}$/.test(address)) {
    return false;
  }

  return true;
}

function readAddressesFromFile(filePath) {
  try {
    const fileContent = fs.readFileSync(filePath, 'utf-8');
    const addresses = fileContent
      .split('\n')
      .map((line) => line.trim())
      .filter((address) => isValidEthereumAddress(address));

    return addresses;
  } catch (error) {
    console.error('Error reading addresses from file:', error);
    return [];
  }
}

async function main() {
  const provider = new ethers.providers.JsonRpcProvider(process.env.NETWORK_PROVIDER_API_URL);
  const wallet = new ethers.Wallet(process.env.PRIVATE_KEY, provider);
  const contract = new ethers.Contract(ALIEN888_ITEMS_CONTRACT_ADDRESS, ALIEN888_ITEMS_ABI, wallet);

  const cmdArgs = getArgs();

  if (cmdArgs["token-id"] === undefined) {
    console.log("token-id must not be null");
    return;
  }

  const tokenId = parseInt(cmdArgs["token-id"]);
  const walletAddress = cmdArgs["wallet-address"];
  const fileName = cmdArgs["file"];

  switch (cmdArgs["cmd"]) {
    //1
    case "add-wallet":
      console.log("[cmd] add wallet");

      if (walletAddress === undefined) {
        console.log("walletAddress must not be null");
        return;
      }

      const tx1 = await contract.addIntoWhiteList(tokenId, walletAddress);
      const receipt1 = tx1.wait();
      console.log(`[ok] tx.hash: ${tx1.hash}`);
      break;

    //2
    case "print":
      console.log("[cmd] print wallets");
      const res2 = await contract.getWhiteListFor(tokenId);
      console.log(res2);
      break;

    //3
    case "remove-wallet":
      console.log("[cmd] remove wallet");

      if (walletAddress === undefined) {
        console.log("walletAddress must not be null");
        return;
      }

      const tx3 = await contract.removeFromWhiteList(tokenId, walletAddress);
      const receipt3 = await tx3.wait();
      console.log(`[ok] tx.hash: ${tx3.hash}`);
      break;

    //4
    case "add-wallets":
      console.log("[cmd] add wallets from a file");

      if (fileName === undefined) {
        console.log("file name must not be null");
        return;
      }

      const addresses = readAddressesFromFile(fileName);
      const tx4 = await contract.batchAddIntoWhiteList(tokenId, addresses);
      const receipt4 = tx4.wait();
      console.log(`[ok] tx.hash: ${tx4.hash}`);
      break;

    default:
      console.log("no command provided; provide: [add-wallet | add-wallets| print | remove-wallet]");
    }
}

main()
  .catch(err => {
    console.log(err);
  })
  .finally(() => {
    process.exit();
});