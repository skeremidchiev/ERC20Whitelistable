// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract ERC20Whitelistable is ERC20, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");
    bytes32 public constant WHITELISTED_ROLE = keccak256("WHITELISTED_ROLE");

    constructor()
        ERC20("LimeChain exam token", "LET")
        public
    {
        _setupRole(DEFAULT_ADMIN_ROLE, _msgSender());

        _setupRole(MINTER_ROLE, _msgSender());
        _setupRole(WHITELISTED_ROLE, _msgSender());
    }

    function mint(
        address to,
        uint256 amount
    ) 
        public
    {
        require(
            hasRole(MINTER_ROLE, _msgSender()),
            "ERC20Whitelistable: must have minter role to mint"
        );
        _mint(to, amount);
    }

    function _beforeTokenTransfer(address from, address to, uint256 amount)
        internal
        override
        virtual
    {
        require(
            hasRole(WHITELISTED_ROLE, to),
            "ERC20Whitelistable: must be whitelisted to recieve tokens"
        );
        super._beforeTokenTransfer(from, to, amount);
    }
}
