// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.18;


import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Strings.sol";
import {Auth, Authority} from "@rari-capital/solmate/src/auth/Auth.sol";
import {ERC2981} from "@openzeppelin/contracts/token/common/ERC2981.sol";

contract Alien888Item is ERC1155, ERC2981, Ownable {
    using Strings for uint256;

    uint public constant TOTAL_SUPPLY = 888;
    string private baseURI;

    constructor(string memory _baseURI) ERC1155(_baseURI) {
        baseURI = _baseURI;
    }

    function getBaseURI() public view returns (string memory) {
        return baseURI;
    }

    function setBaseURI(string memory newBaseURI) external onlyOwner {
        baseURI = newBaseURI;
    }

    function uri(uint256 tokenId) public view override returns (string memory) {
        return bytes(baseURI).length > 0
            ? string(abi.encodePacked(baseURI, tokenId.toString(), ".json"))
            : "";
    }

    //tokenId ==> price
    mapping(uint => uint) private tokenPrices;

    function getPrice(uint tokenId) external view returns (uint) {
        return tokenPrices[tokenId];
    }

    function setPrice(uint tokenId, uint newPrice) external onlyOwner {
        tokenPrices[tokenId] = newPrice;
    }

    //total amount of minted copies per token
    mapping(uint256 => uint256) private totalCopiesMinted;
    // returns the total number of copies minted for a certain token_id
    function getTotalCopiesMinted(uint tokenId) external view returns (uint) {
        return totalCopiesMinted[tokenId];
    }

    // all uniquely minted token ids
    uint256[] private totalUniqueTokenIdsMinted;

    // returns the total number of unique tokens minted
    function getTotalUniqueTokensMinted() external view returns (uint256) {
        return totalUniqueTokenIdsMinted.length;
    }

    uint private MAX_COPIES_PER_TOKEN = 100;
    function getMaxCopiesPerToken() external view returns (uint) {
        return MAX_COPIES_PER_TOKEN;
    }

    //tokenId ==> addresses
    mapping(uint => address[]) private tokenWhiteLists;

    function addIntoWhiteList(uint tokenId, address addr) external onlyOwner {
        bool exists = isInWhiteList(tokenId, addr);
        if (!exists) {
            tokenWhiteLists[tokenId].push(addr);
        }
    }

    function batchAddIntoWhiteList(uint tokenId, address[] memory addresses) external onlyOwner returns (uint256) {
        uint256 count = 0;
        for (uint256 i = 0; i < addresses.length; i++) {
            if (!isInWhiteList(tokenId, addresses[i])) {
                tokenWhiteLists[tokenId].push(addresses[i]);
                count++;
            }
        }
        return count;
    }

    function removeFromWhiteList(uint tokenId, address addr) external onlyOwner {
        require(addr != address(0), "address must not be 0");
        address[] memory wlAddresses = tokenWhiteLists[tokenId];
        if (wlAddresses.length > 0) {
            for (uint i = 0; i < wlAddresses.length; i++) {
                if (wlAddresses[i] == addr) {
                    delete tokenWhiteLists[tokenId][i];
                    break;
                }
            }
        }
    }

    function isInWhiteList(uint tokenId, address addr) public view returns(bool) {
        require(addr != address(0), "address must not be 0");

        address[] memory wlAddresses = tokenWhiteLists[tokenId];
        for (uint i = 0; i < wlAddresses.length; i++) {
            if (wlAddresses[i] == addr) {
                return true;
            }
        }

        return false;
    }

    function getWhiteListFor(uint tokenId) public view returns(address[] memory) {
        return tokenWhiteLists[tokenId];
    }

    uint256 private constant MAX_COPIES_PER_WALLET = 5;

    function getMaxCopiesPerWallet() external pure returns(uint) {
        return MAX_COPIES_PER_WALLET;
    }

    //tokenId ==> (wallet_address => current number of copies)
    mapping(uint256 => mapping(address => uint256)) private copiesPerWallet;

    error ERR_COPIES_PER_TOKEN_MAX_SUPPLY(uint supply);

    function mint(uint tokenId, uint amount) public payable {
        require(tokenId <= TOTAL_SUPPLY, "TOTAL_SUPPLY must not be exceeded");
        require(copiesPerWallet[tokenId][msg.sender] + amount <= MAX_COPIES_PER_WALLET, "Exceeded maximum copies per wallet");
        require(msg.value == tokenPrices[tokenId]);

        bool isFound = isInWhiteList(tokenId, msg.sender);
        require(isFound, "wallet must be in the whitelist");

        (bool success, ) = payable(owner()).call{value: msg.value}("");
        require(success);

        uint copiesNewAmount = totalCopiesMinted[tokenId] + amount;
        if (copiesNewAmount > MAX_COPIES_PER_TOKEN) {
          revert ERR_COPIES_PER_TOKEN_MAX_SUPPLY(MAX_COPIES_PER_TOKEN); 
        }

        _mint(msg.sender, tokenId, amount, "");

        if (!isNewTokenMinted(tokenId)) {
            totalUniqueTokenIdsMinted.push(tokenId);
        }

        totalCopiesMinted[tokenId] = copiesNewAmount;
        copiesPerWallet[tokenId][msg.sender] += amount;
    }


    function ownerMint(uint tokenId, uint amount) public onlyOwner {
        require(tokenId <= TOTAL_SUPPLY, "TOTAL_SUPPLY must not be exceeded");

        uint copiesNewAmount = totalCopiesMinted[tokenId] + amount;
        if (copiesNewAmount > MAX_COPIES_PER_TOKEN) {
          revert ERR_COPIES_PER_TOKEN_MAX_SUPPLY(MAX_COPIES_PER_TOKEN); 
        }

        _mint(msg.sender, tokenId, amount, "");

        // adds the tokenId to the list of totalUniqueTokenIdsMinted only if it's a new one
        if (!isNewTokenMinted(tokenId)) {
            totalUniqueTokenIdsMinted.push(tokenId);
        }

        // update the amount of copies for this tokenId
        totalCopiesMinted[tokenId] = copiesNewAmount;
    }


    // checks whether a token has been minted already
    function isNewTokenMinted(uint256 tokenId) internal view returns (bool) {
        for (uint256 i = 0; i < totalUniqueTokenIdsMinted.length; i++) {
            if (totalUniqueTokenIdsMinted[i] == tokenId) {
                return true;
            }
        }
        return false;
    }

    function burn(address acc, uint tokenId, uint amount) external onlyOwner {
        _burn(acc, tokenId, amount);

        totalCopiesMinted[tokenId] -= amount;
    }

    function setRoyalty(uint id, address receiver, uint96 feeNumerator) external onlyOwner {
        if (receiver == address(0)) {
          return _resetTokenRoyalty(id);
        }

        _setTokenRoyalty(id, receiver, feeNumerator);
    }

    function setRoyalty(address receiver, uint96 feeNumerator) external onlyOwner {
        if (receiver == address(0)) {
          return _deleteDefaultRoyalty();
        }

        _setDefaultRoyalty(receiver, feeNumerator);
    }

    function supportsInterface(bytes4 interfaceId)
        public
        pure
        override(ERC1155, ERC2981)
        returns (bool)
    {
        return
            // ERC165 Interface ID for ERC2981
            interfaceId == 0x2a55205a ||
            // ERC165 Interface ID for ERC1155
            interfaceId == 0xd9b67a26;
    }

    function withdraw() public onlyOwner {
        address receiver = msg.sender;
        (bool res, ) = receiver.call{value: address(this).balance}("");
        if (!res) {
          revert("withdrawal error");
        }
    }
}
