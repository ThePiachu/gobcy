package gobcy

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"strconv"
)

//CreateHDWallet creates a public-address watching HDWallet
//associated with this token/coin/chain, usable anywhere
//in the API where an Address might be used (just use
//the wallet name instead). For example, with checking
//a wallet name balance:
//  addr, err := api.GetAddrBal("your-hd-wallet-name")
func (api *API) CreateHDWallet(c context.Context, req HDWallet) (wal HDWallet, err error) {
	u, err := api.buildURL("/wallets/hd")
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

//ListHDWallets lists all known HDWallets associated with
//this token/coin/chain.
func (api *API) ListHDWallets(c context.Context) (names []string, err error) {
	u, err := api.buildURL("/wallets/hd")
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	jsonResp := new(struct {
		List []string `json:"hd_wallet_names"`
	})
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(jsonResp)
	names = jsonResp.List
	return
}

//GetHDWallet gets a HDWallet based on its name
//and the associated API token/coin/chain.
func (api *API) GetHDWallet(c context.Context, name string) (wal HDWallet, err error) {
	u, err := api.buildURL("/wallets/hd/" + name)
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

//GetAddrHDWallet returns addresses associated with
//a named HDWallet, associated with the API token/coin/chain.
//Offers 4 parameters for customization:
//  "used," if true will return only used addresses
//  "unused," if true will return only unused addresses
//  "zero", if true will return only zero balance addresses
//  "nonzero", if true will return only nonzero balance addresses
//"used" and "unused" cannot be true at the same time; the SDK will throw an error.
//"zero" and "nonzero" cannot be true at the same time; the SDK will throw an error.
func (api *API) GetAddrHDWallet(c context.Context, name string, used bool, unused bool, zero bool, nonzero bool) (addrs HDWallet, err error) {
	params := make(map[string]string)
	if used && unused {
		err = errors.New("GetAddrHDWallet: Unused and used cannot be the same")
		return
	}
	if zero && nonzero {
		err = errors.New("GetAddrHDWallet: Zero and nonzero cannot be the same")
		return
	}
	if used != unused {
		params["used"] = strconv.FormatBool(used)
	}
	if zero != nonzero {
		params["zerobalance"] = strconv.FormatBool(zero)
	}
	u, err := api.buildURLParams("/wallets/hd/"+name+"/addresses", params)
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//decode JSON into result
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&addrs)
	return
}

//DeriveAddrHDWallet derives a new address within the named Wallet,
//associated with the API token/coin/chain. It will only return a partial
//HDWallet, ONLY containing the new address derived. Has two parameters:
// "count," number of addresses to derive. Default is one.
// "isSub," true if deriving an address on a subchain.
//    if false, address will be generated on the first chain in the HDWallet.
// "subchainIndex," Derives address(es) on this specific subchain. Only used
//    if isSubchain is true.
func (api *API) DeriveAddrHDWallet(c context.Context, name string, count int, isSub bool, subchainIndex int) (wal HDWallet, err error) {
	params := map[string]string{"count": strconv.Itoa(count)}
	if isSub {
		params["subchain_index"] = strconv.Itoa(subchainIndex)
	}
	u, err := api.buildURLParams("/wallets/hd/"+name+"/addresses/derive", params)
	resp, err := postResponse(c, u, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//decode JSON into result
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&wal)
	return
}

//DeleteHDWallet deletes a named HDWallet associated with the
//API token/coin/chain.
func (api *API) DeleteHDWallet(c context.Context, name string) (err error) {
	u, err := api.buildURL("/wallets/hd/" + name)
	resp, err := deleteResponse(c, u)
	if err != nil {
		return
	}
	resp.Body.Close()
	return
}
