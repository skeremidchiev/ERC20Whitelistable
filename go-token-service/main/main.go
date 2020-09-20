package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"encoding/json"
	"io/ioutil"
	"net/http"

	"ERC20Whitelistable/go-token-service/whitelistableToken"
)

func Auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		// just a dummy validation
		if !ok || user != "admin" || pass != "pass" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized!"))
			return
		}
		fn(w, r)
	}
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: homePage")

	w.Write([]byte("Welcome to the HomePage!"))
}

func whitelisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint: whitelist")
	var input whitelistableToken.WhitelistInput

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error!"))
		return
	}

	json.Unmarshal(reqBody, &input)

	wlt, err := whitelistableToken.GetWhitelistableToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error!"))
		return
	}

	output, err := wlt.WhitelistAddress(&input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error!"))
		return
	}

	json.NewEncoder(w).Encode(output)
}

func main() {
	// adding configuration file flag
	cfpathFlag := flag.String("cfpath", "", "Configuration file.")
	flag.Parse()

	if len(os.Args) != 2 {
		log.Fatal("Give valid configuration path!")
	}

	whitelistableToken.SetConfigFilePath(*cfpathFlag)

	http.HandleFunc("/", Auth(homePageHandler))
	http.HandleFunc("/whitelist", Auth(whitelisterHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
