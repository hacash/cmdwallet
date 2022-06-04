package gentxs

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/cmdwallet/toolshell/ctx"
	"github.com/hacash/core/actions"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/transactions"
	"github.com/hacash/x16rs"
	"strconv"
	"strings"
)

/*

gentx diamond ${DIAMOND} ${NUMBER} ${PrevHash} ${Nonce} ${Address} ${feeAddress} ${FEE}

passwd 123456
passwd 12345678
gentx diamond NHMYYM 1 000000077790ba2fcdeaef4a4299d9b667135bac577ce204dee8388f1b97f7e6 0100000001552c71 1271438866CSDpJUqrnchoJAiGGBFSQhjd 1EDUeK8NAjrgYhgDFv9NJecn8dNyJJsu3y HCX2:244
gentx diamond_transfer NHMYYM 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 1271438866CSDpJUqrnchoJAiGGBFSQhjd HCX1:244

*/

// Create diamond
func GenTxCreateDiamond(ctx ctx.Context, params []string) {
	if len(params) < 7 {
		fmt.Println("params not enough")
		return
	}

	diamondArgv := params[0]
	numberArgv := params[1]
	prevHashArgv := params[2]
	nonceArgv := params[3]
	addressArgv := params[4]
	feeAddressArgv := params[5]
	feeArgv := params[6]

	// Check field
	_, dddok := x16rs.IsDiamondHashResultString("0000000000" + diamondArgv)
	if !dddok {
		fmt.Printf("%s is not diamond value.\n", diamondArgv)
		return
	}

	number, e3 := strconv.ParseUint(numberArgv, 10, 0)
	if e3 != nil {
		fmt.Printf("number %s is error.\n", numberArgv)
		return
	}

	noncehash, e3 := hex.DecodeString(nonceArgv)
	if e3 != nil {
		fmt.Printf("nonce %s format is error.\n", nonceArgv)
		return
	}

	address := ctx.IsInvalidAccountAddress(addressArgv)
	if address == nil {
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

	blkhash, e0 := hex.DecodeString(prevHashArgv)
	if e0 != nil {
		fmt.Println("block hash format error")
		return
	}

	// Create action
	var dimcreate actions.Action_4_DiamondCreate
	dimcreate.Number = fields.DiamondNumber(number)
	dimcreate.Diamond = fields.DiamondName(diamondArgv)
	dimcreate.PrevHash = blkhash
	dimcreate.Nonce = fields.Bytes8(noncehash)
	dimcreate.Address = *address

	// Create transaction
	newTrs, e5 := transactions.NewEmptyTransaction_2_Simple(*feeAddress)
	newTrs.Timestamp = fields.BlockTxTimestamp(ctx.UseTimestamp()) // Use the timestamp of hold
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}
	newTrs.Fee = *feeAmount // set fee

	// Put in action
	newTrs.AppendAction(&dimcreate)

	// sign
	e6 := newTrs.FillNeedSigns(ctx.GetAllPrivateKeyBytes(), nil)
	if e6 != nil {
		fmt.Println("sign transaction error, " + e6.Error())
		return
	}

	// Check signature
	sigok, sigerr := newTrs.VerifyAllNeedSigns()
	if sigerr != nil || !sigok {
		fmt.Println("transaction VerifyAllNeedSigns fail")
		return
	}

	// Datalization
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
	fmt.Println("-------- TRANSACTION BODY END   --------")

	// record
	ctx.SetTxToRecord(newTrs.Hash(), newTrs)
}

// Transfer diamond
func GenTxDiamondTransfer(ctx ctx.Context, params []string) {
	if len(params) < 4 {
		fmt.Println("params not enough")
		return
	}

	diamondArgv := params[0]
	addressArgv := params[1]
	feeAddressArgv := params[2]
	feeArgv := params[3]
	// Check field
	_, dddok := x16rs.IsDiamondHashResultString("0000000000" + diamondArgv)
	if !dddok {
		fmt.Printf("%s is not diamond value.\n", diamondArgv)
		return
	}

	address := ctx.IsInvalidAccountAddress(addressArgv)
	if address == nil {
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

	// Create action
	var dimtransfer actions.Action_5_DiamondTransfer
	dimtransfer.Diamond = fields.DiamondName(diamondArgv)
	dimtransfer.ToAddress = *address
	// Create transaction
	newTrs, e5 := transactions.NewEmptyTransaction_2_Simple(*feeAddress)
	newTrs.Timestamp = fields.BlockTxTimestamp(ctx.UseTimestamp()) // Use the timestamp of hold
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}

	newTrs.Fee = *feeAmount // set fee
	// Put in action
	newTrs.AppendAction(&dimtransfer)

	// sign
	e6 := newTrs.FillNeedSigns(ctx.GetAllPrivateKeyBytes(), nil)
	if e6 != nil {
		fmt.Println("sign transaction error, " + e6.Error())
		return
	}

	// Check signature
	sigok, sigerr := newTrs.VerifyAllNeedSigns()
	if sigerr != nil || !sigok {
		fmt.Println("transaction VerifyAllNeedSigns fail")
		return
	}

	// Datalization
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
	fmt.Println("-------- TRANSACTION BODY END   --------")

	// record
	ctx.SetTxToRecord(newTrs.Hash(), newTrs)
}

/////////////////////////////////////////////////////////////////////////////////////

// Bulk transfer of diamonds
func GenTxOutfeeQuantityDiamondTransfer(ctx ctx.Context, params []string) {
	if len(params) < 5 {
		fmt.Println("params not enough")
		return
	}

	fromAddressArgv := params[0]
	toAddressArgv := params[1]
	diamondsArgv := params[2]
	feeAddressArgv := params[3]
	feeArgv := params[4]

	// Check field
	diamonds := strings.Split(diamondsArgv, ",")
	if len(diamonds) > 200 {
		fmt.Printf("diamonds number is too much.\n")
		return
	}

	for _, diamond := range diamonds {
		_, dddok := x16rs.IsDiamondHashResultString("0000000000" + diamond)
		if !dddok {
			fmt.Printf("%s is not diamond value.\n", diamond)
			return
		}
	}

	fromaddress := ctx.IsInvalidAccountAddress(fromAddressArgv)
	if fromaddress == nil {
		return
	}

	toaddress := ctx.IsInvalidAccountAddress(toAddressArgv)
	if toaddress == nil {
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

	// Create action
	var dimtransfer actions.Action_6_OutfeeQuantityDiamondTransfer
	dimtransfer.FromAddress = *fromaddress
	dimtransfer.ToAddress = *toaddress
	dimtransfer.DiamondList.Count = fields.VarUint1(len(diamonds))
	dimtransfer.DiamondList.Diamonds = make([]fields.DiamondName, len(diamonds))
	for i, v := range diamonds {
		dimtransfer.DiamondList.Diamonds[i] = fields.DiamondName(v)
	}

	// Create transaction
	newTrs, e5 := transactions.NewEmptyTransaction_2_Simple(*feeAddress)
	newTrs.Timestamp = fields.BlockTxTimestamp(ctx.UseTimestamp()) // Use the timestamp of hold
	if e5 != nil {
		fmt.Println("create transaction error, " + e5.Error())
		return
	}
	newTrs.Fee = *feeAmount // set fee

	// Put in action
	newTrs.AppendAction(&dimtransfer)

	// sign
	e6 := newTrs.FillNeedSigns(ctx.GetAllPrivateKeyBytes(), nil)
	if e6 != nil {
		fmt.Println("sign transaction error, " + e6.Error())
		return
	}

	// Check signature
	sigok, sigerr := newTrs.VerifyAllNeedSigns()
	if sigerr != nil || !sigok {
		fmt.Println("transaction VerifyAllNeedSigns fail")
		return
	}

	// Datalization
	bodybytes, e7 := newTrs.Serialize()
	if e7 != nil {
		fmt.Println("transaction serialize error, " + e7.Error())
		return
	}

	// ok
	ctx.Println("transaction create success! ")
	ctx.Println("hash: <" + hex.EncodeToString(newTrs.Hash()) + ">, hash_with_fee: <" + hex.EncodeToString(newTrs.HashWithFee()) + ">")
	ctx.Println("body length " + strconv.Itoa(len(bodybytes)) + " bytes, hex body is:")
	ctx.Println("-------- TRANSACTION BODY START --------")
	ctx.Println(hex.EncodeToString(bodybytes))
	ctx.Println("-------- TRANSACTION BODY END   --------")

	// record
	ctx.SetTxToRecord(newTrs.Hash(), newTrs)
}
