package gentxs

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/cmdwallet/toolshell/ctx"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"strconv"
)

/*

gentx release_lockbls <lockblsid> <release_amount> <fee_addr> <fee>

// test:
gentx release_lockbls 000000000000000000000000000000000000000000000001 HAC1024:248 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 HAC1:248

*/

// 创建 释放线性锁仓
func GenTxReleaseLockbls(ctx ctx.Context, params []string) {
	if len(params) < 4 {
		fmt.Println("params not enough")
		return
	}

	lockbls_key, e1 := hex.DecodeString(params[0])
	if e1 != nil {
		fmt.Println("param lockbls_id format error.")
		return
	}
	if len(lockbls_key) != 24 {
		fmt.Println("param lockbls_id length error.")
		return
	}

	releaseLockblsAct := &actions.Action_10_LockblsRelease{
		LockblsId: lockbls_key,
	}

	releaseAmount := ctx.IsInvalidAmountString(params[1])
	if releaseAmount == nil {
		fmt.Println("releaseAmount error.")
		return
	}
	releaseLockblsAct.ReleaseAmount = *releaseAmount

	feeAddress := ctx.IsInvalidAccountAddress(params[2])
	if feeAddress == nil {
		fmt.Println("feeAddress error.")
		return
	}
	feeAmount := ctx.IsInvalidAmountString(params[3])
	if feeAmount == nil {
		fmt.Println("feeAmount error.")
		return
	}
	// 创建交易
	newTrs, e5 := transactions.NewEmptyTransaction_2_Simple(*feeAddress)
	newTrs.Timestamp = fields.VarInt5(ctx.UseTimestamp()) // 使用 hold 的时间戳
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}
	newTrs.Fee = *feeAmount // set fee
	// 放入action
	newTrs.AppendAction(releaseLockblsAct)

	// sign
	e6 := newTrs.FillNeedSigns(ctx.GetAllPrivateKeyBytes(), nil)
	if e6 != nil {
		fmt.Println("sign transaction error, " + e6.Error())
		return
	}
	// 检查签名
	sigok, sigerr := newTrs.VerifyNeedSigns(nil)
	if sigerr != nil || !sigok {
		fmt.Println("transaction VerifyNeedSigns fail")
		return
	}

	// 数据化
	bodybytes, e7 := newTrs.Serialize()
	if e7 != nil {
		fmt.Println("transaction serialize error, " + e7.Error())
		return
	}

	// ok
	fmt.Println("transaction create success! ")
	fmt.Println("hash: <" + hex.EncodeToString(newTrs.Hash()) + ">, hash_with_fee: <" + hex.EncodeToString(newTrs.HashWithFee()) + ">")
	fmt.Println("body length " + strconv.Itoa(len(bodybytes)) + " bytes, hex body is:")
	fmt.Println("-------- TRANSACTION BODY START --------")
	fmt.Println(hex.EncodeToString(bodybytes))
	//fmt.Println( hex.EncodeToString( bodybytes2 ) )
	fmt.Println("-------- TRANSACTION BODY END   --------")

	// 记录
	ctx.SetTxToRecord(newTrs.Hash(), newTrs)
}
