//Tests for the BlockCypher Go SDK. Test functions
//try to mirror file names where possible.
package gobcy

import (
	"appengine/aetest"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

var keys1, keys2 AddrKeychain
var txhash1, txhash2 string
var bcy API

func TestMain(m *testing.M) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	//Set Coin/Chain to BlockCypher testnet
	bcy.Coin = "bcy"
	bcy.Chain = "test"
	//Set Token
	bcy.Token = "test-token"
	//Create/fund the test addresses
	keys1, err = bcy.GenAddrKeychain(c)
	keys2, err = bcy.GenAddrKeychain(c)
	if err != nil {
		log.Fatal("Error generating test addresses: ", err)
	}
	txhash1, err = bcy.Faucet(c, keys1, 1e5)
	txhash2, err = bcy.Faucet(c, keys2, 2e5)
	if err != nil {
		log.Fatal("Error funding test addresses: ", err)
	}
	os.Exit(m.Run())
}

//TestsGetTXConf runs first, to test
//Confidence factor
func TestGetTXConf(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	conf, err := bcy.GetTXConf(c, txhash2)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", conf)
	return
}

func TestBlockchain(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	ch, err := bcy.GetChain(c)
	if err != nil {
		t.Error("GetChain error encountered: ", err)
	}
	t.Logf("%+v\n", ch)
	_, err = bcy.GetBlock(c, 187621, "")
	if err != nil {
		t.Error("GetBlock via height error encountered: ", err)
	}
	bl, err := bcy.GetBlock(c, 0, "0000ffeb0031885f2292475eac7f9c6f7bf5057e3b0017a09cd1994e71b431a4")
	if err != nil {
		t.Error("GetBlock via hash error encountered: ", err)
	}
	t.Logf("%+v\n", bl)
	_, err = bcy.GetBlock(c, 187621, "0000ffeb0031885f2292475eac7f9c6f7bf5057e3b0017a09cd1994e71b431a4")
	if err == nil {
		t.Error("Expected error when querying both height and hash in GetBlock, did not receive one")
	}
	err = nil
	bl, err = bcy.GetBlockPage(c, 0, "0000cb69e3c85ec1a4a17d8a66634c1cf136acc9dca9a5a71664a593f92bc46e", 0, 1)
	if err != nil {
		t.Error("GetBlockPage error encountered: ", err)
	}
	t.Logf("%+v\n", bl)
	bl2, err := bcy.GetBlockNextTXs(c, bl)
	if err != nil {
		t.Error("GetBlockNextTXs error encountered: ", err)
	}
	t.Logf("%+v\n", bl2)
	return
}

func TestAddress(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	addr, err := bcy.GetAddrBal(c, keys1.Address)
	if err != nil {
		t.Error("GetAddrBal error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	addr, err = bcy.GetAddr(c, keys1.Address)
	if err != nil {
		t.Error("GetAddr error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	addr, err = bcy.GetAddrFull(c, keys2.Address)
	if err != nil {
		t.Error("GetAddrFull error encountered: ", err)
	}
	t.Logf("%+v\n", addr)
	return
}

func TestGenAddrMultisig(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	pubkeys := []string{
		"02c716d071a76cbf0d29c29cacfec76e0ef8116b37389fb7a3e76d6d32cf59f4d3",
		"033ef4d5165637d99b673bcdbb7ead359cee6afd7aaf78d3da9d2392ee4102c8ea",
		"022b8934cc41e76cb4286b9f3ed57e2d27798395b04dd23711981a77dc216df8ca",
	}
	response, err := bcy.GenAddrMultisig(c, AddrKeychain{PubKeys: pubkeys, ScriptType: "multisig-2-of-3"})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	if response.Address != "De2gwq9GvNgvKgHCYRMKnPqss3pzWGSHiH" {
		t.Error("Response does not match expected address")
	}
	t.Logf("%+v\n", response)
	return
}

func TestWallet(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	wal, err := bcy.CreateWallet(c, Wallet{Name: "testwallet",
		Addresses: []string{keys1.Address}})
	if err != nil {
		t.Error("CreateWallet error encountered: ", err)
	}
	t.Logf("%+v\n", wal)
	list, err := bcy.ListWallets(c)
	if err != nil {
		t.Error("ListWallet error encountered: ", err)
	}
	t.Logf("%+v\n", list)
	wal, err = bcy.AddAddrWallet(c, "testwallet", []string{keys2.Address}, false)
	if err != nil {
		t.Error("AddAddrWallet error encountered: ", err)
	}
	t.Logf("%+v\n", wal)
	err = bcy.DeleteAddrWallet(c, "testwallet", []string{keys1.Address})
	if err != nil {
		t.Error("DeleteAddrWallet error encountered ", err)
	}
	addrs, err := bcy.GetAddrWallet(c, "testwallet", false, false, false, false)
	if err != nil {
		t.Error("GetAddrWallet error encountered: ", err)
	}
	if addrs[0] != keys2.Address {
		t.Error("GetAddrWallet response does not match expected addresses")
	}
	wal, newAddrKeys, err := bcy.GenAddrWallet(c, "testwallet")
	if err != nil {
		t.Error("GenAddrWallet error encountered: ", err)
	}
	t.Logf("%+v\n%+v\n", wal, newAddrKeys)
	err = bcy.DeleteWallet(c, "testwallet")
	if err != nil {
		t.Error("DeleteWallet error encountered: ", err)
	}
	return
}

func TestHDWallet(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	wal, err := bcy.CreateHDWallet(c, HDWallet{Name: "testhdwallet",
		ExtPubKey: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8"})
	if err != nil {
		t.Error("CreateHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", wal)
	//Will reenable once ListHDWallet fixed
	/* list, err := bcy.ListHDWallets()
	if err != nil {
		t.Error("ListHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", list)*/
	addrs, err := bcy.GetAddrHDWallet(c, "testhdwallet", false, false, false, false)
	if err != nil {
		t.Error("GetAddrHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", addrs)
	newhd, err := bcy.DeriveAddrHDWallet(c, "testhdwallet", 1, false, 0)
	if err != nil {
		t.Error("DeriveAddrHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", newhd)
	wal, err = bcy.GetHDWallet(c, "testhdwallet")
	if err != nil {
		t.Error("GetHDWallet error encountered: ", err)
	}
	t.Logf("%+v\n", wal)
	err = bcy.DeleteHDWallet(c, "testhdwallet")
	if err != nil {
		t.Error("DeleteHDWallet error encountered: ", err)
	}
	return
}

func TestTX(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	txs, err := bcy.GetUnTX(c)
	if err != nil {
		t.Error("GetUnTX error encountered: ", err)
	}
	t.Logf("%+v\n", txs)
	tx, err := bcy.GetTX(c, txhash1)
	if err != nil {
		t.Error("GetTX error encountered: ", err)
	}
	t.Logf("%+v\n", tx)
	//Create New TXSkeleton
	temp := TempNewTX(c, keys2.Address, keys1.Address, 45000)
	skel, err := bcy.NewTX(c, temp, true)
	if err != nil {
		t.Error("NewTX error encountered: ", err)
	}
	t.Logf("%+v\n", skel)
	//Sign TXSkeleton
	err = skel.Sign(c, []string{keys2.Private})
	if err != nil {
		t.Error("*TXSkel.Sign error encountered: ", err)
	}
	//Send TXSkeleton
	skel, err = bcy.SendTX(c, skel)
	if err != nil {
		t.Error("SendTX error encountered: ", err)
	}
	t.Logf("%+v\n", skel)
	return
}

func TestMicro(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	result, err := bcy.SendMicro(c, MicroTX{Priv: keys2.Private, ToAddr: keys1.Address, Value: 25000})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("%+v\n", result)
	//Test public key signing method
	micpub, err := bcy.SendMicro(c, MicroTX{Pubkey: keys2.Public, ToAddr: keys1.Address, Value: 15000})
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("First MicroTX Pubkey call: %+v\n", micpub)
	//Sign resultpub
	err = micpub.Sign(keys2.Private)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("Signed MicroTX: %+v\n", micpub)
	//Send signed resultpub
	resultpub, err := bcy.SendMicro(c, micpub)
	if err != nil {
		t.Error("Error encountered: ", err)
	}
	t.Logf("Final MicroTX Pubkey call: %+v\n", resultpub)
	return
}

func TestHook(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	hook, err := bcy.CreateHook(c, Hook{Event: "new-block", URL: "https://my.domain.com/api/callbacks/doublespend?secret=justbetweenus"})
	if err != nil {
		t.Error("PostHook error encountered: ", err)
	}
	t.Logf("%+v\n", hook)
	if err = bcy.DeleteHook(c, hook.ID); err != nil {
		t.Error("DeleteHook error encountered: ", err)
	}
	hooks, err := bcy.ListHooks(c)
	if err != nil {
		t.Error("ListHooks error encountered: ", err)
	}
	//Should be empty
	t.Logf("%+v\n", hooks)
	return
}

func TestPayFwd(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	pay, err := bcy.CreatePayFwd(c, PayFwd{Destination: keys1.Address})
	if err != nil {
		t.Error("CreatePayFwd error encountered: ", err)
	}
	t.Logf("%+v\n", pay)
	pay, err = bcy.GetPayFwd(c, pay.ID)
	if err != nil {
		t.Error("GetPayFwd error encountered: ", err)
	}
	t.Logf("%+v\n", pay)
	if err = bcy.DeletePayFwd(c, pay.ID); err != nil {
		t.Error("DeletePayFwd error encountered: ", err)
	}
	pays, err := bcy.ListPayFwds(c)
	if err != nil {
		t.Error("ListPayFwds error encountered: ", err)
	}
	//Should be empty
	t.Logf("%+v\n", pays)
	return
}

func TestMeta(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	err = bcy.PutMeta(c, keys1.Address, "addr", true, map[string]string{"key": "value"})
	if err != nil {
		t.Error("PutMeta error encountered: ", err)
	}
	data, err := bcy.GetMeta(c, keys1.Address, "addr", true)
	if err != nil {
		t.Error("GetMeta error encountered: ", err)
	}
	if data["key"] != "value" {
		t.Error("GetMeta error encountered, expected data['key']='value', got: ", data["key"])
	}
	err = bcy.DeleteMeta(c, keys1.Address, "addr")
	if err != nil {
		t.Error("DeleteMeta error encountered: ", err)
	}
}

func TestAsset(t *testing.T) {
	c, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	oap1, err := bcy.GenAssetKeychain(c)
	oap2, err := bcy.GenAssetKeychain(c)
	funder, err := bcy.GenAddrKeychain(c)
	_, err = bcy.Faucet(c, funder, 1e6)
	_, err = bcy.Faucet(c, oap1, 1e6)
	if err != nil {
		t.Error("GenAsset/AddrKeychain or Faucet error encountered: ", err)
	}
	tx1, err := bcy.IssueAsset(c, OAPIssue{funder.Private, oap1.OAPAddress, 9000, ""})
	if err != nil {
		t.Error("IssueAsset error encountered: ", err)
	}
	//wait until tx1 is confirmed
	fmt.Printf("Waiting for asset issuance to confirm.")
	for {
		conf, err := bcy.GetTXConf(c, tx1.Hash)
		if err != nil {
			t.Logf("Error polling for Issue Asset tx confirmation: ", err)
			break
		}
		if conf.Confidence == 1 {
			fmt.Printf("\n")
			break
		}
		fmt.Printf(".")
		time.Sleep(2 * time.Second)
	}
	tx2, err := bcy.TransferAsset(c, OAPIssue{oap1.Private, oap2.OAPAddress, 8999, ""}, tx1.AssetID)
	if err != nil {
		t.Error("TransferAsset error encountered: ", err)
		t.Errorf("Returned OAPTX1:%+v\n", tx1)
		t.Errorf("Returned OAPTX2:%+v\n", tx2)
	}
	txs, err := bcy.ListAssetTXs(c, tx1.AssetID)
	if err != nil {
		t.Error("ListAssetTXs error encountered: ", err)
		t.Errorf("Returned TXs:%v\n", txs)
	}
	checktx, err := bcy.GetAssetTX(c, tx1.AssetID, tx1.Hash)
	if err != nil {
		t.Error("GetAssetTX error encountered: ", err)
		t.Errorf("Original OAPTX from first issue: %+v\n", tx1)
		t.Errorf("Returned OAPTX from GetAssetTX endpoint: %+v\n", checktx)
	}
	oapaddr, err := bcy.GetAssetAddr(c, tx1.AssetID, oap1.OAPAddress)
	if err != nil {
		t.Error("GetAssetAddr error encountered: ", err)
		t.Errorf("Original OAPTX from first issue: %+v\n", tx1)
		t.Errorf("Returned Addr from GetAssetAddr endpoint: %+v\n", oapaddr)
	}
	t.Logf("Returned Addr from GetAssetAddr endpoint: %+v\n", oapaddr)
}
