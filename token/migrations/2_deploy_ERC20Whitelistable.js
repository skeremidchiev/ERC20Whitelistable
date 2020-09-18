const ERC20Whitelistable = artifacts.require("ERC20Whitelistable");

module.exports = function (deployer) {
  deployer.deploy(ERC20Whitelistable);
};
