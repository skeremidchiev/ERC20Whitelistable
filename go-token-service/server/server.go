package server

import (
	"log"
	"sync"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"ERC20Whitelistable/go-token-service/token"
)


var (
	wlt *token.WhitelistableToken // token context common for all handlers
)

const (
	internalServerError = "Internal Server Error!"
)

func auth(fn http.HandlerFunc) http.HandlerFunc {
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
	log.Println("Endpoint: homePage")

	w.Write([]byte("Welcome to the HomePage!"))
}

func whitelistHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint: whitelist")
	var input token.WhitelistInput

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	err = json.Unmarshal(reqBody, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	output, _ := wlt.WhitelistAddress(&input)
	json.NewEncoder(w).Encode(output)
}

func whitelistMultipleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint: whitelist multiple")
	var input token.WhitelistMultiInput

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	err = json.Unmarshal(reqBody, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	multiOutput := token.GetTxMultiOutput()
	var wg sync.WaitGroup

	for _, addr := range input.Addresses {
		// skip incorrect inputs
		if addr.Address == "" {
			continue
		}

		wg.Add(1)
		go func(t token.WhitelistInput, wg *sync.WaitGroup) {
			defer wg.Done()
			// no need to handle the error
			output, _ := wlt.WhitelistAddress(&t)
			multiOutput.Add(output)
		}(addr, &wg)
	}
	wg.Wait()

	json.NewEncoder(w).Encode(multiOutput)
}

func mintHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint: mint")
	var input token.MintInput

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	err = json.Unmarshal(reqBody, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	output, _ := wlt.Mint(&input)
	json.NewEncoder(w).Encode(output)
}

func mintMultipleHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint: mint multiple")
	var input token.MintMultiInput

	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	err = json.Unmarshal(reqBody, &input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(internalServerError))
		return
	}

	multiOutput := token.GetTxMultiOutput()
	var wg sync.WaitGroup

	for _, mint := range input.Mints {
		// skip incorrect inputs
		if mint.Address == "" {
			continue
		}

		wg.Add(1)
		go func(t token.MintInput, wg *sync.WaitGroup) {
			defer wg.Done()
			// no need to handle the error
			output, _ := wlt.Mint(&t)
			multiOutput.Add(output)
		}(mint, &wg)
	}
	wg.Wait()

	json.NewEncoder(w).Encode(multiOutput)
}


func Run() {
	// initialize token context
	var err error
	wlt, err = token.GetWhitelistableToken()
	if err != nil {
		log.Println("Can't setup token context: ", err)
		return
	}

	http.HandleFunc("/", auth(homePageHandler))
	http.HandleFunc("/whitelist", auth(whitelistHandler))
	http.HandleFunc("/whitelist/multiple", auth(whitelistMultipleHandler))
	http.HandleFunc("/mint", auth(mintHandler))
	http.HandleFunc("/mint/multiple", auth(mintMultipleHandler))

	log.Println("Server starting ...")
	defer log.Println("Server shutting down ...")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Server failed with error: ", err)
		return
	}
}
