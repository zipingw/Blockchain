# Go Experience

```bash
# 配置代理，依赖源 https://proxy.golang.org在国内无法访问
echo "export GO111MODULE=on" >> ~/.profile
echo "export GOPROXY=https://goproxy.cn" >> ~/.profile
source ~/.profile
# 或者可以不修改proxy 连接VPN再修改go mod tidy 的timeout时间
```

