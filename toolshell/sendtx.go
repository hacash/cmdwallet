package toolshell

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/cmdwallet/toolshell/ctx"
	"io/ioutil"
	"net/http"
)

/*
sendtx <txhx> 127.0.0.1:3338
*/

// Send a transaction to the miner
func sendTxToMiner(ctx ctx.Context, params []string) {
	if len(params) < 2 {
		fmt.Println("params not enough")
		return
	}

	txhash, e0 := hex.DecodeString(params[0])
	if e0 != nil {
		fmt.Println("tx hash format error")
		return
	}

	minerAddress := params[1]

	tx := ctx.GetTxFromRecord(txhash)
	if tx == nil {
		return
	}

	sigok, e2 := tx.VerifyAllNeedSigns()
	if e2 != nil || !sigok {
		fmt.Println("Tx sign error")
		fmt.Println(e2)
		return
	}

	// Post send
	body := new(bytes.Buffer)
	body.Write([]byte{0, 0, 0, 1}) // opcode
	txbytes, e9 := tx.Serialize()
	if e9 != nil {
		fmt.Println("tx serialize error:")
		fmt.Println(e9)
		return
	}

	// Generate transaction
	// transactions.ParseTransaction(txbytes, 0)
	body.Write(txbytes)
	req, e3 := http.NewRequest("POST", "http://"+minerAddress+"/operate", body)
	if e3 != nil {
		fmt.Println("POST NewRequest error:")
		fmt.Println(e3)
		return
	}

	client := &http.Client{}
	resp, e4 := client.Do(req)
	if e4 != nil {
		fmt.Println("POST client.Do(req) error:")
		fmt.Println(e4)
		return
	}

	defer resp.Body.Close()
	// ok
	fmt.Println("add tx to " + minerAddress + ", the response is:")
	resbody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resbody))
}
