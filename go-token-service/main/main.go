package main

import (
	"flag"
	"fmt"

	"ERC20Whitelistable/go-token-service/config"
	"ERC20Whitelistable/go-token-service/whitelistableToken"
)

func main() {
	// adding configuration file flag
	cfpathFlag := flag.String("cfpath", "", "Configuration file.")
	flag.Parse()

	config.SetConfigFilePath(*cfpathFlag)

	// TESTING
	// wlt, err := whitelistableToken.GetWhitelistableToken()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// err = wlt.WhitelistAddress("0x95ec0f286C94534BB6033265e3BF1D6F038E8704")
	// fmt.Println("WhitelistAddress err: ", err)

	// err = wlt.Mint("0x0eC47Eb3645EFf7499fA0C485c2b9D2e1Db2595b", 99)
	// fmt.Println("Mint err: ", err)

	// err = wlt.Mint("0x5A136663a33DaC49362Ad6B4d50bFaE9a8d36002", 99)
	// fmt.Println("Mint err: ", err)
}
