"use sctrict";

const USER_WALLET_ADDRESS_KEY = "user-wallet-address";
window.userWalletAddress = null;
const walletLoginBtn = document.getElementById('wallet-login');
const walletLogin2Btn = document.getElementById('wallet-login2');
const userWallet = document.getElementById('user-wallet');
const walletAddrTopNav = document.getElementById("wallet-address-top-nav");
const pfpImageElem = document.getElementById("pfpImage");

const CONNECT_WALLET = "Connect Wallet";
const CONNECTED = "Connected";
const WALLET_DASHES = "----";

async function loginWithMetaMask(event) {
  const accounts = await window.ethereum.request({
    method: 'eth_requestAccounts'
  })
  .catch((e) => {
    console.error(e.message);
    return;
  });
}

async function isAnyMetaMaskConnected() {
  if ((window.ethereum) && (window.ethersProvider)) {
    const accounts = await window.ethersProvider.listAccounts();
    return accounts.length > 0;
  }

  return false;
}

function walletAddressToShortVersion(addr) {
  return addr.substring(0, 4) + "..." + addr.substring(addr.length - 4, addr.length)
}

//TODO
function clearMetaMaskConnection() {
  window.localStorage.removeItem(USER_WALLET_ADDRESS_KEY);

  walletAddrTopNav.innerText = WALLET_DASHES;
  walletAddrTopNav.title = "";
  walletLoginBtn.innerText = CONNECT_WALLET;
  walletLoginBtn.disabled = false;

  walletLogin2Btn.innerText = CONNECT_WALLET;
  walletLogin2Btn.disabled = false;

  pfpImageElem.classList.remove("header-avatar");
  pfpImageElem.src = pfpImageElem.src.replace("user-avatar", "avatar-anonym");
}

async function signDataWithWallet(provider, data) {
  try {
    const message = JSON.stringify(data);
    const signer = provider.getSigner();
    const address = await signer.getAddress();
    const signature = await signer.signMessage(message);
    return { signature, address };
  } catch (error) {
    console.error('Failed to sign the data:', error);
    return null;
  }
}

window.addEventListener('DOMContentLoaded', () => {
  if (window.ethereum) {
    window.ethersProvider = new ethers.providers.Web3Provider(window.ethereum);

    walletLoginBtn.disabled = false;
    walletLogin2Btn.disabled = false;
    let usrWa = window.localStorage.getItem(USER_WALLET_ADDRESS_KEY);

    fetch("/current_wallet", {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({wallet_address: usrWa}),
    })
    .then(resp => resp.json())
    .then(resp => {
      console.log(`current_wallet > resp: ${resp}`);
    })
    .catch(e => {
      console.log(`error: ${e}`);
    });

    if (usrWa) {
      window.userWalletAddress = usrWa;
      let shortWalletAddr = walletAddressToShortVersion(usrWa);
      walletAddrTopNav.innerText = shortWalletAddr;
      walletAddrTopNav.title = usrWa;

      walletLoginBtn.innerText = CONNECTED;
      walletLoginBtn.disabled = true;
      walletLoginBtn.style.pointerEvents = "none";


      pfpImageElem.classList.add("header-avatar");
      pfpImageElem.src = pfpImageElem.src.replace("avatar-anonym", "user-avatar");

      walletLogin2Btn.innerText = CONNECTED;
      walletLogin2Btn.disabled = true;
      walletLogin2Btn.style.pointerEvents = "none";
    } else {
      //TODO
      walletAddrTopNav.innerText = WALLET_DASHES;
      walletAddrTopNav.title = usrWa;

      pfpImageElem.classList.remove("header-avatar");
      pfpImageElem.src = pfpImageElem.src.replace("user-avatar", "avatar-anonym");

      walletLoginBtn.addEventListener('click', loginWithMetaMask);
      walletLogin2Btn.addEventListener('click', loginWithMetaMask);

      walletLoginBtn.style.pointerEvents = "all";
      walletLogin2Btn.style.pointerEvents = "all";
    }

    window.ethereum.on('accountsChanged', (accounts) => {
      if (accounts && accounts.length > 0) {
        let usrWa = accounts[0];
        window.userWalletAddress = usrWa;
        window.localStorage.setItem(USER_WALLET_ADDRESS_KEY, usrWa);
        let shortWalletAddr = walletAddressToShortVersion(usrWa)
        walletAddrTopNav.innerText = shortWalletAddr;
        walletAddrTopNav.title = usrWa;

        walletLoginBtn.innerText = CONNECTED;
        walletLoginBtn.disabled = true;
        walletLoginBtn.removeEventListener('click', loginWithMetaMask);
        walletLoginBtn.style.pointerEvents = "none";

        walletLogin2Btn.innerText = CONNECTED;
        walletLogin2Btn.disabled = true;
        walletLogin2Btn.removeEventListener('click', loginWithMetaMask);
        walletLogin2Btn.style.pointerEvents = "none";

        pfpImageElem.classList.add("header-avatar");
        pfpImageElem.src = pfpImageElem.src.replace("avatar-anonym", "user-avatar");

        //TODO
        fetch("/current_wallet", {
          method: "POST",
          headers: {"Content-Type": "application/json"},
          body: JSON.stringify({wallet_address: usrWa}),
        })
        .then(resp => resp.json())
        .then(resp => {
          console.log(`current_wallet > resp: ${resp}`);
        })
        .catch(e => {
          console.log(`error: ${e}`);
        });
      } else {
        //TODO extract into a function 'clearMetaMaskConnection()'
        window.localStorage.removeItem(USER_WALLET_ADDRESS_KEY);
        walletAddrTopNav.innerText = WALLET_DASHES;
        walletAddrTopNav.title = "";
        walletLoginBtn.innerText = CONNECT_WALLET;
        walletLoginBtn.disabled = false;
        pfpImageElem.classList.remove("header-avatar");
        pfpImageElem.src = pfpImageElem.src.replace("user-avatar", "avatar-anonym");

        walletLogin2Btn.innerText = CONNECT_WALLET;
        walletLogin2Btn.disabled = false;

        location.reload();
      }
    });
  } else {
    window.localStorage.removeItem(USER_WALLET_ADDRESS_KEY);
    walletAddrTopNav.innerText = WALLET_DASHES;
    walletAddrTopNav.title = usrWa;

    walletLoginBtn.disabled = true;
    walletLoginBtn.innerText = 'MetaMask unavailable';

    walletLogin2Btn.disabled = true;
    walletLogin2Btn.innerText = 'MetaMask unavailable';

    // walletLoginBtn.classList.remove('bg-purple-500', 'text-white')
    // walletLoginBtn.classList.add('bg-gray-500', 'text-gray-100', 'cursor-not-allowed')
    return false;
  }

  //one of the accounts gets disconnected
  window.ethereum.on('disconnect', () => {
    const accounts = window.ethersProvider.listAccounts().
      then((accounts) => {
        if (accounts.length == 0) {
          window.localStorage.removeItem(USER_WALLET_ADDRESS_KEY);
          walletAddrTopNav.innerText = WALLET_DASHES;
          walletAddrTopNav.title = "";
          walletLoginBtn.innerText = CONNECT_WALLET;
          walletLoginBtn.style.pointerEvents = "all";
          walletLoginBtn.disabled = false;
          pfpImageElem.classList.remove("header-avatar");
          pfpImageElem.src = pfpImageElem.src.replace("user-avatar", "avatar-anonym");

          walletLogin2Btn.innerText = CONNECT_WALLET;
          walletLogin2Btn.disabled = false;
          walletLogin2Btn.style.pointerEvents = "all";
        } else {
          window.localStorage.setItem(USER_WALLET_ADDRESS_KEY, accounts[0]);
        }

        location.reload();
      });
  });
});
