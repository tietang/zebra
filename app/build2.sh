#!/usr/bin/env bash
VERSION=0.2

cpath=`pwd`
echo $cpath
PROJECT_PATH=${cpath%src*} #从右向左截取第一个 src 后的字符串
echo ${PROJECT_PATH}


export GOPATH=$GOPATH:${PROJECT_PATH}

SOURCE_FILE_NAME=main
TARGET_FILE_NAME=zebra

go install github.com/tietang/zebra

