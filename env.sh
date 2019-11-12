#!/usr/bin/env bash
#cd ..
#npm run build
#cd brun
#VERSION=0.1
cpath=`pwd`
echo $cpath
PROJECT_PATH=${cpath%src*} #从右向左截取第一个 src 后的字符串
echo ${PROJECT_PATH}


export GOPATH=$GOPATH:${PROJECT_PATH}
export GOPROXY=https://mirrors.aliyun.com/goproxy/