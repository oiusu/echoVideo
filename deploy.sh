#!/usr/bin/env bash
FWDIR="$(cd "`dirname "$0"`"; pwd)"
echo $FWDIR
cd $FWDIR


prodPath=/home/chenc/apps/echoVideo
user_name=root
port=22
#toPath=${prodPath}/upload/oss
#
GOOS=linux GOARCH=amd64 go build -o echoVideo  main.go
ls | grep echoVideo
echo "go build finish..."


HOST=192.168.55.98

#for HOST in {211.159.149.73}
#也可以写成for element in ${array[*]}
#do
echo "deploy $HOST"

ssh  ${user_name}@${HOST} "mkdir -p ${prodPath}/view;mkdir -p ${prodPath}/echoVideo; "
echo "mkdir from server..."
scp   echoVideo  ${user_name}@${HOST}:${prodPath}
scp     ./view/tpl.html  ${user_name}@${HOST}:${prodPath}/view
echo "scp finish..."

#done



# nohup  ./xvideo-cloud-go  -model_server=10.163.20.154:8500  &
# nohup  ./xvideo-cloud-go  -model_server=10.163.20.154:8500 -deploy=prod x &
