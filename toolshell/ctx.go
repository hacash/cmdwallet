package toolshell

import (
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"os"
)

/*

 */

type ctxToolShell struct{
	logfile *os.File
}


func (c *ctxToolShell) Println(argv ...interface{}) {
	fmt.Println(argv)
	strs := make([]string, 0, len(argv))
	for _, a := range argv {
		if str, ok := a.(string); ok {
			strs = append(strs, str)
		}
	}
	c.LogFileWriteln(strs...)
}

func (c *ctxToolShell) LogFileWriteln(strs ...string) {
	for _, str := range strs {
		c.logfile.Write( []byte(str+"\n") )
	}
}

func (c *ctxToolShell) Printf(format string, argv ...interface{}) {
	str := fmt.Sprintf(format, argv)
	c.logfile.Write( []byte(str) )
	fmt.Printf(str)
}

func (c *ctxToolShell) LogFileWritef(format string, argv ...interface{}) {
	str := fmt.Sprintf(format, argv)
	c.logfile.Write( []byte(str) )
}

func (c *ctxToolShell) NotLoadedYetAccountAddress(addr string) bool {
	if _, ok := MyAccounts[addr]; !ok {
		fmt.Println("Account " + addr + " need to be loaded")
		return true
	}
	return false
}

func (c *ctxToolShell) IsInvalidAccountAddress(addr string) *fields.Address {
	address, err := fields.CheckReadableAddress(addr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return address
}

func (c *ctxToolShell) IsInvalidAmountString(amtstr string) *fields.Amount {
	amt, e1 := fields.NewAmountFromFinString(amtstr)
	if e1 != nil {
		fmt.Printf("amount \"%s\" format error or over range, the right example is 'HCX1:248' for one coin\n", amtstr)
		return nil
	}
	return amt
}

func (c *ctxToolShell) GetAllPrivateKeyBytes() map[string][]byte {
	return AllPrivateKeyBytes
}

func (c *ctxToolShell) SetTxToRecord(hash_no_fee []byte, tx interfaces.Transaction) { // 记录交易
	Transactions[string(hash_no_fee)] = tx
}

func (c *ctxToolShell) GetTxFromRecord(hash_no_fee []byte) interfaces.Transaction { // 获取交易
	if tx, ok := Transactions[string(hash_no_fee)]; ok {
		return tx
	} else {
		fmt.Println("Not find tx " + hex.EncodeToString(hash_no_fee))
		return nil
	}
}

func (c *ctxToolShell) UseTimestamp() uint64 { // 当前使用的时间戳
	return uint64(TargetTime.Unix())
}

