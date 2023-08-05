new Vue({
  el: "#frame",
  data: {},
  mounted() {
  },
  methods: {
    setPFP() {
      html2canvas(document.getElementById("capture")).then(
        (canvas) => {
          let imageAvatarInBase64 = canvas
            .toDataURL("image/jpeg")
            .replace("image/jpeg", "image/octet-stream");

          fetch("/current_custom_avatar", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({
              wallet_address: window.userWalletAddress,
              image_in_base64: imageAvatarInBase64,
              //TODO
              main_nft_token_id: 0,
              secondary_nft1_token_id: 0,
              secondary_nft2_token_id: 0,
              secondary_nft3_token_id: 0,
            }),
          })
          .then(resp => resp.json())
          .then(resp => {
                document
                  .getElementById("pfpImage")
                  .setAttribute("src", imageAvatarInBase64);

                alert("avatar updated");
          })
          .catch(e => {
            console.log(`error: ${e}`);
          });
        }
      );
    },

    onDownload() {
      html2canvas(document.getElementById("capture")).then(
        (canvas) => {
          const obNode = document.createElement("a");
          obNode.href = canvas
            .toDataURL("image/jpeg")
            .replace("image/jpeg", "image/octet-stream");

          obNode.download = document.querySelector(".item-drop-btn-main p").textContent + ".jpg";
          obNode.click();
        }
      );
    },
  },
});

// alien
const aliens = document.querySelectorAll(".item-whitelist");
const alienImg = document.getElementById("alien");
const chooseBg = document.getElementById("choose-bg");
const chooseBgREmove = document.querySelector(".choose-bg-remove");
const btnDownload = document.getElementById("download");

//FIXME already defined via 'alienImg' variable
let nftMainId = document.getElementById("alien");
let nftBg = document.getElementById("background");
let nftItem1 = document.getElementById("item1");
let nftItem2 = document.getElementById("item2");


// Dropdown
const itemDropdown = document.querySelector(".item-dropdown");
const itemDropdownMainBtn = document.querySelector(".item-drop-btn-main");
const itemDropdownMenu = document.querySelector(".item-drop-menu");

if (itemDropdownMainBtn) {
  itemDropdownMainBtn.addEventListener("click", () => {
    itemDropdown.classList.toggle("active");
  });
}

aliens.forEach((el) => {
  el.addEventListener("click", (e) => {

    //1
    // let alien = e.target.getAttribute("data-alien");
    // alienImg.setAttribute("src", alien);
    // document.querySelectorAll(".sm-item-alien").forEach((el) => {
    //   el.setAttribute("src", alien);
    // });


    //2
    let itemName = el.querySelector("p").getAttribute("data-name");
    itemDropdownMainBtn.setAttribute("data-name", itemName);
    itemDropdownMainBtn.querySelector("p").textContent = itemName;
    let selectedAlienTokenID = e.target.getAttribute("data-id");
    alienImg.setAttribute("data-id", selectedAlienTokenID);
    itemDropdown.classList.remove("active");

    //TODO
    //window.location.href = currentRelativeURL + "?token_id=" + selectedAlienTokenID;
  });
});


// Select bg
const itemsBg = document.querySelectorAll(".item-select-bg");
const itemsBg1 = document.querySelectorAll(".item-select-1");
const itemsBg2 = document.querySelectorAll(".item-select-2");

const defaultBg = "./assets/img/green-bg.png";
const defaultItem = "./assets/img/placeholder.png";

itemsBg.forEach((el) => {
  el.addEventListener("click", (e) => {
    //TODO
    alienImg.setAttribute("src", transparentImageOfSelectedMainNft);

    let bg = e.target.getAttribute("data-bg");
    let bgId = e.target.getAttribute("data-id");
    nftBg.setAttribute("data-id", bgId);
    nftBg.setAttribute("src", bg);

    document.getElementById("choose-item-1").classList.add("active");
    document.getElementById("choose-item-1-bg").setAttribute("src", bg);

    //TODO
    toggleDownloadButtonEnabledAttribute();
  });
});

if (document.querySelector(".choose-bg-remove")) {
  //TODO
  // alienImg.setAttribute("src", nonTransparentImageOfSelectedMainNft);

  document
    .querySelector(".choose-bg-remove")
    .addEventListener("click", (e) => {
      e.preventDefault();
      document.getElementById("choose-item-1").classList.remove("active");
      // nftBg.setAttribute("src", "");
      nftBg.setAttribute("src", defaultBg);
      removeActiveItem1();

      //TODO
      toggleDownloadButtonEnabledAttribute();
    });
}

// Select item1
itemsBg1.forEach((el) => {
  el.addEventListener("click", (e) => {
    let item = e.target.getAttribute("data-item");
    let item1 = e.target.getAttribute("data-id");
    document.getElementById("item1").setAttribute("data-id", item1);
    document.getElementById("choose-item-2").classList.add("active");
    document.getElementById("item1").setAttribute("src", item);

    //TODO
    toggleDownloadButtonEnabledAttribute();
  });
});

if (document.querySelector(".choose-item1-remove")) {
  document
    .querySelector(".choose-item1-remove")
    .addEventListener("click", (e) => {
      e.preventDefault();
      document.getElementById("choose-item-2").classList.remove("active");
      document.getElementById("item1").setAttribute("src", defaultItem);
      removeActiveItem2();

      //TODO
      toggleDownloadButtonEnabledAttribute();
    });
}

// Select item2
itemsBg2.forEach((el) => {
  el.addEventListener("click", (e) => {
    let item = e.target.getAttribute("data-item");
    let item2 = e.target.getAttribute("data-id");
    document.getElementById("item2").setAttribute("data-id", item2);
    document.getElementById("choose-item-3").classList.add("active");
    document.getElementById("item2").setAttribute("src", item);

    //TODO
    toggleDownloadButtonEnabledAttribute();
  });
});

if (document.querySelector(".choose-item2-remove")) {
  document
    .querySelector(".choose-item2-remove")
    .addEventListener("click", (e) => {
      e.preventDefault();
      document.getElementById("choose-item-3").classList.remove("active");
      document.getElementById("item2").setAttribute("src", defaultItem);
      removeActiveItem3();

      //TODO
      toggleDownloadButtonEnabledAttribute();
    });
}

function toggleDownloadButtonEnabledAttribute() {
  let nftBgId = parseInt(nftBg.getAttribute("data-id"));
  let nftItemFirst = parseInt(nftItem1.getAttribute("data-id"));
  let nftItemSecond = parseInt(nftItem2.getAttribute("data-id"));

  let c1 = initialSelectedTokenIDForBackground !== nftBgId;
  let c2 = initialSelectedTokenIDForItem1 !== nftItemFirst;
  let c3 = initialSelectedTokenIDForItem2 !== nftItemSecond;

  btnDownload.disabled = c1 || c2 || c3;
}

// Modal
const approveBtn = document.querySelector(".approve");
const overlay = document.querySelector(".overlay");
const overlayClose = document.querySelector(".modal-close");

// Data
const attr = {
  selected_tokens: {
    main_alien_id: "",
    item1_bg: "",
    item_2: "",
    item_3: "",
  },
};


approveBtn.addEventListener("click", () => {
  isAnyMetaMaskConnected().then(result => {
    if (result) {
      // overlay.classList.add("is-open");
      let nftMainImgId = nftMainId.getAttribute("data-id");
      let nftBgId = nftBg.getAttribute("data-id");
      let nftItemFirst = nftItem1.getAttribute("data-id");
      let nftItemSecond = nftItem2.getAttribute("data-id");

      attr.selected_tokens.main_alien_id = nftMainImgId;
      attr.selected_tokens.item1_bg = nftBgId;
      attr.selected_tokens.item_2 = nftItemFirst;
      attr.selected_tokens.item_3 = nftItemSecond;

      const data = {
        main_nft_id: attr.selected_tokens.main_alien_id,
        secondary_nft_ids: [
          attr.selected_tokens.item1_bg,
          attr.selected_tokens.item_2,
          attr.selected_tokens.item_3
        ]
      };

      signDataWithWallet(window.ethersProvider, data)
          .then((result) => {
            if (result.signature) {
              let postData = {
                data: data,
                wallet_address: result.address,
                signature: result.signature
              };

              fetch("/wardrobe", {
                method: "POST",
                headers: {"Content-Type": "application/json"},
                body: JSON.stringify(postData),
              })
              .then(resp => resp.json())
              .then(resp => {
                // overlay.classList.add("is-open");
                console.log(JSON.stringify(resp));
                alert("Updated successfully");
                location.reload();
              })
              .catch(e => {
                alert(`error: ${e}`);
              });
            }
          })
          .catch((err) => {
            const errMsg = `Error signing data: ${err}`;
            alert(errMsg);
            console.error(errMsg);
          });
      } else {
        alert("You have to log in via MetaMask first");
      }
  })
});

overlayClose.addEventListener("click", () => {
  overlay.classList.remove("is-open");
});

overlay.addEventListener("click", (e) => {
  if (e.target === overlay) {
    overlay.classList.remove("is-open");
  }
});

function removeActiveItem1() {
  nftItems1.forEach((el) => {
    if (!el.classList.contains("active")) {
      let elemNumber = el
        .querySelector(".balance b")
        .getAttribute("data-value");
      el.querySelector(".balance b").textContent = elemNumber;
    }
  });
}

function removeActiveItem2() {
  nftItems2.forEach((el) => {
    if (!el.classList.contains("active")) {
      let elemNumber = el
        .querySelector(".balance b")
        .getAttribute("data-value");
      el.querySelector(".balance b").textContent = elemNumber;
    }
  });
}

function removeActiveItem3() {
  nftItems2.forEach((el) => {
    if (!el.classList.contains("active")) {
      let elemNumber = el
        .querySelector(".balance b")
        .getAttribute("data-value");
      el.querySelector(".balance b").textContent = elemNumber;
    }
  });
}

const nftItemsNumber1 = document.querySelectorAll(
  ".items-content-1 .balance b"
);
const nftItemsNumber2 = document.querySelectorAll(
  ".items-content-2 .balance b"
);
const nftItemsNumber3 = document.querySelectorAll(
  ".items-content-3 .balance b"
);

const nftItems = document.querySelectorAll(".item-select");
const nftItems1 = document.querySelectorAll(
  ".items-content-1 .item-select"
);
const nftItems1Active = document.querySelectorAll(
  ".items-content-1 .item-select.active"
);
const nftItems2 = document.querySelectorAll(
  ".items-content-2 .item-select"
);
const nftItems2Active = document.querySelectorAll(
  ".items-content-2 .item-select.active"
);
const nftItems3 = document.querySelectorAll(
  ".items-content-3 .item-select"
);
const nftItems3Active = document.querySelectorAll(
  ".items-content-3 .item-select.active"
);

nftItems.forEach((el) => {
  if (el.classList.contains("active")) {
    let curNumber = el
      .querySelector(".balance b")
      .getAttribute("data-value");
    el.querySelector(".balance b").textContent = curNumber - 1;
  }
});

nftItems1.forEach((item) => {
  item.addEventListener("click", (e) => {
    let number = item
      .querySelector(".balance b")
      .getAttribute("data-value");
    if (!item.classList.contains("active")) {
      nftItemsNumber1.forEach((el) => {
        el.textContent = number;
      });
      nftItems1Active.forEach((el) => {
        el.querySelector(".balance b").textContent = number - 1;
      });

      item.querySelector(".balance b").textContent = number - 1;
    } else if (item.classList.contains("active")) {
      nftItemsNumber1.forEach((el) => {
        el.textContent = number;
      });
      item.querySelector(".balance b").textContent = number - 1;
    }
  });
});

nftItems2.forEach((item) => {
  item.addEventListener("click", (e) => {
    let number = item
      .querySelector(".balance b")
      .getAttribute("data-value");
    if (!item.classList.contains("active")) {
      nftItemsNumber2.forEach((el) => {
        el.textContent = number;
      });
      nftItems2Active.forEach((el) => {
        el.querySelector(".balance b").textContent = number - 1;
      });

      item.querySelector(".balance b").textContent = number - 1;
    } else if (item.classList.contains("active")) {
      nftItemsNumber2.forEach((el) => {
        el.textContent = number;
      });
      item.querySelector(".balance b").textContent = number - 1;
    }
  });
});

nftItems3.forEach((item) => {
  item.addEventListener("click", (e) => {
    let number = item
      .querySelector(".balance b")
      .getAttribute("data-value");
    if (!item.classList.contains("active")) {
      nftItemsNumber3.forEach((el) => {
        el.textContent = number;
      });
      nftItems3Active.forEach((el) => {
        el.querySelector(".balance b").textContent = number - 1;
      });

      item.querySelector(".balance b").textContent = number - 1;
    } else if (item.classList.contains("active")) {
      nftItemsNumber3.forEach((el) => {
        el.textContent = number;
      });
      item.querySelector(".balance b").textContent = number - 1;
    }
  });
});

// ITEMS
const imageTabBtns = document.querySelectorAll(".item-btn");
const itemsTabContents = document.querySelectorAll(".items-content");
imageTabBtns.forEach((el) => {
  el.addEventListener("click", (e) => {
    let tab = e.target.getAttribute("data-tab");

    imageTabBtns.forEach((el) => {
      el.classList.remove("selected");
    });
    el.classList.add("selected");

    itemsTabContents.forEach((el) => {
      el.classList.remove("active");
    });
    document.querySelector(`.items-content-${tab}`).classList.add("active");
  });
});
