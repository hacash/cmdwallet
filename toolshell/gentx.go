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
	case "sendcash": // Send transaction
		gentxs.GenTxSimpleTransfer(ctx, bodys)
	case "paychan": // Create payment channel
		gentxs.GenTxCreatePaymentChannel(ctx, bodys)
	case "paychan_close": // Close settlement payment channel
		gentxs.GenTxClosePaymentChannel(ctx, bodys)
	case "diamond": // Create diamond
		gentxs.GenTxCreateDiamond(ctx, bodys)
	case "diamond_transfer": // Transfer diamond
		gentxs.GenTxDiamondTransfer(ctx, bodys)
	case "diamond_transfer_quantity": // Transfer diamond
		gentxs.GenTxOutfeeQuantityDiamondTransfer(ctx, bodys)
	case "btcmove": // Confirm BTC one-way transfer
		gentxs.GenTxCreateSatoshiGenesis(ctx, bodys)
	case "sendsat": // Confirm BTC one-way rotation
		gentxs.GenTxSimpleTransferSatoshi(ctx, bodys)
	case "create_lockbls": // Create lock
		gentxs.GenTxCreateLockbls(ctx, bodys)
	case "release_lockbls": // Release lock
		gentxs.GenTxReleaseLockbls(ctx, bodys)

	default:
		fmt.Println("Sorry, undefined gentx type: " + typename)
	}
}
