#!/bin/bash
supported_os=("linux" "darwin" "windows")
supported_arch=("amd64")
release_root=../../../../../dev/cloudor/releases
for os in "${supported_os[@]}"
do
    for arch in "${supported_arch[@]}"
    do
        CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -o $release_root/$os/$arch/latest/cloudor main.go
        echo "build to $release_root/$os/$arch/latest/cloudor"
    done        

done