import { ethers } from "hardhat";

async function main() {
  const name = "Alien888Item";
  const Alien888Item = await ethers.getContractFactory(name);
  const c = await Alien888Item.deploy("https://alien888.projects.incerteza.one/erc1155-tokens/");
  await c.deployed();
  console.log(`'${name}' deployed to: ${c.address}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
