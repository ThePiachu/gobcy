package gobcy

import (
	"golang.org/x/net/context"
)

//CreateHook creates a new WebHook associated
//with your API.Token, and returns a WebHook
//with a BlockCypher-assigned id.
func (api *API) CreateHook(c context.Context, hook Hook) (result Hook, err error) {
	u, err := api.buildURL("/hooks", nil)
	if err != nil {
		return
	}
	err = postResponse(c, u, &hook, &result)
	return
}

//ListHooks returns a slice of WebHooks
//associated with your API.Token.
func (api *API) ListHooks(c context.Context) (hooks []Hook, err error) {
	u, err := api.buildURL("/hooks", nil)
	if err != nil {
		return
	}
	err = getResponse(c, u, &hooks)
	return
}

//GetHook returns a WebHook by its id.
func (api *API) GetHook(c context.Context, id string) (hook Hook, err error) {
	u, err := api.buildURL("/hooks/"+id, nil)
	if err != nil {
		return
	}
	err = getResponse(c, u, &hook)
	return
}

//DeleteHook deletes a WebHook notification
//from BlockCypher's database, based on its id.
func (api *API) DeleteHook(c context.Context, id string) (err error) {
	u, err := api.buildURL("/hooks/"+id, nil)
	if err != nil {
		return
	}
	err = deleteResponse(c, u)
	return
}
