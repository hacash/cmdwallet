package toolshell

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/cmdwallet/toolshell/ctx"
	"strconv"
)

func getTx(ctx ctx.Context, params []string) {
	if len(params) < 1 {
		fmt.Println("params not enough")
		return
	}
	txhashnofee, err := hex.DecodeString(params[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	newTrs := ctx.GetTxFromRecord(txhashnofee)
	if newTrs == nil {
		fmt.Printf(" tx <%s> not find!", params[0])
		return
	}

	// ok
	ctx.Println("hash: <" + hex.EncodeToString(newTrs.Hash()) + ">, hash_with_fee: <" + hex.EncodeToString(newTrs.HashWithFee()) + ">")

	// Judge whether the signature is completed
	sigok, sigerr := newTrs.VerifyAllNeedSigns()
	nosigntip := ""
	if !sigok || sigerr != nil {
		nosigntip = " [NOT SIGN]"
		fmt.Println("Attention: transaction verify need signs fail!")
		return
	}

	bodybytes, e7 := newTrs.Serialize()
	if e7 != nil {
		fmt.Println("transaction serialize error, " + e7.Error())
		return
	}
	ctx.Println("body length " + strconv.Itoa(len(bodybytes)) + " bytes, hex body is:")
	ctx.Println("-------- TRANSACTION BODY" + nosigntip + " START --------")
	ctx.Println(hex.EncodeToString(bodybytes))
	ctx.Println("-------- TRANSACTION BODY" + nosigntip + " END   --------")

	// record
	ctx.SetTxToRecord(newTrs.Hash(), newTrs)

}
