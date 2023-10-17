

```bash
# 换清华镜像源
sudo -s
cp /etc/apt/sources.list /etc/apt/sources.list.bak
vi /etc/apt/sources.list

# 默认注释了源码镜像以提高 apt update 速度，如有需要可自行取消注释
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse

# 预发布软件源，不建议启用
# deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse


sudo apt-get update

$ sudo apt-get install git
$ sudo apt-get install curl
sudo apt-get -y install docker-compose

sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -a -G docker <username>
newgrp docker # 这行好像又不是一定需要的

sudo apt-get install openssh-server
sudo /etc/init.d/ssh start
```



```bash
# linux下的下载命令
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
# windows下的下载命令
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh" -OutFile "install-fabric.sh"
# install.sh 需要开VPN才能下载，windows上下载后传过去
scp D:\file\install-fabric.sh zpw@192.168.209.133:/home/zpw/go/src/github.com/zipingw/
```

```bash
chmod +x install-fabric.sh
./install-fabric.sh d s b 
```

![image-20231016113042232](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20231016113042232.png)

transaction needs to be signed by multiple organizations and then can be committed to ledger

TX是由智能合约发起的，也是由智能合约endorse的，而智能合约是通过Application调用的

多个智能合约打包后被称为chaincode，chaincode先被install在peer上，再被部署到channel，但是chaincode必须在channel中的所有member同意之后才被允许部署

```bash
cd fabric-samples/test-network
sudo ./network.sh down
sudo ./network.sh up

docker ps -a

# 创建channel 
sudo ./network.sh createChannel [-c channel1]
# 启动network同时创建channel
sudo ./network.sh up createChannel

wget -c https://dl.google.com/go/go1.21.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/bin/go
sudo apt-get install jq
# 注意 go version可以在普通用户下执行但是sudo go version确会找不到go
cd /etc
sudo vim /etc/sudoers
# 在secure_path中新增路径 :/usr/local/bin/go
sudo ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go # 该命令用到jq和go 
# 但这条命令还是会报错，原因是chaincode-go中的go.mod中的版本号要设为1.21，不能是1.21.3，如果设为1.17不知道会不会报错
# 这个地方还需要配置代理 Clash 开启服务模式 打开TUN 使得WSL中可以ping google.com ，但是在执行该命令时还是出现访问proxy.golang.org超时错误 ,可是用ping proxy.golang.org是可以成功的，但是代理配置后 可以成功执行到Chaincode is packaged
# 重复多次(down and up)后发现不再出现访问github失败的问题
```



```bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
# 教程中的该命令有误 
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'
```

```bash
# 报错信息
2023-10-16 14:37:09.465 CST 0001 ERRO [main] InitCmd -> Cannot run peer because error when setting up MSP of type bccsp from directory /home/zpwang/fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp: KeyMaterial not found in SigningIdentityInfo
```

这个地方得通过切换为root再在root里声明环境变量，再invoke链码，这个之后我以为下一条命令peer chaincode query可以直接在普通用户中执行，但是还是会出现一样的MSP KeyMaterial not found in SigningIdentityInfo错误。

解决方法跟上面一样，看来涉及到链码的operation都必须要在root用户下执行，应该是权限问题导致无法访问文件就会使得其误以为文件不存在，***直觉上认为可以通过修改org1.example.com/users/Admin@org1.example.com/msp这个文件夹的权限来解决该问题*** 该文件夹应该是由于涉及安全问题设置了较高权限



wsl2中不能用systemctl来启动docker服务，要用service ,service也会启动失败，在windows中docker desktop中的setting=>Resouces=>WSL integration 选中用的系统，再打开Ubuntu就可以用sudo service start docker启动docker

```bash
# 关闭之前的网络
sudo ./network.sh down
# 重新开启网络 并创建通道
sudo ./network.sh up createChannel
# chaincode 要先安装go并配置环境变量
sudo ./network.sh deployCC -ccn basic -ccp ../asset-transfer-basic/chaincode-go -ccl go
# 如果按照教程不切换为root用户直接export环境变量再执行peer chaincode invoke会报错MSP相关KeyMaterial not found in SigningIdentityInfo，如果加sudo执行会显示peer命令不存在或者仍然报错，这里需要先切换为root用户再export环境变量再执行peeer chaincode invoke就会成功了！！！（搞了一下午）
sudo bash
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'

# 查看Assets
peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'
# 转移Assets
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"TransferAsset","Args":["asset6","Christopher"]}'

# 重新export环境变量
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
# 查看org2上的assets
peer chaincode query -C mychannel -n basic -c '{"Args":["ReadAsset","asset6"]}'
```



sudo bash命令以root用户打开一个终端 ，需要的是zpwang的密码

还有个sudo su也可也进入root用户，直接在同一终端切换用户

su root需要root用户的密码，Ubuntu中root用户的密码是随机的，每次开机自动有新密码，可以用`sudo passwd`来设置

su zpwang切回指定用户 su不带参数默认为su root

```bash
# 下载go.mod中定义的运行代码需要的相关依赖到vendor文件夹中
# 使其在无网络环境下也可也运行
sudo GO111MODULE=on go mod vendor
cd ../../test-network
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
peer lifecycle chaincode package basic.tar.gz --path ../asset-transfer-basic/chaincode-go/ --lang golang --label basic_1.0
# 上述操作可以将链码打包成basic.tar.gz
# 接下来需要将打包后的链码部署到所有需要背书的组织相关的peer上
export PATH=${PWD}/../bin:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
# 这里可以领悟到在测试网络中操作peer的方式是通过控制PATH环境变量来切换
peer lifecycle chaincode install basic.tar.gz
# 这里发现又出现了msp不存在的情况，但是无法用曾经的手段解决
# 发现这里的原因是需要重新启动网络 ./network.sh up 再建立一下channel(不清楚是否必要) 这样就可以成功执行 成功submitInstallProposal

#接下来在Org2上面安装chaincode
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051
# 再次安装
peer lifecycle chaincode install basic.tar.gz

# 接下来通过packageID来控制在peer上install的chaincode
# 可以通过下面的命令查看在peer上安装的chaincode的packageID
peer lifecycle chaincode queryinstalled

# 先对Org2进行操作
export CC_PACKAGE_ID=basic_1.0:69b8a2e2397a9d80e59b51b2ec71cb523f064ea5b184336a413b5ea876a957c1
# 执行以下命令使Org2 approve chaincode , approve操作是以组织为单位的，Org2中的一个peer运行了该命令，其他peer也都会通过分布式系统的gossip approve这个chaincode
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --package-id $CC_PACKAGE_ID --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
# sequence 1 表示该chaincode被定义或被更新的次数总和
# approveformyorg 还可以加参数 --signature-policy , --channel-config-policy 
# endorsement policy可以更改这里的配置

# 之后切换PATH到Org1，再次执行 approveformyorg 
```

![image-20231016194522092](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20231016194522092.png)

approve操作 返回了一个txid 这说明peer安装了chaincode并approve这个智能合约也会提交一个tx给Order，因为这样才能将某Org同意了该chaincode的信息传递给其他Org

```bash
# 可以在某个peer上运行该命令 checkcommitreadiness 来查看各Org 对 chaincode 的同意情况
peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json
```

![image-20231016200401433](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20231016200401433.png)

```bash
# 同意数达到了defination的要求， 则任一Org中的某个Peer可以commit该chaincode to channel
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name basic --version 1.0 --sequence 1 --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
```

![image-20231016200816203](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20231016200816203.png)

```bash
# chaincode defination的上述背书，被channel members提交给Order，Order将该TX放进block分发给channel中的peers,交给peers进行Validation
# 上述执行的commit命令会等待Validation的结果，可以通过querycommitted命令查询Validation结果
peer lifecycle chaincode querycommitted --channelID mychannel --name basic --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"
```



```bash
# 现在chaincode等待被client applications invoke
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n basic --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'

# query Assets
peer chaincode query -C mychannel -n basic -c '{"Args":["GetAllAssets"]}'
```

***到这里相当于拆解了一开始 deployCC命令中打包起来的部署智能合约的具体过程***

Sample application :  makes call to blockchain network

Smart contract : implements the transactions that interact with the ledger

![image-20231017125325033](C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20231017125325033.png)





Fabric Gateway client API      这是什么一个概念  相当于blockchain network对外开放的一个代理人，跟网络中的网关差不多，Client Application通过这个Gateway去提交TX到blockchain network或者查询ledger



