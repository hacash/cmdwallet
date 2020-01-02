package gentxs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	base58check "github.com/hacash/core/account"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"github.com/hacash/cmdwallet/toolshell/ctx"
	"strconv"
)




/*


gentx sendcash ${FROM_ADDRESS} ${TO_ADDRESS} ${AMOUNT} ${FEE}


passwd 123456

gentx sendcash 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 1699oAd32emhfShPDFVs5UY8vJNe2u42Fz 1:248 1:244




 */







// 创建一笔交易
func GenTxSimpleTransfer(ctx ctx.Context, params []string) {
	if len(params) < 4 {
		fmt.Println("params not enough")
		return
	}
	from := params[0]
	to := params[1]
	finamt := params[2]
	finfee := params[3]
	if ctx.NotLoadedYetAccountAddress(from) {
		return
	}
	toAddr := ctx.IsInvalidAccountAddress(to)
	if toAddr == nil {
		return
	}
	amt, e1 := fields.NewAmountFromFinString(finamt)
	if e1 != nil {
		fmt.Println("amount format error or over range, the right example is 'HCX1:248' for one coin")
		return
	}
	fee, e2 := fields.NewAmountFromFinString(finfee)
	if e2 != nil {
		fmt.Println("fee format error or over range")
		return
	}
	masterAddr, e3 := base58check.Base58CheckDecode(from)
	if e3 != nil {
		fmt.Println("from address format error")
		return
	}
	newTrs, e5 := transactions.NewEmptyTransaction_2_Simple(fields.Address(masterAddr))
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}
	newTrs.Timestamp = fields.VarInt5(ctx.UseTimestamp()) // 使用 hold 的时间戳
	newTrs.Fee = *fee // set fee
	tranact := actions.NewAction_1_SimpleTransfer(*toAddr, amt)
	e5 = newTrs.AppendAction(tranact)
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}
	// sign
	e6 := newTrs.FillNeedSigns(ctx.GetAllPrivateKeyBytes(), nil)
	if e6 != nil {
		fmt.Println("sign transaction error, " + e6.Error())
		return
	}
	bodybytes, e7 := newTrs.Serialize()
	if e7 != nil {
		fmt.Println("transaction serialize error, " + e7.Error())
		return
	}

	var trxnew, _, _ = transactions.ParseTransaction(bodybytes, 0)
	bodybytes2, _ := trxnew.Serialize()
	if 0 != bytes.Compare(bodybytes, bodybytes2) {
		fmt.Println("transaction serialize error")
		return
	}

	sigok, sigerr := trxnew.VerifyNeedSigns(nil)
	if sigerr != nil {
		fmt.Println("transaction VerifyNeedSigns error")
		return
	}
	if !sigok {
		fmt.Println("transaction VerifyNeedSigns fail")
		return
	}

	// ok
	ctx.Println("transaction create success! ")
	ctx.Println("hash: <" + hex.EncodeToString(newTrs.Hash()) + ">, hash_with_fee: <" + hex.EncodeToString(newTrs.HashWithFee()) + ">")
	ctx.Println("body length " + strconv.Itoa(len(bodybytes)) + " bytes, hex body is:")
	ctx.Println("-------- TRANSACTION BODY START --------")
	ctx.Println(hex.EncodeToString(bodybytes))
	//fmt.Println( hex.EncodeToString( bodybytes2 ) )
	ctx.Println("-------- TRANSACTION BODY END   --------")

	// 记录
	ctx.SetTxToRecord(newTrs.Hash(), newTrs)

}
