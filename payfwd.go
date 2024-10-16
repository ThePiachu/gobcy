package gobcy

import (
	"strconv"

	"golang.org/x/net/context"
)

//CreatePayFwd creates a new PayFwd forwarding
//request associated with your API.Token, and
//returns a PayFwd with a BlockCypher-assigned id.
func (api *API) CreatePayFwd(c context.Context, payment PayFwd) (result PayFwd, err error) {
	u, err := api.buildURL("/payments", nil)
	if err != nil {
		return
	}
	err = postResponse(c, u, &payment, &result)
	return
}

//ListPayFwds returns a PayFwds slice
//associated with your API.Token.
func (api *API) ListPayFwds(c context.Context) (payments []PayFwd, err error) {
	u, err := api.buildURL("/payments", nil)
	if err != nil {
		return
	}
	err = getResponse(c, u, &payments)
	return
}

//ListPayFwdsPage returns a PayFwds slice
//associated with your API.Token, starting at the start index.
//Useful for paging past the 200 payment forward limit.
func (api *API) ListPayFwdsPage(c context.Context, start int) (payments []PayFwd, err error) {
	params := map[string]string{"start": strconv.Itoa(start)}
	u, err := api.buildURL("/payments", params)
	if err != nil {
		return
	}
	err = getResponse(c, u, &payments)
	return
}

//GetPayFwd returns a PayFwd based on its id.
func (api *API) GetPayFwd(c context.Context, id string) (payment PayFwd, err error) {
	u, err := api.buildURL("/payments/"+id, nil)
	if err != nil {
		return
	}
	err = getResponse(c, u, &payment)
	return
}

//DeletePayFwd deletes a PayFwd request from
//BlockCypher's database, based on its id.
func (api *API) DeletePayFwd(c context.Context, id string) (err error) {
	u, err := api.buildURL("/payments/"+id, nil)
	if err != nil {
		return
	}
	err = deleteResponse(c, u)
	return
}
