#!/bin/bash

echo "create the tool's name !"

read name

echo "start build x64 !!!"

cd ../src

echo -n "windows ? (y,n)"
read yn
if [ $yn == "y" ]
then
    echo "start build windows x64"
    GOOS=windows GOARCH=amd64 go build -o ../build/${name}_windows_x64.exe
else
    echo "no"

fi

echo -n "linux ? (y,n)"
read yn
if [ $yn == "y" ]
then
    echo "start build linux x64"
    GOOS=linux GOARCH=amd64 go build -o ../build/${name}_linux_x64
else
    echo "no"

fi

echo -n "mac ? (y,n)"
read yn
if [ $yn == "y" ]
then
    echo "start build mac x64"
    GOOS=darwin GOARCH=amd64 go build -o ../build/${name}_mac_x64
else
    echo "no"

fi


echo -n "need x86 version. (y,n)"
read yn

if [ $yn == "y" ]
then
    GOOS=windows GOARCH=386 go build -o ../build/${name}_windwos_x86.exe
    GOOS=linux GOARCH=386 go build -o ../build/${name}_linux_x86
    GOOS=darwin GOARCH=386 go build -o ../build/${name}_mac_x86
else
    echo "no"
fi

echo "build finish"