# SSH使用

```bash
# ssh是客户端 sshd是服务端

# 安装 sshd
sudo apt-get install sshd # Unable to locate package sshd

# 安装 ssh
sudo apt-get install openssh-server # succeed
# 启动ssh服务
sudo service ssh start
> * Starting OpenBSD Secure Shell server sshd                      [ OK ]
# 查看ssh服务状态
> * sshd is running  # 为什么启动了ssh但是sshd在running
```





```bash
# Windows 查看 ip 地址
ipconfig
# 这里windows主机的ip为 192.168.137.1
# 从Windows看 wsl中运行的系统的 ip 为 172.22.80.1
```

<img src="C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20231018113325300.png" alt="image-20231018113325300" style="zoom:50%;" />



```bash
# WSL 中 Linux 查看 ip 地址
ip address
# WSL中观察到ip 为 172.22.95.30
```

<img src="C:\Users\zipin\AppData\Roaming\Typora\typora-user-images\image-20231018113413550.png" alt="image-20231018113413550" style="zoom: 50%;" />



```bash
# 从WSL中的Ubuntu通过scp传文件给Windows
# 先查看Windows用户名： 控制面板=>用户账户=>用户账户 王子平 Administrator
scp ./ipfscc.go administrator@192.168.137.1:/d:/ipfscc.go
# failed
```

