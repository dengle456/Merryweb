#!/usr/bin/env bash
frps="150.158.143.171"
port="10010"
token="ssD790416"
echo "1.frp安装install"
echo "2.查看配置并且修改"
echo "3.后台启动frp"
echo "4.进程结束frp"
echo "5.查看日志状态"
read id

if [ $id = "1" ]; then
yum install wget
wget "https://github.com/fatedier/frp/releases/download/v0.39.1/frp_0.39.1_linux_amd64.tar.gz"
tar -zxvf frp_0.39.1_linux_amd64.tar.gz
rm -rf frp_0.39.1_linux_amd64.tar.gz
cd frp_0.39.1_linux_amd64
clear
echo "frps ip地址:"
read ip
sed -r -i "s/(server_addr).*/ \1= $ip/g" frpc.ini
echo "frps 监听端口:"
read ip
sed -r -i "s/(server_port).*/\1= $ip/g" frpc.ini
echo "token:"
read ip
sed -i "3a\token = $ip" frpc.ini
nohup ./frpc -c ./frpc.ini &
fi

if [ $id = "2" ]; then
cd frp_0.39.1_linux_amd64
vim frpc.ini
fi

if [ $id = "3" ]; then
cd frp_0.39.1_linux_amd64
nohup ./frpc -c frpc.ini &
fi

if [ $id = "4" ]; then
cd frp_0.39.1_linux_amd64
pid=`ps -ef|grep frpc|awk 'NR==1{print $2}'`
echo $path
echo `ps -ef|grep frpc`
echo "结束进程: $pid"
kill $pid
fi

if [ $id = "5" ];then
cd frp_0.39.1_linux_amd64
tail -fn 10 nohup.out
fi


if [ $id = "6" ]; then
yum install wget
wget "https://github.com/fatedier/frp/releases/download/v0.39.1/frp_0.39.1_linux_amd64.tar.gz"
tar -zxvf frp_0.39.1_linux_amd64.tar.gz
rm -rf frp_0.39.1_linux_amd64.tar.gz
cd frp_0.39.1_linux_amd64
sed -r -i "s/(server_addr).*/\1= $frps/g" frpc.ini
sed -r -i "s/(server_port).*/\1= $port/g" frpc.ini
sed -i "3a\token = $token" frpc.ini
nohup ./frpc -c ./frpc.ini &
echo "远程地址:$frps:$port"
tail -fn 10 nohup.out
fi
