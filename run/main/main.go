package main

import "github.com/hacash/cmdwallet/toolshell"

/**


go build -o test run/main/main.go && ./test


gentx btcmove 1 1001 1596702752 0 1 1048576 1H1XgkBQdh2Tx3DyetHD2DBSbmzciu5C9s 8deb5180a3388fee4991674c62705041616980e76288a8888b65530e41ccf90d 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 HAC4:244

sendtx 75987b36e4cb7b35c6bd24b61d68e14adfb8ba53aaae74d99539bd0846396ec6 127.0.0.1:33381

gentx sendsat 1H1XgkBQdh2Tx3DyetHD2DBSbmzciu5C9s 1MzNY1oA3kfgYi75zquj3SRUPYztzXHzK9 40000000 HAC1:244



*/

func main() {

	toolshell.RunToolShell()

}
