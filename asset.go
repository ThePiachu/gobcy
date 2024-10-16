package gobcy

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/context"
)

//GenAssetKeychain generates a public/private key pair, alongside
//an associated OAPAddress for use in the Asset API.
func (api *API) GenAssetKeychain(c context.Context) (pair AddrKeychain, err error) {
	u, err := api.buildURL("/oap/addrs")
	resp, err := postResponse(c, u, nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&pair)
	return
}

//IssueAsset issues new assets onto an Open Asset Address,
//using a private key associated with a funded address
//on the underlying blockchain.
func (api *API) IssueAsset(c context.Context, issue OAPIssue) (tx OAPTX, err error) {
	u, err := api.buildURL("/oap/issue")
	if err != nil {
		return
	}
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err = enc.Encode(&issue); err != nil {
		return
	}
	resp, err := postResponse(c, u, &data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tx)
	return
}

//TransferAsset transfers previously issued assets onto a new
//Open Asset Address, based on the assetid and OAPIssue.
func (api *API) TransferAsset(c context.Context, issue OAPIssue, assetID string) (tx OAPTX, err error) {
	u, err := api.buildURL("/oap/" + assetID + "/transfer")
	if err != nil {
		return
	}
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err = enc.Encode(&issue); err != nil {
		return
	}
	resp, err := postResponse(c, u, &data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tx)
	return
}

//ListAssetTXs lists the transaction hashes associated
//with the given assetID.
func (api *API) ListAssetTXs(c context.Context, assetID string) (txs []string, err error) {
	u, err := api.buildURL("/oap/" + assetID + "/txs")
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&txs)
	return
}

//GetAssetTX returns a OAPTX associated with the given
//assetID and transaction hash.
func (api *API) GetAssetTX(c context.Context, assetID, hash string) (tx OAPTX, err error) {
	u, err := api.buildURL("/oap/" + assetID + "/txs/" + hash)
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&tx)
	return
}

//GetAssetAddr returns an Addr associated with the given
//assetID and oapAddr. Note that while it returns an Address,
//anything that would have represented "satoshis" now represents
//"amount of asset."
func (api *API) GetAssetAddr(c context.Context, assetID, oapAddr string) (addr Addr, err error) {
	u, err := api.buildURL("/oap/" + assetID + "/addrs/" + oapAddr)
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&addr)
	return
}
