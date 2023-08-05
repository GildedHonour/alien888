const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Alien888Item", function() {
  let c: Contract;
  let accounts;

  beforeEach(async () => {
    accounts = await ethers.getSigners();

    const CF = await ethers.getContractFactory("Alien888Item");
    c = await CF.deploy();
    await c.deployed();
  });

  it("Should return 0 for non-existing token", async () => {
    let price = await c.getPrice(1111);
    expect(price).to.equal(0);
  });

  it("Should add/remove tokens into/from WhiteList", async () => {
    //1
    const addr1 = "0x10ED43C718714eb63d5aA57B78B54704E256024E";
    const tokenId1 = 1;
    await c.addIntoWhitelist(tokenId1, addr1);
    let a1 = await c.isInWhitelist(tokenId1, addr1);
    expect(a1).to.equal(true);

    await c.removeFromWhitelist(tokenId1, addr1);
    let a1_2 = await c.isInWhitelist(tokenId1, addr1);
    expect(a1_2).to.equal(false);

    //2
    const addr2 = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73";
    const tokenId2 = 999;
    let a2 = await c.isInWhitelist(tokenId2, addr2);
    expect(a2).to.equal(false);

    await c.removeFromWhitelist(tokenId2, addr2);
    let a2_2 = await c.isInWhitelist(tokenId2, addr2);
    expect(a2_2).to.equal(false);
  });

  it("Should update and return correct prices", async () => {
    await c.setPrice(1, 123);
    let p1 = await c.getPrice(1);
    expect(p1).to.equal(123);

    await c.setPrice(1, 124);
    p1 = await c.getPrice(1);
    expect(p1).to.equal(124);
    expect(p1).to.not.equal(123);
  });

  it("Should not permitted to mint if not in WhiteList", async () => {
    const tokenId = 1;
    const amount = 6;
    await expect(c.mint2(tokenId, amount)).to.be.reverted;
  });

  it("Should mint if in WhiteList", async () => {
    const tokenId1 = 1;
    await c.addIntoWhitelist(tokenId1, accounts[0].address);
    let a1 = await c.isInWhitelist(tokenId1, accounts[0].address);
    expect(a1).to.equal(true);

    const amount = 4;
    await c.mint2(tokenId1, amount);
    const balance = await c.balanceOf(accounts[0].address, tokenId1);
    expect(amount).to.equal(Number(balance.toString()));


    const tokenId2 = 555;
    const balance2 = await c.balanceOf(accounts[0].address, tokenId2);
    expect(0).to.equal(Number(balance2.toString()));
  });
});