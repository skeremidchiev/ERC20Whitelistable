package main

import (
    "fmt"
    "flag"

    "ERC20Whitelistable_Token_Service/config"
	"ERC20Whitelistable_Token_Service/whitelistableToken"
)

func main() {
    // adding configuration file flag
    cfpathFlag := flag.String("cfpath", "", "Configuration file.")
    flag.Parse()


    config.SetConfigFilePath(*cfpathFlag)

    // TESTING
	wlt, err := whitelistableToken.GetWhitelistableToken()
	if err != nil {
		fmt.Println(err)
		return
	}

    // err = wlt.WhitelistAddress("0x7b16d00dC38bf023f75116A1e0b405601f0d271C")
    // fmt.Println("WhitelistAddress err: ", err)

    // err = wlt.Mint("0x7b16d00dC38bf023f75116A1e0b405601f0d271C", 99)
    // fmt.Println("Mint err: ", err)

    err = wlt.Mint("0x5A136663a33DaC49362Ad6B4d50bFaE9a8d36002", 99)
    fmt.Println("Mint err: ", err)
}
