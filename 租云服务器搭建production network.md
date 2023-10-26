# 租云服务器搭建production network

## 云服务器

### 华为云

#### 配置ssh登录云服务器

```bash
ssh root@127.0.0.1 
# root用户下容易破坏系统，创建新用户并授予sudo权限
adduser zpwang
usermod -aG sudo zpwang 
# 退出云服务器,在宿主系统中配置zpwang用户的免密登录
cd .ssh
vim config
# 输入以下内容
# Host HwCloud
#        HostName 60.204.247.86
#        User zpwang
ssh HwCloud # 登录云服务器

# 设置云服务器上zpwang用户的免密登录
ssh-keygen # 生成rsa公私钥对，一路回车
cd .ssh
cat id_rsa.pub # 复制公钥
ssh HwCloud    # 登录云服务器
ls -a # 查看有无.ssh文件夹，没有则新建一个 
mkdir .ssh
cd .ssh
vim authorized_keys # 将公钥复制到里面，如果要设置多个主机可以免密登录，每个公钥需要用空行隔开
# 上述复制公钥的步骤也可以在宿主系统上 运行下面命令一步完成
ssh-copy-id HwCloud
```

#### 配置云服务器环境（tmux,docker）

[docker官方Ubuntu安装教程](https://docs.docker.com/engine/install/ubuntu/)

按照上述教程执行命令，2023.10.25命令如下：

```bash
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Add the repository to Apt sources:
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

# To install the latest version, run:
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
# Verify that the Docker Engine installation is successful by running the hello-world image.
docker --version
docker-compose --version
sudo docker run hello-world
# 下面两行命令似乎不需要也已经启动了docker
# sudo systemctl start docker
# sudo systemctl enable docker
```

假设一个云服务器跑一个Org，一个云服务器上开多个docker，每个docker运行一个节点



[docker官方文档](https://docs.docker.com/engine/install/linux-postinstall/)

配置docker的权限

```bash
sudo groupadd docker
sudo usermod -aG docker $USER
newgrp docker
```

## Docker Image

```bash
docker pull ubuntu:20.04
docker images
docker rmi ubuntu:20.04
docker image rm ubuntu:20.04
# 保存docker image到本地,注意需要给ubuntu_20_04.tar文件加可读权限 chmod +r ubuntu_20_04.tar
docker save -o ubuntu_20_04.tar ubuntu:20.04
docker load -i ubuntu_20_04.tar
# 将一个云服务器上的docker image复刻到另一个云服务器
scp server1:ubuntu_20_04.tar .	# 先复制到本地
scp ubuntu_20_04.tar server2:	# 传到另一个云服务器
```

## Docker Container

```bash
docker create -it ubuntu:20.04
docker ps -a # 显示所有容器
docker ps # 显示所有正在运行的容器
# 进入运行
docker start CONTAINER_ID
# 重启
docker restart CONTAINER_ID
# 停止运行
docker stop CONTAINER_ID
docker run -itd CONTAINER_ID # 创建并启动一个容器
docker run -it CONTAINER_ID  # 创建并启动并进入一个容器
# 上述指令创建好了容器
docker attach CONTAINER_NAME # 进入一个正在运行的容器，NAME可以通过 docker ps 命令查看 
# 如何退出容器？ 先按 ctrl+p(挂起容器)， 再按 ctrl+q(退出容器)
# ctrl+d (退出并关闭容器)

# 在容器已经处于运行状态后，可在容器之外令容器内的系统运行命令
docker exec CONTAINER_NAME COMMAND
# 删除容器,只能删除一个不处于运行状态的容器
docker rm CONTAINER_NAME 
# 删除所有容器
docker contrainer prune

# 除了Image可以保存成文件进行转移之外，contrainer也可以导出为一个Image，但是压缩包名称应写为.rar
docker export -o ubuntu_20_04.rar CONTAINER_NAME
# 将.tar中的image导出，导出得到的也是一个Image而不是Contrainer , 可以对Image进行命名
docker import ubuntu_20_04.tar [xxx:xxx] 

# 可以在容器之外调用容器内运行的系统执行命令
docker top CONTAINER_NAME
# 查看所有容器对于系统资源的占用信息
docker stats
# 可以在寄主系统和docker container中运行的系统之间传输文件 且复制文件夹时不需要 -r 
docker cp src dst 
# 如docker cp tmp.rar CONTAINER_NAME:/root
# 如docker cp CONTAINER_NAME:/root .
# 限制CONTAINER的资源 都在docker update命令中
docker update CONTAINER_NAME --memory 1000M
```

## 配置云服务器中的运行在docker container中的Image

下面是在云服务器中配置一个云服务器的操作，并不一定需要

```bash
# 进入AC Terminal
scp /var/lib/acwing/docker/images/docker_lesson_1_0.tar server_name:  # 将镜像上传到自己租的云端服务器
ssh server_name  # 登录自己的云端服务器
docker load -i docker_lesson_1_0.tar  # 将镜像加载到本地
# 创建并运行docker_lesson:1.0镜像,将容器的22端口映射到本地20000端口，因为本地的22端口已经用于ssh登录
docker run -p 20000:22 --name my_docker_server -itd docker_lesson:1.0  
docker attach my_docker_server  # 进入创建的docker容器
passwd  # 设置root密码 
# 去云平台控制台->更多->网络和安全组->安全组配置->配置规则->手动添加->(自定义TCP 20000 源0.0.0.0/0)
# 返回AC Terminal，即可通过ssh登录自己的docker容器：
ssh root@xxx.xxx.xxx.xxx -p 20000  # 将xxx.xxx.xxx.xxx替换成自己租的服务器的IP地址
# 然后，可以仿照上节课内容，创建工作账户acs。
# 最后，可以参考4. ssh——ssh登录配置docker容器的别名和免密登录。
```

按如下步骤制作***Fabric***镜像：

```bash
# 假设刚开始在本地系统
ssh HwCloud
docker ps -a # 查看所有container
# 我们大致需要三种节点，分别制作每一节点需要的Image模板
docker run --name fabric_peer_tmp -itd ubuntu:22.04
docker run --name fabric_order_tmp -itd ubuntu:22.04
docker run --name fabric_ca_tmp -itd ubuntu:22.04
# 下面以fabric_peer_tmp为例
docker attach fabric_peer_tmp
adduser admin
usermod -aG sudo admin
apt-get install sudo
sudo apt-get install tmux
```

```bash
# root用户下
apt-get update
apt-get install sudo # docker跑的ubuntu:22.04镜像中没有sudo,华为云上租下来的服务器中自带sudo
# 重复官网安装docker的步骤
sudo apt-get install ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
docker --version
docker-compose --version
# 补充两行命令启动docker daemon,这里的docker服务依赖于云服务器的docker,云服务器中的docker用sudo systemctl start
sudo service docker start
sudo docker run hello-world

# 不一定可靠安装方式：上述部分只装了docker没装docker-compose，也可以试试下面这部分fabric官方教程提供的命令
sudo apt-get install curl
sudo apt-get -y install docker-compose
sudo usermod -aG docker <username>
newgrp docker
sudo service docker start
```





该思路存在一部分问题，docker in docker不是一种优美的方式，kubo可以在
