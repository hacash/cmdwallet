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

gentx gentx ${FROM_ADDRESS} ${TO_ADDRESS} ${SATOSHI_AMOUNT} ${FEE}

*/

// satoshi 普通转账
func GenTxSimpleTransferSatoshi(ctx ctx.Context, params []string) {
	if len(params) < 4 {
		fmt.Println("params not enough")
		return
	}

	feeAddress := ctx.IsInvalidAccountAddress(params[0])
	if feeAddress == nil {
		return
	}

	targetAddress := ctx.IsInvalidAccountAddress(params[1])
	if targetAddress == nil {
		return
	}

	satoshiAmount, se1 := strconv.ParseInt(params[2], 10, 0)
	if se1 != nil {
		fmt.Println("satoshi number error: " + params[2])
		return
	}

	feeAmount := ctx.IsInvalidAmountString(params[3])
	if feeAmount == nil {
		return
	}

	// 创建action
	newact := actions.NewAction_8_SimpleSatoshiTransfer(*targetAddress, fields.VarInt8(satoshiAmount))

	// 创建交易
	newTrs, e5 := transactions.NewEmptyTransaction_2_Simple(*feeAddress)
	newTrs.Timestamp = fields.VarInt5(ctx.UseTimestamp()) // 使用 hold 的时间戳
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}
	newTrs.Fee = *feeAmount // set fee
	// 放入action
	newTrs.AppendAction(newact)

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

////////////////////////////////////////////////////////

/*

gentx btcmove <trsno> <block_height> <block_timestamp> <prev_btc> <btc> <add_hac> <origin_address> <trs_btc_tx_hx> <fee_addr> <fee>

// test:
gentx btcmove 1 1001 1596702752 0 1 1048576 1EDUeK8NAjrgYhgDFv9NJecn8dNyJJsu3y 8deb5180a3388fee4991674c62705041616980e76288a8888b65530e41ccf90d 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 HAC1:248

*/

func num(str string) fields.VarInt4 {
	n, _ := strconv.ParseInt(str, 10, 0)
	return fields.VarInt4(n)
}

// 创建发布 转移 BTC
func GenTxCreateSatoshiGenesis(ctx ctx.Context, params []string) {
	if len(params) < 10 {
		fmt.Println("params not enough")
		return
	}
	genisisAct := &actions.Action_7_SatoshiGenesis{
		TransferNo:               num(params[0]),
		BitcoinBlockHeight:       num(params[1]),
		BitcoinBlockTimestamp:    num(params[2]),
		BitcoinEffectiveGenesis:  num(params[3]),
		BitcoinQuantity:          num(params[4]),
		AdditionalTotalHacAmount: num(params[5]),
		OriginAddress:            nil,
		BitcoinTransferHash:      nil,
	}

	originAddress := ctx.IsInvalidAccountAddress(params[6])
	if originAddress == nil {
		return
	}
	genisisAct.OriginAddress = *originAddress
	genisisAct.BitcoinTransferHash, _ = hex.DecodeString(params[7])

	feeAddress := ctx.IsInvalidAccountAddress(params[8])
	if feeAddress == nil {
		return
	}
	feeAmount := ctx.IsInvalidAmountString(params[9])
	if feeAmount == nil {
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
	newTrs.AppendAction(genisisAct)

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
