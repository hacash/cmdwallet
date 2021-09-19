package gentxs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/cmdwallet/toolshell/ctx"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"strconv"
)

// gentx paychan ${ADDRESS1} ${AMOUNT1} ${ADDRESS2} ${AMOUNT2} ${FEE}
/*

passwd 123456
passwd 12345678
gentx paychan 1EDUeK8NAjrgYhgDFv9NJecn8dNyJJsu3y HCX1:248 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 HCX1:248 HCX4:244


*/
// 创建支付通道
func GenTxCreatePaymentChannel(ctx ctx.Context, params []string) {
	if len(params) < 5 {
		fmt.Println("params not enough")
		return
	}
	// fmt.Println(strings.Join(params, ","))
	leftAddressArgv := params[0]
	leftAmountArgv := params[1]
	rightAddressArgv := params[2]
	rightAmountArgv := params[3]
	feeArgv := params[4]
	// 检查字段
	leftAddress := ctx.IsInvalidAccountAddress(leftAddressArgv)
	if leftAddress == nil {
		return
	}
	rightAddress := ctx.IsInvalidAccountAddress(rightAddressArgv)
	if rightAddress == nil {
		return
	}
	leftAmount := ctx.IsInvalidAmountString(leftAmountArgv)
	if leftAmount == nil {
		return
	}
	rightAmount := ctx.IsInvalidAmountString(rightAmountArgv)
	if rightAmount == nil {
		return
	}
	fee := ctx.IsInvalidAmountString(feeArgv)
	if fee == nil {
		return
	}
	// 开始拼装 action
	var paychan actions.Action_2_OpenPaymentChannel
	paychan.LeftAddress = *leftAddress
	paychan.LeftAmount = *leftAmount
	paychan.RightAddress = *rightAddress
	paychan.RightAmount = *rightAmount
	pcbts, _ := paychan.Serialize()
	bufs := bytes.NewBuffer(pcbts[16:])
	bufs.Write([]byte(strconv.FormatUint(ctx.UseTimestamp(), 10)))
	hx := fields.CalculateHash(bufs.Bytes())
	paychan.ChannelId = fields.ChannelId(hx[0:16])
	// 创建交易
	newTrs, e5 := transactions.NewEmptyTransaction_2_Simple(*leftAddress)
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}
	newTrs.Timestamp = fields.BlockTxTimestamp(ctx.UseTimestamp()) // 使用 hold 的时间戳
	newTrs.Fee = *fee                                              // set fee
	// 放入action
	e5 = newTrs.AppendAction(&paychan)
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}

	// 打印 hash 签名数据

	// ok
	ctx.Println("transaction create success! ")
	ctx.Println("hash: <" + hex.EncodeToString(newTrs.Hash()) + ">, hash_with_fee: <" + hex.EncodeToString(newTrs.HashWithFee()) + ">")
	ctx.Println("( payment_channel_id = <" + hex.EncodeToString(paychan.ChannelId) + "> )")

	bodybytes, e7 := newTrs.Serialize()
	if e7 != nil {
		fmt.Println("transaction serialize error, " + e7.Error())
		return
	}
	ctx.Println("body length " + strconv.Itoa(len(bodybytes)) + " bytes, hex body is:")
	ctx.Println("-------- TRANSACTION BODY [NOT SIGN] START --------")
	ctx.Println(hex.EncodeToString(bodybytes))
	ctx.Println("-------- TRANSACTION BODY [NOT SIGN] END   --------")

	// 记录
	ctx.SetTxToRecord(newTrs.Hash(), newTrs)

}

/////////////////////////////////////////////////////////////////////////

/*


gentx paychan_close $CHANNELID $FEEADDRESS $FEE


passwd 123456
passwd 12345678
gentx paychan_close d952144400ad6f5ff3da594a357b1dab 1EDUeK8NAjrgYhgDFv9NJecn8dNyJJsu3y HCX1:244




*/

// 关闭结算支付通道
func GenTxClosePaymentChannel(ctx ctx.Context, params []string) {
	if len(params) < 1 {
		fmt.Println("params not enough")
		return
	}
	// fmt.Println(strings.Join(params, ","))
	channelIdArgv := params[0]
	feeAddressArgv := params[1]
	feeArgv := params[2]
	// 检查字段
	channelhash, e3 := hex.DecodeString(channelIdArgv)
	if e3 != nil || len(channelhash) != 16 {
		fmt.Printf("channelId %s format is error.\n", channelhash)
		return
	}
	feeAddress := ctx.IsInvalidAccountAddress(feeAddressArgv)
	if feeAddress == nil {
		return
	}
	feeAmount := ctx.IsInvalidAmountString(feeArgv)
	if feeAmount == nil {
		return
	}
	// 创建支付通道
	var paychanclose actions.Action_3_ClosePaymentChannel
	paychanclose.ChannelId = channelhash

	// 创建交易
	newTrs, e5 := transactions.NewEmptyTransaction_2_Simple(*feeAddress)
	newTrs.Timestamp = fields.BlockTxTimestamp(ctx.UseTimestamp()) // 使用 hold 的时间戳
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}
	newTrs.Fee = *feeAmount // set fee
	// 放入action
	e5 = newTrs.AppendAction(&paychanclose)
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
	// 不检查签名
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
	fmt.Println("-------- TRANSACTION BODY [NOT SIGN COMPLETELY] START --------")
	fmt.Println(hex.EncodeToString(bodybytes))
	//fmt.Println( hex.EncodeToString( bodybytes2 ) )
	fmt.Println("-------- TRANSACTION BODY [NOT SIGN COMPLETELY] END   --------")

	// 记录
	ctx.SetTxToRecord(newTrs.Hash(), newTrs)

}
