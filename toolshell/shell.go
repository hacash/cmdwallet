package toolshell

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/hacash/core/account"
	"github.com/hacash/core/interfaces"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var (
	MyAccounts          = make(map[string]account.Account, 0)
	AllPrivateKeyBytes  = make(map[string][]byte, 0)
	Transactions        = make(map[string]interfaces.Transaction, 0)
	TargetTime          time.Time // 使用的时间
	currentInputContent string
)

////////////////////////////////////

var welcomeContent = `
Welcome to Hacash tool shell, you can:
--------
passwd $XXX $XXX  |  prikey $0xAB123D...  |  newkey  |  accounts  |  update
--------
gentx sendcash $FROM_ADDRESS $TO_ADDRESS $AMOUNT $FEE  |  loadtx $0xTXBODYBYTES  |  txs
--------
sendtx $TXHASH $IP:PORT
--------
exit, quit
--------`

func RunToolShell() {

	abspath, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	abspath = path.Dir(abspath)
	logfilename := path.Join(abspath, time.Now().Format("2006-01-02_15:04:05") + ".log")

	logfile, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY,0660)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	var ctx = &ctxToolShell{
		logfile: logfile,
	}

	fmt.Println(welcomeContent)
	TargetTime = time.Now()
	holdTimeStr := TargetTime.Format("2006/01/02 15:04:05")
	fmt.Println("The use time hold on " + holdTimeStr + ", enter 'update' change to now")
	fmt.Println("Continue to enter anything:")

	ctx.LogFileWriteln(holdTimeStr)

	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}
		//fmt.Scanln(&currentInputContent)
		currentInputContent = strings.TrimSpace(input)
		// empty
		if currentInputContent == "" {
			continue
		}

		// exit
		if currentInputContent == "exit" ||
			currentInputContent == "quit" {
			fmt.Println("Bye")
			break
		}

		ctx.LogFileWriteln("- - - - - - - - - - - - - - - -\n> " + currentInputContent)

		if currentInputContent == "update" {
			TargetTime = time.Now()
			ctx.Println("Hold time change to " + TargetTime.Format("2006/01/02 15:04:05"))
		} else if currentInputContent == "accounts" {
			showAccounts(ctx)
		} else if currentInputContent == "txs" {
			showTxs(ctx)
		}else {
			// other opration
			params := strings.Fields(currentInputContent)
			funcname := params[0]
			parabody := params[1:]
			switch params[0] {
			case "passwd":
				setPrivateKeyByPassword(ctx, parabody)
			case "prikey":
				setPrivateKey(ctx, parabody)
			case "newkey":
				createNewPrivateKey(ctx, parabody)
			case "puttx":
				putTx(ctx, parabody)
			case "gettx":
				getTx(ctx, parabody)
			case "signtx":
				signTx(ctx, parabody)
			case "gentx":
				genTx(ctx, parabody)
			case "sendtx":
				sendTxToMiner(ctx, parabody)
			default:
				fmt.Println("Sorry, undefined instructions: " + funcname)
			}
		}

		// clear
		currentInputContent = ""
	}
}

/////////////////////////////////////////////////

/////////////////////////////////////////////////////////

func showAccounts(ctx *ctxToolShell) {
	if len(MyAccounts) == 0 {
		fmt.Println("none")
		return
	}
	for k, _ := range MyAccounts {
		ctx.Printf(k + " ")
	}
	ctx.Printf("\n")
}

func showTxs(ctx *ctxToolShell) {
	if len(Transactions) == 0 {
		fmt.Println("none")
		return
	}
	for k, _ := range Transactions {
		ctx.Println(hex.EncodeToString([]byte(k)))
	}
}

func setPrivateKey(ctx *ctxToolShell, params []string) {
	for _, hexstr := range params {
		if strings.HasPrefix(hexstr, "0x") {
			hexstr = string([]byte(hexstr)[2:]) // drop 0x
		}
		_, e0 := hex.DecodeString(hexstr)
		if e0 != nil {
			fmt.Println("Private Key '" + hexstr + "' is error")
			return
		}
		acc, e1 := account.GetAccountByPriviteKeyHex(hexstr)
		if e1 != nil {
			fmt.Println("Private Key '" + hexstr + "' is error")
			return
		}
		printLoadAddress(ctx, acc)
	}
}

func setPrivateKeyByPassword(ctx *ctxToolShell, params []string) {
	for _, passwd := range params {
		//fmt.Println(passwd)
		acc := account.CreateAccountByPassword(passwd)
		printLoadAddress(ctx, acc)
	}
}

// 随机创建私钥
func createNewPrivateKey(ctx *ctxToolShell, params []string) {
	acc := account.CreateNewAccount()
	printLoadAddress(ctx, acc)
}

func printLoadAddress(ctx *ctxToolShell, acc *account.Account) {
	MyAccounts[string(acc.AddressReadable)] = *acc // append
	AllPrivateKeyBytes[string(acc.Address)] = acc.PrivateKey
	ctx.Println("Loaded your account private key: 0x" + hex.EncodeToString(acc.PrivateKey) + " address: " + string(acc.AddressReadable))
}
