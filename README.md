# GAE gobcy

A Google App Engine version of the Go wrapper for the [BlockCypher](http://www.blockcypher.com/) API. Targeting support for Bitcoin (main and testnet3) and BlockCypher's internal testnet, but others (Litecoin, Dogecoin) should work too.

## Configuration

Import the package like so:

```go
import "github.com/blockcypher/gobcy"
```

Then initiate an API struct with your credentials:

```go
//explicitly
bc := gobcy.API
bc.Token = "your-api-token-here"
bc.Coin = "btc" //options: "btc","bcy","ltc","doge"
bc.Chain = "main" //depending on coin: "main","test3","test"

//using a struct literal
bc := gobcy.API{"your-api-token-here","btc","main"}

//query away
fmt.Println(bc.GetChain())
```

## Usage

Check the "types.go" file for information on the return types. Almost all API calls are supported, with a few dropped to reduce complexity.

For more information on the API, check out [BlockCypher's documentation](http://dev.blockcypher.com/). We've heavily commented the code following Golang convention, so you might also find [the GoDoc quite useful.](http://godoc.org/github.com/blockcypher/gobcy) The `gobcy_test.go` file also shows most of the API calls in action.

## A Warning for 32-Bit Systems

We assume you use are using a 64-bit architecture for deployment, which automatically makes `int` types 64-bit, the default behavior since [Go 1.1](https://tip.golang.org/doc/go1.1#int). Without 64-bit ints, some values might overflow on certain calls, depending on the blockchain you are querying. If you are using a 32-bit system, you can change all `int` types to `int64` to explicitly work around this issue.

## Testing

The aforementioned `gobcy_test.go` file contains a number of tests to ensure the wrapper is functioning properly. If you run it yourself, you'll have to insert a valid API token; you may also want to generate a new token, as the test POSTs and DELETEs WebHooks and Payment Forwarding requests.

