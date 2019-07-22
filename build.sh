#!/usr/bin/env bash

export GOPATH=$GOPATH:/root/DioClient

export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64
go build -o DioClient_${GOOS}_${GOARCH}

PROJECT_DIR=$(pwd)
PACK_DIR=${PROJECT_DIR}/pack
mkdir -p ${PACK_DIR}/DioClient

cd ${PROJECT_DIR}
cp -R views data assets config logs DioClient_${GOOS}_${GOARCH} ${PACK_DIR}/DioClient
cd ${PACK_DIR}
tar cvf DioClient_${GOOS}_${GOARCH}.tar DioClient
gzip DioClient_${GOOS}_${GOARCH}.tar
rm -rf ${PACK_DIR}/DioClient

