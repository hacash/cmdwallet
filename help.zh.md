Hacash 命令行钱包 使用教程
===

命令行钱包是一个离线的工具，你可以在不联网的桌面电脑上使用它，创建账户、构造转账交易等等。

这个工具是绝对安全的，因为它不会在互联网上传输你的密码或者私钥，一切都是“本地的”操作。

因为本地命令行钱包不联网，所以不支持余额查询和检查。如果你需要查看某个地址的余额，可以去[在线钱包](https://wallet.hacash.org/)上查询。查询余额这个操作由于完全不涉及到密码或私钥，在在线钱包上查询也是安全的。

本教程将为不懂计算机技术的人提供最详细的使用教程，一步一步地教你达成例如创建地址、构造转账交易、提交转账交易等操作。

### 准备工作

你需要去[hacash.org](https://hacash.org/)下载已提前编译好的软件，按照你的操作系统类别（Windows或Ubuntu）下载对应版本的软件。然后解压下载的ZIP文件，然后进入解压的目录，然后运行它们（windows下双击运行，Ubuntu需要在程序目录里按下鼠标右键，选择“在终端打开”，然后输入）。

提示：如果你不会在 Ubuntu 下面运行一个命令行程序，请参照这个[教程](https://zhidao.baidu.com/question/501887268.html)。

### 开始

软件运行后，你将会在窗口看到如下内容：

```shell script
Welcome to Hacash tool shell, you can:
--------
passwd ${XXX} ${XXX}  |  prikey ${0xAB123D...}  |  newkey  |  accounts  |  update | log
--------
gentx sendcash ${FROM_ADDRESS} ${TO_ADDRESS} ${AMOUNT} ${FEE}  |  loadtx ${0xAB123D...}  |  txs
--------
sendtx $TXHASH $IP:PORT
--------
exit, quit
--------
The use time hold on 2020/08/10 15:32:12, enter 'update' change to now
Continue to enter anything:
>

```

这表示工具已成功运行。现在你只需要按输入指定的命令就可以完成诸如创建账户、构造交易等操作。

很简单、我们一步一步来：

#### 1. 通过密码创建一个账户地址

在窗口里输入 `passwd mypassword123456` :
```shell script
Continue to enter anything:
>passwd mypassword123456
```
然后按下回车键，你将会看到如下打印内容：

```shell script
>passwd mypassword123456
Loaded your account 
  private key: 0x9f33a504208cfabc1066f04a1592703465775a89d3fc59bcba54538a354330a9 
  address: 1Nk6Vcq4gZM6hJPhjtJYYWHmXRF8wtD7F5
```
其中 `1Nk6Vcq4gZM6hJPhjtJYYWHmXRF8wtD7F5` 就是通过你的密码 `mypassword123456` 生成的对应的账户地址。一个密码对应唯一一个地址，它可以用来收款或收取挖矿奖励等等。

你需要创造一些尽量复杂，但又不会忘记，且只有你自己知道的密码。如果你的密码过于简单可能被黑客猜中，或者你泄露了密码，则账户里的币就会被别人转走，永远找不回来了。

#### 2. 随机生成一个账户

如果你不想记录密码，而是习惯像比特币一样随机生成一个私钥和对应的账户，则可以输入 `newkey` 命令：

```shell script
>newkey
```
你将会看到如下输出

```shell script
>newkey
Loaded your account
  private key: 0xecf9549e4a6fa172c2d331fa84c0d72b4e06217ff80bf213c9de284cb7c62218
  address: 1Hd1LabwkyFvsZu5MravwR7KBdKripLLDi
```
`0xecf9549e4a6fa172c2d331fa84c0d72b4e06217ff80bf213c9de284cb7c62218` 就是随机生成的私钥，`1Hd1LabwkyFvsZu5MravwR7KBdKripLLDi`即为账户地址。请务必备份好私钥。

#### 3. 登录账户

通过输入 `passwd` 或 `newkey` ， `1Nk6Vcq4gZM6hJPhjtJYYWHmXRF8wtD7F5` 和 `1Hd1LabwkyFvsZu5MravwR7KBdKripLLDi` 两个账户已经在工具内被记录（或者叫登录），也就是说，本工具在运行期间可以从这两个地址生成转账交易。

如果并不是第一次打开运行本工具，但是你有一个备份的私钥，那么输入 `prikey` 命令就可以通过私钥登录到账户，例如你的私钥是 `0xecf9549e4a6fa172c2d331fa84c0d72b4e06217ff80bf213c9de284cb7c62218`， 则在命令行输入：
```shell script
>prikey 0xecf9549e4a6fa172c2d331fa84c0d72b4e06217ff80bf213c9de284cb7c62218
```
可以看到打印：
```shell script
>prikey 0xecf9549e4a6fa172c2d331fa84c0d72b4e06217ff80bf213c9de284cb7c62218
Loaded your account 
  private key: 0xecf9549e4a6fa172c2d331fa84c0d72b4e06217ff80bf213c9de284cb7c62218
  address: 1Hd1LabwkyFvsZu5MravwR7KBdKripLLDi
```
这时 `1Hd1LabwkyFvsZu5MravwR7KBdKripLLDi` 就再次被记录到本工具内了，这时你就可以构造从这个地址转出币的交易。

也就是说，你必须知晓一个账户的 密码 或者 私钥，并通过 `passwd` 或 `prikey` 命令记录到本工具内，在可以构造有效的转账交易。

而且，为了账户安全，本工具不会在某一个磁盘位置上保存你的密码或私钥，你每次运行本工具，都需要再次输入密码或私钥来登录账户。

#### 4. 构造转账交易

构造一笔转账交易，需要输入四个内容： 1.付款账户地址， 2.收款账户地址， 3.付款金额， 4.手续费

比如，地址 `1Nk6Vcq4gZM6hJPhjtJYYWHmXRF8wtD7F5` 需要向地址 `1Hd1LabwkyFvsZu5MravwR7KBdKripLLDi` 转账 200 枚 HAC，并同时设置手续费为 0.0001 枚 HAC，则输入一下命令：

```shell script
>gentx sendcash 1Nk6Vcq4gZM6hJPhjtJYYWHmXRF8wtD7F5 1Hd1LabwkyFvsZu5MravwR7KBdKripLLDi HAC2:250 HAC1:244
```

其中 `HAC2:250` 表示转账金额 200 枚，`HAC1:244` 表示本次转账支付手续费 0.0001 枚。回车后即可看见：

```shell script
>gentx sendcash 1Nk6Vcq4gZM6hJPhjtJYYWHmXRF8wtD7F5 1Hd1LabwkyFvsZu5MravwR7KBdKripLLDi HAC2:250 HAC1:244
transaction create success! 
hash: <9eb1d626393a702766892ddea3e91e2022cb9870ea20e766a76401948e15916d>, hash_with_fee: <4f32b5ae2b2f10d60c190d904bcb54939f09f702f2d246fa82dc57772c2a8273>
body length 159 bytes, hex body is:
-------- TRANSACTION BODY START --------
02005f30f7fc00ee811d0d3efe16c1583cfe35af3bbf080771ec7af401010001000100b65144da73d6063bae8cb8e2fd615223a01a2564fa010200010315af4b3de7ad55bf56198ffaea71fbabb3d98d8a5d6d4cb5933a9574de4096b5504a39427b2adbfdc4de12382a0c28de08e29a553328cc25c2d2256026a62a24779f3543de731892f9b1ce815761e407878e18a0d4e8a2869fd53a969e731d830000
-------- TRANSACTION BODY END   --------
```
成功构造了一笔转账交易，交易哈希为 `9eb1d626393a702766892ddea3e91e2022cb9870ea20e766a76401948e15916d` ，交易内容为 `02005f30f7fc00ee811d0d3efe16c1583cfe35af3bbf080771ec7af401010001000100b65144da73d6063bae8cb8e2fd615223a01a2564fa010200010315af4b3de7ad55bf56198ffaea71fbabb3d98d8a5d6d4cb5933a9574de4096b5504a39427b2adbfdc4de12382a0c28de08e29a553328cc25c2d2256026a62a24779f3543de731892f9b1ce815761e407878e18a0d4e8a2869fd53a969e731d830000`

但此时的交易仍然只存在于本地，并没有广播的区块链主网去打包确认。你需要复制交易内容 `9eb1d626.......31d830000` 去[在线钱包](https://wallet.hacash.org/)上提交。

去在线钱包提交构造好的交易，与查询余额一样，也是安全的，因为交易内容是已经签名好的数据，其中并不包含你的密码或私钥。

打开在线钱包，找到 发送交易 （Send transaction） 项，将拷贝的内容复制进去，点击“发送交易”（Send Tx），向 Hacash 主网提交构造好的转账交易：

![](https://raw.githubusercontent.com/hacash/cmdwallet/master/help/sendtx_online.png)

我们看到此时返回了一个错误： `Send failed: Transaction Add to MemTxPool error: address 1Nk6Vcq4gZM6hJPhjtJYYWHmXRF8wtD7F5 balance ㄜ0:0 not enough， need ㄜ2:250.` 因为 `1Nk6Vcq4gZM6hJPhjtJYYWHmXRF8wtD7F5` 的地址是我们刚刚创建的，里面并没有币（余额为0）。需要通过购买或接受赠送，或者挖矿得到币后，才能转账成功。





