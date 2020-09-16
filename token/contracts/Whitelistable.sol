// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "@openzeppelin/contracts/access/Ownable.sol";

contract Whitelistable is Ownable {
    mapping (address => bool) private whitelist;

    modifier onlyWhitelistedMember(address _address) {
        require(
            whitelist[_address] == true,
            "Whitelistable: caller is not the whitelister"
        );
        _;
    }
      
    function addToWhitelist(address _address)
        public
        onlyOwner
    {
        whitelist[_address] = true;
    }
    
    function removeFromWhitelist(address _address) 
        public
        onlyOwner
    {
        whitelist[_address] = false;
    }

}