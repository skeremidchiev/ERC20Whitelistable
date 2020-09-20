package main

import (
	"flag"
	"log"
	"os"

	"ERC20Whitelistable/go-token-service/token"
	"ERC20Whitelistable/go-token-service/server"
)

func main() {
	// adding configuration file flag
	cfpathFlag := flag.String("cfpath", "", "Configuration file.")
	flag.Parse()

	if len(os.Args) != 2 {
		log.Fatal("Give valid configuration path!")
	}

	token.SetConfigFilePath(*cfpathFlag)

	server.Run()
}
