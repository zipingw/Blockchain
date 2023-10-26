# Application for IPFS & Fabric

## Call IPFS API

python用于调用ipfs api的库ipfshttpclient目前（2023.10.24）只支持0.8.0以下版本的实现

解决方案一： 自行编写ipfsAPI类，通过request包发送http请求的方式调用ipfs api

解决方案二： 使用go语言提供的包go rpc ，尚未做深入探索，不清楚

## Call Fabric Chaincode

主要用SDK去调用chaincode

| fabric-sdk-go         | 需要通过org1下的connection-org1.yaml连接区块链网络 | 成功   |
| --------------------- | -------------------------------------------------- | ------ |
| **fabric-sdk-py**     | 需要network.json连接区块链网络                     | 失败   |
| **fabric-sdk-java**   |                                                    | 未尝试 |
| **fabric-sdk-nodejs** |                                                    | 未尝试 |

connection-org1.json包含的内容是network.json的子集

SDK在调用chaincode时，sdk-go要求chaincode提前使用命令行进行了部署，sdk-py目前遇到问题无法得到channel，也无法成功调用chaincode，并且其sdk实现的invoke中不包含选择chaincode中函数的参数

### 关于fabric test network安装部署新chaincode

```bash
# 使用Go一些通用性的相关配置
# 1. 环境变量配置代理地址
GOPROXY=https://goproxy.io,direct
# 2. 打开mod
sudo GO111MODULE=on
go mod init # 创建go.mod文件管理依赖包
go mod tidy # 在go.mod文件中移除不需要的依赖包，执行后生成go.sum文件（依赖下载条目）
go mod verify # 检查当前模块依赖是否都被下载
go mod vendor # 生成vendor文件夹，存储具体的依赖包，使go程序能够在无网络环境下运行
```

#### 方式一：network.sh脚本文件 deployCC命令

```bash
# 创建新文件夹newcc并进入，newcc相当于打包chaincode的最外层文件夹
mkdir newcc && cd newcc
# 在执行该命令的根目录下生成go.mod文件
go mod init newcc
# 创建chaincode-go文件夹存储具体的智能合约程序
mkdir chaincode-go && cd chaincode-go
# 创建chaincode.go并在其中编写具体的合约
touch chaincode.go && vim chaincode.go
# 处理依赖
go mod tidy
go mod vendor
# 注意go.mod中的go版本号1.21可以，1.21.3不符合格式，只能有1个小数点
# 一步部署chaincode; ccn(chaincode name):basic 是chaincode的name,sdk调用时的参数，ccp(chaincode path) 取到chaincode-go为止，ccl(chaincode language)
./network.sh deployCC -ccn basic -ccp ../ipfscc-basic/chaincode-go -ccl go
```

#### 方式二：

方式一虽然简洁，但是封装了具体部署时的各个细节，如果部署production network，network.sh就脚本并不适用，故需要学会逐步部署chaincode

```bash
# 第一步：Package 打包chaincode为压缩包, 压缩包名称basic.tar.gz，被压缩的chaincode的path，lang为chaincode的编程语言，label为chaincode版本号
peer lifecycle chaincode package basic.tar.gz --path ../asset-transfer-basic/chaincode-go/ --lang golang --label basic_1.0

# 第二步：Install 在test-network中是通过切换环境变量至peer0.org1再切换至peer0.org2，表示在该各个peer上安装chaincode
peer lifecycle chaincode install basic.tar.gz

# 第三步：Query 可查询是否已经在当前peer上安装chaincode，已安装则会输出CC Package ID
peer lifecycle chaincode queryinstalled

# 第四步：approveformyorg 涉及大量参数，每个安装了chaincode的peer需要执行命令提交同意，这里在test-network中每个Org只有一个peer,故只要该peer同意就表示Org同意，若有多个peer，需执行多次，这一步开始都需要提供ordeer的msp中的tlsca证书
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# 第五步：checkcommitreadiness 查看各Org approve状态
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json

# 第六步：commit 当check中的结果全部为true，则任意Org中的任意Peer可以执行commit，注册chaincode至channel,之后chaincode便可以被调用，commit时需要peer的ca证书
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

# 第七步：querycommitted 查询commit结果
peer lifecycle chaincode querycommitted --channelID mychannel --name basic --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

# 第八步：invoke 调用chaincode
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'

# 第九步：query 查询某个channel中某个智能合约的账本，以下是基于asset-transfer-basic的例子
peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'
```

### Gateway or Client

```go
// fabric-sdk-go 依赖包，由于是在github上，故执行go mod vendor获取依赖包时会出现访问超时问题，可以通过先ping github.com，再执行相关命令的方式解决
import (
        "fmt"
        "log"
        "os"
        "path/filepath"
        "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
        "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

func main() {
    // 声明环境变量，后续用到了"CHANNEL_NAME"，"CHAINCODE_NAME"
    err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
    if err != nil {
            log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environment variable: %v", err)
    }
    // Gateway 如何连接至 test-network ?
    ccpPath := filepath.Join(
            "/home",
            "zpwang",
            "fabric",
            "fabric-samples",
            "test-network",
            "organizations",
            "peerOrganizations",
            "org1.example.com",
            "connection-org1.yaml",
    )
    
    // wallet存在的意义还没有搞清楚
    walletPath := "wallet"
    // remove any existing wallet from prior runs
    os.RemoveAll(walletPath)
    wallet, err := gateway.NewFileSystemWallet(walletPath)
    if err != nil {
            log.Fatalf("Failed to create wallet: %v", err)
    }

    if !wallet.Exists("appUser") {
            err = populateWallet(wallet)
            if err != nil {
                    log.Fatalf("Failed to populate wallet contents: %v", err)
            }
    }
    // gateway通过org1的peer0的配置文件连接至peer0,再通过gateway连接network,本质是通过peer0作为桥梁访问整个network
    gw, err := gateway.Connect(
            gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
            gateway.WithIdentity(wallet, "appUser"),
    )
    if err != nil {
            log.Fatalf("Failed to connect to gateway: %v", err)
    }
    defer gw.Close()
	
    channelName := "mychannel"
    if cname := os.Getenv("CHANNEL_NAME"); cname != "" {
            channelName = cname
    }
	// gateway获取network实例
    log.Println("--> Connecting to channel", channelName)
    network, err := gw.GetNetwork(channelName)
    if err != nil {
            log.Fatalf("Failed to get network: %v", err)
    }

    chaincodeName := "basic"
    if ccname := os.Getenv("CHAINCODE_NAME"); ccname != "" {
            chaincodeName = ccname
    }
    // 从network中根据ccn获取chaincode实例
    log.Println("--> Using chaincode", chaincodeName)
    contract := network.GetContract(chaincodeName)
    // 通过chaincode实例调用账本
    // init
    result, err := contract.SubmitTransaction("InitLedger")
    if err != nil {
            log.Fatalf("Failed to Submit transaction: %v", err)
    }
    log.Println(string(result))
    // Create
    log.Println("--> Submit Transaction: CreateRecord, creates new record with HashValue, Source, Time, StoreNodes, Operation arguments")
    result, err = contract.SubmitTransaction("CreateRecord", "Qm001", "Rasperry001", "202310212042", "001,002", "add")
    if err != nil {
            log.Fatalf("Failed to Submit transaction: %v", err)
    }
    log.Println(string(result))
    // GetAll
    log.Println("--> Evaluate Transaction: GetAllRecords, function returns all the current assets on the ledger")
    result, err = contract.EvaluateTransaction("GetAllRecords")
    if err != nil {
            log.Fatalf("Failed to evaluate transaction: %v", err)
    }
    log.Println(string(result))
}
```

