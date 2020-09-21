# ERC20Whitelistable

## TOKEN:
ERC20 Token with Whitelistable and Mintable functionality

## GO-TOKEN-SERVICE:
### SetUp:

```
go get github.com/skeremidchiev/ERC20Whitelistable
```

**Dependencies**:
https://github.com/ethereum/go-ethereum

**Generating ERC20Whitelistable.go:**

```
solc --abi --bin --allow-paths .  contracts/ERC20Whitelistable.sol -o build
abigen --bin=./build/ERC20Whitelistable.bin --abi=./build/ERC20Whitelistable.abi --pkg=token --out=./path-to-go-project/contracts/ERC20Whitelistable.go
```

**solc --version**
*Version: 0.6.9-develop.2020.5.27+commit.9f407fe0.Linux.g++*


### Run:
**config.json** expected with following structure:

```
{
  "privateKey": "user-private-key",
  "network": "main net / ropsten / etc.",
  "infuraKey": "PROJECT ID",
  "contractAddress": "0xa845bE40dd6CF745EAC313837bf7F1eFfBCF0bE4" // contract address deployed on ropsten
}
```

```
go run main.go --cfpath="path-to-config.json"
```

## POSTMAN COLLECTION
https://www.getpostman.com/collections/ad0e43025d5d091519f8
