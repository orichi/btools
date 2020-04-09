#!/bin/bash
if [ -e service.zip ];then
  `rm -f service.zip`
  echo "删除旧的版本版本"
fi

if [ ! -e main ];then
  echo "请先构建项目"
  exit
fi

if [ ! -d service ];then
  `mkdir service`
fi

if [ ! -d service/public ];then
  `mkdir service/public`
fi

cp main service/m3uservice
cp -a templates service/
cp config.ini service/
cp config.example.ini service/

zip service.zip -r service
