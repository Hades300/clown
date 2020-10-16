 # Clown
这是一个运行在局域网内的工具。预期实现功能

## TODO

- [x] 主动ARP攻击
- [ ] 将受害机流量导向指定主机
- [ ] DHCP服务器伪造
  - [ ] DNS劫持
- [ ] 流量导出

## 源码编译

需要go编译环境、c编译器、libpcap开发库

```
# Fedora
sudo dnf install -y libpcap-devel

# Debian/Ubuntu
sudo apt-get install -y libpcap-dev

# OSX
brew install libpcap

# FreeBSD
sudo pkg install libpcap

# Windows
# Install https://www.winpcap.org/
```

安装完libpcap后，使用`go install github.com/hades300/clown`进行安装

```
go install github.com/hades300/clown
clown -h	
```

命令示例

```
clown pretend -i [interface] [target-ip] [pretend-ip]
clown pretend 192.168.3.11 192.168.3.1
clown pretend -i en2 192.168.3.11 192.168.3.1
```

查看本地网卡接口名

```
# for darwin 
ifconfig 
# for linux
ip address
```

