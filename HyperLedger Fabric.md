# HyperLedger Fabric

- Hyperledger Fabric
  - Hyperledger Fabric是由 Linux 基金会支持的开源项目，提供了一种可扩展、可定制和安全的联盟链技术。
  - 官网：[🔗](https://www.hyperledger.org/use/fabric)
  - paper:[🔗](https://arxiv.org/pdf/1801.10228.pdf)
  - 官方文档：[🔗](https://hyperledger-fabric.readthedocs.io/en/release-2.5/)
  - 中文文档：[🔗](https://hyperledgercn.github.io/hyperledgerDocs/)
  - 中文视频Tutorial：[🔗](https://wiki.hyperledger.org/display/TWGC/Fabric+Video+Tutorial)
  - FAQ：[🔗](https://github.com/Hyperledger-TWGC/FAQ)
  - 博客(Hyperledger Fabric 2.0系列)：[🔗](https://blog.csdn.net/qq_28540443/article/details/104265844)
  - 电子书《区块链技术指南》:[🔗](https://github.com/yeasy/blockchain_guide)
  - 电子书《Hyperledger 源码分析之Fabric》:[🔗](https://github.com/yeasy/hyperledger_code_fabric)



- CA 由组织给管理者和网络节点颁发证书，证书可以用来识别属于某组织的组件，也可也为交易提供签名，交易签名是交易入块的前提条件
  - 客户端应用的交易提案
  - 智能合约的交易响应

- **源码分析**
  - 《Hyperledger Fabric源代码分析与深入解读》
  - Hyperledger Fabric v2.x 最新资料汇总[🔗](https://hello2mao.github.io/2020/04/22/hyperledger-fabric-v2.x-info/)
  - Fabric2.2中的Raft共识模块源码分析[🔗](https://www.cnblogs.com/GarrettWale/p/16131853.html)
  - 深入浅出FISCO BCOS[🔗](https://fisco-bcos-documentation.readthedocs.io/zh_CN/latest/docs/articles/index.html)
  - FISCO BCOS核心模块设计解析[🔗](https://fisco-bcos-documentation.readthedocs.io/zh_CN/latest/docs/design/index.html)

## Getting Started

在Ubuntu18.04上的部署过程

### Prerequisites

```bash
sudo apt-get install git   # clone fabric project
sudo apt-get install curl  # ?
sudo apt-get -y install docker-compose # -y 表示自动回答 yes
```

```bash
# 验证docker已经安装
docker --version
docker-compose --version
```

```bash
sudo systemctl enable docker # 设置开机自动启动docker
sudo systemctl start docker # 开启Docker daemon 进程
# 现在只有root用户可以运行docker,需要将用户加入到docker用户组中并切换到docker用户组
sudo usermod -a -G docker <username>
newgrp docker # 当前用户同时属于2个用户组，默认是原用户组，需要切换用户组
```

```bash
# Install Go 如果需要write Go chaincode or SDK applications
 rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
 export PATH=$PATH:/usr/local/go/bin # ~/.profile中添加该声明
 source ~/.profile # 使.profile中的声明生效
 go version 
```

```bash
# Install JQ 仅在与channel configuration transaction相关的tutorial中使用到
```

### Install Fabric and Fabric Samples

官方通过Docker compose创建了a simple Fabric test network，并以一系列应用验证了核心功能

官方预编译了`Fabric CLI tool binaries`和`Fabric Docker images`为我们使用

```bash
# 创建工作目录 
mkdir -p $HOME/go/src/github.com/<your_github_userid>
cd $HOME/go/src/github.com/<your_github_userid>
# cURL工具下载install-fabric.sh并授予执行权限
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
# original script : bootstrap.sh
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/bootstrap.sh| bash -s
#install-fabric.sh 相比 bootstrap.sh 优化了syntax,提供了积极的选择方式
```

接下来可以执行`install-fabric.sh`来安装fabric，但执行该命令时需要选择参数

```bash
# choose compnents and version
./install-fabric.sh docker samples binary
./install-fabric.sh --fabric-version 2.5.4 binary
```

此外如果要安装contributor版本，需要额外参考

[开发者模式]: https://hyperledger-fabric.readthedocs.io/en/latest/dev-setup/devenv.html

### Fabric Contract APIs and Application APIs

#### Fabric Contract APIs

- [Go contract API](https://github.com/hyperledger/fabric-contract-api-go) and [documentation](https://pkg.go.dev/github.com/hyperledger/fabric-contract-api-go).

下面根据Documents中的Tutorial一步一步开发智能合约并部署

##### Prerequisites

- 在安装Fabric时除了Go都已经安装
  - A clone of [fabric-samples](https://github.com/hyperledger/fabric-samples)
  - [Go 1.19.x](https://golang.org/doc/install)
  - [Docker](https://docs.docker.com/install/)
  - [Docker compose](https://docs.docker.com/compose/install/)

##### Housekeeping

```bash
mkdir fabric-samples/chaincode/contract-tutorial
cd fabric-samples/chaincode/contract-tutorial
go mod init github.com/hyperledger/fabric-samples/chaincode/contract-tutorial
go get -u github.com/hyperledger/fabric-contract-api-go
```

##### Declaring a contract

在`contract-tutorial`目录下创建`simple-contract.go`

所有需要在chaincode中使用的合约都需要implement the [contractapi.ContractInterface](https://godoc.org/github.com/hyperledger/fabric-contract-api-go/contractapi#ContractInterface) 实现该接口最简单的方式是将`contractapi.Contract`struct嵌入到合约中,示例如下:

```go
package main

import (
    "errors"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SimpleContract contract for handling writing and reading from the world state
type SimpleContract struct {
    contractapi.Contract
}
```

##### Writing contract functions

合约中的public function需要满足一些规则,否则在chaincode创建时会报错,下面是实现一个contract fucntion的示例:



##### Using contracts in chaincode

同样在`simple-contract.go`目录下创建一个文件`main.go`

```go
package main

import (
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	simpleContract := new(SimpleContract)
    cc, err := contractapi.NewChaincode(simpleContract)
    if err != nil {
        panic(err.Error())
    }
    if err := cc.Start(); err != nil {
        panic(err.Error())
    }
}
```

##### Testing your chaincode as a developer

打开一个terminal,进入到`fabric-samples`目录下的`chaincode-docker-devmode`中,该文件夹下提供了一个定义了一个 simple fabric network的docker compose文件,可以通过该文件运行chaincode

```bash
docker-compose -f docker-compose-simple.yaml.up # 启动fabric network
```

在`chaincode-docker-devmode`目录下打开一个新的terminal

```bash
docker exec -it chaincode sh # enter the chaincode docker container
cd contract-tutorial # 这一步不确定
```

接下来要保证fabric docker image版本在2.x.x,否则`go build`命令会fail , 并且要使用`chmod -R 766`命令使得`contract-tutorial` folder有创建文件的权限

```bash
go mod vendor
go build
# run the chaincode
CORE_CHAINCODE_ID_NAME=mycc:0 CORE_PEER_TLS_ENABLED=false ./contract-tutorial -peer.address peer:7052 
```

##### Interacting with the chaincode

新打开一个终端

```bash
docker exec -it cli sh
peer chaincode install -p chaincodedev/chaincode/contract-tutorial -n mycc -v 0
```



#### Fabric Application APIs

