package gobcy

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/context"
)

//CreateHook creates a new WebHook associated
//with your API.Token, and returns a WebHook
//with a BlockCypher-assigned id.
func (api *API) CreateHook(c context.Context, hook Hook) (result Hook, err error) {
	u, err := api.buildURL("/hooks")
	if err != nil {
		return
	}
	var data bytes.Buffer
	enc := json.NewEncoder(&data)
	if err = enc.Encode(&hook); err != nil {
		return
	}
	resp, err := postResponse(c, u, &data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	return
}

//ListHooks returns a slice of WebHooks
//associated with your API.Token.
func (api *API) ListHooks(c context.Context) (hooks []Hook, err error) {
	u, err := api.buildURL("/hooks")
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//decode JSON into hooks
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&hooks)
	return
}

//GetHook returns a WebHook by its id.
func (api *API) GetHook(c context.Context, id string) (hook Hook, err error) {
	u, err := api.buildURL("/hooks/" + id)
	resp, err := getResponse(c, u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	//decode JSON into hooks
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&hook)
	return
}

//DeleteHook deletes a WebHook notification
//from BlockCypher's database, based on its id.
func (api *API) DeleteHook(c context.Context, id string) (err error) {
	u, err := api.buildURL("/hooks/" + id)
	resp, err := deleteResponse(c, u)
	if err != nil {
		return
	}
	resp.Body.Close()
	return
}
