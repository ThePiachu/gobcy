//Package gobcy implements a wrapper for the http://www.blockcypher.com API.
//You can use it to interact with addresses, transactions, and blocks from
//various blockchains, including Bitcoin's main and test3 chains,
//and the BlockCypher test chain.
//
//Please note: we assume you use are using a 64-bit architecture for deployment,
//which automatically makes `int` types 64-bit. Without 64-bit ints, some values
//might overflow on certain calls, depending on the blockchain you are querying.
//If you are using a 32-bit system, you can change all `int` types to `int64` to
//explicitly work around this issue.
package gobcy

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"io"
	"net/http"
	"net/url"

	"github.com/ThePiachu/Go/Log"
)

const baseURL = "https://api.blockcypher.com/v1/"

//API stores your BlockCypher Token, and the coin/chain
//you're querying. Coins can be "btc","bcy","ltc", and "doge".
//Chains can be "main", "test3", or "test", depending on the Coin.
//Check http://dev.blockcypher.com/ for more information.
//All your credentials are stored within an API struct, as are
//many of the API methods.
//You can allocate an API struct like so:
//	bc = gobcy.API{"your-api-token","btc","main"}
//Then query as you like:
//	chain = bc.GetChain()
type API struct {
	Token, Coin, Chain string
}

//getResponse is a boilerplate for HTTP GET responses.
func getResponse(c context.Context, target *url.URL) (resp *http.Response, err error) {
	tr := urlfetch.Transport{Context: c}
	client := http.Client{Transport: &tr}
	Log.Debugf(c, "Getting %v", target.String())
	resp, err = client.Get(target.String())
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		msg := make(map[string]string)
		dec := json.NewDecoder(resp.Body)
		dec.Decode(&msg)
		resp.Body.Close()
		err = errors.New(resp.Status + ", Message: " + msg["error"])
	}
	return
}

//postResponse is a boilerplate for HTTP POST responses.
func postResponse(c context.Context, target *url.URL, data io.Reader) (resp *http.Response, err error) {
	tr := urlfetch.Transport{Context: c}
	client := http.Client{Transport: &tr}
	Log.Debugf(c, "Posting %v", target.String())
	resp, err = client.Post(target.String(), "application/json", data)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		msg := make(map[string]string)
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&msg)
		if err != nil {
			Log.Debugf(c, "postResponse - %v", err)
		}
		resp.Body.Close()
		err = errors.New(resp.Status + ", Message: " + msg["error"])
	}
	return
}

//putResponse is a boilerplate for HTTP PUT responses.
func putResponse(c context.Context, target *url.URL, data io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("PUT", target.String(), data)
	if err != nil {
		return
	}
	tr := urlfetch.Transport{Context: c}
	client := http.Client{Transport: &tr}
	Log.Debugf(c, "Putting %v", target.String())
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		msg := make(map[string]string)
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&msg)
		if err != nil {
			Log.Debugf(c, "putResponse - %v", err)
		}
		resp.Body.Close()
		err = errors.New(resp.Status + ", Message: " + msg["error"])
	}
	return
}

//deleteResponse is a boilerplate for HTTP DELETE responses.
func deleteResponse(c context.Context, target *url.URL) (resp *http.Response, err error) {
	req, err := http.NewRequest("DELETE", target.String(), nil)
	if err != nil {
		return
	}
	tr := urlfetch.Transport{Context: c}
	client := http.Client{Transport: &tr}
	Log.Debugf(c, "Deleting %v", target.String())
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		msg := make(map[string]string)
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&msg)
		if err != nil {
			Log.Debugf(c, "deleteResponse - %v", err)
		}
		resp.Body.Close()
		err = errors.New(resp.Status + ", Message: " + msg["error"])
	}
	return
}

//constructs BlockCypher URLs for requests
func (api *API) buildURL(u string) (target *url.URL, err error) {
	target, err = url.Parse(baseURL + api.Coin + "/" + api.Chain + u)
	if err != nil {
		return
	}
	//add token to url, if present
	if api.Token != "" {
		values := target.Query()
		values.Set("token", api.Token)
		target.RawQuery = values.Encode()
	}
	return
}

//constructs BlockCypher URLs with parameters for requests
func (api *API) buildURLParams(u string, params map[string]string) (target *url.URL, err error) {
	target, err = url.Parse(baseURL + api.Coin + "/" + api.Chain + u)
	if err != nil {
		return
	}
	values := target.Query()
	//Set parameters
	for k, v := range params {
		values.Set(k, v)
	}
	//add token to url, if present
	if api.Token != "" {
		values.Set("token", api.Token)
	}
	target.RawQuery = values.Encode()
	return
}
