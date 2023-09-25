# IPFS私有网络搭建

## 前言

IPFS具体是指由Protocal Lab提出的一个分布式存储协议，遵循该协议有多个实现版本，可在下面这个网址中查看：https://docs.ipfs.tech/concepts/ipfs-implementations/  其中Kubo（Go）和Helia(JS)是两个比较完善的实现版本，Kubo是第一个IPFS协议的实现版本（基于Go语言实现），官方对Kubo的描述如下：“ Generalist daemon oriented IPFS implementation with an extensive HTTP RPC API.” 此外，还有用于搭建Kubo节点集群工作ipfs-cluster

为了满足我们搭建私有网络的需求，并考虑到系统可靠性与稳定性，我们选择Kubo实现版本（目前I结合IPFS的相关论文中并未提及他们在搭建系统时使用了哪一个实现版本），其次Go语言也适用于分布式系统高并发的场景。

## Kubo安装

https://docs.ipfs.tech/install/command-line/#install-official-binary-distributions

打开上述网址是IPFS官方提供的安装Kubo的教程，安装方式主要有以下三类：

1. 安装官方提供的已编译的二进制文件压缩包，打开https://dist.ipfs.tech/#kubo查看对应的操作系统与CPU架构的二进制文件版本（该网址中还提供了一些IPFS工具的二进制文件）

   ```bash
   # 下面是以linux-amd64为例的安装过程
   wget https://dist.ipfs.tech/kubo/v0.22.0/kubo_v0.22.0_linux-amd64.tar.gz
   tar -xvzf kubo_v0.22.0_linux-amd64.tar.gz
   > x kubo/install.sh
   > x kubo/ipfs
   > x kubo/LICENSE
   > x kubo/LICENSE-APACHE
   > x kubo/LICENSE-MIT
   > x kubo/README.md
   
   cd kubo
   sudo bash install.sh
   > Moved ./ipfs to /usr/local/bin
   
   ipfs --version
   > ipfs version 0.22.0
   ```

2. 下载IPFS源代码文件，在本地进行编译，具体参考https://github.com/ipfs/kubo/blob/v0.22.0/README.md#build-from-source

3. 下载已经安装了IPFS的Docker，具体参考https://docs.ipfs.tech/install/run-ipfs-inside-docker/#set-up

此处最简易的安装方式是直接下载已编译的二进制文件，但如果需要对IPFS系统进行一些修改，可以选择第二种方式，修改后重新编译系统。

## 基本使用

以下基本使用过程都是以linux系统为例，Windows和MacOS的目录结构会有些许不同

```bash
ipfs init # 初始化本地ipfs节点，将会生成".ipfs"文件存储ipfs配置信息
```

![](.\pictures\ipfs_init.png)

```bash
ipfs id # 查看本地节点信息：如peer identity
```

在这里补充对于分布式的IPFS节点如何建立连接构建成分布式的IPFS网络：

```bash
ipfs bootstrap list 
```

正如该命令名称，这个命令将显示一系列的地址，这些地址是由IPFS开发社区提供的运行的IPFS节点，分布式的IPFS本地节点通过与bootstrap list中的初始节点建立连接，从而连入整个IPFS网络。

```bash
ipfs bootstrap rm --all # 该命令删除所有list中的节点
ipfs bootstrap add # 自定义增加想要连接的节点
```

删除了bootstrap list后，当我们启动本地节点时，便不会接入整个IPFS网络

```bash
ipfs daemon # 启动ipfs进程
ipfs swarm peers # 查看已经连接的其他节点
```

但是为了搭建私有网络，仅仅在bootstrap list中设置需要连接的节点这一方法并不稳妥，一旦其中一个节点与其他节点之间建立了连接，该网络便会拓展，我们需要一种更安全的机制

## 搭建限制进入的私网

建立私网引入了一个工具：ipfs-swarm-key，这是IPFS官方提供的工具，可以在github中获取 https://github.com/Kubuxu/go-ipfs-swarm-key-gen/tree/master

该工具基于Go实现，所以我们需要先安装Go，可以选择参考本篇知乎教程https://zhuanlan.zhihu.com/p/656048616 ，官网教程为[https://go.dev/doc/install](https://link.zhihu.com/?target=https%3A//go.dev/doc/install)

### Go Install

下载go压缩包[https://go.dev/dl/](https://link.zhihu.com/?target=https%3A//go.dev/dl/)

官方提供解压命令：

```bash
 rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
```

考虑到以前安装了go的其他版本，故先删除

之后要把bin路径进行声明，可以声明为system级别，也可以是用户级别

选择在/etc/profile 或 $HOME/.profile 文件中末尾增加一行

```bash
export PATH=$PATH:/usr/local/go/bin
```

此时.profile文件的更改尚未生效，如果是在Home中加的，再执行

```bash
source ~/.profile
```

或者

```bash
source $HOME/.profile
```

使得.profile文件内容生效

```bash
go version
```

验证安装成功

### 下载swarm-key-gen工具

go get方式已经不被新版的go所支持，新版本的go需要使用go install 命令，故我们需要下载go并参考下面这篇知乎中的内容进行 https://zhuanlan.zhihu.com/p/656058690 ，其中记录了探索使用该工具的过程。

其实go install和go get的功能相似，简单讲就是获取一个go实现的package，并将其编译，然后将二进制可执行文件加到go路径的$GO_PATH/bin目录下，之后便可以直接调用该功能

```bash
export PATH=$PATH:/usr/local/go/bin # 声明GO_PATH
```

任意目录下执行以下命令（把get改为install，`latest`使得go install下载的是该库的最新版本）

```bash
go install github.com/Kubuxu/go-ipfs-swarm-key-gen/ipfs-swarm-key-gen@latest 
```

在Home目录下找到go工作区，并调用下载的工具生成key文件

```bash
cd ~/go/bin
ipfs-swarm-key-gen > /home/zpwang/.ipfs/swarm.key
```

之后可以在 `/home/zpwang/.ipfs/` 目录下看到swarm.key文件，将该文件通过scp传输到每一个ipfs节点的`~/.ipfs`目录下。

### 声明环境变量修改ipfs启动模式

```bash
cd ~
vim .profile # 在该文件末尾增加下面一行声明
export LIBP2P_FORCE_PNET=1
source ~/.profile # 使环境变量生效
```

完成以上配置后，再以`ipfs daemon`命令启动进程时，将出现以下提示信息

```bash
Swarm is limited to private network of peers with the swarm key
Swarm key fingerprint: c2fc00b19ee671210674155a5cf76ee8
```

## ipfs-cluster

上述方法搭建的私有网络不便于对节点进行统一管理，若需要对集群进行管理，需要额外部署支持Kubo实现版本的ipfs-cluster工具，该工具主要有三部分：

- ipfs-cluster-ctl
- ipfs-cluster-follow
- ipfs-cluster-service
