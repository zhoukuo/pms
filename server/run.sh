#!/bin/bash

# Program to deploy pms server

timestamp=`date +"%Y/%m/%d %H:%M:%S"`
db="$GOPATH/bin/pms.db"

if [ ! -n $GOPATH ]; then
	echo "$GOPATH is NULL, please config golang sdk first!"
	exit
fi

if [ ! -d "$GOPATH/src" ]; then
	echo $timestamp "mkdir src ..."
	mkdir $GOPATH/src
	cd $GOPATH/src
fi

if [ ! -d "$GOPATH/src/pms" ]; then
	echo $timestamp 'checkout source code from svn ...'
	svn checkout http://192.168.131.32/qmd/QMDb/21QDev/PMS/src/server  --username zhoukuo --password 111111
	mv server pms
	cd pms
else
	echo $timestamp 'SVN update ...'
	cd $GOPATH/src/pms
	svn update
fi

echo $timestamp 'Terminate Server ...'
ps -ef |grep -v grep |grep './pms' |cut -c 10-15 |xargs kill 9

echo $timestamp 'Build main.go ...'
go install pms/main

# -f 参数判断 $db 是否存在
echo $timestamp 'Database(pms.db) not exist, generate a new one ...'
cd database;./initdb.sh
cd ..
mv pms.db $db

echo $timestamp 'Start Server ...'
cd $GOPATH/bin
mv main pms
nohup ./pms &

echo -e "=======================================================================================\n"
echo $timestamp "Running Static Code Analysis ..."
echo -e "=======================================================================================\n"

cd $GOPATH/src/pms/business/
go vet ./...

echo -e "=======================================================================================\n"
echo $timestamp "Running Unit Testing ..."
echo -e "=======================================================================================\n"
cd $GOPATH/src/pms/test/api/
rm -fr newman
newman -c cases.json -e env.json -H HTML


echo -e "=======================================================================================\n"
