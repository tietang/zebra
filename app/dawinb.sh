#!/usr/bin/env bash
echo $0 $1
VERSION=
cpath=`pwd`
echo $cpath
PROJECT_PATH=${cpath%src*} #从右向左截取第一个 src 后的字符串
echo ${PROJECT_PATH}


export GOPATH=$GOPATH:${PROJECT_PATH}

SOURCE_FILE_NAME=main
TARGET_FILE_NAME=zebra



build(){
   echo $GOOS $GOARCH
   env  GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${TARGET_FILE_NAME}_${GOOS}_${GOARCH}-${VERSION}${EXT} -v ${SOURCE_FILE_NAME}.go
   mv ${TARGET_FILE_NAME}_${GOOS}_${GOARCH}-${VERSION}${EXT} ${TARGET_FILE_NAME}${EXT}
   if [ ${GOOS} == "windows" ]; then
       zip ${TARGET_FILE_NAME}_${GOOS}_${GOARCH}-${VERSION}.zip ${TARGET_FILE_NAME}${EXT} proxy.ini
   else
       tar -czvf ${TARGET_FILE_NAME}_${GOOS}_${GOARCH}-${VERSION}.tar.gz ${TARGET_FILE_NAME}${EXT} proxy.ini start.sh stop.sh
   fi
   mv  ${TARGET_FILE_NAME}${EXT} ${TARGET_FILE_NAME}_${GOOS}_${GOARCH}-${VERSION}${EXT}
}

CGO_ENABLED=0


# mac osx
GOOS=darwin
GOARCH=amd64
build