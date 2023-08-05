// Amount Dropdown
const amountDrop = document.querySelector(".amount-drop");
const amountDropBtn = document.querySelector(".amount-drop-main");
const amountDropItems = document.querySelectorAll(".amount-drop-items button");
const amountMintBtn = document.querySelector(".mint-free-btn");
const amountMintText = document.querySelector(".mint-free-btn span");

if (amountDrop) {
  amountDropBtn.addEventListener("click", () => {
    amountDrop.classList.toggle("active");
  });

  amountDropItems.forEach((el) => {
    el.addEventListener("click", (e) => {
      let num = Number(e.target.textContent);

      amountDropBtn.setAttribute("value", num);
      amountDrop.classList.remove("active");
      amountDropBtn.value = num;

      if (amountDropBtn.value >= 1) {
        amountMintBtn.style.display = "block";
        amountMintText.textContent = num;
      }
    });
  });

  amountDropBtn.addEventListener("input", (e) => {
    let input = e.target.value;
    amountMintBtn.style.display = "block";
    amountMintText.textContent = input;
  });
  amountDropBtn.addEventListener("focus", (e) => {
    e.target.value = "";
  });
}

const contractAbi = [
  {
    "constant": true,
    "inputs": [
      {"name": "account", "type": "address"},
      {"name": "id", "type": "uint256"}
    ],
    "name": "balanceOf",
    "outputs": [{"name": "", "type": "uint256"}],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [{"name": "id", "type": "uint256"}],
    "name": "uri",
    "outputs": [{"name": "", "type": "string"}],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      {"name": "tokenId", "type": "uint256"},
      {"name": "amount", "type": "uint256"},
    ],
    "name": "mint",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  }
];

async function mintTokens(event) {
  try {
    if (window.ethereum && window.ethereum.selectedAddress) {
      const provider = new ethers.providers.Web3Provider(window.ethereum);
      const currNtw = await provider.getNetwork();
      if (chainID === currNtw.chainId) {
        const signer = provider.getSigner();
        const erc1155Contract = new ethers.Contract(contractAddress, contractAbi, signer);
        const tokenId = currentTokenID;
        let amount = 1;
        let amountElem = document.getElementById("mint-amount");
        if (amountElem) {
          amount = amountElem.options[amountElem.selectedIndex].value;
        }

        erc1155Contract.mint(tokenId, amount)
          .then(res => {
            alert(`Minted successfully: token #${tokenId} (${amount})`);
            // location.reload();
          })
          .catch(err => {
            console.log(`err: ${err}`);
            alert("Error calling mint() method:\n\n" + JSON.stringify(err.error["data"], null, "\t"));
        });
      } else {
        alert("Wrong network in MetaMask; switch to the chain ID: " + chainID);
      }
    } else {
      alert('Log in into or unlock MetaMask');
    }
  } catch (err) {
    console.error('Error calling "mint" method:', err);
    alert(err);
  }
}
