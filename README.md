# ERC20Whitelistable

## TOKEN:
ERC20 Token with Whitelistable and Mintable functionality

## GO-TOKEN-SERVICE:
###SetUp:


###Run:
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
