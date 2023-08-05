import { ethers } from "hardhat";

async function main() {
  const name = "Alien888Hero";
  const Alien888Hero = await ethers.getContractFactory(name);
  const c = await Alien888Hero.deploy(name, "Alien888");
  await c.deployed();
  console.log(`'${name}' deployed to: ${c.address}`);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exitCode = 1;
  });
