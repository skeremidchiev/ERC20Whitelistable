// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./Whitelistable.sol";

contract ERC20Whitelistable is ERC20, Whitelistable {
    constructor()
        ERC20("LimeChain exam token", "LET")
        public
    {
    }

    function mint(
        address account,
        uint256 amount
    ) 
        public
    {
        super._mint(account, amount);
    }

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) 
        internal
        virtual
        override
        onlyWhitelistedMember(to)
    {
        super._beforeTokenTransfer(from, to, amount);
    }
}