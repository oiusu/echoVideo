#!/usr/bin/env bash
FWDIR="$(cd "`dirname "$0"`"; pwd)"
echo $FWDIR
cd $FWDIR


prodPath=/home/chenc/apps/echoVideo
user_name=chenc
port=22
#port=10001
#toPath=${prodPath}/upload/oss
#
GOOS=linux GOARCH=amd64 go build -o echoVideo  main.go
ls | grep echoVideo
echo "go build finish..."


HOST=192.168.55.98
# 台式机
#HOST=172.20.100.44

#for HOST in {211.159.149.73}
#也可以写成for element in ${array[*]}
#do
echo "deploy $HOST $port "


ssh -P ${port} ${user_name}@${HOST} "mkdir -p ${prodPath}/view; "
echo "mkdir from server..."
scp  -P ${port}  echoVideo  ${user_name}@${HOST}:${prodPath}
scp  -P ${port}   ./view/tpl.html  ${user_name}@${HOST}:${prodPath}/view
echo "scp finish..."

#done



# nohup  ./xvideo-cloud-go  -model_server=10.163.20.154:8500  &
# nohup  ./xvideo-cloud-go  -model_server=10.163.20.154:8500 -deploy=prod x &
