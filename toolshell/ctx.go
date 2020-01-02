package toolshell

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/fields"
	"github.com/hacash/core/interfaces"
	"os"
)

/*

 */

func handleArgvToBytes(spx string, end string, argv ...interface{}) []byte {
	buf := bytes.NewBuffer([]byte{})
	for _, a := range argv {
		//fmt.Println(reflect.ValueOf(a).Type().String())
		if str, ok := a.(string); ok {
			buf.Write([]byte(str+spx))
		}else if bts, ok := a.([]byte); ok {
			buf.Write(bts)
		}
	}
	buf.Write([]byte(end))
	return buf.Bytes()
}


type ctxToolShell struct{
	logfile *os.File
}


func (c *ctxToolShell) Println(argv ...interface{}) {
	fmt.Println(argv...)
	c.LogFileWriteln(argv...)
}

func (c *ctxToolShell) LogFileWriteln(strs ...interface{}) {
	c.logfile.Write( handleArgvToBytes(" ", "\n", strs...) )
}

func (c *ctxToolShell) Print(strs ...interface{}) {
	fmt.Print(strs...)
	c.LogFileWrite(strs...)
}

func (c *ctxToolShell) LogFileWrite(strs ...interface{}) {
	c.logfile.Write( handleArgvToBytes("", "", strs...) )
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

