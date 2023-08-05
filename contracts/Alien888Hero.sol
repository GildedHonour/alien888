// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Alien888Hero is ERC721, Ownable {
    constructor(string memory _name, string memory _symbol) ERC721(_name, _symbol) {}

    function mint(address _to, uint256 _tokenId) external onlyOwner {
        _mint(_to, _tokenId);
    }

    function _baseURI() internal pure override returns (string memory) {
        return "https://alien888.projects.incerteza.one/erc721-tokens/";
    }
}
