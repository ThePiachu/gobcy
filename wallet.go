package gobcy

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"strconv"
	"strings"
)

//CreateWallet creates a public-address watching wallet
//associated with this token/coin/chain, usable anywhere
//in the API where an Address might be used (just use
//the wallet name instead). For example, with checking
//a wallet name balance:
//  addr, err := api.GetAddrBal("your-wallet-name")
func (api *API) CreateWallet(c context.Context, req Wallet) (wal Wallet, err error) {
	u, err := api.buildURL("/wallets")
	if err != nil {
		return
	}
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err = enc.Encode(&req); err != nil {
		return
	}
	resp, err := postResponse(c, u, &data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&wal)
	return
}

//ListWallets lists all known Wallets associated with
//this token/coin/chain.
func (api *API) ListWallets(c context.Context) (names []string, err error) {
	u, err := api.buildURL("/wallets")
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	jsonResp := new(struct {
		List []string `json:"wallet_names"`
	})
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(jsonResp)
	names = jsonResp.List
	return
}

//GetWallet gets a Wallet based on its name, the associated
//API token/coin/chain, and whether it's an HD wallet or
//not.
func (api *API) GetWallet(c context.Context, name string) (wal Wallet, err error) {
	u, err := api.buildURL("/wallets/" + name)
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//decode JSON into result
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&wal)
	return
}

//AddAddrWallet adds a slice of addresses to a named Wallet,
//associated with the API token/coin/chain. In addition to your
//list of addresses to add, takes one additional parameter:
//  "omitAddr," if true will omit wallet addresses in your
//  response. Useful to speed up the API call for larger wallets.
func (api *API) AddAddrWallet(c context.Context, name string, addrs []string, omitAddr bool) (wal Wallet, err error) {
	params := map[string]string{"omitWalletAddresses": strconv.FormatBool(omitAddr)}
	u, err := api.buildURLParams("/wallets/"+name+"/addresses", params)
	if err != nil {
		return
	}
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err = enc.Encode(&Wallet{Addresses: addrs}); err != nil {
		return
	}
	resp, err := postResponse(c, u, &data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&wal)
	return
}

//GetAddrWallet returns a slice of addresses associated with
//a named Wallet, associated with the API token/coin/chain.
//Offers 4 parameters for customization:
//  "used," if true will return only used addresses
//  "unused," if true will return only unused addresses
//  "zero", if true will return only zero balance addresses
//  "nonzero", if true will return only nonzero balance addresses
//"used" and "unused" cannot be true at the same time; the SDK will throw an error.
//"zero" and "nonzero" cannot be true at the same time; the SDK will throw an error.
func (api *API) GetAddrWallet(c context.Context, name string, used bool, unused bool, zero bool, nonzero bool) (addrs []string, err error) {
	params := make(map[string]string)
	if used && unused {
		err = errors.New("GetAddrWallet: Unused and used cannot be the same")
		return
	}
	if zero && nonzero {
		err = errors.New("GetAddrWallet: Zero and nonzero cannot be the same")
		return
	}
	if used != unused {
		params["used"] = strconv.FormatBool(used)
	}
	if zero != nonzero {
		params["zerobalance"] = strconv.FormatBool(zero)
	}
	u, err := api.buildURLParams("/wallets/"+name+"/addresses", params)
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//decode JSON into result
	var wal Wallet
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&wal)
	addrs = wal.Addresses
	return
}

//DeleteAddrWallet deletes a slice of addresses associated with
//a named Wallet, associated with the API token/coin/chain.
func (api *API) DeleteAddrWallet(c context.Context, name string, addrs []string) (err error) {
	u, err := api.buildURLParams("/wallets/"+name+"/addresses",
		map[string]string{"address": strings.Join(addrs, ";")})
	resp, err := deleteResponse(c, u)
	if err != nil {
		return
	}
	resp.Body.Close()
	return
}

//GenAddrWallet generates a new address within the named Wallet,
//associated with the API token/coin/chain. Also returns the
//private/WIF/public key of address via an Address Keychain.
func (api *API) GenAddrWallet(c context.Context, name string) (wal Wallet, addr AddrKeychain, err error) {
	u, err := api.buildURL("/wallets/" + name + "/addresses/generate")
	resp, err := postResponse(c, u, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//decode JSON into result
	dec := json.NewDecoder(resp.Body)
	//weird anonymous struct composition FTW
	err = dec.Decode(&struct {
		*Wallet
		*AddrKeychain
	}{&wal, &addr})
	return
}

//DeleteWallet deletes a named wallet associated with the
//API token/coin/chain.
func (api *API) DeleteWallet(c context.Context, name string) (err error) {
	u, err := api.buildURL("/wallets/" + name)
	resp, err := deleteResponse(c, u)
	if err != nil {
		return
	}
	resp.Body.Close()
	return
}
