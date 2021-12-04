package ctx

import (
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
)

type Context interface {
	NotLoadedYetAccountAddress(string) bool         // 检测账户是否已经登录
	IsInvalidAccountAddress(string) *fields.Address // 检测是否为一个合法的账户名
	IsInvalidAmountString(string) *fields.Amount    // 检测是否为一个合法的金额数量
	GetAllPrivateKeyBytes() map[string][]byte       // 获取全部私钥，用于填充签名
	SetTxToRecord([]byte, interfaces.Transaction)   // 记录交易
	GetTxFromRecord([]byte) interfaces.Transaction  // 获取交易
	UseTimestamp() uint64                           // 当前使用的时间戳

	//////////////////////////////////////////////////////////////////

	Println(...interface{})
	Print(...interface{})
}
