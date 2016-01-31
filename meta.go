package gobcy

import (
	"appengine"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

//GetMeta gets BlockCypher-stored metadata associated with
//the hash of the given blockchain object. "Kind" describes
//the blockchain object you're querying:
//  "addr" (for an address)
//  "tx" (for a transaction)
//  "block" (for a block)
//If private is false, will retrieve publicly stored metadata.
//If private is true, will retrieve privately stored metadata
//associated with your token.
func (api *API) GetMeta(c appengine.Context, hash string, kind string, private bool) (meta map[string]string, err error) {
	if kind != "addr" && kind != "tx" && kind != "block" {
		err = errors.New(fmt.Sprintf("Func GetMeta: kind an invalid type: '%v'. Needs to be 'addr', 'tx', or 'block'", kind))
		return
	}
	params := map[string]string{"private": strconv.FormatBool(private)}
	u, err := api.buildURLParams("/"+kind+"s/"+hash+"/meta", params)
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//decode JSON into map[string]string
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&meta)
	return
}

//PutMeta puts BlockCypher-stored metadata associated with
//the hash of the given blockchain object. "Kind" describes
//the blockchain object you're querying:
//  "addr" (for an address)
//  "tx" (for a transaction)
//  "block" (for a block)
//If private is false, will set publicly stored metadata.
//If private is true, will set privately stored metadata
//associated with your token.
func (api *API) PutMeta(c appengine.Context, hash string, kind string, private bool, meta map[string]string) (err error) {
	if kind != "addr" && kind != "tx" && kind != "block" {
		err = errors.New(fmt.Sprintf("Func PutMeta: kind an invalid type: '%v'. Needs to be 'addr', 'tx', or 'block'", kind))
		return
	}
	params := map[string]string{"private": strconv.FormatBool(private)}
	u, err := api.buildURLParams("/"+kind+"s/"+hash+"/meta", params)
	if err != nil {
		return
	}
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err = enc.Encode(&meta); err != nil {
		return
	}
	resp, err := putResponse(c, u, &data)
	if err != nil {
		return
	}
	resp.Body.Close()
	return
}

//DeleteMeta deletes ALL PRIVATE BlockCypher-stored metadata
//associated with the hash of the given blockchain object.
//"Kind" describes the blockchain object you're querying:
//  "addr" (for an address)
//  "tx" (for a transaction)
//  "block" (for a block)
//Public metadata cannot be deleted; it is immutable.
func (api *API) DeleteMeta(c appengine.Context, hash string, kind string) (err error) {
	if kind != "addr" && kind != "tx" && kind != "block" {
		err = errors.New(fmt.Sprintf("Func DeleteMeta: kind an invalid type: '%v'. Needs to be 'addr', 'tx', or 'block'", kind))
		return
	}
	u, err := api.buildURL("/" + kind + "s/" + hash + "/meta")
	resp, err := deleteResponse(c, u)
	if err != nil {
		return
	}
	resp.Body.Close()
	return
}
