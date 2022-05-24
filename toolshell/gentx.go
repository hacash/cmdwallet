package toolshell

import (
	"fmt"
	"github.com/hacash/cmdwallet/toolshell/ctx"
	"github.com/hacash/cmdwallet/toolshell/gentxs"
)

func genTx(ctx ctx.Context, params []string) {
	if len(params) <= 1 {
		fmt.Println("params not enough")
		return
	}

	typename := params[0]
	bodys := params[1:]
	switch typename {
	case "sendcash": // 发送交易
		gentxs.GenTxSimpleTransfer(ctx, bodys)
	case "paychan": // 创建支付通道
		gentxs.GenTxCreatePaymentChannel(ctx, bodys)
	case "paychan_close": // 关闭结算支付通道
		gentxs.GenTxClosePaymentChannel(ctx, bodys)
	case "diamond": // 创建钻石
		gentxs.GenTxCreateDiamond(ctx, bodys)
	case "diamond_transfer": // 转移钻石
		gentxs.GenTxDiamondTransfer(ctx, bodys)
	case "diamond_transfer_quantity": // 转移钻石
		gentxs.GenTxOutfeeQuantityDiamondTransfer(ctx, bodys)
	case "btcmove": // 确认BTC单向转移
		gentxs.GenTxCreateSatoshiGenesis(ctx, bodys)
	case "sendsat": // 确认BTC单向转
		gentxs.GenTxSimpleTransferSatoshi(ctx, bodys)
	case "create_lockbls": // 创建锁仓
		gentxs.GenTxCreateLockbls(ctx, bodys)
	case "release_lockbls": // 释放锁仓
		gentxs.GenTxReleaseLockbls(ctx, bodys)

	default:
		fmt.Println("Sorry, undefined gentx type: " + typename)
	}
}
