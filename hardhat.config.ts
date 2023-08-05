import { config as dotEnvConfig } from "dotenv";
import { HardhatUserConfig } from "hardhat/config";
import "@nomicfoundation/hardhat-toolbox";

dotEnvConfig();

const config: HardhatUserConfig = {
  solidity: "0.8.18",
  defaultNetwork: "arbitrum_goerli",
  networks: {
    ethereum_goerli: {
      url: "https://goerli-rollup.arbitrum.io/rpc",
      chainId: 5,
      accounts: [process.env.ETHEREUM_GOERLI_PRIVATE_KEY]
    },
    ethereum_sepolia: {
      url: `https://sepolia.infura.io/v3/${process.env.INFURA_API_KEY}`,
      chainId: 11155111,
      accounts: [process.env.ETHEREUM_GOERLI_PRIVATE_KEY]
    },
    arbitrum_goerli: {
      url: "https://goerli-rollup.arbitrum.io/rpc",
      chainId: 421613,
      accounts: [process.env.ETHEREUM_GOERLI_PRIVATE_KEY]
    },
    arbitrum_mainnet: {
      url: "https://arb1.arbitrum.io/rpc",
      // accounts: ["todo"]
    },
  }
};

export default config;