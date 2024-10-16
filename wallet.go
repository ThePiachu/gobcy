package gobcy

import (
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

//CreateWallet creates a public-address watching wallet
//associated with this token/coin/chain, usable anywhere
//in the API where an Address might be used (just use
//the wallet name instead). For example, with checking
//a wallet name balance:
//  addr, err := api.GetAddrBal("your-wallet-name", nil)
func (api *API) CreateWallet(c context.Context, req Wallet) (wal Wallet, err error) {
	u, err := api.buildURL("/wallets", nil)
	if err != nil {
		return
	}
	err = postResponse(c, u, &req, &wal)
	return
}

//ListWallets lists all known Wallets associated with
//this token/coin/chain.
func (api *API) ListWallets(c context.Context) (names []string, err error) {
	u, err := api.buildURL("/wallets", nil)
	if err != nil {
		return
	}
	jsonResp := new(struct {
		List []string `json:"wallet_names"`
	})
	err = getResponse(c, u, jsonResp)
	names = jsonResp.List
	return
}

//GetWallet gets a Wallet based on its name, the associated
//API token/coin/chain, and whether it's an HD wallet or
//not.
func (api *API) GetWallet(c context.Context, name string) (wal Wallet, err error) {
	u, err := api.buildURL("/wallets/"+name, nil)
	if err != nil {
		return
	}
	err = getResponse(c, u, &wal)
	return
}

//AddAddrWallet adds a slice of addresses to a named Wallet,
//associated with the API token/coin/chain. In addition to your
//list of addresses to add, takes one additional parameter:
//  "omitAddr," if true will omit wallet addresses in your
//  response. Useful to speed up the API call for larger wallets.
func (api *API) AddAddrWallet(c context.Context, name string, addrs []string, omitAddr bool) (wal Wallet, err error) {
	params := map[string]string{"omitWalletAddresses": strconv.FormatBool(omitAddr)}
	u, err := api.buildURL("/wallets/"+name+"/addresses", params)
	if err != nil {
		return
	}
	err = postResponse(c, u, &Wallet{Addresses: addrs}, &wal)
	return
}

//GetAddrWallet returns a slice of addresses associated with
//a named Wallet, associated with the API token/coin/chain.
//Takes an optionally-nil URL parameter map.
func (api *API) GetAddrWallet(c context.Context, name string, params map[string]string) (addrs []string, err error) {
	u, err := api.buildURL("/wallets/"+name+"/addresses", params)
	if err != nil {
		return
	}
	var wal Wallet
	err = getResponse(c, u, &wal)
	addrs = wal.Addresses
	return
}

//DeleteAddrWallet deletes a slice of addresses associated with
//a named Wallet, associated with the API token/coin/chain.
func (api *API) DeleteAddrWallet(c context.Context, name string, addrs []string) (err error) {
	u, err := api.buildURL("/wallets/"+name+"/addresses",
		map[string]string{"address": strings.Join(addrs, ";")})
	if err != nil {
		return
	}
	err = deleteResponse(c, u)
	return
}

//GenAddrWallet generates a new address within the named Wallet,
//associated with the API token/coin/chain. Also returns the
//private/WIF/public key of address via an Address Keychain.
func (api *API) GenAddrWallet(c context.Context, name string) (wal Wallet, addr AddrKeychain, err error) {
	u, err := api.buildURL("/wallets/"+name+"/addresses/generate", nil)
	if err != nil {
		return
	}
	err = postResponse(c, u, nil, &struct {
		*Wallet
		*AddrKeychain
	}{&wal, &addr})
	return
}

//DeleteWallet deletes a named wallet associated with the
//API token/coin/chain.
func (api *API) DeleteWallet(c context.Context, name string) (err error) {
	u, err := api.buildURL("/wallets/"+name, nil)
	if err != nil {
		return
	}
	err = deleteResponse(c, u)
	return
}
